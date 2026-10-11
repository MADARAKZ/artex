package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/Autumn-27/artex/db"
	actool "github.com/Autumn-27/norma/tool"
)

// assetInterceptCandidates English text English text/IP/URL English text，English text。
// URL English text host English text，English text「English text URL」English text/English text English text/IP English text。
func assetInterceptCandidates(item assetInputItem) (domains, ips, urls []string) {
	add := func(dst *[]string, s string) {
		if s = strings.TrimSpace(s); s != "" {
			*dst = append(*dst, s)
		}
	}
	add(&domains, item.Domain)
	for _, d := range item.BoundDomains {
		add(&domains, d)
	}
	add(&ips, item.IP)
	add(&ips, item.ServiceIP)
	add(&urls, item.URL)
	if item.URL != "" {
		if u, err := url.Parse(item.URL); err == nil {
			if h := u.Hostname(); h != "" {
				if net.ParseIP(h) != nil {
					add(&ips, h)
				} else {
					add(&domains, h)
				}
			}
		}
	}
	return domains, ips, urls
}

// assetInputLabel English text，English text。
func assetInputLabel(item assetInputItem) string {
	typ := strings.TrimSpace(item.Type)
	var target string
	switch {
	case strings.TrimSpace(item.Domain) != "":
		target = strings.TrimSpace(item.Domain)
	case strings.TrimSpace(item.URL) != "":
		target = strings.TrimSpace(item.URL)
	case strings.TrimSpace(item.IP) != "":
		target = strings.TrimSpace(item.IP)
	case strings.TrimSpace(item.ServiceIP) != "":
		target = strings.TrimSpace(item.ServiceIP)
	default:
		target = "(English text)"
	}
	if typ != "" {
		return fmt.Sprintf("[%s] %s", typ, target)
	}
	return target
}

// =====================================================================
// Unified asset insertion tools
// =====================================================================

// SetAssetStore wires the asset store and company store onto this ToolSet
// so the insert_assets, add_company_scope, and list_assets tools are active.
func (t *ToolSet) SetAssetStore(as *db.AssetStore, cs *db.CompanyStore) {
	t.as = as
	t.cs = cs
}

// assetInputItem is one element of the insert_assets "assets" array.
type assetInputItem struct {
	Type string `json:"type"` // root_domain|ip|subdomain|app|service|endpoint

	// ---- root_domain / subdomain ----
	Domain      string   `json:"domain"`
	ICP         string   `json:"icp"`
	RecordType  string   `json:"record_type"`
	RecordValue []string `json:"record_value"`

	// ---- ip ----
	IP           string           `json:"ip"`
	BoundDomains []string         `json:"bound_domains"`
	OpenPorts    []db.PortService `json:"open_ports"`

	// ---- app ----
	AppName     string `json:"app_name"`
	BundleID    string `json:"bundle_id"`
	Category    string `json:"category"`
	Description string `json:"description"`
	AppICP      string `json:"app_icp"`
	CompanyID   *int64 `json:"company_id"` // explicit company link (app only; others auto-attribute via scope)

	// ---- service (http) ----
	URL           string           `json:"url"`
	Technologies  []string         `json:"technologies"`
	StatusCode    *int             `json:"status_code"`
	ContentLength *int64           `json:"content_length"`
	PageTitle     string           `json:"page_title"`
	FaviconMMH3   string           `json:"favicon_mmh3"`
	Auth          []map[string]any `json:"auth"`
	ServiceName   string           `json:"service_name"`
	ServiceIP     string           `json:"service_ip"` // optional enrichment IP

	// ---- service (other) ----
	Port  int    `json:"port"`
	Proto string `json:"proto"`

	// ---- endpoint ----
	Method string           `json:"method"`
	Params []map[string]any `json:"params"`
}

// insertAssets is the unified insert_assets agent tool.
func (t *ToolSet) insertAssets() actool.CoreTool {
	return writeTool(
		"insert_assets",
		"English text，English text（type English text）。\n"+
			"English text：root_domain→domain；ip→ip（English text IPv4/IPv6，English text）；subdomain→domain；app→app_name；service(HTTP)→url；service(English textHTTP)→service_name+port（ip/domain English text）；endpoint→url+method。English text。\n"+
			"auth/technologies/params English text(append)，English text。\n"+
			"English text：{results:[{index,id,type}], errors:[{index,error}]}",
		obj(map[string]any{
			// task_id English text：worker English text task English text SetTaskID English text(English text handler)。
			"assets": map[string]any{
				"type":        "array",
				"description": "English text，English text",
				"items": obj(map[string]any{
					"type": map[string]any{
						"type":        "string",
						"enum":        []string{"root_domain", "ip", "subdomain", "app", "service", "endpoint"},
						"description": "English text",
					},
					// root_domain / subdomain
					"domain":      str("English text（root_domain/subdomain English text）"),
					"icp":         str("ICP English text（English text）"),
					"record_type": str("DNS English text：A/AAAA/CNAME/MX English text（subdomain English text）"),
					"record_value": map[string]any{
						"type":        "array",
						"items":       map[string]any{"type": "string"},
						"description": "DNS English text（subdomain English text，English text [\"1.2.3.4\",\"2.3.4.5\"]）",
					},
					// ip
					"ip": str("IP English text，English text IPv4/IPv6 English text，English text（English text type=subdomain English text domain English text）；ip English text；service/endpoint English text，English text IP"),
					"bound_domains": map[string]any{
						"type":        "array",
						"items":       map[string]any{"type": "string"},
						"description": "English text IP English text（ip English text）",
					},
					"open_ports": map[string]any{
						"type":        "array",
						"description": "English text（ip English text）",
						"items": obj(map[string]any{
							"port":    intp("English text"),
							"service": str("English text，English text http/ssh/mysql English text（English text）"),
						}, "port"),
					},
					// app
					"app_name":    str("English text（app English text）"),
					"bundle_id":   str("Bundle ID（app English text）"),
					"category":    str("English text（English text）"),
					"description": str("English text（English text）"),
					"app_icp":     str("English text ICP English text（English text）"),
					"company_id":  intp("English text id（app English text；app English text scope English text，English text。id English text add_company_scope English text）"),
					// service (http)
					"url":         str("English text URL，English text（HTTP English text；service_type English text http）"),
					"status_code": intp("HTTP English text，English text 200/301/403/404（English text）"),
					"content_length": map[string]any{
						"type":        "integer",
						"description": "HTTP English text（English text）",
					},
					"page_title":   str("English text <title> English text（English text）"),
					"favicon_mmh3": str("favicon MMH3 English text（English text）"),
					"technologies": map[string]any{
						"type":        "array",
						"items":       map[string]any{"type": "string"},
						"description": "English text/English text，English text [\"Nginx\",\"Vue\",\"Bootstrap\"]（English text）",
					},
					"auth": map[string]any{
						"type":        "array",
						"description": "English text，English text type/username/password English text（English text，English text）",
						"items":       map[string]any{"type": "object"},
					},
					// service (other，English text HTTP)
					"service_name": str("English text，English text ssh/mysql/redis（service English text HTTP English text）"),
					"port":         intp("English text（service English text HTTP English text）"),
					// endpoint
					"method": str("HTTP English text：GET/POST/PUT/PATCH/DELETE English text（endpoint English text）"),
					"params": map[string]any{
						"type":        "array",
						"description": "English text，English text location(query/body/header/path)/name/value/type（English text，English text）",
						"items":       map[string]any{"type": "object"},
					},
				}, "type"),
			},
		}, "assets"),
		func(_ context.Context, in json.RawMessage) (actool.Result, error) {
			if t.as == nil {
				return actool.Errorf("insert_assets English text: AssetStore English text"), nil
			}
			var a struct {
				Assets []assetInputItem `json:"assets"`
			}
			if err := json.Unmarshal(in, &a); err != nil {
				return actool.Errorf("invalid input: " + err.Error()), nil
			}
			// task_id English text(worker: SetTaskID)，English text——English text/English text
			// English text。English text(auto/pentest/chat)English text t.taskID=0。
			taskID := t.taskID

			type result struct {
				Index int    `json:"index"`
				ID    int64  `json:"id"`
				Type  string `json:"type"`
			}
			type errEntry struct {
				Index int    `json:"index"`
				Error string `json:"error"`
			}

			var results []result
			var errs []errEntry

			// English text；English text（English text）。
			// English text = English text ∪ English text block；English text = English text allow。
			blockRules, _ := t.as.ListAssetInterceptRules()
			var allowRules []db.AssetInterceptRule
			if t.taskID > 0 {
				if tb, ta, err := t.as.TaskInterceptRulesSplit(t.taskID); err == nil {
					blockRules = append(blockRules, tb...)
					allowRules = ta
				}
			}

			for i, item := range a.Assets {
				// English text：English text，English text（English text Upsert English text）。
				domains, ips, urls := assetInterceptCandidates(item)
				if d := db.EvaluateAssetGate(blockRules, allowRules, domains, ips, urls); !d.Allowed {
					errs = append(errs, errEntry{
						Index: i,
						Error: fmt.Sprintf("English text %s %s，English text", assetInputLabel(item), d.Reason),
					})
					continue
				}

				typ := strings.TrimSpace(item.Type)
				var id int64
				var err error

				switch typ {
				case "root_domain":
					id, err = t.as.UpsertRootDomain(db.UpsertRootDomainReq{
						Domain: item.Domain,
						ICP:    item.ICP,
						TaskID: taskID,
					})

				case "ip":
					id, err = t.as.UpsertIP(db.UpsertIPReq{
						IP:           item.IP,
						BoundDomains: item.BoundDomains,
						OpenPorts:    item.OpenPorts,
						TaskID:       taskID,
					})

				case "subdomain":
					id, err = t.as.UpsertSubdomain(db.UpsertSubdomainReq{
						Domain:      item.Domain,
						RecordType:  item.RecordType,
						RecordValue: item.RecordValue,
						ICP:         item.ICP,
						TaskID:      taskID,
					})

				case "app":
					id, err = t.as.UpsertApp(db.UpsertAppReq{
						Name:        item.AppName,
						BundleID:    item.BundleID,
						Category:    item.Category,
						Description: item.Description,
						ICP:         item.AppICP,
						CompanyID:   item.CompanyID,
						TaskID:      taskID,
					})

				case "service":
					// distinguish HTTP vs other by presence of url
					if item.URL != "" {
						// agent may send "ip" or "service_ip" for the enrichment IP; accept both
						svcIP := item.ServiceIP
						if svcIP == "" {
							svcIP = item.IP
						}
						id, err = t.as.UpsertHTTPService(db.UpsertHTTPServiceReq{
							URL:           item.URL,
							Technologies:  item.Technologies,
							StatusCode:    item.StatusCode,
							ContentLength: item.ContentLength,
							PageTitle:     item.PageTitle,
							FaviconMMH3:   item.FaviconMMH3,
							Auth:          item.Auth,
							IP:            svcIP,
							TaskID:        taskID,
						})
					} else {
						id, err = t.as.UpsertOtherService(db.UpsertOtherServiceReq{
							Domain:      item.Domain,
							IP:          item.IP,
							Port:        item.Port,
							ServiceName: item.ServiceName,
							Auth:        item.Auth,
							TaskID:      taskID,
						})
					}

				case "endpoint":
					id, err = t.as.UpsertEndpoint(db.UpsertEndpointReq{
						URL:    item.URL,
						Method: item.Method,
						Params: item.Params,
						IP:     item.ServiceIP,
						TaskID: taskID,
					})

				default:
					errs = append(errs, errEntry{Index: i, Error: "unknown type: " + typ})
					continue
				}

				if err != nil {
					errs = append(errs, errEntry{Index: i, Error: err.Error()})
					continue
				}
				results = append(results, result{Index: i, ID: id, Type: typ})
				t.writes.Assets++
				t.anchorOwner(id)
				if taskID > 0 {
					var sourceNodeID *int64
					if t.ownerNode > 0 {
						nodeID := t.ownerNode
						sourceNodeID = &nodeID
					}
					summary := "Agent English text insert_assets English text"
					if t.ownerNode > 0 {
						summary = fmt.Sprintf("Worker English text #%d English text insert_assets English text", t.ownerNode)
					}
					_ = t.as.SetTaskAssetSource(taskID, id, "agent", summary, sourceNodeID)
				}
				// English text(source='auto')：English text worker English text，English text
				// English text；side-effect English text，English text。taskID=0 English text。
				// English text：task_scope English text(list/English text)，
				// English text，English text。
				{
					svcIP := item.ServiceIP
					if svcIP == "" {
						svcIP = item.IP
					}
					_ = t.as.AddAutoScope(taskID, typ, item.Domain, item.URL, svcIP)
				}
			}

			return jsonResult(map[string]any{
				"results": results,
				"errors":  errs,
			})
		},
	)
}

// addCompanyScope writes to company_scope table and triggers asset attribution.
func (t *ToolSet) addCompanyScope() actool.CoreTool {
	return writeTool(
		"add_company_scope",
		"English text/IP/CIDR/ICPEnglish text/English text【English text】——English text、English textICPEnglish text，English textAgentEnglish text。\n"+
			"English text：company English text，English text(English text)。\n"+
			"scope English text，English text：English text / URL / English text IP / CIDR English text / ICPEnglish text / English text。\n"+
			"English text reason English text(whois/English text/ASN English text)。\n"+
			"English text：English text TLD English text(IPv4English text/16-/32、IPv6English text/32-/128)，English text errors English text。",
		obj(map[string]any{
			"company": str("English text(English text、English text；English text)"),
			"scope":   str("English text，English text：English text / URL / IP / CIDR / ICPEnglish text / English text"),
			"reason":  str("English text(English text/English text)，English text"),
			"logo":    str("English text URL(English text；English text)"),
		}, "company", "scope"),
		func(_ context.Context, in json.RawMessage) (actool.Result, error) {
			if t.cs == nil {
				return actool.Errorf("add_company_scope English text: CompanyStore English text"), nil
			}
			var a struct {
				Company string `json:"company"`
				Scope   string `json:"scope"`
				Reason  string `json:"reason"`
				Logo    string `json:"logo"`
			}
			if err := json.Unmarshal(in, &a); err != nil {
				return actool.Errorf(err.Error()), nil
			}
			if strings.TrimSpace(a.Company) == "" {
				return actool.Errorf("company English text"), nil
			}
			companyID, _, err := t.cs.UpsertCompany(a.Company, a.Logo)
			if err != nil {
				return actool.Errorf("English text/English text: " + err.Error()), nil
			}
			lines := splitLines(a.Scope)
			added, skipped, invalid, errMsgs := t.cs.AddScope(companyID, lines, a.Reason)
			out := map[string]any{
				"company_id": companyID,
				"added":      added,
				"skipped":    skipped,
				"invalid":    invalid,
			}
			if len(errMsgs) > 0 {
				out["errors"] = errMsgs
			}
			return jsonResult(out)
		},
	)
}

// addTaskScope lets the plan agent add test scope to THE CURRENT TASK — the coverage
// denominator and the task's authorization edge. Worker discoveries are auto-scoped
// (precise host) by insertAssets; this tool is for DELIBERATELY WIDENING: pull a whole
// root domain or whole company into scope, or add a specific subdomain / ip.
func (t *ToolSet) addTaskScope() actool.CoreTool {
	return writeTool(
		"add_task_scope",
		"English text【English text】——English text，English text。\n"+
			"kind English text：company(English text) / root_domain(English text，English text) / subdomain(English text) / ip / cidr / icp / keyword。\n"+
			"English text：worker English text【English text】English text(English text)；English text【English text】——English text/English text，English text/IP。\n"+
			"value：company English text id(English text)；root_domain/subdomain English text；ip/cidr English text IP English text；icp/keyword English text。\n"+
			"English text reason English text(English text)。English text entries English text。",
		obj(map[string]any{
			"entries": map[string]any{"type": "array", "description": "English text：[{kind, value}]。kind∈company/root_domain/subdomain/ip/cidr/icp/keyword。", "items": map[string]any{"type": "object"}},
			"kind":    str("[English text] company / root_domain / subdomain / ip / cidr / icp / keyword"),
			"value":   str("[English text] English textid / English text / IP / CIDR / ICP / English text"),
			"reason":  str("English text(English text)，English text"),
		}),
		func(_ context.Context, in json.RawMessage) (actool.Result, error) {
			if t.as == nil {
				return actool.Errorf("add_task_scope English text: AssetStore English text"), nil
			}
			if t.taskID <= 0 {
				return actool.Errorf("add_task_scope English text(English text task)"), nil
			}
			type scopeEntry struct {
				Kind  string `json:"kind"`
				Value string `json:"value"`
			}
			var a struct {
				Entries    []scopeEntry `json:"entries"`
				scopeEntry              // English text
				Reason     string       `json:"reason"`
			}
			_ = json.Unmarshal(in, &a)
			items := a.Entries
			if len(items) == 0 {
				items = []scopeEntry{a.scopeEntry}
			}
			var added []map[string]any
			errs := map[string]string{}
			for i, e := range items {
				ts, err := t.as.AddAgentScope(t.taskID, strings.TrimSpace(e.Kind), e.Value, a.Reason, "agent")
				if err != nil {
					errs[strconv.Itoa(i)] = err.Error()
					continue
				}
				added = append(added, map[string]any{"kind": ts.Kind, "domain": ts.Domain, "net": ts.Net, "value": ts.Value, "company_id": ts.CompanyID})
			}
			out := map[string]any{"added": added}
			if len(errs) > 0 {
				out["errors"] = errs
			}
			return jsonResult(out)
		},
	)
}

// listUntestedAssets lets the plan agent pull the current + directly inherited
// scope's not-yet-tested assets on demand (filter by type, paginated).
func (t *ToolSet) listUntestedAssets() actool.CoreTool {
	return readTool(
		"list_untested_assets",
		"English text【English text】English text、English text（English text，English text，English text）。\n"+
			"English text：root_domain/subdomain/service/app/endpoint/ip。\n"+
			"English text：page English text 1 English text、page_size English text 10。English text {assets:[{id,type,label}], total, page, page_size}。English text。",
		obj(map[string]any{
			"type":      str("English text（English text）：root_domain/subdomain/service/app/endpoint/ip"),
			"page":      intp("English text，English text 1 English text（English text 1）"),
			"page_size": intp("English text（English text 10）"),
		}),
		func(_ context.Context, in json.RawMessage) (actool.Result, error) {
			if t.as == nil {
				return actool.Errorf("list_untested_assets English text: AssetStore English text"), nil
			}
			if t.taskID <= 0 || t.ts == nil {
				return actool.Errorf("list_untested_assets English text"), nil
			}
			var a struct {
				Type     string `json:"type"`
				Page     int    `json:"page"`
				PageSize int    `json:"page_size"`
			}
			_ = json.Unmarshal(in, &a)
			if a.Page <= 0 {
				a.Page = 1
			}
			if a.PageSize <= 0 {
				a.PageSize = 10
			}
			offset := (a.Page - 1) * a.PageSize
			assets, total, err := t.as.ListUntestedAssetsWithSources(t.taskID, strings.TrimSpace(a.Type), a.PageSize, offset)
			if err != nil {
				return actool.Errorf(err.Error()), nil
			}
			return jsonResult(map[string]any{
				"assets": assets, "total": total, "page": a.Page, "page_size": a.PageSize,
			})
		},
	)
}

// listAssets lets an agent query the asset table.
func (t *ToolSet) listAssets() actool.CoreTool {
	return readTool(
		"list_assets",
		"English text：DSL English text，English text id/ids English text；English text。English text【English text】English text。\n"+
			"DSL：field=value English text(ILIKE) | field==value English text | field!=value English text | English text > >= < <= | English text=English text；AND/OR English text(AND English text)，English text。English text type English text，English text DSL。\n"+
			"English text id/ids English text dsl English text（English text）。\n"+
			"English text：domain(English text/English text/English text)、root_domain、ip、url、page_title、icp、service_name、app_name、method(English text GET/POST)、service_type(http|other)、record_type(English text A/CNAME)、technology(English text，=English text ==English text)、port/status_code/company_id(English text)。\n"+
			"English text：status_code>=400 AND technology=shiro ；(port==80 OR port==443) AND technology=nginx",
		obj(map[string]any{
			"dsl":    str(`DSL English text（English text/English text）。English text id/ids English text。`),
			"type":   str("English text：root_domain|ip|subdomain|app|service|endpoint（English text，English text dsl English text；English text type English text，English text dsl）"),
			"id":     intp("English text id English text（English text，English text dsl/type English text）"),
			"ids":    map[string]any{"type": "array", "items": map[string]any{"type": "integer"}, "description": "English text id English text（English text，English text dsl/type English text）"},
			"limit":  intp("English text，English text 10（English text）"),
			"offset": intp("English text，English text 0（English text）"),
		}),
		func(_ context.Context, in json.RawMessage) (actool.Result, error) {
			if t.as == nil {
				return actool.Errorf("list_assets English text: AssetStore English text"), nil
			}
			var a struct {
				DSL    string  `json:"dsl"`
				Type   string  `json:"type"`
				ID     int64   `json:"id"`
				IDs    []int64 `json:"ids"`
				Limit  int     `json:"limit"`
				Offset int     `json:"offset"`
			}
			_ = json.Unmarshal(in, &a)
			if a.Limit <= 0 {
				a.Limit = 10
			}

			var assets []*db.Asset
			var err error
			switch {
			case a.ID > 0:
				assets, err = t.as.GetByIDsInScope(t.taskID, []int64{a.ID})
			case len(a.IDs) > 0:
				assets, err = t.as.GetByIDsInScope(t.taskID, a.IDs)
			case a.DSL != "":
				assets, err = t.as.QueryDSLInScope(a.DSL, a.Type, t.taskID, a.Limit, a.Offset)
			default:
				return actool.Errorf("English text id/ids English text dsl English text：English text，English text"), nil
			}
			if err != nil {
				return actool.Errorf("DSL English text: " + err.Error()), nil
			}
			return jsonResult(map[string]any{
				"count":  len(assets),
				"assets": assets,
			})
		},
	)
}

// listCompanies lets an agent enumerate companies (English text) with their scope + asset count.
func (t *ToolSet) listCompanies() actool.CoreTool {
	return readTool(
		"list_companies",
		"English text【English text/English text】English text(scope)English text。English text、"+
			"English text company_id（insert_assets English text app、list_assets English text company_id English text）。"+
			"English text search English text(English text)，English text。",
		obj(map[string]any{
			"search": str("English text(English text，English text)；English text"),
		}),
		func(_ context.Context, in json.RawMessage) (actool.Result, error) {
			if t.cs == nil {
				return actool.Errorf("list_companies English text: CompanyStore English text"), nil
			}
			var a struct {
				Search string `json:"search"`
			}
			_ = json.Unmarshal(in, &a)
			cos, err := t.cs.ListCompanies()
			if err != nil {
				return actool.Errorf("English text: " + err.Error()), nil
			}
			q := strings.ToLower(strings.TrimSpace(a.Search))
			type companyOut struct {
				ID         int64    `json:"id"`
				Name       string   `json:"name"`
				AssetCount int      `json:"asset_count"`
				Scope      []string `json:"scope"`
			}
			out := make([]companyOut, 0, len(cos))
			for _, c := range cos {
				if q != "" && !strings.Contains(strings.ToLower(c.Name), q) {
					continue
				}
				scope := make([]string, 0, len(c.Scope))
				for _, r := range c.Scope {
					scope = append(scope, r.Raw)
				}
				out = append(out, companyOut{ID: c.ID, Name: c.Name, AssetCount: c.AssetCount, Scope: scope})
			}
			return jsonResult(map[string]any{"count": len(out), "companies": out})
		},
	)
}

// splitLines splits a multi-line string into non-empty trimmed lines.
func splitLines(s string) []string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

// WorkerTools returns the tool set for a work agent.
func (t *ToolSet) WorkerTools() []actool.CoreTool {
	return []actool.CoreTool{
		// list_findings English text：English text，English text。
		t.listFindings(),
		t.addFinding(), t.recordFact(),
		// asset management (handlers guard nil store internally)。
		// add_company_scope English text worker：English text/English text/Auto English text，worker English text。
		t.insertAssets(), t.listAssets(),
		// English text work English text：worker English text work English text，English text。
		// search_all_worker_traces：English text intent_id，English text；
		// get_worker_trace：English text work English text/English text/English text。
		t.searchAllWorkerTraces(), t.getWorkerTrace(),
		// node_detail：worker English text intent_id/English text id English text（English text）。
		t.nodeDetail(),
		// English text【English text】worker，English text planner/main（English text、English text work English text，
		// worker English text）：list_facts / list_companies / list_worker_traces。
	}
}

// MainAgentTools returns the human-interface tool set.
func (t *ToolSet) MainAgentTools() []actool.CoreTool {
	return []actool.CoreTool{
		t.graphOverview(), t.listFindings(), t.listFacts(), t.nodeDetail(),
		t.expandDigest(), // cold-digest §6.1
		t.getWorkerOutput(), t.getWorkerTrace(), t.searchAllWorkerTraces(), t.addHint(), t.addIntent(),
		// steer_work：English text(work)English text（English text、English text）。
		t.steerWorkTool(),
		// set_goals：English text（English text）。
		t.setGoals(),
		// set_constraints：English text/English text（allow/deny），English text planner/worker English text。
		t.setConstraints(),
		// asset management (handlers guard nil store internally)
		t.insertAssets(), t.addCompanyScope(), t.listAssets(),
		t.addFinding(), t.recordFact(),
		t.addTaskScope(),
		// list_untested_assets：English text(English text+English text)，English text。
		t.listUntestedAssets(),
	}
}

// AllDomainTools returns the union of all domain tools across all agent types,
// deduped by name (mainagent order wins). Used by the server to build a registry
// for injecting domain tools into agents (Auto, custom) that don't own a per-task
// ToolSet. The caller provides real stores; tools are callable at taskID=0 scope.
func (t *ToolSet) AllDomainTools() []actool.CoreTool {
	seen := map[string]bool{}
	var out []actool.CoreTool
	all := append(append(t.MainAgentTools(), t.PlannerTools()...), t.WorkerTools()...)
	for _, tool := range all {
		if !seen[tool.Name()] {
			seen[tool.Name()] = true
			out = append(out, tool)
		}
	}
	return out
}
