package notify

import "testing"

func TestParseFilterMalformedFallsBackToMatchAll(t *testing.T) {
	// English text JSON、English text、English text——English text Filter，
	// English text「English text」。English text「English text」English text：
	// English text，English text。
	cases := []struct {
		name string
		raw  string
	}{
		{"English text", ""},
		{"English text JSON", `{not json`},
		{"English text JSON", `{"min_severity":`},
		{"English text", `{"min_severity": 123, "task_ids": "abc"}`},
		{"English text", `[1,2,3]`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := ParseFilter([]byte(tc.raw))
			if f.MinSeverity != "" || len(f.TaskIDs) != 0 || len(f.AssetIDs) != 0 {
				t.Fatalf("English text Filter，English text %+v", f)
			}
			// English text Filter English text。
			ev := Snapshot{Kind: EventFindingCreated, Severity: "low", VulnClass: "XSS"}
			if !Match(f, ev) {
				t.Fatal("English text Filter English text")
			}
		})
	}
}

func TestMatchSeverityThreshold(t *testing.T) {
	ev := func(sev string) Snapshot {
		return Snapshot{Kind: EventFindingCreated, Severity: sev}
	}
	cases := []struct {
		min    string
		sev    string
		expect bool
	}{
		{"", "low", true},
		{"", "critical", true},
		{"high", "critical", true},
		{"high", "high", true},
		{"high", "medium", false},
		{"high", "low", false},
		{"critical", "high", false},
		{"critical", "critical", true},
		// English text 0，English text（English text）。
		{"low", "", false},
		{"low", "unknown", false},
		{"", "", true},
	}
	for _, tc := range cases {
		got := Match(Filter{MinSeverity: tc.min}, ev(tc.sev))
		if got != tc.expect {
			t.Errorf("min=%q sev=%q: English text %v English text %v", tc.min, tc.sev, tc.expect, got)
		}
	}
}

func TestMatchScopeRestrictions(t *testing.T) {
	ev := Snapshot{
		Kind:      EventFindingCreated,
		Severity:  "high",
		TaskID:    7,
		AssetIDs:  []int64{10, 20},
		VulnClass: "SQLEnglish text",
	}
	cases := []struct {
		name   string
		filter Filter
		expect bool
	}{
		{"English text=English text", Filter{}, true},
		{"English text", Filter{TaskIDs: []int64{7}}, true},
		{"English text", Filter{TaskIDs: []int64{8}}, false},
		{"English text", Filter{TaskIDs: []int64{8, 7}}, true},
		{"English text", Filter{AssetIDs: []int64{20, 99}}, true},
		{"English text", Filter{AssetIDs: []int64{99}}, false},
		{"English text", Filter{TaskIDs: []int64{7}, AssetIDs: []int64{10}}, true},
		{"English text", Filter{TaskIDs: []int64{7}, AssetIDs: []int64{99}}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Match(tc.filter, ev); got != tc.expect {
				t.Errorf("English text %v English text %v", tc.expect, got)
			}
		})
	}
}

func TestMatchVulnClassKeywords(t *testing.T) {
	ev := func(class string) Snapshot {
		return Snapshot{Kind: EventFindingCreated, Severity: "high", VulnClass: class}
	}
	cases := []struct {
		name   string
		filter Filter
		class  string
		expect bool
	}{
		{"include English text=English text", Filter{}, "English text", true},
		{"include English text", Filter{VulnClassInclude: []string{"SQL"}}, "SQLEnglish text", true},
		{"include English text", Filter{VulnClassInclude: []string{"English text"}}, "SQLEnglish text", false},
		{"include English text", Filter{VulnClassInclude: []string{"English text", "SQL"}}, "SQLEnglish text", true},
		{"English text", Filter{VulnClassInclude: []string{"sql"}}, "SQLEnglish text", true},
		{"exclude English text", Filter{VulnClassExclude: []string{"English text"}}, "English text", false},
		{"exclude English text", Filter{VulnClassExclude: []string{"English text"}}, "SQLEnglish text", true},
		// English text：English text。
		{"English text", Filter{
			VulnClassInclude: []string{"SQL"},
			VulnClassExclude: []string{"English text"},
		}, "SQLEnglish text", false},
		// English text，English text「English text」。
		{"English text", Filter{VulnClassInclude: []string{"", "  "}}, "SQLEnglish text", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Match(tc.filter, ev(tc.class)); got != tc.expect {
				t.Errorf("English text %v English text %v", tc.expect, got)
			}
		})
	}
}

func TestMatchStatusChangeRequiresOptIn(t *testing.T) {
	ev := Snapshot{Kind: EventFindingStatusChanged, Severity: "critical", FromStatus: "pending", ToStatus: "fixed"}
	// English text：English text「English text」English text，English text。
	if Match(Filter{MinSeverity: "low"}, ev) {
		t.Fatal("English text")
	}
	if !Match(Filter{OnStatusChange: true}, ev) {
		t.Fatal("English text on_status_change English text")
	}
	// English text on_status_change English text。
	created := Snapshot{Kind: EventFindingCreated, Severity: "critical"}
	if !Match(Filter{MinSeverity: "low"}, created) {
		t.Fatal("English text on_status_change")
	}
}
