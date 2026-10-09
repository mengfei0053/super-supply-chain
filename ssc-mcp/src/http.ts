/**
 * HTTP transports for ssc-mcp.
 *
 * Streamable HTTP (recommended, protocol 2025-11-25) on `/mcp`.
 * Legacy HTTP+SSE (protocol 2024-11-05) on `/sse` and `/messages`.
 *
 * Each MCP session gets its own McpServer. Sessions share one SscClient,
 * so SSC login state is process-wide.
 */

import { randomUUID } from "node:crypto";
import type { Server } from "node:http";
import type { NextFunction, Request, Response } from "express";
import { isInitializeRequest } from "@modelcontextprotocol/sdk/types.js";
import { createMcpExpressApp } from "@modelcontextprotocol/sdk/server/express.js";
import { SSEServerTransport } from "@modelcontextprotocol/sdk/server/sse.js";
import { StreamableHTTPServerTransport } from "@modelcontextprotocol/sdk/server/streamableHttp.js";
import type { SscClient } from "./client.js";
import { createSscMcpServer } from "./server.js";

const IDLE_MS = 30 * 60_000;
const MAX_SESSIONS = 200;

type SessionRecord = {
  transport: StreamableHTTPServerTransport | SSEServerTransport;
  open: number;
  lastActive: number;
};

export type HttpServerHandle = {
  host: string;
  port: number;
  /** Preferred MCP endpoint (Streamable HTTP). */
  mcpUrl: string;
  /** Legacy SSE endpoint. */
  sseUrl: string;
  close: () => Promise<void>;
};

export type StartHttpServerOptions = {
  client: SscClient;
  host: string;
  port: number;
  corsOrigin: string;
  allowedHosts?: string[];
};

export async function startHttpServer(
  options: StartHttpServerOptions,
): Promise<HttpServerHandle> {
  const { client, host, port, corsOrigin, allowedHosts } = options;
  const app = createMcpExpressApp(
    allowedHosts ? { host, allowedHosts } : { host },
  );
  app.use(corsMiddleware(corsOrigin));

  const sessions = new Map<string, SessionRecord>();

  const sweep = setInterval(() => {
    const cutoff = Date.now() - IDLE_MS;
    for (const session of sessions.values()) {
      if (session.open === 0 && session.lastActive < cutoff) {
        void session.transport.close().catch((err: unknown) => {
          console.error("ssc-mcp failed to close idle session:", err);
        });
      }
    }
  }, 60_000);
  sweep.unref();

  app.all("/mcp", async (req: Request, res: Response) => {
    try {
      await handleStreamable(req, res, client, sessions);
    } catch (err) {
      console.error("ssc-mcp streamable HTTP error:", err);
      if (!res.headersSent) {
        jsonRpcError(res, 500, -32603, "Internal server error");
      }
    }
  });

  app.get("/sse", async (_req: Request, res: Response) => {
    try {
      await handleSseConnect(res, client, sessions);
    } catch (err) {
      console.error("ssc-mcp SSE connect error:", err);
      if (!res.headersSent) {
        res.status(500).send("Internal server error");
      }
    }
  });

  app.post("/messages", async (req: Request, res: Response) => {
    try {
      await handleSseMessage(req, res, sessions);
    } catch (err) {
      console.error("ssc-mcp SSE message error:", err);
      if (!res.headersSent) {
        jsonRpcError(res, 500, -32603, "Internal server error");
      }
    }
  });

  const httpServer = await listen(app, port, host);
  const bound = httpServer.address();
  if (!bound || typeof bound === "string") {
    httpServer.close();
    throw new Error("ssc-mcp HTTP server did not bind a TCP port");
  }

  const advertisedHost = clientHost(host);
  const mcpUrl = `http://${advertisedHost}:${bound.port}/mcp`;
  const sseUrl = `http://${advertisedHost}:${bound.port}/sse`;

  let closed = false;
  const close = async (): Promise<void> => {
    if (closed) return;
    closed = true;
    clearInterval(sweep);
    const pending = [...sessions.values()];
    sessions.clear();
    await Promise.all(
      pending.map(async (session) => {
        try {
          await session.transport.close();
        } catch (err) {
          console.error("ssc-mcp session close error:", err);
        }
      }),
    );
    if (typeof httpServer.closeAllConnections === "function") {
      httpServer.closeAllConnections();
    }
    await new Promise<void>((resolve, reject) => {
      httpServer.close((err) => (err ? reject(err) : resolve()));
    });
  };

  return {
    host,
    port: bound.port,
    mcpUrl,
    sseUrl,
    close,
  };
}

async function handleStreamable(
  req: Request,
  res: Response,
  client: SscClient,
  sessions: Map<string, SessionRecord>,
): Promise<void> {
  const sessionId = headerValue(req, "mcp-session-id");
  const existing = sessionId ? sessions.get(sessionId) : undefined;

  let transport: StreamableHTTPServerTransport;

  if (existing) {
    if (!(existing.transport instanceof StreamableHTTPServerTransport)) {
      jsonRpcError(
        res,
        400,
        -32000,
        "Bad Request: Session exists but uses a different transport protocol",
      );
      return;
    }
    existing.lastActive = Date.now();
    trackResponse(existing, res);
    transport = existing.transport;
  } else if (!sessionId && req.method === "POST" && isInitializeRequest(req.body)) {
    if (sessions.size >= MAX_SESSIONS) {
      jsonRpcError(res, 503, -32000, "Too many open sessions");
      return;
    }
    const record: SessionRecord = {
      transport: undefined as unknown as StreamableHTTPServerTransport,
      open: 0,
      lastActive: Date.now(),
    };
    transport = new StreamableHTTPServerTransport({
      sessionIdGenerator: () => randomUUID(),
      onsessioninitialized: (id) => {
        sessions.set(id, record);
        console.error(`ssc-mcp streamable session started ${id}`);
      },
    });
    record.transport = transport;
    transport.onclose = () => {
      const id = transport.sessionId;
      if (id && sessions.delete(id)) {
        console.error(`ssc-mcp streamable session closed ${id}`);
      }
    };
    trackResponse(record, res);
    const server = createSscMcpServer(client);
    await server.connect(transport);
  } else if (sessionId) {
    jsonRpcError(res, 404, -32001, "Session not found");
    return;
  } else {
    jsonRpcError(
      res,
      400,
      -32000,
      "Bad Request: No valid session ID provided",
    );
    return;
  }

  await transport.handleRequest(req, res, req.body);
}

async function handleSseConnect(
  res: Response,
  client: SscClient,
  sessions: Map<string, SessionRecord>,
): Promise<void> {
  if (sessions.size >= MAX_SESSIONS) {
    res.status(503).send("Too many open sessions");
    return;
  }
  const transport = new SSEServerTransport("/messages", res);
  const record: SessionRecord = {
    transport,
    open: 0,
    lastActive: Date.now(),
  };
  sessions.set(transport.sessionId, record);
  trackResponse(record, res);
  transport.onclose = () => {
    if (sessions.delete(transport.sessionId)) {
      console.error(`ssc-mcp SSE session closed ${transport.sessionId}`);
    }
  };
  console.error(`ssc-mcp SSE session started ${transport.sessionId}`);
  const server = createSscMcpServer(client);
  try {
    await server.connect(transport);
  } catch (err) {
    sessions.delete(transport.sessionId);
    throw err;
  }
}

async function handleSseMessage(
  req: Request,
  res: Response,
  sessions: Map<string, SessionRecord>,
): Promise<void> {
  const sessionId = typeof req.query.sessionId === "string" ? req.query.sessionId : undefined;
  if (!sessionId) {
    jsonRpcError(res, 400, -32000, "Bad Request: sessionId query parameter required");
    return;
  }
  const existing = sessions.get(sessionId);
  if (!existing) {
    jsonRpcError(res, 404, -32001, "Session not found");
    return;
  }
  if (!(existing.transport instanceof SSEServerTransport)) {
    jsonRpcError(
      res,
      400,
      -32000,
      "Bad Request: Session exists but uses a different transport protocol",
    );
    return;
  }
  existing.lastActive = Date.now();
  await existing.transport.handlePostMessage(req, res, req.body);
}

function trackResponse(session: SessionRecord, res: Response): void {
  if (res.destroyed) return;
  session.open++;
  res.on("close", () => {
    session.open = Math.max(0, session.open - 1);
    session.lastActive = Date.now();
  });
}

function headerValue(req: Request, name: string): string | undefined {
  const value = req.headers[name];
  if (Array.isArray(value)) return value[0];
  return value;
}

function jsonRpcError(
  res: Response,
  status: number,
  code: number,
  message: string,
): void {
  res.status(status).json({
    jsonrpc: "2.0",
    error: { code, message },
    id: null,
  });
}

function clientHost(host: string): string {
  if (host === "0.0.0.0") return "127.0.0.1";
  if (host === "::" || host === "[::]") return "[::1]";
  return host;
}

function corsMiddleware(corsOrigin: string) {
  const allowAny = corsOrigin.trim() === "*";
  const allowed = new Set(
    corsOrigin
      .split(",")
      .map((origin) => origin.trim())
      .filter((origin) => origin.length > 0),
  );

  return (req: Request, res: Response, next: NextFunction): void => {
    const requestOrigin = headerValue(req, "origin");
    const allow = allowAny
      ? "*"
      : requestOrigin && allowed.has(requestOrigin)
        ? requestOrigin
        : undefined;
    if (allow) {
      res.setHeader("Access-Control-Allow-Origin", allow);
      res.setHeader("Vary", "Origin");
      res.setHeader(
        "Access-Control-Allow-Methods",
        "GET, POST, DELETE, OPTIONS",
      );
      res.setHeader(
        "Access-Control-Allow-Headers",
        "Content-Type, Accept, Authorization, Mcp-Session-Id, Mcp-Protocol-Version, Last-Event-ID",
      );
      res.setHeader(
        "Access-Control-Expose-Headers",
        "Mcp-Session-Id, Mcp-Protocol-Version",
      );
      res.setHeader("Access-Control-Max-Age", "600");
    }
    if (req.method === "OPTIONS") {
      res.status(204).end();
      return;
    }
    next();
  };
}

function listen(app: { listen: ServerListen }, port: number, host: string): Promise<Server> {
  return new Promise((resolve, reject) => {
    const bound: { server?: Server } = {};
    const done = (err?: Error): void => {
      // Express may report listen errors through this callback. Defer so the
      // assignment below has completed if the callback ever runs inline.
      queueMicrotask(() => {
        if (err) {
          reject(err);
          return;
        }
        if (!bound.server) {
          reject(new Error("HTTP server failed to bind"));
          return;
        }
        resolve(bound.server);
      });
    };
    bound.server = app.listen(port, host, done);
  });
}

type ServerListen = (
  port: number,
  host: string,
  callback: (err?: Error) => void,
) => Server;
