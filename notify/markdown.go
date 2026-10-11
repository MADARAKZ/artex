package notify

import (
	"fmt"
	"strings"
)

// English text「Markdown English text」English text（English text、English text）English text。
// English text JSON、Telegram English text HTML、English text HTML，English text。

// maxAssetsShown English text。English text，
// English text——English text 4 English text IM English text。
const maxAssetsShown = 3

// maxSummaryRunes English text。IM English text「English text」，
// English text，English text。
const maxSummaryRunes = 120

// markdownReservedBytes English text（English text + English text + English text）
// English text（English text）。English text，English text——
// English text，English text「English text、English text」English text。
const markdownReservedBytes = 320

// markdownEscape English text markdown English text。
//
// English text：English text、English text、English text、English text**English text**——
// English text（English text），English text url English text
// English text URL（English text）。English text，English text
//
//	English text SQL English text\n[English text：English text](http://attacker.tld)
//
// English text/English text**English text**；English text
// `![](http://attacker.tld/beacon)` English text，English text
// 「English text」English text IP。English text，English text
// English text。
//
// English text/English text/English text/English text/English text/English text
// English text。`\` English text，English text。
func markdownEscape(s string) string {
	replacer := strings.NewReplacer(
		`\`, `\\`,
		"`", "\\`",
		"*", `\*`,
		"_", `\_`,
		"[", `\[`,
		"]", `\]`,
		"(", `\(`,
		")", `\)`,
		"!", `\!`,
		"#", `\#`,
		">", `\>`,
		"|", `\|`,
		"~", `\~`,
	)
	return replacer.Replace(s)
}

// markdownText English text，English text markdown English text。
// English text：English text，
// English text。
func markdownText(s string, maxRunes int) string {
	return markdownEscape(OneLine(s, maxRunes))
}

// markdownTitle English text（IM English text/English text），English text**English text**。
//
// English text：English text——markdown English text、Telegram English text
// HTML、English text plain_text、English text Webhook English text JSON English text。English text
// English text（markdown English text HTML English text，English text JSON English text
// English text），English text，English text writeItem / feishuItemLines /
// telegramEscape。English text markdown English text，English text Telegram English text
// English text `\(1\)` English text。
func markdownTitle(m Message) string {
	if m.Batch {
		return fmt.Sprintf("English text · English text %d English text", len(m.Items))
	}
	if len(m.Items) == 0 {
		return "English text"
	}
	it := m.Items[0]
	return fmt.Sprintf("[%s] %s", SeverityLabel(it.Severity), OneLine(it.Title(), 0))
}

// markdownBody English text，English text**English text**。
//
// English text kept English text，English text kept English text
// English text——English text，English text
// English text。English text「English text」English text：English text，English text，
// English text。
//
// maxBytes<=0 English text。
func markdownBody(m Message, maxBytes int) (string, int) {
	if !m.Batch {
		if len(m.Items) == 0 {
			return "", 0
		}
		var b strings.Builder
		writeItem(&b, m.Items[0], "", true)
		// English text（English text）：English text
		// English text。
		return TruncateBytes(b.String(), maxBytes), 1
	}

	footer := ""
	if m.HomeURL != "" {
		footer = fmt.Sprintf("\n[English text](%s)\n", m.HomeURL)
	}
	kept := packItemCount(m.Items, maxBytes, markdownReservedBytes, footer, byteSize, func(it Item, idx int) string {
		var b strings.Builder
		writeItem(&b, it, fmt.Sprintf("%d. ", idx+1), false)
		return b.String()
	})

	items := m.Items[:kept]
	var b strings.Builder
	b.WriteString(markdownBatchIntro(m, items, len(m.Items)))
	for i, it := range items {
		writeItem(&b, it, fmt.Sprintf("%d. ", i+1), false)
	}
	b.WriteString(footer)
	return TruncateBytes(b.String(), maxBytes), kept
}

// markdownBatchIntro English text：English text、English text。
// English text，English text。
//
// items English text**English text**English text，total English text。English text
// 「English text」——English text，
// English text。
func markdownBatchIntro(m Message, items []Item, total int) string {
	var b strings.Builder
	if m.WindowMinutes > 0 {
		fmt.Fprintf(&b, "**English text %d English text %d English text**", m.WindowMinutes, total)
	} else {
		fmt.Fprintf(&b, "**English text %d English text**", total)
	}
	if extra := total - len(items); extra > 0 {
		fmt.Fprintf(&b, "（English text %d English text，English text %d English text）", len(items), extra)
	}
	// English text，English text。English text**English text**English text
	// English text，English text「English text 3」English text。
	counts := map[string]int{}
	for _, it := range items {
		counts[it.Severity]++
	}
	var parts []string
	for _, sev := range []string{"critical", "high", "medium", "low"} {
		if n := counts[sev]; n > 0 {
			parts = append(parts, fmt.Sprintf("%s %d", SeverityLabel(sev), n))
		}
	}
	if len(parts) > 0 {
		b.WriteString("\n" + strings.Join(parts, " · "))
	}
	b.WriteString("\n\n")
	return b.String()
}

// writeItem English text。
//
// prefix English text；single=true English text（English text），
// English text——English text 50 English text。
//
// English text（English text/English text/English text/English text）English text markdownText：
// English text + English text。English text public_base_url English text，English text，
// English text，English text。
func writeItem(b *strings.Builder, it Item, prefix string, single bool) {
	line := fmt.Sprintf("%s**%s · %s**", prefix, SeverityLabel(it.Severity), markdownText(it.Title(), 0))
	if !single {
		// English text：English text，English text。
		var extras []string
		if a := assetLine(it.Assets, maxAssetsShown); a != "" {
			extras = append(extras, markdownText(a, 0))
		}
		if it.Summary != "" {
			extras = append(extras, markdownText(it.Summary, 60))
		}
		if len(extras) > 0 {
			line += " — " + strings.Join(extras, " · ")
		}
		b.WriteString(line + "\n")
		return
	}
	b.WriteString(line + "\n")
	if it.IsStatusChange() {
		fmt.Fprintf(b, "**English text**：%s → %s\n",
			markdownText(StatusLabel(it.FromStatus), 0), markdownText(StatusLabel(it.ToStatus), 0))
	}
	if it.VulnClass != "" && it.VulnClass != it.Title() {
		fmt.Fprintf(b, "**English text**：%s\n", markdownText(it.VulnClass, 0))
	}
	if a := assetLine(it.Assets, maxAssetsShown); a != "" {
		fmt.Fprintf(b, "**English text**：%s\n", markdownText(a, 0))
	}
	if it.Summary != "" {
		if s := markdownText(it.Summary, maxSummaryRunes); s != "" {
			fmt.Fprintf(b, "**English text**：%s\n", s)
		}
	}
	if it.DetailURL != "" {
		fmt.Fprintf(b, "[English text](%s)\n", it.DetailURL)
	}
}
