package notify

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"
)

// dingTalkChannel English text。
//
// English text（English text）：
//   - English text 20 English text/English text，English text（HTTP English text 200），
//     English text，English text DefaultRatePerMin。
//   - English text：English text / English text / IP English text。English text
//     English text，English text（English text webhook）。
//   - English text/English text HTTP 200，English text body English text errcode English text——English text errcode
//     English text。
type dingTalkChannel struct{}

func (dingTalkChannel) Kind() string { return KindDingTalk }

func (dingTalkChannel) DefaultRatePerMin() int { return 20 }

// English text Webhook English text access_token，English text，English text。
func (dingTalkChannel) SecretKeys() []string { return []string{"webhook", "secret"} }

// English text Webhook English text；English text。
func (dingTalkChannel) DestinationKeys() []string { return []string{"webhook"} }

func (dingTalkChannel) Validate(cfg map[string]any) error {
	hook := cfgString(cfg, "webhook")
	if hook == "" {
		return errors.New("English text Webhook English text")
	}
	if err := validateHTTPURL(hook); err != nil {
		return fmt.Errorf("Webhook English text: %w", err)
	}
	return nil
}

// Send English text。English text ActionCard（English text），English text markdown。
func (c dingTalkChannel) Send(ctx context.Context, cfg map[string]any, m Message) (int, error) {
	hook := cfgString(cfg, "webhook")
	if err := c.Validate(cfg); err != nil {
		return 0, Permanent(err)
	}
	endpoint, err := dingTalkSignedURL(hook, cfgString(cfg, "secret"), time.Now())
	if err != nil {
		return 0, Permanent(err)
	}

	title := markdownTitle(m)
	// English text markdown English text，English text，English text。
	text, kept := markdownBody(m, 20000)

	var payload any
	if !m.Batch && len(m.Items) == 1 && m.Items[0].DetailURL != "" {
		payload = map[string]any{
			"msgtype": "actionCard",
			"actionCard": map[string]any{
				"title":          title,
				"text":           text,
				"btnOrientation": "0",
				"singleTitle":    "English text",
				"singleURL":      m.Items[0].DetailURL,
			},
		}
	} else {
		payload = map[string]any{
			"msgtype":  "markdown",
			"markdown": map[string]any{"title": title, "text": text},
		}
	}

	raw, err := doJSON(ctx, "POST", endpoint, nil, payload)
	if err != nil {
		return 0, err
	}
	// English text 200 English text。
	var res struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.Unmarshal(raw, &res); err != nil {
		return 0, fmt.Errorf("English text: %w (%s)", err, snippet(raw))
	}
	if res.ErrCode != 0 {
		// 301000 English text、310000 English text——English text，
		// English text。
		return 0, Permanent(fmt.Errorf("English text %d: %s", res.ErrCode, res.ErrMsg))
	}
	return kept, nil
}

// dingTalkSignedURL English text webhook English text timestamp English text sign English text。
//
// English text：English text = timestamp + "\n" + secret，HMAC-SHA256 English text**English text secret**，
// English text base64 English text URL English text。timestamp English text。secret English text，
// English text。
func dingTalkSignedURL(hook, secret string, now time.Time) (string, error) {
	if secret == "" {
		return hook, nil
	}
	ts := strconv.FormatInt(now.UnixMilli(), 10)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(ts + "\n" + secret))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	u, err := url.Parse(hook)
	if err != nil {
		// English text err：url.Parse English text（English text access_token）。
		return "", fmt.Errorf("English text Webhook English text: %s", redactRequestTarget(hook))
	}
	q := u.Query()
	q.Set("timestamp", ts)
	q.Set("sign", sign)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// validateHTTPURL English text、English text，English text IP English text。
//
// English text：
//
//  1. **English text**。url.Parse English text *url.Error，English text Error() English text
//     **English text**，English text（English text access_token、
//     English text key、Telegram English text bot token、English text hook id）。English text `return err`，
//     English text「English text」English text，English text 400 English text、
//     English text last_error、English text。
//
//  2. **English text IP English text**，English text（blockInternalDial English text
//     English text，English text DNS English text）。English text，
//     English text。
//
// English text：file:///gopher:// English text http.Client English text
// English text（English text scheme English text，English text）。
func validateHTTPURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("English text（%s）", redactRequestTarget(raw))
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("English text http/https，English text %q", u.Scheme)
	}
	if u.Host == "" {
		return errors.New("English text")
	}
	if ip := net.ParseIP(u.Hostname()); ip != nil && isBlockedDialIP(ip) && !allowLocalTargets() {
		return fmt.Errorf("English text/English text %s（English text，English text %s=1）", ip, AllowLocalTargetsEnv)
	}
	return nil
}
