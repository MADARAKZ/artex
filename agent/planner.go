package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/intercept"
	"github.com/Autumn-27/norma/agentcore"
	"github.com/Autumn-27/norma/llm"
	"github.com/Autumn-27/norma/permission"
	actool "github.com/Autumn-27/norma/tool"
	"github.com/Autumn-27/norma/transcript"
)

// Planner is the event-driven LLM planner (docs §4.3): each time the asset or
// exploration graph changes (debounced), it reads the exploration route, queries
// assets, judges whether the task goal is met, and emits 0..N exploration intents
// into the frontier. It is the sole intent generator.
type Planner struct {
	findingRecorder   FindingRecorder
	prov              llm.Provider
	model             string
	tx                *transcript.Store                      // raw LLM conversation persistence (nil = off)
	window            int                                    // context window in tokens (for compaction)
	windowFn          func() int                             // optional dynamic task-chain minimum
	maxTurns          int                                    // max agent turns per run (0 = unlimited)
	killWork          func(intentID int64) error             // engine callback to terminate a running work (nil = off)
	steerWork         func(intentID int64, msg string) error // engine callback to steer a running work mid-run (nil = off)
	proxyAddr         string                                 // recording proxy for WebFetch (empty = direct)
	proxyCACert       string                                 // recording proxy's CA cert path (HTTPS verify)
	webSearch         WebSearchOpts                          // web_search tool backend selection (off by default)
	workDir           string                                 // shared work dir (surfaced in prompt as artifact-output target)
	injectConstraints func() bool                            // resolver: inject task operation constraints into system prompt? (nil = yes)
	nonStreamingFn    func() bool                            // resolver: use non-streaming (Complete) path? (nil = streaming)
	noaEnabledFn      func() bool                            // resolver: use experimental noa compaction? (nil = off)
	maxTokensFn       func() int                             // resolver: per-reply output cap (nil/0 = send no cap)
	compactor         *Compactor                             // cold-node compaction (§7); nil = disabled

	// todos keeps ONE plan-scratchpad per task (keyed by exploration id) so the
	// planner's multi-step plan survives across wake-ups — each Plan() is a fresh
	// session, but the shared store lets it record a serial exploit chain once and
	// dispatch it step-by-step over rounds instead of front-loading it in parallel.
	todoMu sync.Mutex
	todos  map[int64]*actool.TodoStore
}

func NewPlanner(prov llm.Provider, model, workDir string, tx *transcript.Store, window, maxTurns int) *Planner {
	return &Planner{prov: prov, model: model, workDir: workDir, tx: tx, window: window, maxTurns: maxTurns, todos: map[int64]*actool.TodoStore{}}
}

func (p *Planner) SetCompactionWindowResolver(fn func() int) { p.windowFn = fn }

// SetCompactor wires the cold-node compactor (cold-digest §7). Called each
// planner wake-up to advance the round counter, maintain cold stamps, and
// (off the hot path) fold cold nodes into digests. nil = feature disabled.
func (p *Planner) SetCompactor(c *Compactor) { p.compactor = c }

// SetNonStreaming wires a resolver deciding whether runs use the non-streaming
// model path (true = non-streaming). nil/unset = streaming (default).
func (p *Planner) SetNonStreaming(fn func() bool) { p.nonStreamingFn = fn }

func (p *Planner) nonStreaming() bool { return p.nonStreamingFn != nil && p.nonStreamingFn() }

// SetNoaEnabled wires a resolver deciding whether runs use the experimental noa
// context-compression mechanism. nil/unset = off (built-in compaction). Read per
// run so the settings toggle takes effect without rebuilding the agent.
func (p *Planner) SetNoaEnabled(fn func() bool) { p.noaEnabledFn = fn }

// SetMaxTokens wires a resolver for the per-reply output cap. nil/unset or 0 =
// send no cap and let the endpoint decide. Read per run, like nonStreaming.
func (p *Planner) SetMaxTokens(fn func() int) { p.maxTokensFn = fn }

func (p *Planner) maxTokens() int {
	if p.maxTokensFn == nil {
		return 0
	}
	return p.maxTokensFn()
}

func (p *Planner) compactionWindow() int {
	if p.windowFn != nil {
		return p.windowFn()
	}
	return p.window
}

// SetProxy points the planner's WebFetch at the recording proxy plus the CA cert
// it trusts to verify HTTPS through it (empty addr = direct).
func (p *Planner) SetProxy(addr, caCert string) { p.proxyAddr, p.proxyCACert = addr, caCert }

// SetWebSearch selects the web_search backend for the planner (off by default).
func (p *Planner) SetWebSearch(o WebSearchOpts) { p.webSearch = o }

// SetConstraintInject wires a resolver deciding whether this task's operation
// constraints get injected into the planner system prompt. Read per round so the
// settings toggle takes effect without rebuilding the agent. nil = inject (default).
func (p *Planner) SetConstraintInject(fn func() bool) { p.injectConstraints = fn }

// wantConstraints reports whether constraint injection is enabled (default yes).
func (p *Planner) wantConstraints() bool { return p.injectConstraints == nil || p.injectConstraints() }

// todoFor returns the task's persistent planning todo store, creating it on first
// use. Shared across all of this task's planner wake-ups.
func (p *Planner) todoFor(expID int64) *actool.TodoStore {
	p.todoMu.Lock()
	defer p.todoMu.Unlock()
	s := p.todos[expID]
	if s == nil {
		s = actool.NewTodoStore()
		p.todos[expID] = s
	}
	return s
}

// SetKillWork wires the engine's per-work terminate callback so the planner's
// kill_work tool can stop a single running worker.
func (p *Planner) SetKillWork(fn func(intentID int64) error) { p.killWork = fn }

// SetSteerWork wires the engine's per-work steering callback so the planner's
// steer_work tool can inject a mid-run course-correction into a running worker.
func (p *Planner) SetSteerWork(fn func(intentID int64, msg string) error) { p.steerWork = fn }

// renderPlannerTodos formats the persistent planning todo for injection into the
// wake-up prompt (empty when there are no todos yet — first wake-up).
func renderPlannerTodos(items []actool.Todo) string {
	if len(items) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("\n\n【English text（English text，English text）】：\n")
	for _, it := range items {
		mark := map[actool.TodoStatus]string{actool.TodoPending: "☐", actool.TodoInProgress: "▶", actool.TodoCompleted: "✔"}[it.Status]
		if mark == "" {
			mark = "☐"
		}
		b.WriteString(fmt.Sprintf("  %s %s\n", mark, it.Content))
	}
	b.WriteString("English text：English text【English text / English text fact English text】English text；English text TodoWrite English text（English text fact English text completed）。English text pending/in_progress English text。")
	return b.String()
}

// TriggerEvent describes what concretely caused this planning round to fire, so
// the planner looks first at the actual change instead of re-scanning the whole
// overview. Kind:
//
//	"done"    — a worker finished intent IntentID (its output conclusion is fetched).
//	"finding" — a worker reported a finding on intent IntentID (Detail = English text).
//	"goal"    — the human (via English text agent English text set_goals) added one OR MORE goals in a
//	            single call (Goals = English text，1+ English text；set_goals English text).
//	"goal_deleted" — the human deleted a goal from English text (Detail = English text).
//	"goal_edited"  — the human edited a goal from English text (OldGoal→NewGoal English text).
//	"cancelled" — the human deleted intent IntentID (Detail = English text). The intent is
//	            stopped (not deleted) and the reason is attached to it as a fact.
type TriggerEvent struct {
	Kind     string
	IntentID int64
	Detail   string
	Summary  string   // Kind=="cancelled" English text：English text（English text，English text）
	Goals    []string // Kind=="goal" English text：English text set_goals English text（1 English text）
	OldGoal  string   // Kind=="goal_edited" English text：English text
	NewGoal  string   // Kind=="goal_edited" English text：English text
	Hints    []string // Kind=="hint" English text：English text add_hint English text（1 English text）
}

// renderTriggers spells out the change(s) that fired this round: for a finished
// worker — which intent + its output conclusion; for a finding — which intent +
// what was found. Empty for time/heartbeat wakes. Reads the store (best-effort;
// a blank field never blocks the round).
func renderTriggers(ts *db.ExplorationStore, evs []TriggerEvent) string {
	if len(evs) == 0 || ts == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString("\n\n【English text（English text，English text）】：")
	for _, ev := range evs {
		switch ev.Kind {
		case "goal":
			if len(ev.Goals) == 1 {
				b.WriteString(fmt.Sprintf("\n- English text（English text agent）English text：%s —— English text，English text（English text）。", ev.Goals[0]))
			} else {
				b.WriteString(fmt.Sprintf("\n- English text（English text agent）English text %d English text：%s —— English text，English text。", len(ev.Goals), strings.Join(ev.Goals, "；")))
			}
		case "hint":
			if len(ev.Hints) == 1 {
				b.WriteString(fmt.Sprintf("\n- English text（English text agent）English text：%s —— English text，English text/English text（English text）。", ev.Hints[0]))
			} else {
				b.WriteString(fmt.Sprintf("\n- English text（English text agent）English text %d English text：%s —— English text，English text/English text。", len(ev.Hints), strings.Join(ev.Hints, "；")))
			}
		case "goal_deleted":
			b.WriteString(fmt.Sprintf("\n- English text：%s —— English text，English text/English text（English text）。", ev.Detail))
		case "goal_edited":
			b.WriteString(fmt.Sprintf("\n- English text，English text「%s」English text「%s」—— English text（English text）。", ev.OldGoal, ev.NewGoal))
		case "finding":
			b.WriteString(fmt.Sprintf("\n- English text #%d（%s）English text worker English text finding：%s", ev.IntentID, intentSummary(ts, ev.IntentID), ev.Detail))
		case "cancelled":
			// English text Summary（English text，intentSummary English text）。
			sm := ev.Summary
			if sm == "" {
				sm = intentSummary(ts, ev.IntentID)
			}
			b.WriteString(fmt.Sprintf("\n- English text #%d English text，English text：%s、English text：%s。English text（English text）；English text。", ev.IntentID, sm, ev.Detail))
		default: // "done"
			b.WriteString(fmt.Sprintf("\n- English text #%d（%s）English text worker English text，English text：%s", ev.IntentID, intentSummary(ts, ev.IntentID), workerOutput(ts, ev.IntentID)))
			if fids := factIDsYielded(ts, ev.IntentID); fids != "" {
				b.WriteString(fmt.Sprintf("；English text id：%s ", fids))
			}
		}
	}
	b.WriteString("\n（English text node_detail / get_worker_output / list_findings English text。）")
	return b.String()
}

// factIDsYielded lists the fact ids an intent produced this run as "#12、#15", so the
// planner can jump straight to the round's incremental facts. Empty (best-effort) when
// the intent yielded no facts or the lookup fails.
func factIDsYielded(ts *db.ExplorationStore, id int64) string {
	ids, err := ts.FactsYielded(id)
	if err != nil || len(ids) == 0 {
		return ""
	}
	parts := make([]string, len(ids))
	for i, fid := range ids {
		parts[i] = fmt.Sprintf("#%d", fid)
	}
	return strings.Join(parts, "、")
}

// intentSummary reads an intent node's one-line summary (best-effort, "?" on miss).
func intentSummary(ts *db.ExplorationStore, id int64) string {
	n, err := ts.GetNode(id)
	if err != nil || n == nil {
		return "?"
	}
	var p map[string]any
	if json.Unmarshal(n.Payload, &p) == nil {
		if s, ok := p["summary"].(string); ok && s != "" {
			return s
		}
	}
	return "?"
}

// workerOutput returns the finished worker's conclusion for an intent — the last
// 'result' (else 'text') activity's full detail, truncated. Same source get_worker_output uses.
func workerOutput(ts *db.ExplorationStore, id int64) string {
	acts, _, err := ts.ActivityList(&id, 0, 1000)
	if err != nil {
		return "(English text)"
	}
	var pick *db.Activity
	for i := range acts {
		if acts[i].Kind == "result" {
			pick = &acts[i]
		} else if acts[i].Kind == "text" && pick == nil {
			pick = &acts[i]
		}
	}
	if pick == nil {
		return "(English text work English text)"
	}
	out, _ := ts.ActivityDetail(pick.ID)
	if out == "" {
		out = pick.Summary
	}
	return truncOutput(out, 800)
}

// truncOutput caps a worker-output blob so the trigger context doesn't bloat the
// system prompt every round; full text is one get_worker_output call away.
func truncOutput(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + " …（English text，English text get_worker_output）"
}

// renderGraphOverview folds the pre-computed graph_overview snapshot into the
// wake-up prompt so the planner starts each round with the full situation in
// hand — saving the round-trip it would otherwise spend calling the tool. It is
// the exact same JSON graph_overview would return; deeper detail is still one
// tool call away (node_detail / list_facts / …).
func renderGraphOverview(data map[string]any) string {
	b, err := json.Marshal(data)
	if err != nil {
		return "" // fall back to the model calling graph_overview itself
	}
	return "\n\n【English text（graph_overview English text，English text；English text node_detail/list_facts English text）】：\n" + string(b)
}

// plannerDefaultTmpl is the built-in EDITABLE body (English text [A]) of the planner prompt,
// seeded into agent_prompts. Goal is a {{.Goal}} template var; the English text
// tail is code-owned (artifactSpec) and appended by plannerSystem after rendering.
const plannerDefaultTmpl = `English text"English text"，English text（English text）。English text：English text → English text → **English text**English text。English text、English text：English text【English text/English text】English text【English text】，English text plan English text。

English text：{{.Goal}}

**English text（English text）**：
- **English text（English text）**：English text【English text】English text【English text open English text running English text】（frontier_open=0 English text running_intents English text），English text【English text】English text——English text work English text、English text，English text 0 English text=English text；English text recent_done English text，English text done/exhausted/blocked English text。
- English text，**English text 0 English text，English text**（English text"English text"English text）：①**English text**——English text open/running English text（English text）；②**English text**——English text work English text、English text（English text，English text）。
- English text：English text【English text、English text work】English text，English text，English text——English text 0 English text。

**English text**：

1. **English text**（English text graph_overview English text，English text）：task（English text+English text/English text）、English text、goals+English text、open/running/recent_done English text、sites_without_endpoints（English text，English text）、facts（English text，English text）、recent_facts（{id,summary,confidence?}）。
   - **English text**：English text（goals/English text/facts/findings）English text；**English text**（English text，English text、English text）——English text。
   - **English text**：English text parents（English text：English text/English text）English text yields（English text：English text/English text），recent_facts English text from_intent；English text"English text、English text"。
   - **English text/English text**（recent_facts English text"English text/English text"English text）English text worker English text、English text：English text node_detail(id) English text evidence——evidence English text、confidence=observed English text；evidence English text、English text"English text/English text"、English text confidence=inferred English text，English text【English text】English text，English text，English text（**English text**；English text、English text，English text、English text）。
   - **English text**：list_facts（English text，English text，English text 20，English text q English text、before English text，English text total/has_more）、list_findings（English text）、node_detail(id)（English text/English text；English text/recent_facts English text）、list_assets（pull：q English text、type/company_id/task_id English text、English text，English text id/ids English text）、asset_neighbors。English text，English text。

2. **English text（English text）**：goals English text；English text/English text，English text prove_goal(goal_id, evidence_id, reason) English text met。**English text，English text**——English text prove_goal English text，English text"English text"English text。
   - ⚠️ **English text（English text）**：English text（English text X%、English text N English text flag、English text）English text，prove_goal English text【English text】English text graph_overview English text（coverage.pct、findings_total English text）：English text【English text】prove_goal，English text；English text"English text/English text"English text met。English text：English text 100% English text coverage.pct=40% → English text，English text。

3. **（English text，English text、English text）English text**：English text fact（recent_facts English text、English text）、English text，English text Bash English text、English text（English text 1–2 English text curl English text/English text）。**English text**——English text/English text/English text，English text/English text/English text（English text worker English text，English text）。English text：
   - English text worker English text fact（facts>0 / recent_facts English text）→【English text】English text，English text fact，English text"English text"English text"English text"；English text → English text worker English text，English text curl。
   - English text ≤3 English text，English text；English text"English text"English text"English text"（English text/English text、English text id、English text、English text、English text/English text/English text——English text worker English text），English text。
   - English text/English text，English text。

4. **English text**：**English text"English text"English text【English text】，English text"English text"**——English text，English text"English text，English text、English text、English text"，English text"English text"。English text【English text】（English text/English text），English text、English text、English text，English text open + running + recent_done English text：
   - English text open/running English text → English text（English text）。
   - English text recent_done English text → **English text state（English text）English text，English text**：
     · **done（English text）**：English text → English text；English text yields English text fact English text、English text state；English text【English text】（English text/English text/English text/English text）English text，English text summary English text；English text、"English text"English text，English text。
     · **exhausted（English text、English text，English text）/ blocked（English text、English text）**：English text、English text——English text get_worker_trace / get_worker_output English text、English text，English text：English text → English text"English text"；English text（blocked English text）→ English text；English text → English text/English text。English text trace English text，English text state English text。
   - English text → English text。
   - English text open/running English text → English text、English text（English text/English text work，English text）；English text recent_done English text、English text open/running English text → English text。
   - **English text**：coverage English text/English text、English text；English text（English text RCE/English text/English text）English text，English text【English text】，English text、English text。
   - **English text、English text**：English text，English text/English text，English text【English text】English text（English text/English text/English text），English text，English text（English text，English text）；English text，English text。English text 2–3 English text（English text"English text"English text"English text"），English text【English text】English text。**English text【English text】**：English text/English text/English text/English text，English text。

   **English text：English text，English text。** English text（①→②→③，English text）：English text（English text/English text）；English text TodoWrite English text（English text），English text"English text"English text（English text），English text fact English text（English text）English text completed。"English text"English text（"English text"English text"English text"English text）；English text【English text、English text】English text（English text）English text。

5. **English text**：English text【English text】add_intent English text（intents English text，English text 4 English text，English text）：
   - **summary**：English text（English text + English text + English text），English text；English text。
   - **asset_ids**：English text/English text id（English text，0/1/English text，English text list_assets）——English text（English text/English text/English text/English text）English text，English text、English text，English text；English text。
   - **parent_ids**：English text（English text，0/1/English text）——English text，English text/English text id，English text。

English text、English text；English text、English text，English text。English text、English text、English text。`

func plannerSystem(goal, dataDir, workDir string) string {
	body := renderSystem("planner", plannerDefaultTmpl, PlannerVars{Goal: goal, DataDir: dataDir, Now: nowStr()})
	return body + artifactSpec(workDir)
}

// Plan runs one planning round. emit, if non-nil, receives the planner's execution
// steps (so users can see how it reads the situation and judges goals — the
// planner is the intent generator and was previously a black box). Returns whether
// the planner judged the goal met.
// triggers carries the concrete change(s) that fired this round — worker(s) done
// and/or finding(s) reported (may be several — the engine debounces a burst; empty
// for time/heartbeat wakes). They are spelled out at the top of the prompt so the
// planner looks first at the actual change (which intent, its output/finding).
func (p *Planner) Plan(ctx context.Context, taskID int64, as *db.AssetStore, ts *db.ExplorationStore, goal string, triggers []TriggerEvent, emit func(db.Activity)) (met bool, reason string, err error) {
	// cold-digest §2.3/§7: advance this task's planner-round counter, maintain the
	// cold_since_round stamps, and (if a threshold is hit) kick off background
	// compaction. Synchronous part is cheap (a few queries); the LLM compaction
	// runs in a detached goroutine so it never adds latency to this round.
	p.compactor.OnPlannerRound(ctx, ts)
	tsx := NewToolSet(ts, "planner")
	tsx.SetFindingRecorder(p.findingRecorder)
	if as != nil {
		tsx.SetAssetStore(as, as.Companies())
	}
	tsx.SetTaskID(taskID)
	tsx.SetCoverageEnabled(as == nil || as.CoverageEnabled(taskID))
	tsx.killWork = p.killWork   // enable kill_work tool (nil = unavailable)
	tsx.steerWork = p.steerWork // enable steer_work tool (nil = unavailable)
	if origin, _ := ts.OriginFactID(); origin > 0 {
		tsx.SetOwnerNode(origin) // planner-side anchors default to the task root (origin fact)
	}
	// English text + English text（Read/Write/Edit/MultiEdit/LS/Glob/Grep/Bash）
	// English text add_task_scope/list_untested_assets（English text prompt）。
	base := append(tsx.DropCoverageTools(tsx.PlannerTools()), actool.DefaultTools()...)
	ctx = WithRunInfo(ctx, RunInfo{TaskID: taskID, ExplorationID: explorationID(ts)})
	tools, def, cleanup := AugmentTools(ctx, "planner", base)
	defer cleanup()
	// English text（English text + English text）English text【English text user English text】(English text input)，system
	// English text。move-out English text system English text、English text；English text，English text
	// English text compaction English text（planner English text，English text）。situational English text input。
	situational := renderTriggers(ts, triggers) + renderGraphOverview(tsx.graphOverviewData())
	// English text deadline / English text(English text ctx English text,English text taskclock.go)。English text
	// planner English text【English text】English text user English text(English text situational),English text
	// English text、English text。
	tc := taskClockFrom(ctx)
	if tc.Final {
		situational += "\n\n【English text（English text，English text）】：" + resolveTaskTimeoutWrapup("planner")
	}
	// English text <workDir>/tasks/<taskID>，English text。
	taskDir := ensureRunDir(p.workDir, taskID, 0)
	ctx = intercept.WithReviewContext(ctx, taskDir, intercept.ReviewBackground{})
	sysBody := plannerSystem(goal, p.workDir, taskDir)
	if p.wantConstraints() {
		sysBody += constraintBlock(ts) // English text(English text)English text,English text
	}
	system, boundary := deferredSystem(sysBody, def)
	// planner English text;English text deadline English text MaxDuration English text,English text
	// English text(English text→English text,English text→per-run English text)。
	maxDur, clamped := clampMaxDuration(tc.DeadlineUnix, 0)
	settle := wrapupSettlement("planner", nil)
	if tc.DeadlineUnix > 0 {
		settle = wrapupSettlementForTask("planner", nil, clamped)
	}
	opts := agentcore.Options{
		Provider:        p.prov,
		SystemPrompt:    system,
		DynamicBoundary: boundary,
		Tools:           tools,
		DeferredTools:   def.Deferred,
		UnlockSet:       def.Unlock,
		PermissionMode:  permission.ModeBypass,
		EnableWebFetch:  true, // English text；English text CA English text MITM English text HTTPS English text
		WebFetchProxy:   p.proxyAddr,
		WebFetchCACert:  p.proxyCACert,
		// English text(English text)。ddgs English text key；brave-free English text BraveKey；tavily English text TavilyKey。
		// WebSearchProxy English text(http/https/socks5)，English text MITM English text；English text。
		EnableWebSearch:       p.webSearch.Enabled,
		WebSearchBackend:      p.webSearch.Backend,
		BraveSearchAPIKey:     p.webSearch.BraveKey,
		TavilySearchAPIKey:    p.webSearch.TavilyKey,
		DeepSeekSearchBaseURL: p.webSearch.DeepSeekBaseURL,
		DeepSeekSearchAPIKey:  p.webSearch.DeepSeekAPIKey,
		DeepSeekSearchModel:   p.webSearch.DeepSeekModel,
		WebSearchProxy:        p.webSearch.Proxy,
		BashEnv:               proxyEnv(p.proxyAddr, p.proxyCACert), // Bash English text+English text CA
		WorkingDir:            taskDir,                              // English text <workDir>/tasks/<taskID>
		ToolOutputDir:         cmdOutDir(taskDir),
		MaxTurns:              p.maxTurns, // 0 = unlimited (configurable in agent management)
		MaxDuration:           maxDur,     // 0=English text;English text deadline English text=English text deadline English text
		Compaction:            compactionConfig(p.compactionWindow()),
		// English text：English text（session English text，store English text）。
		Todos: p.todoFor(ts.ID()),
		// English text【English text】English text→ SDK English text:English text(English text add_intent、
		// English text prove_goal、English text TodoWrite),English text——planner English text。
		// clamped(English text deadline English text)English text PromptByReason(English text wrapupSettlementForTask)。
		Settlement:   settle,
		NonStreaming: p.nonStreaming(), // English text profile English text Provider.Complete
		MaxTokens:    p.maxTokens(),    // 0 = English text,English text
	}
	if p.tx != nil { // persist raw LLM conversation; one accumulating file per task's planner
		opts.Transcript = p.tx
		opts.SessionID = fmt.Sprintf("exp%d-planner", ts.ID())
	}
	// English text:English text noa English text(English text <workDir>/noa/<SessionID> English text,English text)。
	noaSession := fmt.Sprintf("exp%d-planner", ts.ID())
	enableNoa(&opts, p.noaEnabledFn, p.workDir, noaSession, noaWarn(noaSession))
	// English text（English text + English text）English text user English text（English text input）。user English text
	// English text + English text（todo English text，English text，English text user English text）。
	// English text「English text」English text：English text → English text【English text】English text；English text
	// (English text / hint / English text) → English text"English text",English text。
	lead := "English text（English text【English text】），English text："
	if len(triggers) == 0 {
		lead = "English text**English text（English text）/English text**English text——English text。English text：English text steer_work English text、English text kill_work English text；English text、English text："
		// English text/English text,English text open English text running English text → English text(English text worker English text、
		// English text)。English text planner English text,English text。
		if active, err := ts.HasActiveIntent(); err == nil && !active {
			lead = "English text**English text（English text）**English text,English text**English text open English text running English text**——English text worker English text、English text,English text。English text**English text**English text、English text**English text**English text(English text 0 English text);English text,English text："
		}
	}
	input := lead + situational + "\n\nEnglish text，English text。English text【English text】（English text/English text）English text prove_goal English text。**English text：English text、English text open English text running English text（frontier_open=0 English text running_intents English text），English text——English text work English text、English text，English text 0 English text=English text。English text open/running English text、English text，English text。**" +
		renderPlannerTodos(opts.Todos.List())
	// MaxDuration English text(English text ctx English text),English text
	// English text,English text ctx English text。ctx English text pause / kill / shutdown。
	_, _, err = captureRun(ctx, opts, input,
		func(r db.Activity) {
			if emit != nil {
				r.Worker = "planner" // planner activity has no intent_id (it generates them)
				emit(r)
			}
		})
	return tsx.GoalMet, tsx.Reason, err
}
