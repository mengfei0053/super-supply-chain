import assert from "node:assert/strict";
import { createServer } from "node:http";
import { once } from "node:events";
import { mkdtemp, readFile, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";
import test from "node:test";
import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { StdioClientTransport } from "@modelcontextprotocol/sdk/client/stdio.js";
import { StreamableHTTPClientTransport } from "@modelcontextprotocol/sdk/client/streamableHttp.js";
import { SscClient } from "../src/client.js";
import { startHttpServer } from "../src/http.js";

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
const expectedTools = [
  "ssc_upload_excel",
  "ssc_delete_excel_row",
  "ssc_export_excel",
  "ssc_list_excel",
  "ssc_search_companies",
  "ssc_create_company",
  "ssc_update_company",
  "ssc_delete_company",
];

test("streamable HTTP exposes upload and export tools", async () => {
  const client = new SscClient({ baseUrl: "http://127.0.0.1:9" });
  const httpServer = await startHttpServer({
    client,
    host: "127.0.0.1",
    port: 0,
    corsOrigin: "*",
  });
  const mcp = new Client({ name: "ssc-mcp-test", version: "0.0.0" });
  const transport = new StreamableHTTPClientTransport(new URL(httpServer.mcpUrl));
  try {
    await mcp.connect(transport);
    const listed = await mcp.listTools();
    const names = listed.tools.map((tool) => tool.name);
    for (const name of expectedTools) {
      assert.ok(names.includes(name), `missing tool ${name}`);
    }
    assert.equal(names.length, expectedTools.length);
    const denied = await mcp.callTool({
      name: "ssc_delete_excel_row",
      arguments: {
        tableName: "dynamic_settlement_statement_suqian",
        id: "896",
      },
    });
    assert.equal(denied.isError, true);
    const text = JSON.stringify(denied);
    assert.match(text, /Not authenticated/);
    assert.equal(text.includes("Bearer "), false);
  } finally {
    await mcp.close();
    await httpServer.close();
  }
});

test("HTTP tool call uploads and exports through the SSC client", async () => {
  const xlsx = Buffer.concat([
    Buffer.from([0x50, 0x4b, 0x03, 0x04]),
    Buffer.from("via-mcp"),
  ]);
  const dir = await mkdtemp(path.join(tmpdir(), "ssc-mcp-e2e-"));
  const filePath = path.join(dir, "in.xlsx");
  await writeFile(filePath, xlsx);

  const api = createServer(async (req, res) => {
    if (req.method === "POST") {
      res.writeHead(200, { "Content-Type": "application/json" });
      res.end(JSON.stringify({ inserted: 1 }));
      return;
    }
    if (req.method === "DELETE") {
      res.writeHead(200, { "Content-Type": "application/json" });
      res.end(JSON.stringify({ message: "Delete successfully" }));
      return;
    }
    res.writeHead(200, {
      "Content-Type": "application/zip",
      "Content-Disposition": 'attachment; filename="export.xlsx"',
    });
    res.end(xlsx);
  });
  api.listen(0, "127.0.0.1");
  await once(api, "listening");
  const apiAddress = api.address();
  if (!apiAddress || typeof apiAddress === "string") {
    throw new Error("expected API port");
  }

  const ssc = new SscClient({
    baseUrl: `http://127.0.0.1:${apiAddress.port}`,
    token: "unit-test-token-not-a-jwt",
  });
  const httpServer = await startHttpServer({
    client: ssc,
    host: "127.0.0.1",
    port: 0,
    corsOrigin: "*",
  });
  const mcp = new Client({ name: "ssc-mcp-test", version: "0.0.0" });
  const transport = new StreamableHTTPClientTransport(new URL(httpServer.mcpUrl));

  try {
    await mcp.connect(transport);
    const uploaded = await mcp.callTool({
      name: "ssc_upload_excel",
      arguments: {
        tableName: "dynamic_settlement_statement_suqian",
        filePath,
        name: "宿迁结算.xlsx",
      },
    });
    assert.equal(uploaded.isError, undefined);
    const uploadText = textOf(uploaded);
    assert.match(uploadText, /宿迁结算\.xlsx/);
    assert.equal(uploadText.includes("unit-test-token"), false);

    const exported = await mcp.callTool({
      name: "ssc_export_excel",
      arguments: {
        tableName: "dynamic_settlement_statement_suqian",
        ids: "896",
        type: "shortHaulInvoice",
        outputPath: dir,
      },
    });
    assert.equal(exported.isError, undefined);
    const exportText = textOf(exported);
    const parsed = JSON.parse(exportText) as {
      path: string;
      xlsx: boolean;
      fileName: string;
    };
    assert.equal(parsed.xlsx, true);
    assert.equal(parsed.fileName, "export.xlsx");
    assert.match(parsed.path, /\.xlsx$/i);
    assert.equal(parsed.path.toLowerCase().endsWith(".zip"), false);
    assert.deepEqual(await readFile(parsed.path), xlsx);
    assert.equal(exportText.includes("unit-test-token"), false);

    const deleted = await mcp.callTool({
      name: "ssc_delete_excel_row",
      arguments: {
        tableName: "dynamic_settlement_statement_suqian",
        id: 896,
      },
    });
    assert.equal(deleted.isError, undefined);
    assert.match(textOf(deleted), /Delete successfully/);
    assert.equal(textOf(deleted).includes("unit-test-token"), false);
  } finally {
    await mcp.close();
    await httpServer.close();
    api.close();
    await rm(dir, { recursive: true, force: true });
  }
});

test("stdio transport exposes upload and export tools", async () => {
  const tsx = path.join(root, "node_modules", ".bin", "tsx");
  const mcp = new Client({ name: "ssc-mcp-test", version: "0.0.0" });
  const transport = new StdioClientTransport({
    command: tsx,
    args: [path.join(root, "src", "index.ts")],
    cwd: root,
    stderr: "pipe",
    env: {
      PATH: process.env.PATH ?? "",
      HOME: process.env.HOME ?? "",
      SSC_BASE_URL: "http://127.0.0.1:9",
    },
  });
  try {
    await mcp.connect(transport);
    const listed = await mcp.listTools();
    const names = listed.tools.map((tool) => tool.name);
    for (const name of expectedTools) {
      assert.ok(names.includes(name), `missing tool ${name}`);
    }
    assert.equal(names.length, expectedTools.length);
    const denied = await mcp.callTool({
      name: "ssc_export_excel",
      arguments: {
        tableName: "dynamic_settlement_statement_suqian",
        ids: "896",
        type: "invoice_freight",
      },
    });
    assert.equal(denied.isError, true);
    assert.match(JSON.stringify(denied), /Not authenticated/);
  } finally {
    await mcp.close();
  }
});

function textOf(result: { content?: unknown }): string {
  const content = result.content;
  if (!Array.isArray(content)) return JSON.stringify(result);
  return content
    .map((item) =>
      item && typeof item === "object" && "text" in item
        ? String((item as { text: unknown }).text)
        : "",
    )
    .join("\n");
}
