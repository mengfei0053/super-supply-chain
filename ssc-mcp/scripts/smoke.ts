/**
 * Smoke test against live SSC API (not over MCP stdio).
 * Uses SSC_TOKEN or SSC_USERNAME+SSC_PASSWORD.
 * Does not print tokens.
 */
import {
  SscClient,
  SscApiError,
  loadConfigFromEnv,
  redactSecrets,
} from "../src/client.js";

async function main(): Promise<void> {
  const cfg = loadConfigFromEnv();
  const client = new SscClient(cfg);
  console.log("SSC_BASE_URL:", client.baseUrl);

  const hasToken = Boolean(cfg.token);
  const hasPass = Boolean(cfg.username && cfg.password);
  console.log("Auth mode:", hasToken ? "token" : hasPass ? "login" : "none");

  if (!hasToken && !hasPass) {
    console.error(
      "FAIL: set SSC_TOKEN or SSC_USERNAME+SSC_PASSWORD before smoke test. See README.md Manual verification. Tokens are not printed.",
    );
    process.exit(2);
  }

  try {
    if (!hasToken && hasPass) {
      const user = await client.login(cfg.username!, cfg.password!);
      console.log("login: OK as", user.username, "(token not printed)");
    } else {
      await client.ensureAuth();
      console.log(
        "token: present (not printed). Set SSC_SMOKE_EXPORT_* to verify it against the API.",
      );
    }

    const exportTable = process.env.SSC_SMOKE_EXPORT_TABLE?.trim();
    const exportIds = process.env.SSC_SMOKE_EXPORT_IDS?.trim();
    const exportType = process.env.SSC_SMOKE_EXPORT_TYPE?.trim();
    if (exportTable && exportIds && exportType) {
      const file = await client.exportExcel({
        tableName: exportTable,
        ids: exportIds,
        type: exportType,
      });
      console.log(
        "GET excel-exports: OK",
        JSON.stringify({
          path: file.path,
          bytes: file.bytes,
          fileName: file.fileName,
          xlsx: file.xlsx,
        }),
      );
    } else {
      console.log(
        "export smoke skipped (set SSC_SMOKE_EXPORT_TABLE, SSC_SMOKE_EXPORT_IDS, SSC_SMOKE_EXPORT_TYPE)",
      );
    }

    const uploadFile = process.env.SSC_SMOKE_UPLOAD_FILE?.trim();
    const uploadTable = process.env.SSC_SMOKE_UPLOAD_TABLE?.trim();
    if (uploadFile && uploadTable) {
      const uploaded = await client.uploadExcel({
        tableName: uploadTable,
        filePath: uploadFile,
        name:
          process.env.SSC_SMOKE_UPLOAD_NAME?.trim() ||
          uploadFile.split(/[/\\]/).pop() ||
          "upload.xlsx",
      });
      console.log(
        "POST excel upload: OK",
        JSON.stringify({
          tableName: uploaded.tableName,
          fileName: uploaded.fileName,
          bytes: uploaded.bytes,
        }),
      );
    } else {
      console.log(
        "upload smoke skipped (set SSC_SMOKE_UPLOAD_TABLE and SSC_SMOKE_UPLOAD_FILE; this inserts a row)",
      );
    }

    console.log("SMOKE: PASS");
  } catch (err) {
    if (err instanceof SscApiError) {
      console.error("SMOKE: FAIL", err.message);
      // Avoid dumping large bodies; truncate
      console.error("body:", redactSecrets(err.body.slice(0, 300)));
    } else {
      console.error("SMOKE: FAIL", err);
    }
    process.exit(1);
  }
}

main();
