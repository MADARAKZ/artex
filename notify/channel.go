package notify

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"
)

// Channel English text。English text**English text**：English text
// English text，English text cfg English text。
type Channel interface {
	// Kind English text，English text。
	Kind() string
	// Validate English text，English text。English text
	// English text，English text「English text」English text「English text」。
	Validate(cfg map[string]any) error
	// Send English text，English text**English text**English text。
	//
	// English text：English text，English text。
	// English text，English text——English text、
	// English text，English text。English text kept English text，
	// English text kept English text，English text。
	//
	// English text，English text *PermanentError English text。
	// English text kept English text，English text。
	Send(ctx context.Context, cfg map[string]any, m Message) (int, error)
	// DefaultRatePerMin English text，English text
	// English text。English text 0 English text。
	DefaultRatePerMin() int
	// SecretKeys English text。API English text，
	// English text。English text
	// （English text Webhook English text，English text secret），
	// English text，English text。
	SecretKeys() []string
	// DestinationKeys English text「English text」English text。
	//
	// English text SecretKeys English text：English text，
	// English text「English text、English text」，English text
	// English text，English text。
	// English text PrepareConfigUpdate。
	DestinationKeys() []string
}

// registry English text。English text init() English text：English text「English text」
// English text，English text，English text。
var registry = map[string]Channel{
	KindDingTalk: dingTalkChannel{},
	KindFeishu:   feishuChannel{},
	KindWeCom:    weComChannel{},
	KindWebhook:  webhookChannel{},
	KindTelegram: telegramChannel{},
	KindEmail:    emailChannel{},
}

// Get English text。
func Get(kind string) (Channel, bool) {
	c, ok := registry[kind]
	return c, ok
}

// ValidKind English text kind English text。
func ValidKind(kind string) bool {
	_, ok := registry[kind]
	return ok
}

// Kinds English text，English text（English text UI English text）。
func Kinds() []string {
	out := make([]string, 0, len(registry))
	for k := range registry {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// PermanentError English text：English text、English text、English text。
// English text（English text、English text、English text 5xx）English text；English text
// English text，English text。
type PermanentError struct{ Err error }

func (e *PermanentError) Error() string { return e.Err.Error() }
func (e *PermanentError) Unwrap() error { return e.Err }

// Permanent English text err English text。err English text nil English text nil，
// English text `return Permanent(someCheck())`。
func Permanent(err error) error {
	if err == nil {
		return nil
	}
	return &PermanentError{Err: err}
}

// IsPermanent English text err English text。
func IsPermanent(err error) bool {
	var pe *PermanentError
	return errors.As(err, &pe)
}

// ---- English texthelper ----
//
// English text JSONB English text，English text encoding/json English text map[string]any，
// English text float64、English text []any。English text helper English text，English text
// English text UI English text（English text）。

// cfgString English text，English text——English text。
func cfgString(cfg map[string]any, key string) string {
	v, ok := cfg[key]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(s)
}

// cfgInt English text，English text float64（JSON English text）English text。
func cfgInt(cfg map[string]any, key string) int {
	switch v := cfg[key].(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		n, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return 0
		}
		return n
	default:
		return 0
	}
}

// cfgBool English text，English text "true"/"1"。
func cfgBool(cfg map[string]any, key string) bool {
	switch v := cfg[key].(type) {
	case bool:
		return v
	case string:
		s := strings.ToLower(strings.TrimSpace(v))
		return s == "true" || s == "1" || s == "yes"
	default:
		return false
	}
}

// cfgStrings English text，English text。
func cfgStrings(cfg map[string]any, key string) []string {
	raw, ok := cfg[key].([]any)
	if !ok {
		// English text，English text。
		if s := cfgString(cfg, key); s != "" {
			return []string{s}
		}
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		s, ok := v.(string)
		if !ok {
			continue
		}
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}

// cfgMap English text（English text HTTP English text），English text，English text。
func cfgMap(cfg map[string]any, key string) map[string]string {
	raw, ok := cfg[key].(map[string]any)
	if !ok {
		return nil
	}
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		out[k] = s
	}
	return out
}
