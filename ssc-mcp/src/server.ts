/**
 * SSC MCP tools shared by stdio and Streamable HTTP.
 *
 * Commonly used Excel upload, list, delete, and export calls, plus company
 * keyword search. Auth is environment-only: SSC_TOKEN, or SSC_USERNAME + SSC_PASSWORD.
 */

import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import * as z from "zod/v4";
import {
  EXCEL_EXPORT_TYPES,
  SscClient,
  errorResult,
  jsonResult,
} from "./client.js";

const excelSortSchema = z
  .union([
    z.object({
      field: z
        .string()
        .describe("Sort field, e.g. id. Sent as {\"field\",\"order\"}."),
      order: z.enum(["ASC", "DESC"]),
    }),
    z.tuple([z.string(), z.enum(["ASC", "DESC"])]).describe(
      'Array form ["id","ASC"], also accepted by the list API.',
    ),
  ])
  .optional();

export function createSscMcpServer(client: SscClient): McpServer {
  const server = new McpServer({
    name: "ssc-mcp",
    version: "1.2.0",
  });

  server.registerTool(
    "ssc_upload_excel",
    {
      description:
        "常用. Upload an Excel workbook into a dynamic table (POST /api/admin/excel/:tableName). Multipart fields: file (local .xlsx/.xls path on the MCP host) and name. Inserts a parsed row. Example tableName: dynamic_settlement_statement_suqian.",
      inputSchema: {
        tableName: z
          .string()
          .describe(
            "Dynamic table name, e.g. dynamic_settlement_statement_suqian",
          ),
        filePath: z
          .string()
          .describe(
            "Absolute or relative path to an .xlsx/.xls file readable by this MCP process",
          ),
        name: z
          .string()
          .describe(
            "Upload name sent as the multipart name field. When it ends in .xls or .xlsx it is also the stored file name.",
          ),
      },
    },
    async ({ tableName, filePath, name }) => {
      try {
        const uploaded = await client.uploadExcel({ tableName, filePath, name });
        return jsonResult({
          ok: true,
          tableName: uploaded.tableName,
          name: uploaded.name,
          fileName: uploaded.fileName,
          bytes: uploaded.bytes,
          response: boundPayload(uploaded.response),
        });
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_delete_excel_row",
    {
      description:
        "常用. Hard-delete one dynamic Excel row (DELETE /api/admin/excel/:tableName/:id). Example tableName: dynamic_settlement_statement_suqian, id: 896.",
      inputSchema: {
        tableName: z
          .string()
          .describe(
            "Dynamic table name, e.g. dynamic_settlement_statement_suqian",
          ),
        id: z
          .union([z.string(), z.number()])
          .describe("Single row id, e.g. 896"),
      },
    },
    async ({ tableName, id }) => {
      try {
        return jsonResult({
          ok: true,
          tableName,
          id: String(id),
          response: await client.deleteExcelRow(tableName, id),
        });
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_export_excel",
    {
      description: `常用. Download an Excel export (GET /api/admin/excel-exports/:tableName?ids=&type=). Writes one workbook as .xlsx (name from Content-Disposition; Content-Type application/zip is still .xlsx) and returns that path. Types: ${EXCEL_EXPORT_TYPES.join(", ")}. Example: tableName dynamic_settlement_statement_suqian, ids 896, type shortHaulInvoice.`,
      inputSchema: {
        tableName: z
          .string()
          .describe(
            "Dynamic table name, e.g. dynamic_settlement_statement_suqian",
          ),
        ids: z
          .union([
            z.string(),
            z.number(),
            z.array(z.union([z.string(), z.number()])),
          ])
          .describe(
            'Row id(s). Example: "896", "896,897", or ["896","897"]. Sent as repeated ids query params.',
          ),
        type: z
          .enum(EXCEL_EXPORT_TYPES)
          .describe(`Export type: ${EXCEL_EXPORT_TYPES.join(", ")}`),
        outputPath: z
          .string()
          .optional()
          .describe(
            "Optional file path or directory for the downloaded workbook. A path ending in .zip is written as .xlsx when the response is one workbook. Default: OS temp dir / ssc-mcp-exports (or SSC_EXPORT_DIR).",
          ),
      },
    },
    async ({ tableName, ids, type, outputPath }) => {
      try {
        const file = await client.exportExcel({
          tableName,
          ids,
          type,
          outputPath,
        });
        return jsonResult({
          ok: true,
          tableName: file.tableName,
          ids: file.ids,
          type: file.type,
          path: file.path,
          bytes: file.bytes,
          contentType: file.contentType,
          fileName: file.fileName,
          contentDisposition: file.contentDisposition,
          xlsx: file.xlsx,
        });
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_list_excel",
    {
      description:
        "常用. List dynamic Excel rows (GET /api/admin/excel/:tableName). Sends filter={\"start\",\"end\"} (created_at; required for non-empty 宿迁结算 results), optional sort, and repeated range=start&range=end (default 0 and 50). Example tableName: dynamic_settlement_statement_suqian.",
      inputSchema: {
        tableName: z
          .string()
          .describe(
            "Dynamic table name, e.g. dynamic_settlement_statement_suqian",
          ),
        filterStart: z
          .string()
          .describe("created_at start date YYYY-MM-DD, inclusive"),
        filterEnd: z
          .string()
          .describe("created_at end date YYYY-MM-DD, inclusive"),
        sort: excelSortSchema.describe(
          'Optional React Admin sort. Object {"field":"id","order":"ASC"} or array ["id","ASC"]. The list handler accepts either and does not ORDER BY it.',
        ),
        range: z
          .tuple([z.number().int().min(0), z.number().int().min(0)])
          .optional()
          .describe(
            "Pagination [start, end). Sent as repeated range=start&range=end. Default [0, 50]. The API also accepts range=[start,end].",
          ),
      },
    },
    async ({ tableName, filterStart, filterEnd, sort, range }) => {
      try {
        const listed = await client.listExcel({
          tableName,
          filterStart,
          filterEnd,
          sort,
          range,
        });
        return jsonResult({
          ok: true,
          tableName,
          contentRange: listed.contentRange,
          rows: listed.rows,
        });
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_search_companies",
    {
      description:
        "常用. Keyword search on base_companies_infos (GET /api/admin/companies?keyword=). Matches name or alias with a literal substring (LIKE %keyword%). Returns every column. Soft-deleted rows are omitted unless includeDeleted is true.",
      inputSchema: {
        keyword: z
          .string()
          .describe(
            "Substring matched against name and alias, e.g. 南阳. % and _ are literal.",
          ),
        includeDeleted: z
          .boolean()
          .optional()
          .describe(
            "When true, include rows with deleted_at set. Default false (deleted_at IS NULL).",
          ),
      },
    },
    async ({ keyword, includeDeleted }) => {
      try {
        const rows = await client.searchCompanies({ keyword, includeDeleted });
        return jsonResult({
          ok: true,
          keyword: keyword.trim(),
          includeDeleted: includeDeleted === true,
          rows,
        });
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  return server;
}

function boundPayload(data: unknown, max = 100_000): unknown {
  const text = typeof data === "string" ? data : JSON.stringify(data);
  if (text.length <= max) return data;
  return {
    truncated: true,
    preview: text.slice(0, max),
  };
}
