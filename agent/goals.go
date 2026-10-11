package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/norma/agentcore"
	"github.com/Autumn-27/norma/llm"
	acperm "github.com/Autumn-27/norma/permission"
	actool "github.com/Autumn-27/norma/tool"
	"github.com/Autumn-27/norma/transcript"
)

// goalsDefaultTmpl is the built-in EDITABLE body (English text [A]) of the goals-decomposer
// prompt, seeded into agent_prompts. No template vars are used today.
const goalsDefaultTmpl = `English text。English text**English text**，English text。

**English text（English text）：English text**
English text「English text / English text」English text【English text、English text】English text，English text set_constraints English text（English text、English text）：
- type=deny：English text（English text「English text」「English text/English text」「English text」「English text」）。
- type=allow：English text/English text（English text「English text」「English text」）。
- English text ≠ English text，English text ≠ English text：English text。
- **English text【English text、English text】**：English text「English text/English text/English textIP/English text/English text」English text**English text**English text/English text**English text**。English text，English text。
  English text：English text https://abc.example.net → English text「English text abc.example.net」English text「English text」；「English text 443，English text」English text「English text」。English text「English text」English text，English text。
- **English text/English text【English text】English text，English text**；English text deny（English text）。
- English text/English text，English text**English text**English text set_constraints。
English text（English text）English text，English text。

**English text = English text/English text**

**English text（English text）**：
- English text、English text、English text
- English text
- English text、English text
- English text

**English text**：
- English text → English text
- English text**English text**English text → English text
- English text vulnclass；English text/English text
- English text

English text set_goals English text。`

// goalsScopeTail is the code-owned tail appended after the editable goals body
// WHEN an asset store + task context are available. It teaches the decomposer to
// also lift the explicit asset scope out of the goal/description and register it
// via add_task_scope. Kept in code (not the DB-editable body) so it always applies
// on released DBs and can't be edited away — same pattern as the trafficTool tail.
const goalsScopeTail = `

**English text：English text**
English text，English text「English text / English text」English text**English text**，English text add_task_scope English text（English text，English text）。**English text：English text，English text。**
- English text URL English text（English text https://xxx.example.com/path、app.example.com）→ English text**English text**，kind=subdomain，value=English text。
  English text：English text https://a1b2c3.lab.example.net/path → kind=subdomain，value=a1b2c3.lab.example.net（**English text** example.net）。
  **English text**English text——English text xxx.example.com English text example.com English text，English text。
- English text**English text、English text**（English text example.com），English text“English text / English text / English text” → English text kind=root_domain，value=example.com。
- English text IP English text → kind=ip / cidr，value=IP English text CIDR。
- **English text**English text（company）——English text、English text，English text，English text plan English text。
English text：
- English text**English text/English text**English text；English text/IP。
- reason English text，English text。
- English text/English text，English text**English text**English text add_task_scope。
English text add_task_scope English text（English text），English text set_goals English text。`

// GoalSpec is one decomposed objective.
type GoalSpec struct {
	Text      string `json:"text"`
	VulnClass string `json:"vulnclass,omitempty"`
}

// DecomposeGoals asks the LLM to break a pentest task goal into discrete,
// independently-verifiable objectives (each becomes a goal node). Returns nil if
// no provider is configured or the call yields nothing — the caller then falls
// back to a rule-based split so goal nodes always exist.
//
// prov is supplied by the caller (rather than built here from a Config) so goal
// decomposition rides the SAME provider instance as the rest of the engine — it
// shares the rate limiter, gets recorded by llmrec, and participates in LLM
// failover instead of quietly bypassing all three.
//
// desc is the task's free-text description (English text：English text/flag English text/English text).
// It is fed alongside the goal so the decomposer no longer splits blind — the
// prompt still forbids inventing anything the two texts don't state.
//
// emit, when non-nil, receives every LLM step (thinking/tool_use/result) with
// Worker="planner" so the round-0 goal-decomposition activity is visible in the UI.
//
// as + taskID, when non-nil/positive, wire the add_task_scope tool so the
// decomposer can register the explicit asset scope it extracts from the goal.
//
// ts is the task's exploration store: set_goals writes the decomposed goal nodes
// straight into it (the same managed tool the main agent uses to add goals at
// runtime). The returned specs are read back from the store so callers can emit
// per-goal activity and detect the "LLM produced nothing" case for their fallback.
func DecomposeGoals(ctx context.Context, prov llm.Provider, dataDir, goalText, desc string, as *db.AssetStore, ts *db.ExplorationStore, taskID int64, emit func(db.Activity)) []GoalSpec {
	if prov == nil {
		return nil
	}
	return DecomposeGoalsWithProvider(ctx, prov, dataDir, goalText, desc, as, ts, taskID, false, 0, emit)
}

// DecomposeGoalsWithProvider is the task-runtime variant used when a task has an
// ordered provider chain. It preserves the same tools and write behavior while
// letting the caller own provider selection/failover. maxTokens is the profile's
// per-reply output cap (0 = send none).
func DecomposeGoalsWithProvider(ctx context.Context, prov llm.Provider, dataDir, goalText, desc string, as *db.AssetStore, ts *db.ExplorationStore, taskID int64, nonStreaming bool, maxTokens int, emit func(db.Activity)) []GoalSpec {
	if prov == nil {
		return nil
	}
	// English text：English text transcript store，English text agentcore English text ctx English text
	// session id（English text writer English text，English text agentcore.Prompt）。English text session-id English text
	// English text/English text（opencode zen English text x-opencode-session English text 400
	// MissingSessionID）English text ctx English text——English text「English text、English text 400」。
	// English text id：English text（English text），English text
	// planner/worker English text，English text llmrec.parseSession English text。
	if ts != nil {
		ctx = transcript.WithSessionID(ctx, fmt.Sprintf("exp%d-goals", ts.ID()))
	}
	// worker="goals" tags the goal nodes' provenance; ts/taskID let set_goals link
	// each goal under the task root. This is the catalog's real set_goals tool, so a
	// web-edited description/schema on it applies here too.
	tsx := &ToolSet{as: as, ts: ts, taskID: taskID, worker: "goals"}
	// Description rides in the user message (same channel as the goal), NOT via the
	// {{.EngagementDescription}} template var — else a prompt that references the var
	// would inject the description twice. System prompt stays pure static instructions.
	sys := renderSystem("goals", goalsDefaultTmpl, GoalsVars{DataDir: dataDir, Now: nowStr()})
	// set_constraints English text(English text asset store):English text「English text」English text
	// (English text agent English text),English text。
	tools := []actool.CoreTool{tsx.setGoals(), tsx.setConstraints()}
	// Wire add_task_scope only when we have a real asset store + task to write to.
	// The scope-extraction tail is appended in lockstep so the prompt never asks for
	// a tool that isn't present.
	if as != nil && taskID > 0 {
		tools = append(tools, tsx.addTaskScope())
		sys += goalsScopeTail
	}
	userMsg := "English text：\n" + goalText
	if d := strings.TrimSpace(desc); d != "" {
		userMsg += "\n\nEnglish text（English text，English text/flag English text/English text；English text，English text）：\n" + d
	}
	// Use captureRun so every LLM step is emitted as an activity record (visible in
	// the plan tab under the round-0 marker). Falls back gracefully when emit is nil.
	captureEmit := func(r db.Activity) {
		if emit != nil {
			r.Worker = "planner"
			emit(r)
		}
	}
	captureRun(ctx, agentcore.Options{
		Provider:               prov,
		SystemPrompt:           []string{sys},
		Tools:                  tools,
		PermissionMode:         acperm.ModeBypass,
		DisableBackgroundTasks: true,
		// 3 English text(English text → English text → English text)English text,English text set_goals。
		MaxTurns:     8,
		NonStreaming: nonStreaming, // English text profile English text Provider.Complete
		MaxTokens:    maxTokens,    // 0 = English text,English text
	}, userMsg, captureEmit)
	// set_goals persisted the goals directly; read them back so the caller sees what
	// was written (empty slice ⇒ the LLM produced nothing ⇒ caller falls back).
	if ts == nil {
		return nil
	}
	nodes, _ := ts.ListByKind(db.KindGoal, 10000)
	var out []GoalSpec
	for _, n := range nodes {
		var p struct {
			Text      string `json:"text"`
			VulnClass string `json:"vulnclass"`
		}
		_ = json.Unmarshal(n.Payload, &p)
		if strings.TrimSpace(p.Text) != "" {
			out = append(out, GoalSpec{Text: p.Text, VulnClass: p.VulnClass})
		}
	}
	return out
}
