#!/usr/bin/env node
/**
 * Super Supply Chain (SSC) MCP server.
 *
 * Default transport is stdio. HTTP (Streamable HTTP + legacy SSE) is opt-in
 * via SSC_MCP_TRANSPORT=http or --transport http.
 */

import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { SscClient, loadConfigFromEnv } from "./client.js";
import { loadMcpRuntimeConfig, printHelp } from "./config.js";
import { startHttpServer } from "./http.js";
import { createSscMcpServer } from "./server.js";

async function main(): Promise<void> {
  const runtime = loadMcpRuntimeConfig();
  if (runtime.help) {
    printHelp();
    return;
  }

  const client = new SscClient(loadConfigFromEnv());

  if (runtime.transport === "stdio") {
    const server = createSscMcpServer(client);
    await server.connect(new StdioServerTransport());
    return;
  }

  const http = await startHttpServer({
    client,
    host: runtime.host,
    port: runtime.port,
    corsOrigin: runtime.corsOrigin,
    allowedHosts: runtime.allowedHosts,
  });
  console.error(
    [
      `ssc-mcp HTTP listening on ${http.host}:${http.port}`,
      `  Streamable HTTP  ${http.mcpUrl}`,
      `  Legacy SSE       ${http.sseUrl}`,
      `  CORS origin      ${runtime.corsOrigin}`,
    ].join("\n"),
  );

  let shuttingDown = false;
  const shutdown = async (signal: string): Promise<void> => {
    if (shuttingDown) return;
    shuttingDown = true;
    console.error(`ssc-mcp shutting down (${signal})`);
    try {
      await http.close();
    } catch (err) {
      console.error("ssc-mcp shutdown error:", err);
      process.exit(1);
    }
    process.exit(0);
  };
  process.on("SIGINT", () => {
    void shutdown("SIGINT");
  });
  process.on("SIGTERM", () => {
    void shutdown("SIGTERM");
  });
}

main().catch((err) => {
  console.error("ssc-mcp failed to start:", err);
  process.exit(1);
});
