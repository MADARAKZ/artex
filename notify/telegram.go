package notify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// telegramTextLimit English text Telegram sendMessage English text text English text（English text）。
const telegramTextLimit = 4096

// telegramChannel English text Telegram Bot API。
//
// English text：
//   - English text URL path English text（/bot<token>/sendMessage），English text。
//   - English text HTML English text MarkdownV2：MarkdownV2 English text `_*[]()~`>#+-=|{}.!`
//     English text 18 English text，English text；HTML English text & < > English text。
//   - English text HTTP 200 English text，English text ok English text。
type telegramChannel struct{}

func (telegramChannel) Kind() string { return KindTelegram }

// Telegram English text 1 English text/English text、English text 20 English text/English text。English text。
func (telegramChannel) DefaultRatePerMin() int { return 20 }

// Bot Token English text；chat_id English text，English text（English text Token English text）。
func (telegramChannel) SecretKeys() []string { return []string{"bot_token"} }

// base_url English text Token English text API English text（English text），English text Token。
func (telegramChannel) DestinationKeys() []string { return []string{"base_url"} }

func (telegramChannel) Validate(cfg map[string]any) error {
	if cfgString(cfg, "bot_token") == "" {
		return errors.New("English text Bot Token")
	}
	if cfgString(cfg, "chat_id") == "" {
		return errors.New("English text Chat ID")
	}
	if base := cfgString(cfg, "base_url"); base != "" {
		if err := validateHTTPURL(base); err != nil {
			return fmt.Errorf("API English text: %w", err)
		}
	}
	return nil
}

func (c telegramChannel) Send(ctx context.Context, cfg map[string]any, m Message) (int, error) {
	if err := c.Validate(cfg); err != nil {
		return 0, Permanent(err)
	}
	endpoint, err := telegramEndpoint(cfg)
	if err != nil {
		return 0, Permanent(err)
	}
	text, kept := telegramHTML(m)
	payload := map[string]any{
		"chat_id":                  cfgString(cfg, "chat_id"),
		"text":                     text,
		"parse_mode":               "HTML",
		"disable_web_page_preview": false,
	}
	raw, err := doJSON(ctx, "POST", endpoint, nil, payload)
	if err != nil {
		return 0, err
	}
	var res struct {
		OK          bool   `json:"ok"`
		ErrorCode   int    `json:"error_code"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(raw, &res); err != nil {
		return 0, fmt.Errorf("English text Telegram English text: %w (%s)", err, snippet(raw))
	}
	if res.OK {
		return kept, nil
	}
	// 429 English text，English text；English text（400 English text、401 token English text、403 English text、
	// 404 chat English text）English text，English text。
	if res.ErrorCode == 429 {
		return 0, fmt.Errorf("Telegram English text: %s", res.Description)
	}
	return 0, Permanent(fmt.Errorf("Telegram English text %d: %s", res.ErrorCode, res.Description))
}

// telegramEndpoint English text sendMessage English text。base_url English text API，
// English text Bot API English text（English text）。
func telegramEndpoint(cfg map[string]any) (string, error) {
	base := cfgString(cfg, "base_url")
	if base == "" {
		base = "https://api.telegram.org"
	}
	base = strings.TrimSuffix(base, "/")
	token := cfgString(cfg, "bot_token")
	raw := base + "/bot" + token + "/sendMessage"
	u, err := url.Parse(raw)
	if err != nil {
		// English text err：English text Bot Token，English text addr English text。
		return "", fmt.Errorf("English text API English text（API English text：%s）", redactRequestTarget(base))
	}
	return u.String(), nil
}

// telegramHTML English text HTML English text，English text（English text Channel.Send）。
func telegramHTML(m Message) (string, int) {
	var b strings.Builder
	b.WriteString("<b>" + telegramEscape(markdownTitle(m)) + "</b>\n")
	if m.Batch {
		// Telegram English text**English text**，English text（runeSize）。
		footer := ""
		if m.HomeURL != "" {
			footer = fmt.Sprintf("\n\n<a href=\"%s\">English text</a>", telegramEscapeAttr(m.HomeURL))
		}
		kept := packItemCount(m.Items, telegramTextLimit, telegramReservedRunes, footer, runeSize, func(it Item, idx int) string {
			return telegramBatchLine(it, idx+1)
		})
		items := m.Items[:kept]
		b.Reset()
		b.WriteString("<b>" + telegramEscape(telegramBatchTitle(m, items, len(m.Items))) + "</b>")
		for i, it := range items {
			b.WriteString("\n" + telegramEscape(telegramBatchLine(it, i+1)))
		}
		b.WriteString(footer)
		return TruncateHTML(b.String(), telegramTextLimit), kept
	}
	if len(m.Items) == 0 {
		return b.String(), 0
	}
	it := m.Items[0]
	if it.IsStatusChange() {
		b.WriteString(fmt.Sprintf("\n<b>English text</b>：%s → %s",
			telegramEscape(StatusLabel(it.FromStatus)), telegramEscape(StatusLabel(it.ToStatus))))
	}
	if it.VulnClass != "" && it.VulnClass != it.Title() {
		b.WriteString("\n<b>English text</b>：" + telegramEscape(it.VulnClass))
	}
	if a := assetLine(it.Assets, maxAssetsShown); a != "" {
		b.WriteString("\n<b>English text</b>：" + telegramEscape(a))
	}
	if s := OneLine(it.Summary, maxSummaryRunes); s != "" {
		b.WriteString("\n<b>English text</b>：" + telegramEscape(s))
	}
	if it.DetailURL != "" {
		b.WriteString(fmt.Sprintf("\n\n<a href=\"%s\">English text</a>", telegramEscapeAttr(it.DetailURL)))
	}
	return TruncateHTML(b.String(), telegramTextLimit), 1
}

// telegramReservedRunes English text（English text）。
const telegramReservedRunes = 160

// telegramBatchLine English text（English text，English text）。
func telegramBatchLine(it Item, idx int) string {
	if a := assetLine(it.Assets, maxAssetsShown); a != "" {
		return fmt.Sprintf("%d. %s · %s — %s", idx, SeverityLabel(it.Severity), it.Title(), a)
	}
	return fmt.Sprintf("%d. %s · %s", idx, SeverityLabel(it.Severity), it.Title())
}

// telegramBatchTitle English text。English text**English text**English text，
// English text——English text。
func telegramBatchTitle(m Message, items []Item, total int) string {
	title := fmt.Sprintf("English text · English text %d English text", total)
	if extra := total - len(items); extra > 0 {
		title += fmt.Sprintf("（English text %d English text，English text %d English text）", len(items), extra)
	}
	if m.WindowMinutes > 0 {
		title = fmt.Sprintf("English text %d English text · %s", m.WindowMinutes, title)
	}
	return title
}

// telegramEscape English text HTML English text。
// Telegram English text，English text &amp; English text——English text
// English text：English text，English text HTML。
func telegramEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

// telegramEscapeAttr English text HTML English text。English text——
// URL English text href English text，English text。
func telegramEscapeAttr(s string) string {
	s = telegramEscape(s)
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
