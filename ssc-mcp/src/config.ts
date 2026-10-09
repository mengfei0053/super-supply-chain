/**
 * Runtime transport settings for the SSC MCP server.
 *
 * Precedence: CLI flags, then environment variables, then defaults.
 * stdio is the default so existing Cursor configs keep working.
 */

export type TransportMode = "stdio" | "http";

export type McpRuntimeConfig = {
  transport: TransportMode;
  host: string;
  port: number;
  /**
   * CORS allow-list. "*" allows any browser origin.
   * Otherwise a comma-separated list of origins (for example
   * "http://localhost:5173,https://cursor.com").
   */
  corsOrigin: string;
  /** Optional Host allow-list passed to the SDK DNS-rebinding middleware. */
  allowedHosts?: string[];
  help: boolean;
};

const DEFAULT_HOST = "0.0.0.0";
const DEFAULT_PORT = 3100;
const DEFAULT_CORS = "*";

export function loadMcpRuntimeConfig(
  argv: string[] = process.argv.slice(2),
  env: NodeJS.ProcessEnv = process.env,
): McpRuntimeConfig {
  const flags = parseArgs(argv);
  const transport = parseTransport(
    stringFlag(flags.transport) ?? env.SSC_MCP_TRANSPORT,
  );
  const host =
    stringFlag(flags.host) ?? env.SSC_MCP_HOST?.trim() ?? DEFAULT_HOST;
  const port = parsePort(stringFlag(flags.port) ?? env.SSC_MCP_PORT);
  const corsOrigin =
    stringFlag(flags["cors-origin"]) ??
    env.SSC_MCP_CORS_ORIGIN?.trim() ??
    DEFAULT_CORS;
  const allowedHosts = parseHostList(
    stringFlag(flags["allowed-hosts"]) ?? env.SSC_MCP_ALLOWED_HOSTS,
  );

  if (!host) {
    throw new Error("SSC_MCP_HOST / --host must be a non-empty hostname");
  }

  return {
    transport,
    host,
    port,
    corsOrigin,
    allowedHosts,
    help: flags.help === true,
  };
}

export function printHelp(): void {
  console.error(`ssc-mcp — Super Supply Chain MCP server

Usage:
  ssc-mcp [--transport stdio|http] [--host HOST] [--port PORT]
          [--cors-origin ORIGIN] [--allowed-hosts host1,host2]

Defaults:
  transport     stdio   (env SSC_MCP_TRANSPORT)
  host          ${DEFAULT_HOST}   (env SSC_MCP_HOST)
  port          ${DEFAULT_PORT}      (env SSC_MCP_PORT)
  cors-origin   ${DEFAULT_CORS}       (env SSC_MCP_CORS_ORIGIN)
  allowed-hosts (unset) (env SSC_MCP_ALLOWED_HOSTS)

HTTP mode serves:
  Streamable HTTP   http://HOST:PORT/mcp     (recommended)
  Legacy SSE        http://HOST:PORT/sse
                    POST http://HOST:PORT/messages?sessionId=...
`);
}

function parseArgs(argv: string[]): Record<string, string | boolean> {
  const out: Record<string, string | boolean> = {};
  for (let i = 0; i < argv.length; i++) {
    const arg = argv[i];
    if (!arg) continue;
    if (arg === "--help" || arg === "-h") {
      out.help = true;
      continue;
    }
    if (!arg.startsWith("--")) {
      throw new Error(`Unexpected argument: ${arg}`);
    }
    const body = arg.slice(2);
    const eq = body.indexOf("=");
    if (eq !== -1) {
      out[body.slice(0, eq)] = body.slice(eq + 1);
      continue;
    }
    const next = argv[i + 1];
    if (next !== undefined && !next.startsWith("-")) {
      out[body] = next;
      i++;
    } else {
      out[body] = true;
    }
  }
  return out;
}

function stringFlag(value: string | boolean | undefined): string | undefined {
  if (typeof value !== "string") return undefined;
  const trimmed = value.trim();
  return trimmed.length > 0 ? trimmed : undefined;
}

function parseTransport(raw: string | undefined): TransportMode {
  const value = (raw ?? "stdio").trim().toLowerCase();
  if (value === "stdio" || value === "http") return value;
  throw new Error(
    `Invalid transport "${raw}". Use SSC_MCP_TRANSPORT=stdio|http or --transport stdio|http`,
  );
}

function parsePort(raw: string | undefined): number {
  if (raw === undefined || raw.trim() === "") return DEFAULT_PORT;
  const port = Number(raw);
  if (!Number.isInteger(port) || port < 0 || port > 65535) {
    throw new Error(
      `Invalid port "${raw}". SSC_MCP_PORT / --port must be an integer from 0 to 65535`,
    );
  }
  return port;
}

function parseHostList(raw: string | undefined): string[] | undefined {
  if (!raw || !raw.trim()) return undefined;
  const hosts = raw
    .split(",")
    .map((host) => host.trim())
    .filter((host) => host.length > 0);
  return hosts.length > 0 ? hosts : undefined;
}
