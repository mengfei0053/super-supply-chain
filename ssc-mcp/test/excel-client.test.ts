import assert from "node:assert/strict";
import { createServer, type IncomingMessage } from "node:http";
import { mkdtemp, readFile, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import path from "node:path";
import { once } from "node:events";
import test from "node:test";
import {
  SscClient,
  filenameFromContentDisposition,
  parseExcelIds,
  redactSecrets,
  sanitizeDownloadName,
} from "../src/client.js";

const TOKEN = "unit-test-token-not-a-jwt";

test("parseExcelIds accepts a single id, commas, and arrays", () => {
  assert.deepEqual(parseExcelIds("896"), ["896"]);
  assert.deepEqual(parseExcelIds("896, 897"), ["896", "897"]);
  assert.deepEqual(parseExcelIds([896, "897", "896"]), ["896", "897"]);
  assert.throws(() => parseExcelIds(""), /ids is required/);
  assert.throws(() => parseExcelIds("abc"), /Invalid excel row id/);
});

test("filenameFromContentDisposition prefers RFC 5987 filename*", () => {
  const header =
    "attachment; filename=\"export.xlsx\"; filename*=UTF-8''%E5%AF%BC%E5%87%BA%E5%8F%91%E7%A5%A8.xlsx";
  assert.equal(
    filenameFromContentDisposition(header, "fallback.xlsx"),
    "导出发票.xlsx",
  );
  assert.equal(
    sanitizeDownloadName("../../etc/passwd.xlsx"),
    "passwd.xlsx",
  );
  assert.equal(sanitizeDownloadName("plain"), "plain.xlsx");
});

test("redactSecrets strips bearer tokens and JWTs", () => {
  const jwt = "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxIn0.signature";
  const text = `Authorization: Bearer ${TOKEN} token=${jwt}`;
  const redacted = redactSecrets(text);
  assert.equal(redacted.includes(TOKEN), false);
  assert.equal(redacted.includes(jwt), false);
  assert.match(redacted, /Bearer \[redacted\]/);
  assert.match(redacted, /\[redacted-jwt\]/);
});

test("uploadExcel posts multipart file and name", async () => {
  const dir = await mkdtemp(path.join(tmpdir(), "ssc-upload-"));
  const filePath = path.join(dir, "statement.xlsx");
  const payload = Buffer.concat([
    Buffer.from([0x50, 0x4b, 0x03, 0x04]),
    Buffer.from("fake-xlsx"),
  ]);
  await writeFile(filePath, payload);

  let seenAuth = "";
  let seenType = "";
  let seenUrl = "";
  let body = Buffer.alloc(0);

  const server = createServer(async (req, res) => {
    seenAuth = req.headers.authorization ?? "";
    seenType = req.headers["content-type"] ?? "";
    seenUrl = req.url ?? "";
    body = await readBody(req);
    res.writeHead(200, { "Content-Type": "application/json" });
    res.end(JSON.stringify({ inserted: true }));
  });
  const port = await listen(server);

  try {
    const client = new SscClient({
      baseUrl: `http://127.0.0.1:${port}`,
      token: TOKEN,
    });
    const result = await client.uploadExcel({
      tableName: "dynamic_settlement_statement_suqian",
      filePath,
      name: "宿迁结算.xlsx",
    });
    assert.equal(result.fileName, "宿迁结算.xlsx");
    assert.equal(result.bytes, payload.length);
    assert.deepEqual(result.response, { inserted: true });
    assert.equal(seenAuth, `Bearer ${TOKEN}`);
    assert.match(seenType, /^multipart\/form-data;/);
    assert.equal(
      seenUrl,
      "/api/admin/excel/dynamic_settlement_statement_suqian",
    );
    const raw = body.toString("utf8");
    assert.match(raw, /name="name"/);
    assert.match(raw, /宿迁结算\.xlsx/);
    assert.match(raw, /name="file"/);
    assert.ok(body.includes(payload));
    assert.equal(raw.includes(TOKEN), false);
  } finally {
    server.close();
    await rm(dir, { recursive: true, force: true });
  }
});

test("exportExcel writes a temp xlsx and repeats ids", async () => {
  const xlsx = Buffer.concat([
    Buffer.from([0x50, 0x4b, 0x03, 0x04]),
    Buffer.from("export-body"),
  ]);
  let seenUrl = "";
  let seenAuth = "";
  const server = createServer((req, res) => {
    seenUrl = req.url ?? "";
    seenAuth = req.headers.authorization ?? "";
    res.writeHead(200, {
      "Content-Type":
        "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
      "Content-Disposition":
        "attachment; filename=\"export.xlsx\"; filename*=UTF-8''%E5%AF%BC%E5%87%BA-%E7%9F%AD%E9%A9%B3%E8%B4%B9-%E5%8F%91%E7%A5%A8.xlsx",
    });
    res.end(xlsx);
  });
  const port = await listen(server);
  const outDir = await mkdtemp(path.join(tmpdir(), "ssc-export-out-"));

  try {
    const client = new SscClient({
      baseUrl: `http://127.0.0.1:${port}`,
      token: TOKEN,
    });
    const file = await client.exportExcel({
      tableName: "dynamic_settlement_statement_suqian",
      ids: "896,897",
      type: "shortHaulInvoice",
      outputPath: outDir,
    });
    const url = new URL(seenUrl, "http://127.0.0.1");
    assert.deepEqual(url.searchParams.getAll("ids"), ["896", "897"]);
    assert.equal(url.searchParams.get("type"), "shortHaulInvoice");
    assert.equal(
      url.pathname,
      "/api/admin/excel-exports/dynamic_settlement_statement_suqian",
    );
    assert.equal(seenAuth, `Bearer ${TOKEN}`);
    assert.equal(file.xlsx, true);
    assert.equal(file.bytes, xlsx.length);
    assert.equal(file.fileName, "导出-短驳费-发票.xlsx");
    assert.equal(path.dirname(file.path), outDir);
    assert.deepEqual(await readFile(file.path), xlsx);
    assert.equal(JSON.stringify(file).includes(TOKEN), false);
  } finally {
    server.close();
    await rm(outDir, { recursive: true, force: true });
  }
});

test("exportExcel surfaces JSON errors and does not write a file", async () => {
  const server = createServer((_req, res) => {
    res.writeHead(500, { "Content-Type": "application/json" });
    res.end(JSON.stringify({ error: "missing rows" }));
  });
  const port = await listen(server);
  const missing = path.join(
    tmpdir(),
    `ssc-should-not-exist-${Date.now()}.xlsx`,
  );

  try {
    const client = new SscClient({
      baseUrl: `http://127.0.0.1:${port}`,
      token: TOKEN,
    });
    await assert.rejects(
      () =>
        client.exportExcel({
          tableName: "dynamic_settlement_statement_suqian",
          ids: 896,
          type: "invoice_freight",
          outputPath: missing,
        }),
      /HTTP 500/,
    );
    await assert.rejects(() => readFile(missing));
  } finally {
    server.close();
  }
});

test("uploadExcel rejects a missing file before calling the API", async () => {
  let called = false;
  const server = createServer((_req, res) => {
    called = true;
    res.writeHead(500);
    res.end();
  });
  const port = await listen(server);
  try {
    const client = new SscClient({
      baseUrl: `http://127.0.0.1:${port}`,
      token: TOKEN,
    });
    await assert.rejects(
      () =>
        client.uploadExcel({
          tableName: "dynamic_settlement_statement_suqian",
          filePath: path.join(tmpdir(), "missing-statement.xlsx"),
          name: "missing.xlsx",
        }),
      /File not found/,
    );
    assert.equal(called, false);
  } finally {
    server.close();
  }
});

test("deleteExcelRow sends one id and rejects extra export types", async () => {
  let seenUrl = "";
  let seenMethod = "";
  const server = createServer((_req, res) => {
    seenUrl = _req.url ?? "";
    seenMethod = _req.method ?? "";
    res.writeHead(200, { "Content-Type": "application/json" });
    res.end(JSON.stringify({ message: "Delete successfully" }));
  });
  const port = await listen(server);
  try {
    const client = new SscClient({
      baseUrl: `http://127.0.0.1:${port}`,
      token: TOKEN,
    });
    const result = await client.deleteExcelRow(
      "dynamic_settlement_statement_suqian",
      896,
    );
    assert.equal(seenMethod, "DELETE");
    assert.equal(
      seenUrl,
      "/api/admin/excel/dynamic_settlement_statement_suqian/896",
    );
    assert.deepEqual(result, { message: "Delete successfully" });
    await assert.rejects(
      () =>
        client.exportExcel({
          tableName: "dynamic_settlement_statement_suqian",
          ids: "896",
          type: "shortHaul",
        }),
      /Unknown export type/,
    );
  } finally {
    server.close();
  }
});

function listen(server: ReturnType<typeof createServer>): Promise<number> {
  server.listen(0, "127.0.0.1");
  return once(server, "listening").then(() => {
    const address = server.address();
    if (!address || typeof address === "string") {
      throw new Error("expected a TCP port");
    }
    return address.port;
  });
}

function readBody(req: IncomingMessage): Promise<Buffer> {
  return new Promise((resolve, reject) => {
    const chunks: Buffer[] = [];
    req.on("data", (chunk) => {
      chunks.push(Buffer.isBuffer(chunk) ? chunk : Buffer.from(chunk));
    });
    req.on("end", () => resolve(Buffer.concat(chunks)));
    req.on("error", reject);
  });
}
