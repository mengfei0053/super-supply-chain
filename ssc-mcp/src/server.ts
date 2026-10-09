/**
 * SSC MCP tool server. One instance is bound to a single transport
 * (stdio session, Streamable HTTP session, or legacy SSE session).
 * The shared {@link SscClient} keeps SSC auth state for the process.
 */

import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import * as z from "zod/v4";
import {
  SscClient,
  errorResult,
  jsonResult,
} from "./client.js";

export function createSscMcpServer(client: SscClient): McpServer {
  const server = new McpServer({
    name: "ssc-mcp",
    version: "1.0.0",
  });

  const rangeSchema = z
    .tuple([z.number().int().min(0), z.number().int().min(0)])
    .optional()
    .describe(
      'React-Admin style pagination range [start, endInclusive]. Default [0, 49]. Sent as range=[start,end].',
    );

  function asRange(
    range: [number, number] | undefined,
  ): [number, number] {
    return range ?? [0, 49];
  }

  server.registerTool(
    "ssc_login",
    {
      description:
        "Log in to SSC via POST /api/login. Stores JWT in-process for later tools. Prefer SSC_TOKEN env when possible.",
      inputSchema: {
        username: z.string().describe("SSC account username"),
        password: z.string().describe("SSC account password"),
      },
    },
    async ({ username, password }) => {
      try {
        const user = await client.login(username, password);
        return jsonResult({
          ok: true,
          id: user.id,
          username: user.username,
          fullName: user.fullName,
          email: user.email,
          tokenSet: true,
          note: "Token stored in MCP process memory only; not printed.",
        });
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_status",
    {
      description:
        "Check auth status: verifies JWT by calling GET /api/admin/menus (same as CLI `ssc status`).",
      inputSchema: {},
    },
    async () => {
      try {
        await client.ensureAuth();
        const menus = await client.getMenus();
        const menuCount = Array.isArray(menus) ? menus.length : undefined;
        return jsonResult({
          ok: true,
          baseUrl: client.baseUrl,
          authenticated: true,
          menuCount,
        });
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_list_menus",
    {
      description:
        "List dynamic Excel menus (GET /api/admin/menus). Each item has id, menuName, dynamicTableName — use dynamicTableName with excel table tools.",
      inputSchema: {},
    },
    async () => {
      try {
        return jsonResult(await client.getMenus());
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_list_orders",
    {
      description:
        "List settlement / order summary rows (GET /api/admin/settlement-form-entry). Returns id, orderNumber, arrivalDate, arrivalPort.",
      inputSchema: {},
    },
    async () => {
      try {
        return jsonResult(await client.listOrders());
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_get_order",
    {
      description: "Get one settlement/order summary by id (GET /api/admin/settlement-form-entry/:id).",
      inputSchema: {
        id: z.union([z.string(), z.number()]).describe("Order id"),
      },
    },
    async ({ id }) => {
      try {
        return jsonResult(await client.getOrder(id));
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_list_dicts",
    {
      description: "List dictionary entries (GET /api/admin/dict-manage). Requires range query.",
      inputSchema: {
        range: rangeSchema,
      },
    },
    async ({ range }) => {
      try {
        return jsonResult(await client.listDicts(asRange(range)));
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_get_dict",
    {
      description: "Get one dictionary entry (GET /api/admin/dict-manage/:id).",
      inputSchema: {
        id: z.union([z.string(), z.number()]),
      },
    },
    async ({ id }) => {
      try {
        return jsonResult(await client.getDict(id));
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_get_dict_map",
    {
      description:
        "Get key→value map for a dictionary type (GET /api/admin/dict-manage/map/:type).",
      inputSchema: {
        type: z.string().describe("Dictionary type string"),
      },
    },
    async ({ type }) => {
      try {
        return jsonResult(await client.getDictMap(type));
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_create_dict",
    {
      description: "Create a dictionary entry (POST /api/admin/dict-manage).",
      inputSchema: {
        key: z.string(),
        value: z.string(),
        type: z.string(),
      },
    },
    async (args) => {
      try {
        return jsonResult(await client.createDict(args));
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_update_dict",
    {
      description: "Update a dictionary entry (PUT /api/admin/dict-manage/:id).",
      inputSchema: {
        id: z.union([z.string(), z.number()]),
        key: z.string().optional(),
        value: z.string().optional(),
        type: z.string().optional(),
      },
    },
    async ({ id, key, value, type }) => {
      try {
        const body: { key?: string; value?: string; type?: string } = {};
        if (key !== undefined) body.key = key;
        if (value !== undefined) body.value = value;
        if (type !== undefined) body.type = type;
        return jsonResult(await client.updateDict(id, body));
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_delete_dict",
    {
      description: "Delete a dictionary entry (DELETE /api/admin/dict-manage/:id).",
      inputSchema: {
        id: z.union([z.string(), z.number()]),
      },
    },
    async ({ id }) => {
      try {
        return jsonResult(await client.deleteDict(id));
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_list_excel_read_rules",
    {
      description:
        "List Excel read-rule definitions / dynamic table registry (GET /api/admin/excel-read-rules).",
      inputSchema: {
        range: rangeSchema,
      },
    },
    async ({ range }) => {
      try {
        return jsonResult(await client.listExcelReadRules(asRange(range)));
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_get_excel_read_rule",
    {
      description: "Get one Excel read rule (GET /api/admin/excel-read-rules/:id).",
      inputSchema: {
        id: z.union([z.string(), z.number()]),
      },
    },
    async ({ id }) => {
      try {
        return jsonResult(await client.getExcelReadRule(id));
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_list_excel_rows",
    {
      description:
        "List rows from a dynamic Excel table (GET /api/admin/excel/:tableName). tableName comes from menus/read-rules (e.g. dynamic_customs_declaration_form). Date filter defaults to 2000-01-01 .. 2099-12-31.",
      inputSchema: {
        tableName: z
          .string()
          .describe("Dynamic MySQL table name, e.g. dynamic_yifan_cost_cal"),
        range: rangeSchema,
        filterStart: z
          .string()
          .optional()
          .describe("created_at start date YYYY-MM-DD"),
        filterEnd: z
          .string()
          .optional()
          .describe("created_at end date YYYY-MM-DD"),
      },
    },
    async ({ tableName, range, filterStart, filterEnd }) => {
      try {
        return jsonResult(
          await client.listExcelRows({
            tableName,
            range: asRange(range),
            filterStart,
            filterEnd,
          }),
        );
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_get_excel_row",
    {
      description: "Get one dynamic Excel table row (GET /api/admin/excel/:tableName/:id).",
      inputSchema: {
        tableName: z.string(),
        id: z.union([z.string(), z.number()]),
      },
    },
    async ({ tableName, id }) => {
      try {
        return jsonResult(await client.getExcelRow(tableName, id));
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_update_excel_row",
    {
      description:
        "Update a dynamic Excel table row (PUT /api/admin/excel/:tableName/:id). Body should match DynamicExcelTable fields (fileName, datas, etc.).",
      inputSchema: {
        tableName: z.string(),
        id: z.union([z.string(), z.number()]),
        body: z
          .record(z.string(), z.unknown())
          .describe("Partial DynamicExcelTable JSON body"),
      },
    },
    async ({ tableName, id, body }) => {
      try {
        return jsonResult(await client.updateExcelRow(tableName, id, body));
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_delete_excel_row",
    {
      description:
        "Hard-delete a dynamic Excel table row (DELETE /api/admin/excel/:tableName/:id).",
      inputSchema: {
        tableName: z.string(),
        id: z.union([z.string(), z.number()]),
      },
    },
    async ({ tableName, id }) => {
      try {
        return jsonResult(await client.deleteExcelRow(tableName, id));
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_list_export_rules",
    {
      description:
        "List Excel export templates for a table (GET /api/admin/excel-export-rule/template/:tableName).",
      inputSchema: {
        tableName: z.string(),
        range: rangeSchema,
      },
    },
    async ({ tableName, range }) => {
      try {
        return jsonResult(
          await client.listExportRules(tableName, asRange(range)),
        );
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_get_export_rule",
    {
      description:
        "Get one export template (GET /api/admin/excel-export-rule/template/:tableName/:id).",
      inputSchema: {
        tableName: z.string(),
        id: z.union([z.string(), z.number()]),
      },
    },
    async ({ tableName, id }) => {
      try {
        return jsonResult(await client.getExportRule(tableName, id));
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  server.registerTool(
    "ssc_list_export_template_options",
    {
      description:
        "List export-template select options for a table (GET /api/admin/options/export-templates?associated_table=...).",
      inputSchema: {
        associatedTable: z
          .string()
          .describe("Associated dynamic table name"),
      },
    },
    async ({ associatedTable }) => {
      try {
        return jsonResult(
          await client.listExportTemplateOptions(associatedTable),
        );
      } catch (err) {
        return errorResult(err);
      }
    },
  );

  return server;
}
