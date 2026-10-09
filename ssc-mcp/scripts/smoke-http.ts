/**
 * HTTP + stdio transport smoke test.
 *
 * Starts the MCP server in HTTP mode, performs the MCP handshake, and lists tools
 * over Streamable HTTP and legacy SSE. Then spawns the default stdio server and
 * lists tools again. Does not call the SSC API and does not need credentials.
 */
import { fileURLToPath } from "node:url";
import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { SSEClientTransport } from "@modelcontextprotocol/sdk/client/sse.js";
import { StdioClientTransport } from "@modelcontextprotocol/sdk/client/stdio.js";
import { StreamableHTTPClientTransport } from "@modelcontextprotocol/sdk/client/streamableHttp.js";
import { SscClient } from "../src/client.js";
import { loadMcpRuntimeConfig } from "../src/config.js";
import { startHttpServer, type HttpServerHandle } from "../src/http.js";

const EXPECTED_TOOLS = [
  "ssc_login",
  "ssc_status",
  "ssc_list_menus",
  "ssc_list_orders",
  "ssc_get_order",
  "ssc_list_dicts",
  "ssc_get_dict",
  "ssc_get_dict_map",
  "ssc_create_dict",
  "ssc_update_dict",
  "ssc_delete_dict",
  "ssc_list_excel_read_rules",
  "ssc_get_excel_read_rule",
  "ssc_list_excel_rows",
  "ssc_get_excel_row",
  "ssc_update_excel_row",
  "ssc_delete_excel_row",
  "ssc_list_export_rules",
  "ssc_get_export_rule",
  "ssc_list_export_template_options",
] as const;

const root = fileURLToPath(new URL("..", import.meta.url));

function assert(condition: unknown, message: string): asserts condition {
  if (!condition) throw new Error(message);
}

function textContent(content: unknown): string {
  if (!Array.isArray(content)) return "";
  return content
    .map((item) => {
      if (item && typeof item === "object" && "text" in item && typeof item.text === "string") {
        return item.text;
      }
      return "";
    })
    .join("\n");
}

function checkConfig(): void {
  const defaults = loadMcpRuntimeConfig([], {});
  assert(defaults.transport === "stdio", "default transport should be stdio");
  assert(defaults.host === "0.0.0.0", `default host should be 0.0.0.0, got ${defaults.host}`);
  assert(defaults.port === 3100, `default port should be 3100, got ${defaults.port}`);
  assert(defaults.corsOrigin === "*", "default CORS origin should be *");

  const fromEnv = loadMcpRuntimeConfig([], {
    SSC_MCP_TRANSPORT: "http",
    SSC_MCP_HOST: "127.0.0.1",
    SSC_MCP_PORT: "3201",
    SSC_MCP_CORS_ORIGIN: "https://allowed.example",
  });
  assert(fromEnv.transport === "http", "env transport");
  assert(fromEnv.host === "127.0.0.1", "env host");
  assert(fromEnv.port === 3201, "env port");
  assert(fromEnv.corsOrigin === "https://allowed.example", "env cors");

  const cliOverrides = loadMcpRuntimeConfig(
    ["--transport", "stdio", "--port=0", "--host", "127.0.0.1"],
    { SSC_MCP_TRANSPORT: "http", SSC_MCP_PORT: "9999" },
  );
  assert(cliOverrides.transport === "stdio", "CLI should override env transport");
  assert(cliOverrides.port === 0, "CLI port");
  assert(cliOverrides.host === "127.0.0.1", "CLI host");

  let rejected = false;
  try {
    loadMcpRuntimeConfig(["--transport", "websocket"], {});
  } catch {
    rejected = true;
  }
  assert(rejected, "unknown transport should fail");
  console.log("config: OK (stdio default, http via env/CLI)");
}

async function listToolNames(label: string, connect: () => Promise<Client>): Promise<void> {
  const mcp = await connect();
  try {
    const listed = await mcp.listTools();
    const names = new Set(listed.tools.map((tool) => tool.name));
    const missing = EXPECTED_TOOLS.filter((name) => !names.has(name));
    assert(
      missing.length === 0,
      `${label} missing tools: ${missing.join(", ")} (got ${listed.tools.length})`,
    );
    console.log(`${label}: handshake OK, tools = ${listed.tools.length}`);
  } finally {
    await mcp.close().catch(() => undefined);
  }
}

async function connectHttp(url: string, mode: "streamable" | "sse"): Promise<Client> {
  const mcp = new Client({ name: "ssc-mcp-smoke", version: "1.0.0" });
  const transport =
    mode === "streamable"
      ? new StreamableHTTPClientTransport(new URL(url))
      : new SSEClientTransport(new URL(url));
  await mcp.connect(transport);
  return mcp;
}

async function checkToolCall(mcpUrl: string): Promise<void> {
  const mcp = await connectHttp(mcpUrl, "streamable");
  try {
    const result = await mcp.callTool({ name: "ssc_status", arguments: {} });
    assert(result.isError === true, "ssc_status without credentials should return isError");
    const text = textContent(result.content);
    assert(
      text.includes("Not authenticated"),
      `ssc_status body did not describe missing auth: ${text}`,
    );
    console.log("streamable ssc_status: OK (tool dispatched, no SSC credentials)");
  } finally {
    await mcp.close().catch(() => undefined);
  }
}

async function checkCors(handle: HttpServerHandle): Promise<void> {
  const preflight = await fetch(handle.mcpUrl, {
    method: "OPTIONS",
    headers: {
      Origin: "https://example.test",
      "Access-Control-Request-Method": "POST",
      "Access-Control-Request-Headers": "content-type, mcp-session-id",
    },
  });
  assert(preflight.status === 204, `CORS preflight status ${preflight.status}`);
  assert(
    preflight.headers.get("access-control-allow-origin") === "*",
    "default CORS should allow any origin",
  );
  const allowHeaders = (preflight.headers.get("access-control-allow-headers") ?? "").toLowerCase();
  assert(allowHeaders.includes("mcp-session-id"), `CORS allow-headers missing session id: ${allowHeaders}`);
  assert(allowHeaders.includes("content-type"), "CORS allow-headers missing content-type");
  console.log("CORS preflight: OK (*)");

  const restricted = await startHttpServer({
    client: new SscClient({ baseUrl: "https://ssc.mengfei.tech" }),
    host: "127.0.0.1",
    port: 0,
    corsOrigin: "https://allowed.example",
  });
  try {
    const denied = await fetch(restricted.mcpUrl, {
      method: "OPTIONS",
      headers: {
        Origin: "https://evil.example",
        "Access-Control-Request-Method": "POST",
      },
    });
    assert(
      denied.headers.get("access-control-allow-origin") == null,
      "unlisted origin should not be echoed",
    );
    const allowed = await fetch(restricted.mcpUrl, {
      method: "OPTIONS",
      headers: {
        Origin: "https://allowed.example",
        "Access-Control-Request-Method": "POST",
      },
    });
    assert(
      allowed.headers.get("access-control-allow-origin") === "https://allowed.example",
      "listed origin should be echoed",
    );
    console.log("CORS allow-list: OK");
  } finally {
    await restricted.close();
  }
}

function stringEnv(source: NodeJS.ProcessEnv): Record<string, string> {
  const env: Record<string, string> = {};
  for (const [key, value] of Object.entries(source)) {
    if (typeof value === "string") env[key] = value;
  }
  delete env.SSC_MCP_TRANSPORT;
  delete env.SSC_MCP_HOST;
  delete env.SSC_MCP_PORT;
  delete env.SSC_MCP_CORS_ORIGIN;
  delete env.SSC_MCP_ALLOWED_HOSTS;
  return env;
}

async function checkStdio(): Promise<void> {
  let stderr = "";
  const transport = new StdioClientTransport({
    command: process.execPath,
    args: ["node_modules/tsx/dist/cli.mjs", "src/index.ts"],
    cwd: root,
    env: stringEnv(process.env),
    stderr: "pipe",
  });
  transport.stderr?.on("data", (chunk: Buffer | string) => {
    stderr += String(chunk);
  });
  const mcp = new Client({ name: "ssc-mcp-smoke-stdio", version: "1.0.0" });
  try {
    await mcp.connect(transport);
    const listed = await mcp.listTools();
    const names = new Set(listed.tools.map((tool) => tool.name));
    const missing = EXPECTED_TOOLS.filter((name) => !names.has(name));
    assert(missing.length === 0, `stdio missing tools: ${missing.join(", ")}`);
    assert(
      !stderr.includes("HTTP listening"),
      "default stdio process should not start the HTTP server",
    );
    console.log(`stdio: handshake OK, tools = ${listed.tools.length}`);
  } catch (err) {
    if (stderr.trim()) console.error("stdio stderr:\n", stderr);
    throw err;
  } finally {
    await mcp.close().catch(() => undefined);
  }
}

async function main(): Promise<void> {
  checkConfig();
  const handle = await startHttpServer({
    client: new SscClient({ baseUrl: "https://ssc.mengfei.tech" }),
    host: "0.0.0.0",
    port: 0,
    corsOrigin: "*",
  });
  console.log(`HTTP listening ${handle.mcpUrl} (bound ${handle.host})`);
  try {
    await listToolNames("streamable", () => connectHttp(handle.mcpUrl, "streamable"));
    await checkToolCall(handle.mcpUrl);
    await listToolNames("sse", () => connectHttp(handle.sseUrl, "sse"));
    await checkCors(handle);
  } finally {
    await handle.close();
  }
  await checkStdio();
  console.log("SMOKE HTTP: PASS");
}

main().catch((err: unknown) => {
  console.error("SMOKE HTTP: FAIL", err);
  process.exit(1);
});
