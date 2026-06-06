# Graph Report - .  (2026-06-03)

## Corpus Check
- Corpus is ~19,711 words - fits in a single context window. You may not need a graph.

## Summary
- 498 nodes · 694 edges · 60 communities (34 shown, 26 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 76 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Invoice Templates|Invoice Templates]]
- [[_COMMUNITY_Dynamic Excel API|Dynamic Excel API]]
- [[_COMMUNITY_Frontend Package Scripts|Frontend Package Scripts]]
- [[_COMMUNITY_Backend Utility Tests|Backend Utility Tests]]
- [[_COMMUNITY_Excel Export Utilities|Excel Export Utilities]]
- [[_COMMUNITY_Frontend Dependencies|Frontend Dependencies]]
- [[_COMMUNITY_App TypeScript Config|App TypeScript Config]]
- [[_COMMUNITY_CLI Login Flow|CLI Login Flow]]
- [[_COMMUNITY_Backend Startup Config|Backend Startup Config]]
- [[_COMMUNITY_Dictionary API|Dictionary API]]
- [[_COMMUNITY_Admin Auth Middleware|Admin Auth Middleware]]
- [[_COMMUNITY_WebDAV Upload Tests|WebDAV Upload Tests]]
- [[_COMMUNITY_Invoice Generator Tests|Invoice Generator Tests]]
- [[_COMMUNITY_Frontend App Pages|Frontend App Pages]]
- [[_COMMUNITY_Node TypeScript Config|Node TypeScript Config]]
- [[_COMMUNITY_Excel Read Rules API|Excel Read Rules API]]
- [[_COMMUNITY_Frontend Auth Data|Frontend Auth Data]]
- [[_COMMUNITY_Tool Page Fields|Tool Page Fields]]
- [[_COMMUNITY_Input Export Forms|Input Export Forms]]
- [[_COMMUNITY_Web Manifest|Web Manifest]]
- [[_COMMUNITY_Auth Controllers|Auth Controllers]]
- [[_COMMUNITY_Navigation Menu|Navigation Menu]]
- [[_COMMUNITY_Clearance Pricing Model|Clearance Pricing Model]]
- [[_COMMUNITY_Read Rule Models|Read Rule Models]]
- [[_COMMUNITY_Settlement List UI|Settlement List UI]]
- [[_COMMUNITY_Export Rule Management|Export Rule Management]]
- [[_COMMUNITY_Settlement Entry UI|Settlement Entry UI]]
- [[_COMMUNITY_Select Options API|Select Options API]]
- [[_COMMUNITY_Accounts Model|Accounts Model]]
- [[_COMMUNITY_Dynamic Menus API|Dynamic Menus API]]
- [[_COMMUNITY_Companies Model|Companies Model]]
- [[_COMMUNITY_Dynamic Excel Model|Dynamic Excel Model]]
- [[_COMMUNITY_Generic Map Helper|Generic Map Helper]]
- [[_COMMUNITY_Companies API|Companies API]]
- [[_COMMUNITY_Dictionary Model|Dictionary Model]]
- [[_COMMUNITY_Export Template Model|Export Template Model]]
- [[_COMMUNITY_Freight Model|Freight Model]]
- [[_COMMUNITY_Shipping Order Model|Shipping Order Model]]
- [[_COMMUNITY_Upload File Model|Upload File Model]]
- [[_COMMUNITY_JSON Logging Helper|JSON Logging Helper]]
- [[_COMMUNITY_Safe Array Helper|Safe Array Helper]]
- [[_COMMUNITY_Claude Permissions|Claude Permissions]]
- [[_COMMUNITY_Date Range Input|Date Range Input]]
- [[_COMMUNITY_Settlement Create UI|Settlement Create UI]]
- [[_COMMUNITY_Dev Script|Dev Script]]
- [[_COMMUNITY_Rule Input|Rule Input]]
- [[_COMMUNITY_Data Input|Data Input]]
- [[_COMMUNITY_Frontend Claude Settings|Frontend Claude Settings]]
- [[_COMMUNITY_TypeScript References|TypeScript References]]
- [[_COMMUNITY_Run Script|Run Script]]
- [[_COMMUNITY_Settlement Create Page|Settlement Create Page]]
- [[_COMMUNITY_Order Model|Order Model]]
- [[_COMMUNITY_Vite Alias Config|Vite Alias Config]]
- [[_COMMUNITY_Product Info Model|Product Info Model]]
- [[_COMMUNITY_User Fixture|User Fixture]]

## God Nodes (most connected - your core abstractions)
1. `compilerOptions` - 19 edges
2. `httpClient()` - 12 edges
3. `T` - 11 edges
4. `setupInvoiceGeneratorTest()` - 10 edges
5. `T` - 10 edges
6. `GetListQueryParams()` - 10 edges
7. `CreateCostCalculation()` - 10 edges
8. `GetCompanyInfo()` - 10 edges
9. `GetExcelExportFilePath()` - 9 edges
10. `UploadToNas()` - 9 edges

## Surprising Connections (you probably didn't know these)
- `GetDicts()` --calls--> `SetContentRange()`  [INFERRED]
  backend/controllers/base-dict.go → backend/utils/SetContentRange.go
- `GetDynamicExcelTableList()` --calls--> `GetListQueryParams()`  [INFERRED]
  backend/controllers/dynamic-excel-tables.go → backend/utils/GetQueryParams.go
- `CreateDynamicExcelTable()` --calls--> `GetBaoguanDan()`  [INFERRED]
  backend/controllers/dynamic-excel-tables.go → backend/utils/GetBaoguanDan.go
- `CreateDynamicExcelTable()` --calls--> `GetExcelData()`  [INFERRED]
  backend/controllers/dynamic-excel-tables.go → backend/utils/ParseExcel.go
- `CreateDynamicExcelTable()` --calls--> `UploadToNas()`  [INFERRED]
  backend/controllers/dynamic-excel-tables.go → backend/utils/upload-to-nas.go

## Import Cycles
- None detected.

## Communities (60 total, 26 thin omitted)

### Community 0 - "Invoice Templates"
Cohesion: 0.12
Nodes (29): DynamicExcelTable, ExcelData, DynamicExcelTable, T, BaseCompaniesInfos, generateFileFromChangjiuTemplate(), generateFileFromInvoiceTemplate(), generateFileFromTemplate() (+21 more)

### Community 1 - "Dynamic Excel API"
Cohesion: 0.12
Nodes (23): Context, Context, Context, Context, CreateDynamicExcelTable(), DeleteDynamicExcelTable(), ExportDynamicExcel(), GetDynamicExcelTableDetail() (+15 more)

### Community 2 - "Frontend Package Scripts"
Cohesion: 0.08
Nodes (25): devDependencies, eslint, eslint-config-prettier, eslint-plugin-react, eslint-plugin-react-hooks, prettier, @types/node, @types/qs (+17 more)

### Community 3 - "Backend Utility Tests"
Cohesion: 0.15
Nodes (21): T, ExcelData, ExcelData, TestGenericHelpers(), TestGetListQueryParamsKeepsRepeatedRangeCompatibility(), TestGetListQueryParamsParsesReactAdminQuery(), TestGetPkgCountExtractsContainerCount(), TestGetPriceExtractsCurrencyCode() (+13 more)

### Community 4 - "Excel Export Utilities"
Cohesion: 0.15
Nodes (17): Context, DynamicExcelTable, ExportExcel(), SingleExportExcel(), ExportExcelReq, CostType, FreightInfo, PortInfo (+9 more)

### Community 5 - "Frontend Dependencies"
Cohesion: 0.10
Nodes (21): dependencies, dayjs, @emotion/react, @emotion/styled, localforage, @mui/icons-material, @mui/material, @mui/x-date-pickers (+13 more)

### Community 6 - "App TypeScript Config"
Cohesion: 0.10
Nodes (20): compilerOptions, allowImportingTsExtensions, composite, isolatedModules, jsx, lib, module, moduleDetection (+12 more)

### Community 7 - "CLI Login Flow"
Cohesion: 0.21
Nodes (19): claimsPayload, config, loginRequest, loginResponse, checkAuthenticatedEndpoint(), defaultConfigPath(), envOrDefault(), Time (+11 more)

### Community 8 - "Backend Startup Config"
Cohesion: 0.15
Nodes (13): Engine, main(), HandlerFunc, Logger, Logger, IsDev(), LoadConfigFile(), LoadStatic() (+5 more)

### Community 9 - "Dictionary API"
Cohesion: 0.18
Nodes (16): Context, Context, CreateDict(), DeleteDict(), GetDictDeltail(), GetDictMap(), GetDicts(), UpdateDict() (+8 more)

### Community 10 - "Admin Auth Middleware"
Cohesion: 0.19
Nodes (15): Context, HandlerFunc, Engine, T, Time, T, AuthMiddleware(), handleNoAuth() (+7 more)

### Community 11 - "WebDAV Upload Tests"
Cohesion: 0.30
Nodes (12): T, T, findRepoDotEnv(), loadDotEnvOrSkip(), TestUploadToNasEndToEndAgainstDotEnvWebDAV(), TestUploadToNasReturnsErrorOnServerFailure(), TestUploadToNasReturnsErrorWhenLocalFileMissing(), TestUploadToNasUploadsFileAndReturnsURL() (+4 more)

### Community 12 - "Invoice Generator Tests"
Cohesion: 0.46
Nodes (13): T, File, assertCell(), chdirForInvoiceTemplate(), openWorkbook(), seedChengxinBaoguandan(), seedInvoiceCompanies(), setupInvoiceGeneratorTest() (+5 more)

### Community 13 - "Frontend App Pages"
Cohesion: 0.20
Nodes (5): Client, IDashboardProps, Layout(), App(), router

### Community 14 - "Node TypeScript Config"
Cohesion: 0.18
Nodes (10): compilerOptions, allowSyntheticDefaultImports, composite, module, moduleResolution, noEmit, skipLibCheck, strict (+2 more)

### Community 15 - "Excel Read Rules API"
Cohesion: 0.29
Nodes (9): Context, CreateExcelReadRules(), DeleteExcelReadRules(), GetExcelReadRule(), GetExcelReadRulesList(), UpdateExcelReadRules(), ExcelReadRuleItem, MapRule (+1 more)

### Community 16 - "Frontend Auth Data"
Cohesion: 0.24
Nodes (5): CreateTemplateParams, authProvider, User, baseDataProvider, dataProvider

### Community 19 - "Input Export Forms"
Cohesion: 0.36
Nodes (4): AdSelectInput(), IAdSelectInputProps, ToolPageCreate(), httpClient()

### Community 21 - "Web Manifest"
Cohesion: 0.25
Nodes (7): background_color, display, icons, name, short_name, start_url, theme_color

### Community 22 - "Auth Controllers"
Cohesion: 0.33
Nodes (6): Context, Login(), Register(), Claims, Credentials, RegisteredClaims

### Community 23 - "Navigation Menu"
Cohesion: 0.29
Nodes (3): IMenuProps, MenuName, Props

### Community 24 - "Clearance Pricing Model"
Cohesion: 0.40
Nodes (5): Model, UpdateDb(), ClearancePriceBase, CostType, WholeOrBulkType

### Community 25 - "Read Rule Models"
Cohesion: 0.60
Nodes (5): Model, ExcelReadRuleInfos, IterateRule, MappingRule, Rules

### Community 29 - "Select Options API"
Cohesion: 0.60
Nodes (4): Context, ExcelMappingRuleOptions, GetExportTemplates(), GetOptions()

### Community 30 - "Accounts Model"
Cohesion: 0.40
Nodes (3): Model, Time, BaseAccountsInfos

### Community 31 - "Dynamic Menus API"
Cohesion: 0.50
Nodes (3): Context, DynamicExcelMenu, GetDynamicExcelMenus()

### Community 32 - "Companies Model"
Cohesion: 0.50
Nodes (3): Model, Time, BaseCompaniesInfos

### Community 33 - "Dynamic Excel Model"
Cohesion: 0.67
Nodes (3): Model, DynamicExcelTable, ExcelData

### Community 34 - "Generic Map Helper"
Cohesion: 0.50
Nodes (3): T, U, Map()

## Knowledge Gaps
- **151 isolated node(s):** `allow`, `Credentials`, `RegisteredClaims`, `Context`, `MapRule` (+146 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **26 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `GetUploadTmpDir()` connect `Dynamic Excel API` to `Invoice Templates`, `Backend Utility Tests`, `Excel Export Utilities`?**
  _High betweenness centrality (0.055) - this node is a cross-community bridge._
- **Why does `GetExcelExportFilePath()` connect `Invoice Templates` to `Dynamic Excel API`?**
  _High betweenness centrality (0.054) - this node is a cross-community bridge._
- **Why does `setupInvoiceGeneratorTest()` connect `Invoice Generator Tests` to `Admin Auth Middleware`?**
  _High betweenness centrality (0.050) - this node is a cross-community bridge._
- **What connects `allow`, `Credentials`, `RegisteredClaims` to the rest of the system?**
  _151 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Invoice Templates` be split into smaller, more focused modules?**
  _Cohesion score 0.11764705882352941 - nodes in this community are weakly interconnected._
- **Should `Dynamic Excel API` be split into smaller, more focused modules?**
  _Cohesion score 0.1168091168091168 - nodes in this community are weakly interconnected._
- **Should `Frontend Package Scripts` be split into smaller, more focused modules?**
  _Cohesion score 0.07692307692307693 - nodes in this community are weakly interconnected._