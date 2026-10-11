// Package notify English text IM / English text。
//
// English text：English text**English text**，English text。English text、English text server。English text
// English text map[string]any English text（English text notification_channels.config English text JSONB English text），
// English text Message English text。English text：English text、UTF-8 English text、English text
// English text PostgreSQL English text，English text server English text。
//
// English text：Channel English text**English text**。English text Channel English text
// （English text）English text，English text cfg English text，
// English text webhook URL English text。
package notify

// English text。English text notification_channels.kind English text，English text server English text
// English text（English text findings.status English text，English text DB CHECK，English text）。
const (
	KindDingTalk = "dingtalk" // English text
	KindFeishu   = "feishu"   // English text(English text Lark)English text
	KindWeCom    = "wecom"    // English text
	KindWebhook  = "webhook"  // English text Webhook：English text/English text/JSON English text
	KindTelegram = "telegram" // Telegram Bot API
	KindEmail    = "email"    // SMTP English text
)

// English text，English text notification_events.kind。
const (
	EventFindingCreated       = "finding_created"
	EventFindingStatusChanged = "finding_status_changed"
)

// InitKind English text config English text kind English text。
const InitKind = KindDingTalk

// severityRank English text。English text 0，English text
// min_severity English text——English text，English text。
var severityRank = map[string]int{
	"low":      1,
	"medium":   2,
	"high":     3,
	"critical": 4,
}

// SeverityRank English text；English text 0。
func SeverityRank(severity string) int { return severityRank[severity] }

// SeverityLabel English text emoji English text，English text。
// English text，English text。
func SeverityLabel(severity string) string {
	switch severity {
	case "critical":
		return "🔴 English text"
	case "high":
		return "🟠 English text"
	case "medium":
		return "🟡 English text"
	case "low":
		return "🔵 English text"
	default:
		return severity
	}
}

// StatusLabel English text，English text。
func StatusLabel(status string) string {
	switch status {
	case "pending":
		return "English text"
	case "in_progress":
		return "English text"
	case "confirmed":
		return "English text"
	case "resolved":
		return "English text"
	case "fixed":
		return "English text"
	case "false_positive":
		return "English text"
	case "ignored":
		return "English text"
	case "duplicate":
		return "English text"
	case "risk_accepted":
		return "English text"
	default:
		return status
	}
}

// AtLeast English text severity English text min English text。min English text，English text。
// English text severity English text 0，English text min English text（English text severityRank English text）。
func AtLeast(severity, min string) bool {
	if min == "" {
		return true
	}
	return SeverityRank(severity) >= SeverityRank(min)
}
