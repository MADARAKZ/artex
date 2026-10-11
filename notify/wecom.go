package notify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// weComMarkdownLimit English text markdown content English text（English text，English text）。
// English text，English text TruncateBytes English text。
const weComMarkdownLimit = 4096

// weComChannel English text。
//
// English text：
//   - English text URL English text key English text，English text——English text webhook English text。
//   - markdown content English text 4096 **English text**，English text（English text）。English text 3 English text/English text，
//     English text，English text。
//   - English text 20 English text/English text，English text。
type weComChannel struct{}

func (weComChannel) Kind() string { return KindWeCom }

func (weComChannel) DefaultRatePerMin() int { return 20 }

// English text Webhook English text（URL English text key），English text——
// English text，English text。
func (weComChannel) SecretKeys() []string { return []string{"webhook"} }

// English text Webhook English text，English text，English text「English text」English text。
func (weComChannel) DestinationKeys() []string { return []string{"webhook"} }

func (weComChannel) Validate(cfg map[string]any) error {
	hook := cfgString(cfg, "webhook")
	if hook == "" {
		return errors.New("English text Webhook English text")
	}
	if err := validateHTTPURL(hook); err != nil {
		return fmt.Errorf("Webhook English text: %w", err)
	}
	return nil
}

func (c weComChannel) Send(ctx context.Context, cfg map[string]any, m Message) (int, error) {
	if err := c.Validate(cfg); err != nil {
		return 0, Permanent(err)
	}
	// English text（50 English text × English text + English text），4096 English text。
	// English text：English text，English text。
	content, kept := markdownBody(m, weComMarkdownLimit)
	payload := map[string]any{
		"msgtype":  "markdown",
		"markdown": map[string]any{"content": content},
	}
	raw, err := doJSON(ctx, "POST", cfgString(cfg, "webhook"), nil, payload)
	if err != nil {
		return 0, err
	}
	var res struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.Unmarshal(raw, &res); err != nil {
		return 0, fmt.Errorf("English text: %w (%s)", err, snippet(raw))
	}
	if res.ErrCode != 0 {
		// 45009 English text——English text，English text，
		// English text。English text rate_per_min English text，
		// English text，English text。
		if res.ErrCode == 45009 {
			return 0, fmt.Errorf("English text %d: %s", res.ErrCode, res.ErrMsg)
		}
		// 93000 English text webhook key English text——English text，English text。
		return 0, Permanent(fmt.Errorf("English text %d: %s", res.ErrCode, res.ErrMsg))
	}
	return kept, nil
}
