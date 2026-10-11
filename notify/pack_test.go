package notify

import (
	"strings"
	"testing"
	"unicode/utf8"
)

// English text「English text」English text：English text，English text**English text**
// English text，English text。
//
// English text、English text：English text，
// English text——English text，English text。

func TestMarkdownBodyPacksWholeItemsWithinByteLimit(t *testing.T) {
	// 200 English text，English text 4096 English text。
	m := batchMsg(200)
	body, kept := markdownBody(m, weComMarkdownLimit)

	if len(body) > weComMarkdownLimit {
		t.Fatalf("English text %d English text %d", len(body), weComMarkdownLimit)
	}
	if !utf8.ValidString(body) {
		t.Fatal("English text UTF-8")
	}
	if kept <= 0 || kept >= len(m.Items) {
		t.Fatalf("English text（0 < kept < %d），English text %d", len(m.Items), kept)
	}
	// English text、English text——English text
	// English text。
	if !strings.Contains(body, "English text") || !strings.Contains(body, "English text") {
		t.Fatalf("English text:\n%s", body[:minInt(400, len(body))])
	}
	// English text kept English text。
	for i := 0; i < kept; i++ {
		if !strings.Contains(body, "English text"+itoa(i+1)) {
			t.Fatalf("English text %d English text:\n%s", i+1, body)
		}
	}
	if strings.Contains(body, "English text"+itoa(kept+1)) {
		t.Fatalf("English text %d English text（English text）", kept+1)
	}
}

func TestMarkdownBodyKeepsEverythingWhenUnderLimit(t *testing.T) {
	m := batchMsg(3)
	body, kept := markdownBody(m, 0) // 0 = English text
	if kept != len(m.Items) {
		t.Fatalf("English text，English text kept=%d", kept)
	}
	if strings.Contains(body, "English text") {
		t.Fatalf("English text:\n%s", body)
	}
}

func TestMarkdownBodyAlwaysKeepsAtLeastOneItem(t *testing.T) {
	// English text，English text（English text）。
	// English text：English text、English text。
	m := batchMsg(5)
	_, kept := markdownBody(m, 50)
	if kept != 1 {
		t.Fatalf("English text 1 English text，English text %d", kept)
	}
}

func TestMarkdownBodySingleReturnsOne(t *testing.T) {
	_, kept := markdownBody(singleMsg(), 4096)
	if kept != 1 {
		t.Fatalf("English text 1 English text，English text %d", kept)
	}
	// English text。
	if _, k := markdownBody(Message{}, 4096); k != 0 {
		t.Fatalf("English text 0 English text，English text %d", k)
	}
}

func TestTelegramPackingUsesRuneBudget(t *testing.T) {
	m := batchMsg(200)
	text, kept := telegramHTML(m)
	// Telegram English text**English text**English text；English text。
	if n := utf8.RuneCountInString(text); n > telegramTextLimit {
		t.Fatalf("English text %d English text %d", n, telegramTextLimit)
	}
	if kept <= 0 || kept >= len(m.Items) {
		t.Fatalf("English text，English text %d", kept)
	}
	if !strings.Contains(text, "English text") {
		t.Fatalf("English text:\n%.300s", text)
	}
}

func TestFeishuPackingReportsKept(t *testing.T) {
	m := batchMsg(2000)
	_, kept := feishuCard(m)
	if kept <= 0 || kept >= len(m.Items) {
		t.Fatalf("English text，English text %d", kept)
	}
}

func TestWebhookAndEmailReportAllItems(t *testing.T) {
	// English text，English text。
	m := batchMsg(7)
	if n := len(m.Items); n != 7 {
		t.Fatal("English text")
	}
	// English text：markdownBody(0) English text。
	if _, k := markdownBody(m, 0); k != len(m.Items) {
		t.Fatalf("English text，English text %d", k)
	}
}

// TestMarkdownEscapesUntrustedContent English text「English text」English text。
// English text（English text），English text URL。
func TestMarkdownEscapesUntrustedContent(t *testing.T) {
	cases := []struct {
		name  string
		item  Item
		must  []string // English text（English text）
		wrong []string // English text（English text）
	}{
		{
			name: "English text + English text",
			item: Item{
				Severity: "high",
				Name:     "English text SQL English text\n[English text：English text](http://attacker.tld)",
			},
			// English text（English text/English text）；
			// English text（English text）。
			must:  []string{`\[English text：English text\]`, `\(http://attacker.tld\)`},
			wrong: []string{"\n[English text", "\n\n[English text"},
		},
		{
			name: "English text",
			item: Item{
				Severity: "high",
				Name:     "English text ![](http://attacker.tld/beacon)",
			},
			must:  []string{`\!`, `\(http://attacker.tld/beacon\)`},
			wrong: []string{"![]("},
		},
		{
			name: "English text",
			item: Item{
				Severity: "high",
				Name:     "English text",
				Assets:   []string{"a.com/*English text*>English text"},
			},
			must:  []string{`\*English text\*`, `\>`},
			wrong: []string{"*English text*"},
		},
		{
			name: "English text",
			item: Item{
				Severity: "high",
				Name:     "English text",
				Summary:  "`code` | English text",
			},
			must:  []string{"\\`code\\`", `\|`},
			wrong: []string{"`code`"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := Message{Items: []Item{tc.item}}
			// English textItem English text markdown English text。
			var b strings.Builder
			writeItem(&b, tc.item, "", true)
			got := b.String()
			for _, want := range tc.must {
				if !strings.Contains(got, want) {
					t.Errorf("English text %q:\n%s", want, got)
				}
			}
			for _, bad := range tc.wrong {
				if strings.Contains(got, bad) {
					t.Errorf("English text %q（English text）:\n%s", bad, got)
				}
			}
			_ = m
		})
	}
}

// TestMarkdownEscapeBackslashFirst English text：English text，
// English text，English text。
func TestMarkdownEscapeBackslashFirst(t *testing.T) {
	if got := markdownEscape(`a\b*c`); got != `a\\b\*c` {
		t.Fatalf("English text，English text %q", got)
	}
}

// TestTelegramTitleHasNoMarkdownEscapes English text：
// markdown English text Telegram English text HTML English text（English text
// English text，English text Telegram English text `\(1\)` English text）。
func TestTelegramTitleHasNoMarkdownEscapes(t *testing.T) {
	m := Message{Items: []Item{{Severity: "high", Name: "alert(1) *English text*"}}}
	text, _ := telegramHTML(m)
	if strings.Contains(text, `\(`) || strings.Contains(text, `\*`) {
		t.Fatalf("Telegram English text markdown English text:\n%s", text)
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
