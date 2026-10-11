package agent

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Autumn-27/norma/harness"
)

// runTrace retains the latest tool call so an interrupted run can identify the
// operation that was still in flight.
type runTrace struct {
	startedAt time.Time
	id        string
	name      string
	input     string
	at        time.Time
	pending   bool
}

func (t *runTrace) start(id, name, input string) {
	t.id, t.name, t.input, t.at, t.pending = id, name, input, time.Now(), true
}

func (t *runTrace) done(id string) {
	if id == t.id {
		t.pending = false
	}
}

var reasonHint = map[harness.TerminalReason]string{
	harness.ReasonCompleted:         "English text，English text；English text",
	harness.ReasonMaxTurns:          "English text(MaxTurns)：SDK English text，English text exhausted，English text，English text",
	harness.ReasonTimeout:           "English text(MaxDuration)：English text，English text，English text exhausted",
	harness.ReasonModelError:        "English text API English text（English text、English text、English text、English text 5xx English text），English text blocked——English text；English text（get_worker_trace）English text",
	harness.ReasonBlockingLimit:     "English text，English text；English text",
	harness.ReasonPromptTooLong:     "English text，English text",
	harness.ReasonImageError:        "English text；English text",
	harness.ReasonStopHookPrevented: "Stop English text，English text；English text Guard English text",
	harness.ReasonHookStopped:       "English text，English text；English text tool_result English text",
	harness.ReasonAbortedStreaming:  "English text",
	harness.ReasonAbortedTools:      "English text",
}

// terminalText renders a terminal event with no final text into a compact summary
// and a Markdown detail block.
func terminalText(ctx context.Context, term *harness.Terminal, tr *runTrace) (string, string) {
	reason := term.Reason
	aborted := reason == harness.ReasonAbortedStreaming || reason == harness.ReasonAbortedTools
	// Prompt may return ctx.Err directly without a terminal event. Preserve the
	// cancellation cause instead of falling back to an empty/unknown terminal reason.
	if reason == "" && ctx.Err() != nil {
		aborted = true
	}

	var sum string
	if aborted {
		_, short, _, ok := AbortReason(ctx)
		if !ok {
			short = "English text"
		}
		stage := "English text"
		switch reason {
		case harness.ReasonAbortedStreaming:
			stage = "English text"
		case harness.ReasonAbortedTools:
			stage = "English text"
		}
		sum = "（English text：" + short + "；English text" + stage + progressSuffix(term, tr) + "，English text）"
	} else if reason == harness.ReasonMaxTurns || reason == harness.ReasonTimeout {
		sum = "（English text(" + string(reason) + ")，English text" + progressSuffix(term, tr) + "；English text）"
	} else {
		hint := terminalReasonHint(reason)
		sum = "（English text，English text " + terminalReasonLabel(reason) + "：" + firstLine(hint, 80) + "）"
	}

	var b strings.Builder
	b.WriteString(sum)
	b.WriteString("\n\n")
	displayReason := terminalReasonLabel(reason)
	fmt.Fprintf(&b, "- **English text**: `%s` - %s\n", displayReason, terminalReasonHint(reason))
	if aborted {
		code, _, why, ok := AbortReason(ctx)
		if ok {
			fmt.Fprintf(&b, "- **English text** (`%s`): %s\n", code, why)
		} else {
			b.WriteString("- **English text**: English text；English text context.WithCancelCause English text\n")
		}
	}
	if term.Err != nil {
		fmt.Fprintf(&b, "- **English text**: `%v`\n", term.Err)
	}
	if aborted && strings.TrimSpace(term.Text) != "" {
		b.WriteString("- **English text**:\n\n")
		b.WriteString(term.Text)
		b.WriteString("\n\n")
	}
	if term.Turns > 0 {
		fmt.Fprintf(&b, "- **English text**: %d English text\n", term.Turns)
	}
	if !tr.startedAt.IsZero() {
		fmt.Fprintf(&b, "- **English text**: %s\n", roundDur(time.Since(tr.startedAt)))
	}
	if u := term.Usage; u.InputTokens+u.OutputTokens+u.CacheReadTokens+u.CacheWriteTokens > 0 {
		fmt.Fprintf(&b, "- **English text token**: English text %d / English text %d / English text %d / English text %d\n",
			u.InputTokens, u.OutputTokens, u.CacheReadTokens, u.CacheWriteTokens)
	}
	if tr.name == "" {
		b.WriteString("- **English text**: English text\n")
	} else if tr.pending {
		fmt.Fprintf(&b, "- **English text**: `%s`（English text %s，**English text**）\n\n  ```json\n  %s\n  ```\n",
			tr.name, roundDur(time.Since(tr.at)), firstLine(tr.input, 300))
	} else {
		fmt.Fprintf(&b, "- **English text**: `%s`（English text）\n", tr.name)
	}
	return sum, b.String()
}

func terminalReasonLabel(reason harness.TerminalReason) string {
	if reason == "" {
		return "context_canceled"
	}
	return string(reason)
}

func terminalReasonHint(reason harness.TerminalReason) string {
	if hint := reasonHint[reason]; hint != "" {
		return hint
	}
	if reason == "" {
		return "English text context English text，English text Terminal English text"
	}
	return "English text；harness English text TerminalReason，English text reasonHint"
}

func progressSuffix(term *harness.Terminal, tr *runTrace) string {
	var parts []string
	if term.Turns > 0 {
		parts = append(parts, fmt.Sprintf("%d English text", term.Turns))
	}
	if !tr.startedAt.IsZero() {
		parts = append(parts, roundDur(time.Since(tr.startedAt)))
	}
	if len(parts) == 0 {
		return ""
	}
	return "，English text " + strings.Join(parts, " / ")
}

func roundDur(d time.Duration) string {
	switch {
	case d < time.Minute:
		return d.Round(100 * time.Millisecond).String()
	case d < time.Hour:
		return d.Round(time.Second).String()
	default:
		return d.Round(time.Minute).String()
	}
}
