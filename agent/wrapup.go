package agent

import (
	"strings"

	"github.com/Autumn-27/norma/harness"
)

// English text(wrap-up / settlement prompt):English text agent English text【English text(MaxTurns)】English text
// 【English text(run_seconds/MaxDuration)】English text,SDK English text settlement English text,
// English text agent English text、English text,English text。
//
// English text agent English text(English text agents.wrapup_prompt),English text
// English text。English text【English text】English text;English text、English text。

// WrapupOverride, if set, returns the stored wrap-up prompt for an agent key and
// whether a non-empty one exists. Wired by the server to the agents table (like
// PromptOverride for system prompts). nil / empty → the built-in default is used.
var WrapupOverride func(agentKey string) (string, bool)

// WrapupMaxTurnsOverride, if set, returns the admin-configured turn budget for the
// wrap-up phase of an agent and whether a positive one exists. Wired to the agents
// table. nil / ≤0 → the built-in per-agent default (wrapupTurnDefaults) is used.
var WrapupMaxTurnsOverride func(agentKey string) (int, bool)

// English text,English text agent key English text。worker English text settleWrapUpPrompt
// (English text worker.go),planner/mainagent English text;English text(English text agent)English text。
var wrapupDefaults = map[string]string{
	"worker":    settleWrapUpPrompt,
	"planner":   plannerWrapUpDefault,
	"mainagent": mainAgentWrapUpDefault,
}

// wrapupTurnDefaults: English text agent English text【English text】English text(English text >0 English text)。
// English text 10 English text,English text。English text genericWrapupTurns。
var wrapupTurnDefaults = map[string]int{
	"worker":    10,
	"planner":   10,
	"mainagent": 10,
}

const genericWrapupTurns = 10

const plannerWrapUpDefault = "English text——English text【English text】English text,English text,English text,English text。English text、English text,English text【English text】(English text 0 English text)：(1) English text【English text】English text,English text add_intent English text(English text);(2) English text/English text,English text prove_goal English text met(English text);(3) English text,English text TodoWrite English text,English text。English text,English text。"

const mainAgentWrapUpDefault = "English text,English text。English text/English text。English text**English text**English text、English text,English text。"

const genericWrapUpDefault = "English text。English text,English text**English text**English text、English text(English text)。"

// WrapupDefault returns the built-in default wrap-up prompt for an agent key —
// used by the admin UI as the "restore default" value and empty-field placeholder.
func WrapupDefault(agentKey string) string {
	if d, ok := wrapupDefaults[agentKey]; ok {
		return d
	}
	return genericWrapUpDefault
}

// WrapupTurnsDefault returns the built-in wrap-up turn budget for an agent key —
// used by the admin UI as the "0 = default N" hint.
func WrapupTurnsDefault(agentKey string) int {
	if n, ok := wrapupTurnDefaults[agentKey]; ok {
		return n
	}
	return genericWrapupTurns
}

// resolveWrapup returns the effective wrap-up prompt: the DB override (if set and
// non-empty) over the built-in default.
func resolveWrapup(agentKey string) string {
	if WrapupOverride != nil {
		if t, ok := WrapupOverride(agentKey); ok && strings.TrimSpace(t) != "" {
			return t
		}
	}
	return WrapupDefault(agentKey)
}

// resolveWrapupTurns returns the effective wrap-up turn budget: a positive DB
// override over the built-in per-agent default.
func resolveWrapupTurns(agentKey string) int {
	if WrapupMaxTurnsOverride != nil {
		if v, ok := WrapupMaxTurnsOverride(agentKey); ok && v > 0 {
			return v
		}
	}
	return WrapupTurnsDefault(agentKey)
}

// wrapupSettlement builds the settlement config for an agent's run. Prompt and the
// turn budget are admin-editable per agent; disabled tools are code-owned policy so
// a user can't edit away the "stop probing" guardrail. Resolved fresh each run
// (reads DB live), so edits apply on the next run without a restart.
func wrapupSettlement(agentKey string, disabledTools []string) *harness.Settlement {
	return &harness.Settlement{
		Prompt:        resolveWrapup(agentKey),
		DisabledTools: disabledTools,
		MaxTurns:      resolveWrapupTurns(agentKey),
	}
}

// ---------- English text（English text docs/English text.md）----------
//
// English text per-run English text【English text】：per-run English text"English text run English text"；English text
// "English text、English text"。English text（English text planner：per-run English text"English text"，
// English text"English text、English text"）。English text worker/planner English text。

// WrapupTaskTimeoutOverride / …TurnsOverride：English text DB English text
// （wire English text agents.task_timeout_wrapup_prompt / _max_turns，English text worker/planner）。
var (
	WrapupTaskTimeoutOverride      func(agentKey string) (string, bool)
	WrapupTaskTimeoutTurnsOverride func(agentKey string) (int, bool)
)

var taskTimeoutWrapupDefaults = map[string]string{
	"worker":  workerTaskTimeoutDefault,
	"planner": plannerTaskTimeoutDefault,
}

const workerTaskTimeoutDefault = "**English text，English text**（English text run English text，English text）。English text：(1) English text【English text】English text——English text insert_assets、English text/English text record_fact、English text report_finding；(2) English text/English text；(3) **English text**English text。"

const plannerTaskTimeoutDefault = "**English text，English text**（English text，English text）。English text【English text】English text，English text：English text prove_goal English text met（English text）。**English text**（English text）。English text，English text。"

// TaskTimeoutWrapupDefault English text agent English text（English text/English text）。
func TaskTimeoutWrapupDefault(agentKey string) string {
	return taskTimeoutWrapupDefaults[agentKey] // English text(mainagent/chat)English text
}

// resolveTaskTimeoutWrapup：DB English text(English text) > English text。English text agent English text
// （English text worker/planner），English text per-run English text。
func resolveTaskTimeoutWrapup(agentKey string) string {
	if WrapupTaskTimeoutOverride != nil {
		if t, ok := WrapupTaskTimeoutOverride(agentKey); ok && strings.TrimSpace(t) != "" {
			return t
		}
	}
	return TaskTimeoutWrapupDefault(agentKey)
}

func resolveTaskTimeoutTurns(agentKey string) int {
	if WrapupTaskTimeoutTurnsOverride != nil {
		if v, ok := WrapupTaskTimeoutTurnsOverride(agentKey); ok && v > 0 {
			return v
		}
	}
	return resolveWrapupTurns(agentKey) // English text per-run English text
}

// wrapupSettlementForTask builds settlement for a worker/planner run that is aware
// of the task deadline. See §5 of the design doc:
//   - clamped=true  → English text run English text deadline English text：English text Timeout English text=English text→English text；
//     English text MaxTurns English text=English text、English text→English text per-run English text。
//   - clamped=false → English text：English text reason English text per-run English text（English text wrapupSettlement）。
//
// English text harness English text PromptByReason English text【English text】reason English text，English text build English text。
func wrapupSettlementForTask(agentKey string, disabledTools []string, clamped bool) *harness.Settlement {
	perRun := resolveWrapup(agentKey)
	st := &harness.Settlement{
		Prompt:        perRun, // English text(English text clamped English text reason English text)
		DisabledTools: disabledTools,
		MaxTurns:      resolveWrapupTurns(agentKey),
	}
	if clamped {
		if tt := resolveTaskTimeoutWrapup(agentKey); tt != "" {
			st.PromptByReason = map[harness.TerminalReason]string{
				harness.ReasonTimeout:  tt,     // English text
				harness.ReasonMaxTurns: perRun, // English text、English text
			}
			st.MaxTurns = resolveTaskTimeoutTurns(agentKey)
		}
	}
	return st
}
