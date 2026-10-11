package agent

import (
	"context"
	"encoding/json"

	actool "github.com/Autumn-27/norma/tool"
)

// English text「English text」English text、English text DB English text：
//   - BuiltinToolSeeds()：English text agent English text seed English text（key +
//     English text + English text schema + English text agent），English text tools English text。
//   - ToolResolve English text：English text DB English text tools English text「English text agent English text +
//     English text/schema + English text」。key/handler English text，DB English text「English text」。
// handler（Call English text）English text——DB English text，English text。

// ToolSeed English text：key English text CoreTool.Name()（English text handler English text，
// UI English text），Desc/Schema English text，Agents English text agent。
type ToolSeed struct {
	Key    string         // = CoreTool.Name()，English text，English text
	Desc   string         // English text（English text UI English text）
	Schema map[string]any // English text JSON-Schema（English text，description/default English text UI English text）
	Agents []string       // English text agent key（worker/planner/mainagent）
}

// builtinToolsByAgent English text「English text」ToolSet（nil stores）English text agent English text
// English text。English text Spec、English text store，English text nil English text——
// English text Name()/Description()/InputSchema()，English text Call。
//
// English text【English text】SDK English text actool.DefaultTools()（Read/Write/Edit/MultiEdit/LS/Glob/
// Grep/Bash）：English text agent English text、English text「English text」English text，English text Prompt()
// English text（English text Description()，English text）。English text seed → English text DB English text → ToolResolve
// English text、English text，English text。English text artex English text。
func builtinToolsByAgent() map[string][]actool.CoreTool {
	ts := NewToolSet(nil, "")
	return map[string][]actool.CoreTool{
		"mainagent": ts.MainAgentTools(),
		"planner":   ts.PlannerTools(),
		"worker":    ts.WorkerTools(),
		// goals（English text）English text set_goals + set_constraints：English text、
		// English text。English text mainagent English text，web English text/schema、English text agent English text。
		"goals": {ts.setGoals(), ts.setConstraints()},
		// auto English text + English text，English text UI English text。
		// English text seed English text；English text seedAutoDefaultBindings English text。
		"auto": {ts.addFinding(), ts.insertAssets(), ts.addCompanyScope(), ts.listAssets(), ts.listCompanies()},
		// pentest（English text agent）English text：English text / English text / English text / English text / English text。
		// English text seed English text；English text seedPentestDefaultBindings English text。
		"pentest": {ts.listAssets(), ts.insertAssets(), ts.addFinding(), ts.listFindings(), ts.listCompanies()},
	}
}

// defaultUnbound：English text system English text（web English text、English text agent English text），English text
// English text【English text agent】——ToolResolve English text agent English text,English text opt-in。
// English text agent English text base English text(English text goal_met English text PlannerTools):English text seed English text
// English text desc/schema,English text base English text、ToolResolve English text。
//
// goal_met：English text prove_goal、English text【English text】,English text,English text
// 「prove_goal English text → English text」English text,English text agent,English text。
var defaultUnbound = map[string]bool{"goal_met": true}

// BuiltinToolSeeds English text agent English text seed English text：English text（English text list_assets
// English text agent English text）English text，Agents English text；defaultUnbound English text。
func BuiltinToolSeeds() []ToolSeed {
	byAgent := builtinToolsByAgent()
	order := []string{"mainagent", "goals", "planner", "worker", "auto", "pentest"}

	type acc struct {
		tool   actool.CoreTool
		agents []string
	}
	m := map[string]*acc{}
	var keys []string
	for _, ak := range order {
		for _, t := range byAgent[ak] {
			a, ok := m[t.Name()]
			if !ok {
				a = &acc{tool: t}
				m[t.Name()] = a
				keys = append(keys, t.Name())
			}
			a.agents = append(a.agents, ak)
		}
	}

	out := make([]ToolSeed, 0, len(keys))
	for _, k := range keys {
		a := m[k]
		agents := a.agents
		if defaultUnbound[k] {
			agents = []string{} // English text、English text，English text agent（English text [] English text null，English text）
		}
		out = append(out, ToolSeed{
			Key:    k,
			Desc:   a.tool.Description(),
			Schema: a.tool.InputSchema(),
			Agents: agents,
		})
	}
	return out
}

// ToolResolve, if set, post-processes an agent's fully-assembled tool list against
// the DB tools table: it drops tools not bound to this agent (or globally disabled)
// and wraps the rest so the model sees the DB-overridden description/schema and
// English text get injected. Tools with no matching DB row (MCP/skill/host tools like
// traffic) pass through untouched. nil = tools unchanged. Wired in server/assembly.go.
var ToolResolve func(ctx context.Context, agentKey string, tools []actool.CoreTool) []actool.CoreTool

// DecorateTool wraps t so Description()/InputSchema() report the DB overrides and
// Call() injects scalar parameter defaults (from schema's "default" props) whenever
// the model omitted them. Name/Prompt/permission/scheduler flags delegate to t, so
// the tool's identity and handler are unchanged. Empty desc/schema fall back to t's.
func DecorateTool(t actool.CoreTool, desc string, schema map[string]any) actool.CoreTool {
	if desc == "" {
		desc = t.Description()
	}
	if len(schema) == 0 {
		schema = t.InputSchema()
	}
	return &overriddenTool{CoreTool: t, desc: desc, schema: schema}
}

// overriddenTool is a CoreTool decorator: it embeds the original (so all behavioral
// methods — Prompt/IsReadOnly/IsConcurrencySafe/CheckPermissions/Name — delegate)
// and overrides only the model-facing description/schema plus default injection.
type overriddenTool struct {
	actool.CoreTool
	desc   string
	schema map[string]any
}

func (o *overriddenTool) Description() string         { return o.desc }
func (o *overriddenTool) InputSchema() map[string]any { return o.schema }

func (o *overriddenTool) Call(ctx context.Context, in json.RawMessage, tc *actool.ToolContext) (actool.Result, error) {
	return o.CoreTool.Call(ctx, injectDefaults(in, o.schema), tc)
}

// injectDefaults fills scalar parameter defaults declared in the (possibly edited)
// schema into the input JSON whenever the model omitted the field or left it empty/
// null. Structure (names/types/required) is untouched — onlyEnglish text are merged in.
func injectDefaults(in json.RawMessage, schema map[string]any) json.RawMessage {
	defs := scalarDefaults(schema)
	if len(defs) == 0 {
		return in
	}
	m := map[string]json.RawMessage{}
	if len(in) > 0 {
		if err := json.Unmarshal(in, &m); err != nil {
			return in // non-object input: don't touch it
		}
	}
	changed := false
	for k, dv := range defs {
		if cur, ok := m[k]; !ok || isEmptyJSON(cur) {
			m[k] = dv
			changed = true
		}
	}
	if !changed {
		return in
	}
	b, err := json.Marshal(m)
	if err != nil {
		return in
	}
	return b
}

// scalarDefaults extracts properties[k]["default"] for scalar params (string/
// integer/number/boolean). Array/object defaults are skipped: merging them is
// ambiguous and not worth the surprise.
func scalarDefaults(schema map[string]any) map[string]json.RawMessage {
	props, _ := schema["properties"].(map[string]any)
	if len(props) == 0 {
		return nil
	}
	out := map[string]json.RawMessage{}
	for name, raw := range props {
		p, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		dv, ok := p["default"]
		if !ok || dv == nil {
			continue
		}
		switch p["type"] {
		case "string", "integer", "number", "boolean":
			if b, err := json.Marshal(dv); err == nil {
				out[name] = b
			}
		}
	}
	return out
}

func isEmptyJSON(raw json.RawMessage) bool {
	s := string(raw)
	return s == "null" || s == `""`
}
