# api-recon — English text

Grep English text、`config.json` English text。English text grep English text `js/` English text。bundle English text `js-beautify` English text `sed 's/}/}\n/g'`，English text raw grep English text。

## English text

`scripts/` English text**English text**，English text。English text：

| English text | English text |
|---|---|
| `harvest_static.py` | endpoint English text、webpack/Vite manifest English text、English text publicPath、English text/English text |
| `runtime_harvest.js` | neutralize English text、stub English text body English text、routes English text、WS English text、`waitUntil`/`routeTimeout`/`proxy` |
| `preload.js` | `loginPathRe`、L1 stubs、`neutralize.fields`、`apiPattern`、English text L3、`recordDetail`、`observe.*`、`neutralizeVueRouter` |
| `spider_mpa.py` | `--exclude` English text、cookie、depth/max、English text |
| `extract_route_map.py` | `routeMap` / `routeLink` English text、KEY English text |
| `build_perm_tree.py` | `userRouteAuth` English text、`ROOTS`/`PREFIX_PARENT` English text、stub English text |
| `config.json` | English text |

English text（English text `recon/`），English text。

---

## A. English text

### A1. English text — 「English text？」

```bash
grep -rhoaE '.{0,40}(isLogin|isAuthenticated|loggedIn|hasLogin|requireAuth)\b.{0,80}' js | head
grep -rhoaE 'function (getUser|getToken|getAuth)[0-9]?\([^)]*\)\{.{0,200}' js | head
grep -rhoaE '(localStorage|sessionStorage)\.getItem\("[^"]+"\)' js | sort -u
grep -rhoaE '(Cookies?|cookie)\.(get|load)\("[^"]+"\)' js | sort -u
grep -rhoaE '\batob\(|JSON\.parse\(|jwt|decode' js | head
```

English text `isLogin = f(getUser())` → `getUser = decode(storage.read(KEY))`，English text **English text**、**English text**（Cookie vs localStorage）、**English text**：

| English text | config English text |
|---|---|
| English text / `"1"` / token | `"value": "anything-truthy"` |
| `JSON.parse(x)` | `"value": "json:{\"id\":1,\"username\":\"admin\"}"` |
| `JSON.parse(atob(x))` | `"value": "b64json:{\"id\":1,\"username\":\"admin\"}"` |
| JWT | English text/`alg:none` JWT，English text bundle English text |
| English text（SM2/AES/RSA） | English text；English text blob English text forge；English text |

→ English text `cookies` / `localStorage`。

### A2. English text — 「English text /login？」

```bash
grep -rhoaE '.{0,60}(interceptors\.response|axios|request\.use).{0,120}' js | head
grep -rhoaE '.{0,40}(response_code|errcode|errno|\bcode\b|\bret\b|\bstatus\b)\s*[=!]==?\s*[\-0-9]{1,4}.{0,60}' js | head -20
grep -rhoaE '.{0,40}(English text|English text|English text|unauthorized|English text|English text|token.{0,10}invalid).{0,40}' js | head
grep -rhoaE '.{0,30}(location\.href|router\.(push|replace)|navigate)\([^)]*login[^)]*\)' js | head
```

English text：**English text**、**English text**（English text `0` English text `200`）、**English text**。English text junk session English text：

```bash
curl -sk -X POST -H 'Cookie: <fakekey>=junk' https://target/api/<protected> -d '{}' -H 'Content-Type: application/json'
```

→ English text `neutralize.fields` + `neutralize.success`。

### A3. English text — 「English text/English text？」

```bash
grep -rhoaE '"/api[^"]*(permission|perm|role|menu|acl|resource|nav)[^"]*"' js | sort -u
grep -rhoaE '.{0,30}(menus|permissions|menuList|routeList|authList|role_permissions)\b.{0,120}' js | head
grep -rhoaE 'userRouteAuth|getResultTree|routeMap|routeLink|hasPermission|checkAuth' js | head
grep -rhoaE '([A-Z_][A-Z0-9_]*):\{name:"[^"]*",link:"/[^"]+"\}' js | head
```

**English text**（English text）：

| API | English text payload | English text |
|---|---|---|
| `.../role_permissions` | `{ permissions: string[], role_type }` | English text、English text ACL |
| `.../permissions/all` | `tree[{ code, position, children }]` | English text |
| bundle English text `userRouteAuth` | `{ CODE: { url, name? } }` | code → English text path |
| bundle English text `routeMap` | `{ KEY: { name, link } }` | English text（webpack `o.DASHBOARD`） |

English text：`getResultTree(tree, permissions)` English text、`v-if` / `hasAuth(code)` English text。

**English text forge**（English text）：English text permissive payload → `stubs`。

**English text**（English text，English text/English text）：English text **I English text**。

---

## B. config.json English text

```json
{
  "baseUrl": "https://target/",
  "runtimeMode": "both",
  "chromium": "/usr/bin/chromium",

  "cookies": [
    { "name": "auth", "value": "b64json:{\"id\":1,\"username\":\"admin\",\"role\":\"admin\",\"func\":{},\"permissions\":[\"*\"]}" }
  ],
  "localStorage": { "token": "faketoken", "isLogin": "1" },

  "neutralize": {
    "fields": ["response_code", "code", "errno", "ret", "status"],
    "success": 0,
    "flags": { "success": true, "message": "ok" }
  },
  "forward": true,
  "loginUrlPattern": "/login",
  "apiPattern": "/api/|/rest/|/graphql",

  "mockTier": "L1+L2",
  "recordDetail": true,
  "observe": {
    "storageReads": false,
    "cookieReads": false,
    "xhrHeaders": true
  },
  "neutralizeVueRouter": true,
  "stubs": [
    {
      "match": "permissions/all|/menu|role_permissions",
      "body": {
        "response_code": 0, "code": 0,
        "data": {
          "permissions": ["*"],
          "menus": [
            { "name": "dashboard", "path": "/dashboard", "show": true, "children": [] },
            { "name": "alert", "path": "/alert", "show": true, "children": [] }
          ]
        }
      }
    }
  ],

  "explore": {
    "clickTabs": true,
    "clickTables": true,
    "pushStateFallback": true,
    "maxMenuItems": 50
  },

  "routes": ["/dashboard", "/alert", "/asset", "/device", "/report", "/config", "/system"],
  "waitMs": 1500, "perRouteMs": 900, "headless": true,
  "waitUntil": "domcontentloaded",
  "routeTimeout": 12000,
  "proxy": "",

  "captureResponses": true, "recordWs": true, "respMax": 600
}
```

English text：
- `runtimeMode`：`depth`（Puppeteer）、`coverage`（browser MCP）、`both`
- `cookies[].value` English text：`b64json:` → base64(JSON)；`json:` → English text JSON；English text → English text
- `forward: true` English text；`false` English text stub
- `mockTier`：coverage English text preload English text，English text `L1+L2`、`L1+L2+L3`
- `routes` English text `routes.txt`；forge English text harness English text `<a href>`
- `captureResponses` / `recordWs` English text depth English text
- `waitUntil`：English text SPA English text `domcontentloaded`，English text `networkidle2` English text
- `routeTimeout`：English text `page.goto` English text（English text）
- `proxy`：Puppeteer `--proxy-server`；English text `HTTP_PROXY` / `HTTPS_PROXY`

### B1. English text stub English text（role_permissions + permissions/all）

```json
"stubs": [
  {
    "match": "role_permissions",
    "body": {
      "response_code": 0,
      "data": {
        "permissions": ["MONITOR", "MONITOR_ALERT", "THREAT", "ASSETS_RISK"],
        "role_type": "SUPER_ADMIN"
      }
    }
  },
  {
    "match": "permissions/all",
    "body": {
      "response_code": 0,
      "data": [
        {
          "code": "MONITOR",
          "position": 1,
          "children": [
            { "code": "MONITOR_ALERT", "position": 1, "children": [] }
          ]
        }
      ]
    }
  }
]
```

English text（`response_code` / `code` / `data`）English text A2 English text；`permissions` English text tree English text leaf code。

---

## C. coverage English text：preload English text

English text `scripts/preload.js` English text `CONFIG` English text，English text CDP English text：

```javascript
const CONFIG = {
  loginPathRe: /\/(login|signin)(\/|$|\?)/i,
  mockTier: 'L1+L2',
  forward: true,
  recordDetail: true,
  extractUrlsFromResponse: true,
  neutralizeVueRouter: true,
  observe: { storageReads: false, cookieReads: false, xhrHeaders: true },
  neutralize: { fields: ['response_code', 'code'], success: 0 },
  stubs: [ /* English text config.json stubs */ ],
  apiPattern: /\/(api|apis|v\d+|dev|internal|graphql)\//i,
};
```

English text：`window.__API_RECON_PRELOAD__ === true` English text pathname English text。

English text：

```javascript
JSON.stringify({
  apis: [...window.__API_RECON_LOG__],
  detail: window.__API_RECON_DETAIL__,
  routes: [...(window.__API_RECON_ROUTES__ || [])],
  observe: window.__API_RECON_OBSERVE__,
}, null, 2)
```

---

## D. preload / runtime Hook English text

preload（coverage）English text runtime_harvest（depth）English text Hook English text：

| Hook English text | English text API English text | English text |
|---|---|---|
| Hook fetch / XHR.open | English text URL/English text | ✅ `recordDetail` + `__API_RECON_LOG__` |
| Hook XHR.setRequestHeader | English text Authorization English text | ✅ `observe.xhrHeaders` |
| Hook localStorage/cookie English text | English text | ⚠️ English text `observe.storageReads/cookieReads` |
| Vue English text | English text frontendRoutes | ✅ `__API_RECON_ROUTES__`（English text） |
| Vue English text / English text | English text API | ✅ `neutralizeVueRouter` + English text |
| React English text | English text | ⚠️ English text + English text；English text Hook |
| English text（English text path） | English text | ⚠️ English text path，English text |
| Hook English text（CryptoJS/SM English text） | English text → English text API body | ❌ English text Hook English text；English text config |
| English text bypass | English text runtime English text API | ❌ English text；English text |

---

## E. Endpoint English text（English text）

English text `harvest_static.py` English text `extract_endpoints` English text，English text：

```bash
grep -rhoaE '"/[a-z][A-Za-z0-9_/\-]{3,}"' js | sort -u
grep -rhoaE '/api/[a-zA-Z0-9_./-]+' js | sort -u
```

---

## F. English text

| English text | English text → English text |
|---|---|
| English text API English text | endpoint English text → English text（D English text） |
| chunk English text ≪ manifest | CSS-only English text chunk；404 English text |
| runtime English text | English text → English text A1：English text、English text、English text、domain |
| English text | English text → forge English text（A3）；`routes` path English text |
| English text bootstrap/locale | English text → I English text；English text `role_permissions` + `permissions/all` English text stub |
| English text | tree English text intermediate English text code English text `userRouteAuth` English text |
| English text API English text | English text → English text `neutralize`；English text walk English text |
| WS English text 0 | English text subscribe；English text `perRouteMs` |
| English text | English text `forward: true` English text |
| Chromium English text | English text chromium English text `config.chromium` / `CHROMIUM` |
| Mock English text | Hook English text `location.href` setter → document-start + preload |
| English text | L3 English text；English text Tab/English text/English text |
| English text Redux action English text | English text get/set/change/clear/toggle/upload English text path |
| Vue English text | preload English text document-start → English text；English text `neutralizeVueRouter: false` English text |
| English text URL English text log | English text `extractUrlsFromResponse`；English text `__API_RECON_DETAIL__` English text |
| English text Authorization English text | English text `observe.xhrHeaders` English text DevTools English text |
| runtime English text / English text | English text `waitUntil: domcontentloaded`；English text `routeTimeout`；English text `networkidle2` |
| English text | English text `proxy` / English text；Puppeteer English text curl English text |

---

## G. hardened English text

English text（English text forge English text cookie、English text stub English text）English text，runtime English text shell English text。English text：

- **English text endpoint English text** — English text path English text
- English text，English text**English text**English text harness：`forward: true`、English text neutralize，English text methods/params/responses

---

## H. English text

1. English text
2. **English text** `scripts/harvest_static.py` → English text → English text → English text `api_static.txt`、`routes.txt`
3. **Phase 1b**：path English text + English text → `param_candidates.json`（J English text）
4. English text A1/A2/A3 → English text `config.json`
5. **English text** `runtime_harvest.js` / `preload.js` English text
6. `runtimeMode=depth`：`npm install` → English text harvest English text
7. `runtimeMode=coverage/both`：document-start English text preload → browser MCP English text + **English text**
8. English text → **I English text** → patch stubs → English text
9. English text diff + English text → `params_merged.json`
10. English text → `site_map.json` + `api_merged.txt`，English text、English text

---

## I. English text（Phase 4 English text）

English text forge English text `menus: [{ path, show: true }]` English text、English text mount English text。

### I1. English text auth English text

```bash
grep -l 'userRouteAuth' js/*.js
grep -l 'routeMap\|routeLink' js/*.js
grep -rhoaE 'getResultTree|role_permissions|permissions/all' js | head
```

English text：**English text API path**、**English text**、**English text chunk English text**。

### I2. English text routeMap

```bash
python3 scripts/extract_route_map.py recon/js recon/
# English text recon/route_map.json
```

English text `[!] no routeMap pattern found`：English text `extract_route_map.py` English text，English text grep：

```bash
grep -rhoaE '([A-Z_][A-Z0-9_]*):\{name:"[^"]*",link:"/[^"]+"\}' js | head -20
```

### I3. English text + stub

```bash
python3 scripts/build_perm_tree.py recon/js recon/ --config recon/config.json
```

English text：
1. English text `userRouteAuth={MONITOR:{url:...},...}`（English text webpack English text `He=o.DASHBOARD`）
2. English text `route_map.json` English text alias → English text path
3. English text code English text parent（`MONITOR_ALERT` → `MONITOR`）
4. English text `permissions_tree.json`、`permissions_all_stub.json`、`role_permissions_stub.json`
5. `--config` English text `config.json` English text `stubs` English text `routes`

**English text**（English text）：
- `DEFAULT_ROOTS`：English text code English text
- `DEFAULT_PREFIX_PARENT`：`PREFIX_` → parent English text
- `DEFAULT_EXTRA_PARENT`：English text orphan English text

### I4. English text stub English text

```bash
# permissions English text ≈ userRouteAuth English text
wc -l recon/perm_codes_all.txt
# routes English text route_map English text link
python3 -c "import json; m=json.load(open('recon/route_map.json')); r=set(json.load(open('recon/config.json'))['routes']); print('missing', [v['link'] for v in m.values() if v['link'] not in r])"
```

### I5. English text runtime English text

```bash
node recon/runtime_harvest.js recon/config.json
# English text forge English text runtime_api.json English text；English text /attack、/asset English text API
```

| forge English text | forge English text（English text） |
|---|---|
| English text 3–5 English text bootstrap | English text module API |
| English text `/api/locale/language` | English text `/api/web/...` English text endpoint |
| `routes.txt` English text | `routes` 80–110+ English text route_map |

### I6. English text

- **coverage English text**：English text + Tab，English text gating English text
- **stub English text**：English text API（curl + English text session）English text stub English text nesting
- **English text**：grep `hasPermission|checkRole|func.` English text，English text `role_permissions.permissions`
- **English text**：English text API path English text `api_static.txt`，runtime English text METHOD/body；English text `param_candidates.json` + English text

---

## J. English text（Phase 1b / 5b / 5c）

**English text，English text。** English text path English text；English text + UI English text + English text diff + English text。

### J1. English text — English text path English text

```bash
# English text Phase 1 English text path English text
grep -n '"/api/user/list"' js/*.js
grep -rhoaE '.{0,120}("/api[^"]+").{0,200}' js | head
grep -rhoaE '(params|data|body|payload)\s*:\s*\{' js | head
grep -rhoaE '(get|post|put|delete|patch)\([^,]+,\s*\{' js | head
```

### J2. English text

```bash
# axios / English text request
grep -rhoaE '(axios|request)\.(get|post|put|delete|patch)\(' js | head
grep -rhoaE 'interceptors\.(request|response)' js | head

# GraphQL
grep -rhoaE '(query|mutation)\s+\w+|gql`|graphql\(' js | head
grep -rhoaE '\$[a-zA-Z_]+\s*:\s*(Int|String|Boolean|\[)' js | head

# FormData / multipart
grep -rhoaE 'FormData|\.append\(' js | head

# English text
grep -rhoaE 'path:\s*"/[^"]*:[^"]+"' js | head
grep -rhoaE 'useParams|route\.params|\$route\.params' js | head
```

### J3. English text — English text / English text / English text

```bash
grep -rhoaE '(required|message|pattern|enum|validator)\s*:' js | head
grep -rhoaE 'yup\.|zod\.|async-validator|Form\.Item|a-form-item|el-form-item' js | head
grep -rhoaE 'rules\s*:\s*\[|name:\s*["\'][a-zA-Z_]+["\']' js | head
grep -rhoaE 'label.*value|options\s*:\s*\[' js | head
```

### J4. English text — English text → API

```bash
grep -rhoaE 'onFinish|handleSubmit|getFieldsValue|validateFields' js | head
grep -rhoaE '(pick|omit|transform|dayjs|moment)\(' js | head
```

runtime English text：DevTools → Network → English text → **English text**（call stack）English text `fetch`/`send` English text。

### J5. English text

```bash
grep -rhoaE 'encrypt|decrypt|sign|CryptoJS|sm2|sm3|sm4|RSA|AES' js | head
```

**English text** — Hook English text**English text**，English text plaintext payload；English text `config.json` / `param_candidates.json`。

### J6. English text（Phase 3 English text）

English text，diff English text body/query：

| English text | English text |
|---|---|
| English text | English text |
| English text | keyword、filters |
| English text | optional English text |
| English text/English text | English text entity |
| English text/English text | `ids[]`、`exportType` |
| English text/English text | `sortField`、`order` |

English text `param_samples.json`：`[{ "path", "method", "action": "search", "body", "query", "headers" }]`

### J7. English text

| English text | English text |
|---|---|
| **English text** | English text callsite + runtime ≥2 English text |
| **English text** | English text，English text 1 English text runtime |
| **English text** | English text/English text，English text |
| **English text** | English text，UI/English text |

### J8. English text

| English text | English text |
|---|---|
| REST English text | J1 English text → J6 English text diff → J3 rules |
| English text/English text | J3 Form name → J4 submit English text → runtime English text + English text 400 |
| GraphQL | J2 variables English text → runtime English text operation English text variables |
| English text body | J5 Hook English text → English text params |

### J9. English text api-recon English text

| api-recon | English text recon |
|---|---|
| Phase 1 English text | J1 English text |
| Phase 2 A2 English text | English text（tenantId、sign） |
| Phase 3 runtime | J6 English text + `param_samples.json` |
| Phase 4 English text | English text → English text |
| Phase 5 English text | `params_merged.json` + English text；English text |

### J10. English text

| English text | English text |
|---|---|
| English text runtime English text | English text「English text」；English text / English text / English text select English text option |
| English text path English text body English text | English text — English text `action` English text，English text schema |
| stub English text params | **English text outbound English text** body/headers，English text stub English text |
| 400 English text nested field | English text `data`/`bizData`/`variables` |
| GraphQL English text operation English text | English text `variables` JSON；English text `$var: Type` |

---
