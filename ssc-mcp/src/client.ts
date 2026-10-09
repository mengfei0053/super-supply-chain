/**
 * Thin HTTP client for Super Supply Chain (SSC) backend.
 *
 * API origin: SSC_BASE_URL (default https://ssc.mengfei.tech)
 * Public:    POST /api/login, POST /api/register
 * Protected: /api/admin/* with Authorization: Bearer <jwt>
 */

export type SscConfig = {
  baseUrl: string;
  token?: string;
  username?: string;
  password?: string;
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

  getToken(): string | undefined {
    return this.token;
  }

  setToken(token: string | undefined): void {
    this.token = token;
  }

  /** Ensure we have a JWT: use SSC_TOKEN, or login with username/password. */
  async ensureAuth(): Promise<string> {
    if (this.token) return this.token;
    if (!this.username || !this.password) {
      throw new Error(
        "Not authenticated. Set SSC_TOKEN, or SSC_USERNAME + SSC_PASSWORD, or call ssc_login.",
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

  async getMenus(): Promise<unknown> {
    return this.requestJson("GET", "/api/admin/menus");
  }

  async listOrders(): Promise<unknown> {
    return this.requestJson("GET", "/api/admin/settlement-form-entry");
  }

  async getOrder(id: string | number): Promise<unknown> {
    return this.requestJson(
      "GET",
      `/api/admin/settlement-form-entry/${encodeURIComponent(String(id))}`,
    );
  }

  async listDicts(range: [number, number] = [0, 49]): Promise<unknown> {
    return this.requestJson("GET", "/api/admin/dict-manage", {
      query: { range: JSON.stringify(range) },
    });
  }

  async getDict(id: string | number): Promise<unknown> {
    return this.requestJson(
      "GET",
      `/api/admin/dict-manage/${encodeURIComponent(String(id))}`,
    );
  }

  async getDictMap(type: string): Promise<unknown> {
    return this.requestJson(
      "GET",
      `/api/admin/dict-manage/map/${encodeURIComponent(type)}`,
    );
  }

  async createDict(body: {
    key: string;
    value: string;
    type: string;
  }): Promise<unknown> {
    return this.requestJson("POST", "/api/admin/dict-manage", { body });
  }

  async updateDict(
    id: string | number,
    body: { key?: string; value?: string; type?: string },
  ): Promise<unknown> {
    return this.requestJson(
      "PUT",
      `/api/admin/dict-manage/${encodeURIComponent(String(id))}`,
      { body },
    );
  }

  async deleteDict(id: string | number): Promise<unknown> {
    return this.requestJson(
      "DELETE",
      `/api/admin/dict-manage/${encodeURIComponent(String(id))}`,
    );
  }

  async listExcelReadRules(
    range: [number, number] = [0, 49],
  ): Promise<unknown> {
    return this.requestJson("GET", "/api/admin/excel-read-rules", {
      query: { range: JSON.stringify(range) },
    });
  }

  async getExcelReadRule(id: string | number): Promise<unknown> {
    return this.requestJson(
      "GET",
      `/api/admin/excel-read-rules/${encodeURIComponent(String(id))}`,
    );
  }

  async listExcelRows(opts: {
    tableName: string;
    range?: [number, number];
    filterStart?: string;
    filterEnd?: string;
  }): Promise<unknown> {
    const range = opts.range ?? [0, 49];
    const filterStart = opts.filterStart ?? "2000-01-01";
    const filterEnd = opts.filterEnd ?? "2099-12-31";
    return this.requestJson(
      "GET",
      `/api/admin/excel/${encodeURIComponent(opts.tableName)}`,
      {
        query: {
          range: JSON.stringify(range),
          "filter.start": filterStart,
          "filter.end": filterEnd,
        },
      },
    );
  }

  async getExcelRow(
    tableName: string,
    id: string | number,
  ): Promise<unknown> {
    return this.requestJson(
      "GET",
      `/api/admin/excel/${encodeURIComponent(tableName)}/${encodeURIComponent(String(id))}`,
    );
  }

  async updateExcelRow(
    tableName: string,
    id: string | number,
    body: Record<string, unknown>,
  ): Promise<unknown> {
    return this.requestJson(
      "PUT",
      `/api/admin/excel/${encodeURIComponent(tableName)}/${encodeURIComponent(String(id))}`,
      { body },
    );
  }

  async deleteExcelRow(
    tableName: string,
    id: string | number,
  ): Promise<unknown> {
    return this.requestJson(
      "DELETE",
      `/api/admin/excel/${encodeURIComponent(tableName)}/${encodeURIComponent(String(id))}`,
    );
  }

  async listExportRules(
    tableName: string,
    range: [number, number] = [0, 49],
  ): Promise<unknown> {
    return this.requestJson(
      "GET",
      `/api/admin/excel-export-rule/template/${encodeURIComponent(tableName)}`,
      { query: { range: JSON.stringify(range) } },
    );
  }

  async getExportRule(
    tableName: string,
    id: string | number,
  ): Promise<unknown> {
    return this.requestJson(
      "GET",
      `/api/admin/excel-export-rule/template/${encodeURIComponent(tableName)}/${encodeURIComponent(String(id))}`,
    );
  }

  async listExportTemplateOptions(associatedTable: string): Promise<unknown> {
    return this.requestJson("GET", "/api/admin/options/export-templates", {
      query: { associated_table: associatedTable },
    });
  }

  private async requestJson<T = unknown>(
    method: string,
    path: string,
    opts: {
      body?: unknown;
      query?: Record<string, string>;
      auth?: boolean;
    } = {},
  ): Promise<T> {
    const needAuth = opts.auth !== false;
    if (needAuth) {
      await this.ensureAuth();
    }

    const url = new URL(this.baseUrl + path);
    if (opts.query) {
      for (const [k, v] of Object.entries(opts.query)) {
        url.searchParams.set(k, v);
      }
    }

    const headers: Record<string, string> = {
      Accept: "application/json",
    };
    if (opts.body !== undefined) {
      headers["Content-Type"] = "application/json";
    }
    if (needAuth && this.token) {
      headers.Authorization = `Bearer ${this.token}`;
    }

    const res = await fetch(url, {
      method,
      headers,
      body: opts.body !== undefined ? JSON.stringify(opts.body) : undefined,
    });

    const text = await res.text();
    if (!res.ok) {
      throw new SscApiError(
        `SSC ${method} ${path} failed: HTTP ${res.status}`,
        res.status,
        text,
      );
    }

    if (!text) return undefined as T;
    try {
      return JSON.parse(text) as T;
    } catch {
      return text as unknown as T;
    }
  }
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
  return { content: [{ type: "text", text }], isError: true };
}
