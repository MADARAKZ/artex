package notify

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestTruncateBytesKeepsValidUTF8(t *testing.T) {
	// English text。English text**English text**English text，English text 3 English text/English text，
	// English text，English text UTF-8 English text。
	// English text。
	inputs := []string{
		"English text",
		"English text mixed English text content",
		"aEnglish textbEnglish textcEnglish textdEnglish texte",
		"🔴🟠🟡🔵", // 4 English text emoji，English text
		strings.Repeat("English text", 100),
	}
	for _, in := range inputs {
		for max := 1; max <= len(in)+2; max++ {
			got := TruncateBytes(in, max)
			if !utf8.ValidString(got) {
				t.Fatalf("English text %q max=%d: English text UTF-8 %q", in, max, got)
			}
			if len(got) > max {
				t.Fatalf("English text %q max=%d: English text %d English text", in, max, len(got))
			}
			// English text。
			if len(in) <= max && got != in {
				t.Fatalf("English text %q max=%d: English text -> %q", in, max, got)
			}
		}
	}
}

func TestTruncateBytesZeroMeansUnlimited(t *testing.T) {
	long := strings.Repeat("x", 10000)
	if got := TruncateBytes(long, 0); got != long {
		t.Fatal("max=0 English text")
	}
	if got := TruncateBytes(long, -5); got != long {
		t.Fatal("max<0 English text")
	}
}

func TestTruncateBytesEllipsisBudget(t *testing.T) {
	// max English text，English text。
	got := TruncateBytes("abcdefgh", 1)
	if len(got) > 1 {
		t.Fatalf("max=1 English text %q English text %d English text", got, len(got))
	}
	// English text。
	if got := TruncateBytes("abcdefgh", 5); !strings.HasSuffix(got, ellipsis) {
		t.Fatalf("English text，English text %q", got)
	}
}

func TestTruncateRunesCountsCharactersNotBytes(t *testing.T) {
	// English text TruncateBytes English text：Telegram English text，
	// English text。
	s := "English text"
	got := TruncateRunes(s, 5)
	if n := utf8.RuneCountInString(got); n != 5 {
		t.Fatalf("English text 5 English text，English text %d English text (%q)", n, got)
	}
	// English text。
	if utf8.RuneCountInString(TruncateBytes(s, 5)) >= 5 {
		t.Fatal("English text")
	}
}

func TestOneLineCollapsesWhitespace(t *testing.T) {
	got := OneLine("English text\n\nEnglish text\tEnglish text   English text", 0)
	if strings.ContainsAny(got, "\n\t") {
		t.Fatalf("English text，English text %q", got)
	}
	if strings.Contains(got, "  ") {
		t.Fatalf("English text，English text %q", got)
	}
	// English text。
	got = OneLine("English text", 4)
	if n := utf8.RuneCountInString(got); n != 4 {
		t.Fatalf("English text 4 English text，English text %d (%q)", n, got)
	}
}

func TestTruncateHTMLNeverCutsTagInHalf(t *testing.T) {
	// English text HTML English text `<a href="htt` English text，English text。
	s := `<b>English text</b>English text<a href="https://example.com/very/long/path">English text</a>`
	for max := 1; max <= utf8.RuneCountInString(s)+2; max++ {
		got := TruncateHTML(s, max)
		if n := utf8.RuneCountInString(got); max > 0 && n > max {
			t.Fatalf("max=%d: English text %d English text", max, n)
		}
		// English text `<`（English text `<` English text `>`）。
		if lt := strings.LastIndex(got, "<"); lt >= 0 && !strings.Contains(got[lt:], ">") {
			t.Fatalf("max=%d: English text -> %q", max, got)
		}
	}
}

func TestAssetLineOmitsExcess(t *testing.T) {
	if got := assetLine(nil, 3); got != "" {
		t.Fatalf("English text，English text %q", got)
	}
	if got := assetLine([]string{"a", "b"}, 3); got != "a、b" {
		t.Fatalf("English text，English text %q", got)
	}
	// English text，English text。
	got := assetLine([]string{"a", "b", "c", "d", "e"}, 2)
	if !strings.Contains(got, "English text 5 English text") {
		t.Fatalf("English text 5，English text %q", got)
	}
}

func TestSeverityAndStatusLabels(t *testing.T) {
	if AtLeast("", "low") {
		t.Fatal("English text 0，English text")
	}
	if !AtLeast("critical", "") {
		t.Fatal("English text")
	}
	if got := StatusLabel("fixed"); got != "English text" {
		t.Fatalf("English text，English text %q", got)
	}
	// English text，English text。
	if got := StatusLabel("weird_status"); got != "weird_status" {
		t.Fatalf("English text，English text %q", got)
	}
}

// TestTruncateHTMLNeverCutsEntity English text：English text
// English text，English text HTML English text。
//
// `&amp;` English text `&amp` English text，English text**English text**English text——
// English text，English text。
func TestTruncateHTMLNeverCutsEntity(t *testing.T) {
	s := "aaaa&amp;bbbb&lt;cccc&quot;dddd"
	for max := 1; max <= utf8.RuneCountInString(s)+2; max++ {
		got := TruncateHTML(s, max)
		// English text「English text & English text ;」English text。
		if amp := strings.LastIndex(got, "&"); amp >= 0 && !strings.Contains(got[amp:], ";") {
			t.Fatalf("max=%d: English text %q", max, got[amp:])
		}
		if strings.Contains(got, "&amp\x00") {
			t.Fatalf("max=%d: English text", max)
		}
	}
}
