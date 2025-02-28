import { Admin, Resource } from "react-admin";
import { dataProvider } from "./dataProvider";
import { authProvider } from "./authProvider";
import Dashboard from "./pages/dashboard";
import { Layout } from "./layout";
import SettlementFormEntry from "./pages/settlement-form-entry";
import toolPages from "./pages/tool-pages/excel-mapping-rules";
import ExcelReadRules from "./pages/tool-pages/excel-read-rules";
import { Route } from "react-router-dom";
import ExcelPage from "./pages/excels";
import dictManage from "./pages/tool-pages/dict-manage";
import ExcelEditPage from "./pages/excels/Edit";
import YifanCostCalculation from "./pages/yifan/cost-calculation";

export const App = () => (
  <Admin
    layout={Layout}
    dashboard={Dashboard}
    dataProvider={dataProvider}
    authProvider={authProvider}
  >
    <Resource name="settlement-form-entry" {...SettlementFormEntry}></Resource>
    <Resource name="excel-mapping-rule" {...toolPages}></Resource>
    <Resource name="excel-read-rules" {...ExcelReadRules}></Resource>
    <Resource name="dict-manage" {...dictManage}></Resource>
    <Resource
      name="yifan/cost-calculation"
      {...YifanCostCalculation}
    ></Resource>
    <Resource name="excel">
      <Route path=":tableName" element={<ExcelPage></ExcelPage>}></Route>
      <Route
        path=":tableName/:id"
        element={<ExcelEditPage></ExcelEditPage>}
      ></Route>
    </Resource>
  </Admin>
);
