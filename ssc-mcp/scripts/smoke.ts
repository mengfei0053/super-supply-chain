/**
 * Smoke test against live SSC API (not over MCP stdio).
 * Uses SSC_TOKEN or SSC_USERNAME+SSC_PASSWORD.
 * Does not print tokens.
 */
import {
  SscClient,
  SscApiError,
  loadConfigFromEnv,
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
      "FAIL: set SSC_TOKEN or SSC_USERNAME+SSC_PASSWORD before smoke test",
    );
    process.exit(2);
  }

  try {
    if (!hasToken && hasPass) {
      const user = await client.login(cfg.username!, cfg.password!);
      console.log("login: OK as", user.username, "(token not printed)");
    } else {
      await client.ensureAuth();
      console.log("token: present (not printed)");
    }

    const menus = await client.getMenus();
    const n = Array.isArray(menus) ? menus.length : -1;
    console.log("GET /api/admin/menus: OK, count =", n);

    const orders = await client.listOrders();
    const on = Array.isArray(orders) ? orders.length : -1;
    console.log("GET /api/admin/settlement-form-entry: OK, count =", on);

    const dicts = await client.listDicts([0, 5]);
    const dn = Array.isArray(dicts) ? dicts.length : -1;
    console.log("GET /api/admin/dict-manage: OK, page size =", dn);

    console.log("SMOKE: PASS");
  } catch (err) {
    if (err instanceof SscApiError) {
      console.error("SMOKE: FAIL", err.message);
      // Avoid dumping large bodies; truncate
      console.error("body:", err.body.slice(0, 300));
    } else {
      console.error("SMOKE: FAIL", err);
    }
    process.exit(1);
  }
}

main();
