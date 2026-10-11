package agent

import (
	"context"
	"errors"
	"fmt"
)

// AbortCause names why an agent run's context was cancelled. Every cancellation
// site should attach one so the activity trace can report the real initiator.
type AbortCause struct {
	Code  string
	Short string
	Text  string
}

func (c *AbortCause) Error() string { return c.Text }

func cause(code, short, text string) *AbortCause {
	return &AbortCause{Code: code, Short: short, Text: text}
}

// Causef builds a cause that includes runtime-specific detail.
func Causef(code, short, format string, args ...any) *AbortCause {
	return &AbortCause{Code: code, Short: short, Text: fmt.Sprintf(format, args...)}
}

var (
	// Task-level execution context.
	AbortPausedByUser = cause("paused_by_user", "English text",
		"English text（POST /api/tasks/{id}/control，action=pause）English text。English text Planner/Worker English text；English text frontier(open)，English text")
	AbortPausedByOrchestrator = cause("paused_by_orchestrator", "English text Agent English text",
		"English text Agent English text pause_task English text。English text Planner/Worker English text；English text frontier(open)，English text")
	AbortTaskDeleted = cause("task_deleted", "English text",
		"English text（DELETE /api/tasks/{id}），English text Planner、Worker English text Agent；English text")
	AbortPausedOnReload = cause("paused_on_reload", "English text",
		"English text。English text；English text Agent")
	AbortGoalMet = cause("goal_met", "English text",
		"English text done，English text Worker；English text stopped，English text")
	AbortSettleDrainTimeout = cause("settle_drain_timeout", "English text",
		"English text timeout English text Worker English text，English text 90 English text drain English text，English text；English text exhausted，English text")

	// Per-work context.
	AbortKilledByPlanner = cause("killed_by_planner", "English text",
		"English text kill_work English text，English text；English text stopped，English text")
	AbortWorkPausedByUser = cause("work_paused_by_user", "English text Worker English text",
		"English text Worker。English text，English text paused；English text、English text、English text，English text")
	AbortWorkCancelledByUser = cause("work_cancelled_by_user", "English text Worker English text",
		"English text Worker。English text；Worker English text，English text——English text，English text")
	AbortWorkFinished = cause("work_finished", "Worker English text context",
		"Worker English text，English text detachWork English text context English text。English text；English text，English text")
	AbortPausedRaceGuard = cause("paused_race_guard", "English text",
		"English text，English text context，English text claim English text Worker English text；English text frontier")

	// Main Agent and standalone conversation contexts.
	AbortChatStoppedByUser = cause("chat_stopped_by_user", "English text",
		"English text，English text Agent English text Agent English text。English text，English text")
	AbortChatPausedWithTask = cause("chat_paused_with_task", "English text Agent English text",
		"English text，English text Agent English text。English text；English text")
	AbortChatTurnFinished = cause("chat_turn_finished", "English text context",
		"English text，English text context English text。English text；English text，English text")

	// Process-level and per-run hard backstop.
	AbortShutdown = cause("shutdown", "English text",
		"English text SIGINT English text SIGTERM，English text、English text。English text Agent English text；English text running English text open English text")
	AbortRunHardTimeout = cause("run_hard_timeout", "English text",
		"English text，English text，English text。English text")
)

// AbortReason resolves the named cause attached to a cancelled run context.
func AbortReason(ctx context.Context) (code, short, text string, ok bool) {
	c := context.Cause(ctx)
	if c == nil {
		return "", "", "", false
	}
	var ac *AbortCause
	if errors.As(c, &ac) {
		return ac.Code, ac.Short, ac.Text, true
	}
	switch {
	case errors.Is(c, context.DeadlineExceeded):
		return "deadline_exceeded", "English text context English text deadline",
			"English text context English text deadline，English text WithTimeoutCause English text: " + c.Error(), true
	case errors.Is(c, context.Canceled):
		return "canceled_no_cause", "English text",
			"English text context English text，English text context.WithCancelCause English text；English text agent/cancelcause.go English text", true
	default:
		return "other", firstLine(c.Error(), 80), c.Error(), true
	}
}
