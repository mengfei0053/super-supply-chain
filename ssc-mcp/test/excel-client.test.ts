import assert from "node:assert/strict";
import { createServer, type IncomingMessage } from "node:http";
import { mkdtemp, readFile, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import path from "node:path";
import { once } from "node:events";
import { crc32 } from "node:zlib";
import test from "node:test";
import {
  SscClient,
  classifyExportBody,
  filenameFromContentDisposition,
  parseExcelIds,
  redactSecrets,
  resolveExportDownloadName,
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

test("disposition .xlsx wins over Content-Type application/zip", async () => {
  const xlsx = storedZip([
    { name: "[Content_Types].xml", data: Buffer.from("<Types/>") },
    { name: "xl/workbook.xml", data: Buffer.from("<workbook/>") },
    { name: "xl/worksheets/sheet1.xml", data: Buffer.from("<worksheet/>") },
    { name: "docProps/core.xml", data: Buffer.from("<core/>") },
  ]);
  assert.equal(classifyExportBody(xlsx), "ooxml-workbook");
  const header =
    "attachment; filename=\"export.xlsx\"; filename*=UTF-8''%E5%AF%BC%E5%87%BA%E5%8F%91%E7%A5%A8.xlsx";
  assert.equal(
    resolveExportDownloadName(header, "shortHaulInvoice.xlsx", xlsx),
    "导出发票.xlsx",
  );
  assert.equal(
    resolveExportDownloadName(null, "shortHaulInvoice.xlsx", xlsx),
    "shortHaulInvoice.xlsx",
  );
  assert.equal(
    resolveExportDownloadName(
      'attachment; filename="invoice.zip"',
      "shortHaulInvoice.xlsx",
      xlsx,
    ),
    "invoice.xlsx",
  );
  assert.equal(
    resolveExportDownloadName(
      'attachment; filename="invoice"',
      "shortHaulInvoice.xlsx",
      xlsx,
    ),
    "invoice.xlsx",
  );

  const server = createServer((_req, res) => {
    res.writeHead(200, {
      "Content-Type": "application/zip",
      "Content-Disposition": header,
    });
    res.end(xlsx);
  });
  const port = await listen(server);
  const zipPath = path.join(
    tmpdir(),
    `ssc-invoice-${Date.now()}-${Math.random().toString(16).slice(2)}.zip`,
  );

  try {
    const client = new SscClient({
      baseUrl: `http://127.0.0.1:${port}`,
      token: TOKEN,
    });
    const file = await client.exportExcel({
      tableName: "dynamic_settlement_statement_suqian",
      ids: "896",
      type: "invoice_unpacking",
      outputPath: zipPath,
    });
    assert.equal(file.xlsx, true);
    assert.equal(file.fileName, "导出发票.xlsx");
    assert.equal(file.path, zipPath.replace(/\.zip$/i, ".xlsx"));
    assert.match(file.path, /\.xlsx$/i);
    assert.equal(file.path.toLowerCase().endsWith(".zip"), false);
    assert.deepEqual(await readFile(file.path), xlsx);
    await assert.rejects(() => readFile(zipPath));
  } finally {
    server.close();
    await rm(zipPath, { force: true });
    await rm(zipPath.replace(/\.zip$/i, ".xlsx"), { force: true });
  }
});

test("a zip of separate workbooks keeps .zip when the name has no spreadsheet extension", () => {
  const archive = storedZip([
    { name: "clearance.xlsx", data: Buffer.from("one") },
    { name: "freight.xlsx", data: Buffer.from("two") },
  ]);
  assert.equal(classifyExportBody(archive), "multi-workbook-zip");
  assert.equal(
    resolveExportDownloadName(null, "invoice_freight.xlsx", archive),
    "invoice_freight.zip",
  );
  assert.equal(
    resolveExportDownloadName(
      'attachment; filename="bundle"',
      "invoice_freight.xlsx",
      archive,
    ),
    "bundle.zip",
  );
  assert.equal(
    resolveExportDownloadName(
      "attachment; filename*=UTF-8''%E5%AF%BC%E5%87%BA.xlsx",
      "invoice_freight.xlsx",
      archive,
    ),
    "导出.xlsx",
  );
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

test("listExcel sends filter JSON, object sort, and repeated range", async () => {
  let seenUrl = "";
  let seenAuth = "";
  const server = createServer((_req, res) => {
    seenUrl = _req.url ?? "";
    seenAuth = _req.headers.authorization ?? "";
    res.writeHead(200, {
      "Content-Type": "application/json",
      "Content-Range": "2",
    });
    res.end(
      JSON.stringify([
        {
          id: 896,
          fileName: "statement.xlsx",
          datas: { baseData: { port: "上海" }, list: [] },
        },
      ]),
    );
  });
  const port = await listen(server);
  try {
    const client = new SscClient({
      baseUrl: `http://127.0.0.1:${port}`,
      token: TOKEN,
    });
    const listed = await client.listExcel({
      tableName: "dynamic_settlement_statement_suqian",
      filterStart: "2026-10-09",
      filterEnd: "2026-10-09",
      sort: { field: "id", order: "ASC" },
    });
    const url = new URL(seenUrl, "http://127.0.0.1");
    assert.equal(
      url.pathname,
      "/api/admin/excel/dynamic_settlement_statement_suqian",
    );
    assert.equal(
      url.searchParams.get("filter"),
      '{"start":"2026-10-09","end":"2026-10-09"}',
    );
    assert.equal(
      url.searchParams.get("sort"),
      '{"field":"id","order":"ASC"}',
    );
    assert.deepEqual(url.searchParams.getAll("range"), ["0", "50"]);
    assert.equal(seenAuth, `Bearer ${TOKEN}`);
    assert.equal(listed.contentRange, "2");
    assert.ok(Array.isArray(listed.rows));
    assert.equal(JSON.stringify(listed).includes(TOKEN), false);

    await client.listExcel({
      tableName: "dynamic_settlement_statement_suqian",
      filterStart: "2026-10-01",
      filterEnd: "2026-10-09",
      sort: ["id", "DESC"],
      range: [0, 10],
    });
    const next = new URL(seenUrl, "http://127.0.0.1");
    assert.equal(next.searchParams.get("sort"), '["id","DESC"]');
    assert.deepEqual(next.searchParams.getAll("range"), ["0", "10"]);
  } finally {
    server.close();
  }
});

test("searchCompanies queries name and alias keyword", async () => {
  let seenUrl = "";
  const server = createServer((_req, res) => {
    seenUrl = _req.url ?? "";
    res.writeHead(200, { "Content-Type": "application/json" });
    res.end(
      JSON.stringify([
        {
          id: 1,
          created_at: "2026-01-01T00:00:00Z",
          updated_at: "2026-01-01T00:00:00Z",
          deleted_at: null,
          name: "南阳食品有限公司",
          addr_country: "中国",
          addr_province: "河南",
          addr_city: "南阳",
          addr_street: "示例路 1 号",
          unified_social_credit_code: "91330000NAME",
          bank_code: "BANK001",
          phone_num: "0377-0000000",
          alias: "食品",
          target_addr: "南阳",
        },
      ]),
    );
  });
  const port = await listen(server);
  try {
    const client = new SscClient({
      baseUrl: `http://127.0.0.1:${port}`,
      token: TOKEN,
    });
    const rows = await client.searchCompanies({ keyword: "南阳" });
    const url = new URL(seenUrl, "http://127.0.0.1");
    assert.equal(url.pathname, "/api/admin/companies");
    assert.equal(url.searchParams.get("keyword"), "南阳");
    assert.equal(url.searchParams.get("includeDeleted"), null);
    assert.ok(Array.isArray(rows));
    const first = (rows as Array<Record<string, unknown>>)[0];
    for (const key of [
      "id",
      "created_at",
      "updated_at",
      "deleted_at",
      "name",
      "addr_country",
      "addr_province",
      "addr_city",
      "addr_street",
      "unified_social_credit_code",
      "bank_code",
      "phone_num",
      "alias",
      "target_addr",
    ]) {
      assert.ok(key in first, `missing ${key}`);
    }

    await client.searchCompanies({ keyword: "南阳", includeDeleted: true });
    const next = new URL(seenUrl, "http://127.0.0.1");
    assert.equal(next.searchParams.get("includeDeleted"), "true");
    await assert.rejects(
      () => client.searchCompanies({ keyword: "  " }),
      /keyword is required/,
    );
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

function storedZip(files: Array<{ name: string; data: Buffer }>): Buffer {
  const locals: Buffer[] = [];
  const centrals: Buffer[] = [];
  let offset = 0;
  for (const file of files) {
    const name = Buffer.from(file.name, "utf8");
    const sum = crc32(file.data);
    const local = Buffer.alloc(30 + name.length + file.data.length);
    local.writeUInt32LE(0x04034b50, 0);
    local.writeUInt16LE(20, 4);
    local.writeUInt16LE(0, 8);
    local.writeUInt32LE(sum, 14);
    local.writeUInt32LE(file.data.length, 18);
    local.writeUInt32LE(file.data.length, 22);
    local.writeUInt16LE(name.length, 26);
    name.copy(local, 30);
    file.data.copy(local, 30 + name.length);
    locals.push(local);

    const central = Buffer.alloc(46 + name.length);
    central.writeUInt32LE(0x02014b50, 0);
    central.writeUInt16LE(20, 4);
    central.writeUInt16LE(20, 6);
    central.writeUInt32LE(sum, 16);
    central.writeUInt32LE(file.data.length, 20);
    central.writeUInt32LE(file.data.length, 24);
    central.writeUInt16LE(name.length, 28);
    central.writeUInt32LE(offset, 42);
    name.copy(central, 46);
    centrals.push(central);
    offset += local.length;
  }
  const localAll = Buffer.concat(locals);
  const centralDir = Buffer.concat(centrals);
  const eocd = Buffer.alloc(22);
  eocd.writeUInt32LE(0x06054b50, 0);
  eocd.writeUInt16LE(files.length, 8);
  eocd.writeUInt16LE(files.length, 10);
  eocd.writeUInt32LE(centralDir.length, 12);
  eocd.writeUInt32LE(localAll.length, 16);
  return Buffer.concat([localAll, centralDir, eocd]);
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
