/**
 * Thin HTTP client for Super Supply Chain (SSC) backend.
 *
 * API origin: SSC_BASE_URL (default https://ssc.mengfei.tech)
 * Public:    POST /api/login, POST /api/register
 * Protected: /api/admin/* with Authorization: Bearer <jwt or personal access token>
 */

import { mkdir, readFile, stat, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import path from "node:path";

export type SscConfig = {
  baseUrl: string;
  token?: string;
  username?: string;
  password?: string;
};

/** Export types exposed by ssc_export_excel. */
export const EXCEL_EXPORT_TYPES = [
  "shortHaulInvoice",
  "invoice_unpacking",
  "invoice_clearance_only",
  "invoice_freight",
] as const;

export type ExcelExportType = (typeof EXCEL_EXPORT_TYPES)[number];

const XLSX_MIME =
  "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet";

export type ExcelExportFile = {
  path: string;
  bytes: number;
  contentType: string | null;
  fileName: string;
  contentDisposition: string | null;
  xlsx: boolean;
};

export class SscApiError extends Error {
  constructor(
    message: string,
    readonly status: number,
    readonly body: string,
  ) {
    super(message);
    this.name = "SscApiError";
  }
}

function trimSlash(url: string): string {
  return url.replace(/\/+$/, "");
}

export function assertTableName(tableName: string): string {
  const name = tableName.trim();
  if (!/^[A-Za-z0-9_]+$/.test(name)) {
    throw new Error(
      "tableName must be a dynamic table identifier (letters, digits, underscore), e.g. dynamic_settlement_statement_suqian",
    );
  }
  return name;
}

/** Accept "896", "896,897", 896, or ["896", 897]. */
export function parseExcelIds(
  ids: string | number | Array<string | number>,
): string[] {
  const parts = Array.isArray(ids) ? ids : [ids];
  const out: string[] = [];
  for (const part of parts) {
    for (const piece of String(part).split(/[,\s]+/)) {
      const id = piece.trim();
      if (!id) continue;
      if (!/^\d+$/.test(id)) {
        throw new Error(`Invalid excel row id: ${id}`);
      }
      if (!out.includes(id)) out.push(id);
    }
  }
  if (out.length === 0) {
    throw new Error("ids is required (example: 896)");
  }
  return out;
}

export function isExcelExportType(value: string): value is ExcelExportType {
  return (EXCEL_EXPORT_TYPES as readonly string[]).includes(value);
}

export function filenameFromContentDisposition(
  header: string | null,
  fallback: string,
): string {
  if (!header) return fallback;
  const star = header.match(/filename\*\s*=\s*(?:UTF-8''|utf-8'')([^;\s]+)/i);
  if (star?.[1]) {
    try {
      const decoded = decodeURIComponent(star[1].trim());
      if (decoded) return decoded;
    } catch {
      // Fall through to filename=
    }
  }
  const quoted = header.match(/filename\s*=\s*"([^"]+)"/i);
  if (quoted?.[1]) return quoted[1];
  const plain = header.match(/filename\s*=\s*([^;\s]+)/i);
  if (plain?.[1]) return plain[1].replace(/^"(.*)"$/, "$1");
  return fallback;
}

export function sanitizeDownloadName(name: string): string {
  const base = path
    .basename(name)
    .replace(/[\r\n\0]/g, "")
    .trim()
    .replace(/^\.+/, "");
  const cleaned = base || "export.xlsx";
  const limited = cleaned.slice(0, 180);
  return /\.xlsx?$/i.test(limited) ? limited : `${limited}.xlsx`;
}

function multipartFileName(name: string, filePath: string): string {
  const fromName = path.basename(name.trim());
  if (/\.xlsx?$/i.test(fromName)) return fromName;
  return path.basename(filePath);
}

async function resolveExportPath(
  outputPath: string | undefined,
  fileName: string,
): Promise<string> {
  if (outputPath && outputPath.trim()) {
    const resolved = path.resolve(outputPath.trim());
    try {
      const st = await stat(resolved);
      if (st.isDirectory()) return path.join(resolved, fileName);
    } catch {
      // Caller asked for a file path that does not exist yet.
    }
    return resolved;
  }
  const dir = process.env.SSC_EXPORT_DIR?.trim()
    ? path.resolve(process.env.SSC_EXPORT_DIR.trim())
    : path.join(tmpdir(), "ssc-mcp-exports");
  const stamp = new Date().toISOString().replace(/[:.]/g, "-");
  return path.join(dir, `${stamp}-${fileName}`);
}

export function loadConfigFromEnv(): SscConfig {
  const baseUrl = trimSlash(
    process.env.SSC_BASE_URL?.trim() || "https://ssc.mengfei.tech",
  );
  return {
    baseUrl,
    token: process.env.SSC_TOKEN?.trim() || undefined,
    username: process.env.SSC_USERNAME?.trim() || undefined,
    password: process.env.SSC_PASSWORD || undefined,
  };
}

export class SscClient {
  private token?: string;
  readonly baseUrl: string;
  private username?: string;
  private password?: string;

  constructor(config: SscConfig) {
    this.baseUrl = trimSlash(config.baseUrl);
    this.token = config.token;
    this.username = config.username;
    this.password = config.password;
  }

  /** Ensure we have a bearer token: use SSC_TOKEN (JWT or PAT), or login with username/password. */
  async ensureAuth(): Promise<string> {
    if (this.token) return this.token;
    if (!this.username || !this.password) {
      throw new Error(
        "Not authenticated. Set SSC_TOKEN, or SSC_USERNAME + SSC_PASSWORD.",
      );
    }
    const result = await this.login(this.username, this.password);
    return result.token;
  }

  async login(
    username: string,
    password: string,
  ): Promise<{
    id: number;
    username: string;
    fullName: string;
    email: string;
    token: string;
    avatar: string;
  }> {
    const data = await this.requestJson<{
      id: number;
      username: string;
      fullName: string;
      email: string;
      token: string;
      avatar: string;
    }>("POST", "/api/login", {
      body: { username, password },
      auth: false,
    });
    if (!data.token) {
      throw new Error("Login response did not include token");
    }
    this.token = data.token;
    this.username = username;
    return data;
  }

  async deleteExcelRow(
    tableName: string,
    id: string | number,
  ): Promise<unknown> {
    const table = assertTableName(tableName);
    const ids = parseExcelIds(id);
    if (ids.length !== 1) {
      throw new Error("delete accepts exactly one row id");
    }
    return this.requestJson(
      "DELETE",
      `/api/admin/excel/${encodeURIComponent(table)}/${encodeURIComponent(ids[0])}`,
    );
  }

  /**
   * POST /api/admin/excel/:tableName
   * multipart fields: file (xlsx bytes), name (display name; also used as the file part filename when it ends in .xls/.xlsx).
   */
  async uploadExcel(opts: {
    tableName: string;
    filePath: string;
    name: string;
  }): Promise<{
    tableName: string;
    name: string;
    fileName: string;
    bytes: number;
    response: unknown;
  }> {
    const tableName = assertTableName(opts.tableName);
    const name = opts.name.trim();
    if (!name) throw new Error("name is required");

    const abs = path.resolve(opts.filePath);
    let info;
    try {
      info = await stat(abs);
    } catch {
      throw new Error(`File not found: ${abs}`);
    }
    if (!info.isFile()) throw new Error(`Not a file: ${abs}`);
    const ext = path.extname(abs).toLowerCase();
    if (ext !== ".xlsx" && ext !== ".xls") {
      throw new Error(
        `Expected an .xlsx or .xls file, got ${ext || "no extension"}`,
      );
    }

    const bytes = await readFile(abs);
    const fileName = multipartFileName(name, abs);
    const form = new FormData();
    form.append(
      "file",
      new Blob([new Uint8Array(bytes)], { type: XLSX_MIME }),
      fileName,
    );
    form.append("name", name);

    const apiPath = `/api/admin/excel/${encodeURIComponent(tableName)}`;
    const res = await this.request("POST", apiPath, { body: form });
    const response = await readJsonBody(res, "POST", apiPath);
    return {
      tableName,
      name,
      fileName,
      bytes: bytes.length,
      response,
    };
  }

  /**
   * GET /api/admin/excel-exports/:tableName?ids=&type=
   * Writes the workbook to a temp file (or outputPath) and returns that path.
   */
  async exportExcel(opts: {
    tableName: string;
    ids: string | number | Array<string | number>;
    type: string;
    outputPath?: string;
  }): Promise<ExcelExportFile & { tableName: string; ids: string[]; type: ExcelExportType }> {
    const tableName = assertTableName(opts.tableName);
    const ids = parseExcelIds(opts.ids);
    const type = opts.type.trim();
    if (!isExcelExportType(type)) {
      throw new Error(
        `Unknown export type "${type}". Expected one of: ${EXCEL_EXPORT_TYPES.join(", ")}`,
      );
    }

    const apiPath = `/api/admin/excel-exports/${encodeURIComponent(tableName)}`;
    const res = await this.request("GET", apiPath, {
      query: { ids, type },
      accept: "*/*",
    });
    const buf = Buffer.from(await res.arrayBuffer());
    const contentType = res.headers.get("content-type");
    const disposition = res.headers.get("content-disposition");
    const declaredError =
      (contentType?.includes("json") ?? false) ||
      (contentType?.includes("text/") ?? false);
    const zip = buf.length >= 2 && buf[0] === 0x50 && buf[1] === 0x4b;

    if (!res.ok || declaredError || (looksLikeJsonError(buf) && !zip)) {
      throw new SscApiError(
        `SSC GET ${apiPath} failed: HTTP ${res.status}`,
        res.status,
        buf.toString("utf8").slice(0, 2000),
      );
    }
    if (buf.length === 0) {
      throw new Error("Export response was empty");
    }

    const fileName = sanitizeDownloadName(
      filenameFromContentDisposition(disposition, `${type}.xlsx`),
    );
    const outPath = await resolveExportPath(opts.outputPath, fileName);
    await mkdir(path.dirname(outPath), { recursive: true });
    await writeFile(outPath, buf);

    return {
      tableName,
      ids,
      type,
      path: outPath,
      bytes: buf.length,
      contentType,
      fileName,
      contentDisposition: disposition,
      xlsx: zip,
    };
  }

  private async request(
    method: string,
    path: string,
    opts: {
      jsonBody?: unknown;
      body?: BodyInit;
      query?: Record<string, string | string[]>;
      auth?: boolean;
      accept?: string;
    } = {},
  ): Promise<Response> {
    const needAuth = opts.auth !== false;
    if (needAuth) {
      await this.ensureAuth();
    }

    const url = new URL(this.baseUrl + path);
    if (opts.query) {
      for (const [key, value] of Object.entries(opts.query)) {
        const values = Array.isArray(value) ? value : [value];
        for (const item of values) url.searchParams.append(key, item);
      }
    }

    const headers: Record<string, string> = {};
    if (opts.accept) headers.Accept = opts.accept;
    else headers.Accept = "application/json";
    if (opts.jsonBody !== undefined) {
      headers["Content-Type"] = "application/json";
    }
    if (needAuth && this.token) {
      headers.Authorization = `Bearer ${this.token}`;
    }

    return fetch(url, {
      method,
      headers,
      body:
        opts.jsonBody !== undefined
          ? JSON.stringify(opts.jsonBody)
          : opts.body,
    });
  }

  private async requestJson<T = unknown>(
    method: string,
    path: string,
    opts: {
      body?: unknown;
      query?: Record<string, string | string[]>;
      auth?: boolean;
    } = {},
  ): Promise<T> {
    const res = await this.request(method, path, {
      jsonBody: opts.body,
      query: opts.query,
      auth: opts.auth,
    });
    return readJsonBody<T>(res, method, path);
  }
}

async function readJsonBody<T>(
  res: Response,
  method: string,
  path: string,
): Promise<T> {
  const text = await res.text();
  if (!res.ok) {
    throw new SscApiError(
      `SSC ${method} ${path} failed: HTTP ${res.status}`,
      res.status,
      text.slice(0, 2000),
    );
  }
  if (!text) return undefined as T;
  try {
    return JSON.parse(text) as T;
  } catch {
    return text as unknown as T;
  }
}

function looksLikeJsonError(buf: Buffer): boolean {
  const start = buf.subarray(0, 1).toString("utf8");
  return start === "{" || start === "[";
}

export function jsonResult(data: unknown): {
  content: Array<{ type: "text"; text: string }>;
} {
  return {
    content: [
      {
        type: "text",
        text: typeof data === "string" ? data : JSON.stringify(data, null, 2),
      },
    ],
  };
}

export function redactSecrets(text: string): string {
  return text
    .replace(/Bearer\s+[A-Za-z0-9\-._~+/]+=*/gi, "Bearer [redacted]")
    .replace(
      /eyJ[A-Za-z0-9_-]{8,}\.[A-Za-z0-9_-]{8,}\.[A-Za-z0-9_-]{8,}/g,
      "[redacted-jwt]",
    );
}

export function errorResult(err: unknown): {
  content: Array<{ type: "text"; text: string }>;
  isError: true;
} {
  let text: string;
  if (err instanceof SscApiError) {
    text = `${err.message}\n${err.body}`;
  } else if (err instanceof Error) {
    text = err.message;
  } else {
    text = String(err);
  }
  return { content: [{ type: "text", text: redactSecrets(text) }], isError: true };
}
