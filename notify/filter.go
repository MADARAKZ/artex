package notify

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// Filter English text notification_channels.filter English text JSONB English text：English text。
// English text，English text「English text」——English text，English text ParseFilter。
type Filter struct {
	// MinSeverity English text（low/medium/high/critical），English text=English text。
	MinSeverity string `json:"min_severity"`
	// TaskIDs / AssetIDs English text；English text。
	TaskIDs  []int64 `json:"task_ids"`
	AssetIDs []int64 `json:"asset_ids"`
	// VulnClassInclude English text；English text vulnclass English text。
	// VulnClassExclude English text（English text）。
	// English text——English text：English text。
	VulnClassInclude []string `json:"vulnclass_include"`
	VulnClassExclude []string `json:"vulnclass_exclude"`
	// OnStatusChange English text（English text realtime English text）。
	OnStatusChange bool `json:"on_status_change"`
}

// ParseFilter English text。
//
// **English text error。** English text：English text
// Filter（= English text = English text），English text，**English text
// English text**。English text「English text」，English text、
// English text——English text。
func ParseFilter(raw []byte) Filter {
	var f Filter
	if len(raw) == 0 {
		return f
	}
	// English text f English text，English text。
	_ = json.Unmarshal(raw, &f)
	return f
}

// ValidMinSeverity English text s English text（English text）。
func ValidMinSeverity(s string) bool {
	if s == "" {
		return true
	}
	_, ok := severityRank[s]
	return ok
}

// Validate English text**English text**English text，English text。
//
// English text：Match English text `rank >= 0`，English text——
// English text min_severity English text（"hgih"），English text**English text**English text
// 「English text」。English text「English text」English text（English text），
// English text、English text，
// English text。English text「English text」English text。
//
// English text Validate English text**English text**English text。English text ParseFilter English text，
// English text。
func (f Filter) Validate() error {
	if !ValidMinSeverity(f.MinSeverity) {
		return fmt.Errorf("English text %q English text，English text：low / medium / high / critical，English text", f.MinSeverity)
	}
	return nil
}

// Match English text。
//
// **English text error**，English text ParseFilter：English text「English text」English text。
// English text：English text → English text → English text/English text → English text。
func Match(f Filter, s Snapshot) bool {
	// English text。English text，English text
	// English text「English text」English text「English text」，English text。
	if s.Kind == EventFindingStatusChanged && !f.OnStatusChange {
		return false
	}
	if !AtLeast(s.Severity, f.MinSeverity) {
		return false
	}
	if len(f.TaskIDs) > 0 && !slices.Contains(f.TaskIDs, s.TaskID) {
		return false
	}
	if len(f.AssetIDs) > 0 && !intersectsInt(f.AssetIDs, s.AssetIDs) {
		return false
	}
	// English text：English text，English text。
	if len(f.VulnClassExclude) > 0 && containsAnyFold(s.VulnClass, f.VulnClassExclude) {
		return false
	}
	if len(f.VulnClassInclude) > 0 && !containsAnyFold(s.VulnClass, f.VulnClassInclude) {
		return false
	}
	return true
}

func intersectsInt(a, b []int64) bool {
	// English text；English text「English text」，
	// English text map English text。
	for _, v := range b {
		if slices.Contains(a, v) {
			return true
		}
	}
	return false
}

// containsAnyFold English text s English text keywords English text（English text）。
func containsAnyFold(s string, keywords []string) bool {
	lower := strings.ToLower(s)
	for _, kw := range keywords {
		kw = strings.ToLower(strings.TrimSpace(kw))
		if kw != "" && strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}
