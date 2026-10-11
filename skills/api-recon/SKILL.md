---
name: api-recon
description: English textAPIEnglish textskill。
---

# API Recon（English text）

English text**English text**English text，English text：**English text API**（English text、English text、English text、English text）、**English text**、**UI English text**（Tab、English text、English text）。

---

## English text（Agent English text · English text）

English text skill **English text API / English text**，English text。

### English text

| English text | English text | English text |
|---|---|---|
| **English text** | English text path、method、English text、English text、UI English text | SQLi/XSS/English text/English text/fuzz English text、English text、English text |
| **English text** | Hook + stub/mock English text**English text**English text | English text；English text |
| **English text** | English text hook English text，English text mock English text SPA English text | English text |

### English text（Phase 3 English text）

1. English text `preload.js` / `runtime_harvest.js` **English text stub** English text、English text、English text bootstrap English text；
2. English text **English text、English text、English text** English text mock body；
3. English text 401 English text，English text XHR/fetch/WebSocket；
4. **English text、English text、English text UI English text**——English text。

**English text**：English text mock English text，**English text outbound English text**；English text，English text**English text**。

### English text

| English text | English text |
|---|---|
| Phase 1 English text grep/curl/Read English text entry `index-*.js` English text API path | English text `OUTDIR/harvest_static.py` |
| English text `extract_apis.py` English text harvest English text | English text `OUTDIR/harvest_static.py` English text |
| English text grep/English text ≥2 English text | English text：English text tool_logs、English text harvest、English text reference |
| English text A/B，English text `scripts/` English text | English text OUTDIR English text |
| English text/English text、OTP、OAuth English text | stub/mock（English text） |
| English text「English text」English text stub，English text/English text | English text outbound，English text recon English text |
| English text、English text、English text | coverage English text |
| English text runtime + English text，English text | English text「English text」English text |
| English text + diff，English text | Phase 3b English text + Phase 5 diff |
| English text runtime English text/English text | English text diff English text/English text |

---

## English text + English text

| English text | English text | English text |
|---|---|---|
| **English text**（JS bundle） | English text endpoint English text、English text、English text | English text HTTP English text；English text Phase 1b；English text URL |
| **English text**（English text） | English text + body + English text + English text URL + WS/SSE；English text diff English text | English text；English text/English text |

| English text | English text | English text |
|---|---|---|
| **depth** | `runtime_harvest.js`（Puppeteer） | API English text、METHOD/params/English text、WS/SSE、English text |
| **coverage** | browser + `preload.js` | English text Tab/English text/English text，English text |
| **both** | English text depth English text coverage | English text，English text |

**English text**（English text）：path English text harvest/English text；English text **English text + UI English text + English text diff + English text**（grep English text [reference.md](reference.md) J English text）。

---

## English text

English text recon English text：

- [ ] **English text**：Phase 1 harvest English text `api_static.txt`、`routes.txt`、`js/`
- [ ] **English text**：English text depth English text coverage English text；coverage/both English text **Hook English text + English text**
- [ ] **English text**：English text path English text `/login`（English text hash English text）
- [ ] **English text**：coverage/both English text + `param_samples.json`；Phase 5 English text `params_merged.json`
- [ ] **English text**（English text）：Phase 4 English text，English text **module English text API**（English text locale/bootstrap）
- [ ] **English text**：Phase 5 English text（English text Phase 5 English text）；`insert_assets` English text

---

## English text

`scripts/` English text，**English text**English text。

**English text**：English text → English text → English text `OUTDIR`（English text `recon/`）→ English text `CHANGES.md`；English text，English text。

| English text | English text | English text → OUTDIR English text | English text |
|---|---|---|---|
| **A（English text）** | Phase 0 English text、**English text**English text harvest/spider English text | `harvest_static.py` / `spider_mpa.py` | **English text regex English text**；English text manifest/English text endpoint English text、webpack/Vite `publicPath`、MPA exclude/cookie |
| **B（English text）** | Phase 2 English text、English text depth/coverage English text | `runtime_harvest.js` / `preload.js` + `config.json` | Cookie/localStorage English text、neutralize English text、stubs、login English text、api English text、hash/history |

**SPA English text**（English text；Phase English text「English text」）：

| English text | English text | English text |
|---|---|---|
| Phase 0 English text | English text Bash = `python3 OUTDIR/harvest_static.py <URL> OUTDIR` | curl/grep/Read English text entry `index-*.js`（English text >500KB） |
| English text A | English text → English text → **English text** | English text API English text harvest |
| Phase 1 English text | `wc -l` English text；404 English text harvest English text | English text extract English text；English text URL English text grep |
| Phase 1b English text | grep English text `OUTDIR/js/*.js` | English text bundle English text harvest |

- ✅ English text `harvest_static.py` → （English text）English text regex → **English text**
- ❌ curl English text bundle → grep English text → English text extract → English text harvest
- **MPA**：Phase 0 English text Bash = `python3 OUTDIR/spider_mpa.py ...`

---

## English text

| English text | English text |
|---|---|
| English text | >100KB English text `index-*.js` **English text** Read/grep English text；English text OUTDIR English text |
| grep English text | English text `\| head -20` English text `-m 5`；English text path English text，English text bundle English text |
| English text | English text `wc -l`、`ls \| wc -l`；English text Read English text |
| regex English text | English text、≤1 English text、English text ≤50KB English text chunk English text HTML；English text harvest English text |
| reference | English text/English text/English text [reference.md](reference.md)，English text inline English text |

---

## English text

```
Phase 0 English text + OUTDIR
  → English text A → Phase 1 harvest（★ English text ★）
  → Phase 1b English text
  → Phase 2 English text → config.json
  → English text B → Phase 3 English text + English text
  → Phase 4 English text（English text）→ English text Phase 3
  → Phase 5 English text + insert_assetsEnglish text、English textapiEnglish text，English text
```

English text；**English text Phase**。

1. [ ] **Phase 0**：English text SPA/MPA；English text `OUTDIR` → [Phase 0](#phase-0--English text)
2. [ ] **English text A + Phase 1**：English text → **English text** harvest → `wc -l` English text → [Phase 1](#phase-1--English text)
3. [ ] **Phase 1b**：English text + English text → `param_candidates.json` → [Phase 1b](#phase-1b--English text)
4. [ ] **Phase 2**：English text → `config.json` → [Phase 2](#phase-2--English text)
5. [ ] **English text B**：English text runtime English text → [Phase 3](#phase-3--English text)
6. [ ] **Phase 3**：depth / coverage / both；English text；English text → `param_samples.json`
7. [ ] **Phase 4**（English text）：English text → patch stubs → English text Phase 3 → [Phase 4](#phase-4--English text)
8. [ ] **Phase 5**：English text + English text + `insert_assets` → [Phase 5](#phase-5--English text)

---

## Phase 0 — English text

English text HTML，**English text `OUTDIR`**（English text skill English text `scripts/`）：

- **SPA**：English text + `<div id=app>` + chunk → Phase 1–5
- **MPA**：SSR + `<form>`、English text endpoint bundle → English text A English text：

```bash
python3 recon/spider_mpa.py <BASE_URL> <OUTDIR> [--cookie "session=..."] [--max 300] [--depth 5] [--exclude "logout|delete|destroy"]
```

English text `forms.txt`、`links.txt`、`api_inline.txt`。SPA English text forms ≈ 0 → English text Phase 1。

---

## Phase 1 — English text

English text [English text](#English text) · [English text](#English text)。

```bash
python3 recon/harvest_static.py <BASE_URL> <OUTDIR>
```

harvest：English text HTML script → webpack/Vite manifest → English text lazy chunk → English text `js/`、`api_static.txt`、`routes.txt`、`chunkmap.txt`。

```bash
wc -l OUTDIR/api_static.txt OUTDIR/routes.txt
ls OUTDIR/js | wc -l
```

- chunk English text vs manifest：404 English text harvest English text，English text curl English text chunk
- `api_static.txt` English text → English text OUTDIR English text endpoint English text（English text reference）

### Phase 1b — English text

path English text Phase 1；English text recon。grep English text [English text](#English text)。

**English text**：English text——English text、English text、English text、English text、English text、English text。

#### 1b.0 — English text

| English text | English text | English text |
|---|---|---|
| REST JSON | body + query | path English text `(params\|data\|body)\s*:\s*\{` |
| GraphQL | `variables` | gql English text、`$page: Int` |
| English text form | urlencoded | `<form>`、`FormData` |
| English text | multipart | `FormData.append` |
| English text | `/user/:id` | English text + `useParams` / `$route.params` |
| English text/English text | English text `sign`/`data` | Hook English text（reference D English text） |

English text：English text `transport: query|json|form|graphql|encrypted`。

#### 1b.1 — English text

English text path English text，English text：

```bash
grep -n '"/api/user/list"' OUTDIR/js/*.js | head -20
grep -rhoaE '.{0,120}("/api[^"]+").{0,200}' OUTDIR/js/*.js | head -20
grep -rhoaE '(params|data|body|payload)\s*:\s*\{' OUTDIR/js/*.js | head -20
```

| English text | English text |
|---|---|
| axios English text | `data` / `params` |
| English text request | English text |
| OpenAPI English text | English text method English text |
| React Query / SWR | hook English text |
| Vue composable | composable English text |

English text：`yup`/`zod`/rules、`Form.Item name=`、English text Swagger。

→ `param_candidates.json`：`{ path, fields[], source: "static-callsite", confidence }`

#### 1b.2 — English text

```
Form field → onFinish/handleSubmit → transform → API payload
```

| English text | English text |
|---|---|
| English text submit | English text submit → transform → API |
| English text | `getFieldsValue()` → `params` |
| English text | `:id` / `?tab=` |
| English text | English text `tenantId`、English text、sign |
| English text select | `options` → API English text |

DevTools call stack English text `fetch`/`XHR.send` English text。

#### 1b.3 — English text（≠ Phase 2 English text）

| English text | English text |
|---|---|
| **English text** | payload English text build、transform English text |
| **English text** | required、pattern、enum |
| **English text** | path / query / body / multipart / English text |

English text（Phase 2）English text（Authorization、`X-Tenant-Id`、sign）。

#### 1b.4 — English text Phase 3 English text

English text/English text；**English text/English text/English text**English text Phase 3 English text + diff + Phase 5 English text。

---

## Phase 2 — English text

English text `OUTDIR/js/` grep（English text `head`），English text `config.json`（English text reference）：

| English text | English text | English text |
|---|---|---|
| **English text** | English text？ | `isLogin`、`getToken`、Cookie/localStorage |
| **English text** | English text `/login`？ | `response_code`、`errno`、axios interceptor |
| **English text** | English text/English text？ | `menu`、`permission`、`role`、`acl`、`routes` |

English text localStorage English text——English text chunk/English text。

**English text = English text B**：English text `config.json`，English text `OUTDIR/runtime_harvest.js` / `preload.js`。

### Phase 2b — API English text（English text）

English text OUTDIR English text `preload.js` English text、Authorization、English text API URL：

| English text | English text |
|---|---|
| `recordDetail: true` | `__API_RECON_DETAIL__` |
| `observe.xhrHeaders: true` | headers English text |
| `extractUrlsFromResponse: true` | English text API |
| `observe.storageReads/cookieReads: true` | English text config |
| `neutralizeVueRouter: true` | `__API_RECON_ROUTES__` |

coverage English text：`__API_RECON_LOG__`、`__API_RECON_DETAIL__`、`__API_RECON_ROUTES__`、`__API_RECON_OBSERVE__`。

---

## Phase 3 — English text

English text B；English text [English text](#English textagent-English text--English text) · English text mock English text。

`config.json` English text `"runtimeMode": "depth" | "coverage" | "both"`（English text reference）。

### Hook English text stub（depth + coverage English text）

| English text | English text | English text |
|---|---|---|
| L1 English text | auth/English text/bootstrap stub | English text |
| L2 English text | English text JSON English text | English text → English text |
| L3 English text | English text L1 English text `/api` English text | English text，English text UI |

- **depth**：fake auth + `forward` English text + `stubs`；English text `routes`（hash/history）；English text `runtime_api.json`
- **coverage**：**document-start** English text `preload.js`（CDP `addScriptToEvaluateOnNewDocument` English text Userscript）

English text：`window.__API_RECON_PRELOAD__` English text；English text path English text `/login`。

```bash
cd recon && npm install
node runtime_harvest.js config.json
```

### 3b — coverage English text（English text）

1. English text/English text — English text，English text 1–3s
2. Tab — `role=tab`、`.ant-tabs-tab`
3. English text — English text/English text/English text
4. English text — English text、English text、English text（**English text**）
5. English text — English text API/English text
6. SPA — English text `routes.txt` English text path English text `pushState`（MPA English text）

**English text**（English text）：English text，**diff English text**：

| English text | English text |
|---|---|
| English text | English text + English text |
| English text | keyword、filter |
| English text | English text optional |
| English text/English text | English text entity |
| English text/English text/English text | `ids[]`、`exportType`、`sortField` |

**stub English text outbound body/headers English text**——English text。English text → `scan_raw.json`、`param_samples.json`、`api_detail.json`。

- **Vue**：`neutralizeVueRouter: true` + document-start preload
- **React**：`routes.txt` + English text + `pushState`
- **both**：English text 3a depth，English text 3b coverage

---

## Phase 4 — English text

**English text**：English text / English text bootstrap（English text locale）→ English text。

| English text | English text |
|---|---|
| English text | English text + English text |
| English text/English text | stub shape English text |
| English text API English text | `v-if permission` English text |
| `routes.txt` English text bundle | English text auth English text |

```bash
grep -rhoaE '"/api[^"]*(permission|perm|role|menu|acl)[^"]*"' OUTDIR/js/*.js | sort -u | head -30
grep -rhoaE 'userRouteAuth|getResultTree|routeMap|routeLink|menuList|authList' OUTDIR/js/*.js | head -20
```

English text：`role_permissions`（flat codes）+ `permissions/all`（tree）→ `getResultTree` → `userRouteAuth[CODE].url`。

```bash
python3 recon/extract_route_map.py recon/js recon/
python3 recon/build_perm_tree.py recon/js recon/ --config recon/config.json
```

English text：`route_map.json`、`userRouteAuth.json`、`permissions_tree.json`、`*_stub.json`、`perm_codes_all.txt`。

stub English text：English text `response_code` English text；flat codes English text tree English text；`routes` English text `route_map` English text link。

English text `config.json` English text**English text Phase 3**。English text SPA English text `waitUntil`、`routeTimeout`、`perRouteMs`（English text reference A3/I English text）。

---

## Phase 5 — English text

### English text

| English text | English text | English text |
|---|---|---|
| `js/`、`api_static.txt`、`routes.txt`、`chunkmap.txt` | 1 | English text bundle English text path |
| `param_candidates.json` | 1b | English text |
| `config.json` | 2 | English text + runtime English text |
| `runtime_api.json` | 3a | depth English text（English text WS/SSE） |
| `param_samples.json`、`scan_raw.json`、`api_detail.json` | 3b | English text、English text、detail |
| `route_map.json` English text | 4 | English text（English text） |
| `params_merged.json` | 5 | English text + English text |
| `api_merged.txt` | 5 | `METHOD /path [params] [static\|runtime\|both]` |
| `site_map.json` | 5 | English text、API、params、English text、English text |
| **insert_assets** | 5 | English text、English text |

### 5b — English text

English text `param_samples.json` diff，**English text**。English text reference J7（English text/English text/English text/English text）。

### 5c — English text

English text 400（**English text recon，English text**）：`field 'x' is required`、English text。English text `data` English text、`variables`、English text `bizData`。

English text：runtimeMode、English text/English text API English text、English text、English text、English text `CHANGES.md` English text。

`site_map.json` English text：

```json
{
  "site": "https://example.com",
  "runtimeMode": "both",
  "appType": "vue-spa",
  "routeGuardStrategy": ["nav-neutralize", "L1-auth", "L2-patch", "forward"],
  "apisFromStatic": [],
  "apisFromRuntime": [],
  "apis": [],
  "params": [{ "method": "POST", "path": "/api/user/list", "transport": "json", "fields": [] }],
  "frontendRoutes": [],
  "routesVerifiedByClick": [],
  "featuresTriggered": [],
  "limitations": ""
}
```

English text grep English text [reference.md](reference.md)。

---

## English text

- **English text**：webpack/Vite/Angular lazy load English text
- **English text**：REST/JSON、GraphQL、WebSocket、SSE；gRPC-web English text
- **SSR**：English text fetch English text；RSC/Server Actions English text
- **English text**：JSVMP、WASM、HMAC/mTLS English text → English text + English text
- **English text**：English text、hidden params、WASM English text → 「English text」/「English text」
- **English text**：runtime English text endpoint

---

## English text

- Grep English text、`config.json` English text、English text、Hook、English text J English text、site_map English text：**[reference.md](reference.md)**
- English text [English text](#English text) English text
