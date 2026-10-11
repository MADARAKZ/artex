package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"syscall"
	"time"
)

// allowLocalTargets English text / English text。
//
// English text。English text IM English text，English text
// English text：English text、English text
// （169.254.169.254，English text）。English text，English text XSS/CSRF
// English text、English text JWT English text，English text
// ——doJSON English text 4xx/5xx English text 200 English text last_error，English text
// English text，English text。
//
// English text「English text SMTP English text」（127.0.0.1:25 English text postfix）English text，
// English text。English text：
// English text ARTEX_NOTIFY_ALLOW_LOCAL=1 English text。
//
// English text AllowLocalTargetsEnv English text——English text server English text
// English text 127.0.0.1 English text httptest English text，English text。
const AllowLocalTargetsEnv = "ARTEX_NOTIFY_ALLOW_LOCAL"

func allowLocalTargets() bool {
	v := strings.TrimSpace(os.Getenv(AllowLocalTargetsEnv))
	return v == "1" || strings.EqualFold(v, "true")
}

// isBlockedDialIP English text IP English text「English text」English text。
//
// English text、English text（English text 169.254.169.254）、English text。
// **English text** RFC1918 English text：English text Mattermost / SMTP English text，
// English text。English text——
// English text，English text。
func isBlockedDialIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	// IPv4-mapped IPv6（::ffff:127.0.0.1）English text IPv4 English text，English text。
	if v4 := ip.To4(); v4 != nil {
		ip = v4
	}
	return ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() ||
		ip.IsInterfaceLocalMulticast() || ip.IsUnspecified() || ip.IsMulticast()
}

// blockInternalDial English text http.Transport English text Control English text，English text**English text**
// English text。
//
// English text：English text。
// English text——DNS English text（English text IP、
// English text）English text（English text，English text
// English text）。
func blockInternalDial(_, address string, _ syscall.RawConn) error {
	if allowLocalTargets() {
		return nil
	}
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return fmt.Errorf("English text %q", host)
	}
	if isBlockedDialIP(ip) {
		return fmt.Errorf("English text/English text %s（English text，English text %s=1）", ip, AllowLocalTargetsEnv)
	}
	return nil
}

// notifyTransport English text Transport English text。
// English text Clone English text（English text、HTTP/2、English text、proxy English text），
// English text。
var notifyTransport = func() *http.Transport {
	t, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return &http.Transport{}
	}
	clone := t.Clone()
	clone.DialContext = (&net.Dialer{Timeout: 10 * time.Second, Control: blockInternalDial}).DialContext
	return clone
}()

// httpClient English text。
//
// English text**English text**English text（server English text GlobalProxy）：English text
// English text，English text，English text。
// IM English text。English text 15 English text——English text。
//
// English text：English text「English text endpoint」English text，English text
// English text；English text（English text access_token、English text key、Telegram English text
// bot token）**English text URL English text**，English text。English text
// English text（English text）English text。
var httpClient = &http.Client{
	Timeout:   15 * time.Second,
	Transport: notifyTransport,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return errors.New("English text")
		}
		if len(via) > 0 && req.URL.Host != via[0].URL.Host {
			return fmt.Errorf("English text（%s → %s）", via[0].URL.Host, req.URL.Host)
		}
		return nil
	},
}

// respBodyLimit English text。English text，English text
// English text。
const respBodyLimit = 8 << 10

// doJSON English text（English text）。
//
// payload English text nil English text body（English text GET English text body English text）。
// headers English text，English text Webhook English text。
//
// English text：English text 5xx/408/429 English text「English text」，
// English text 4xx English text「English text」——English text 403 English text 3 English text。
func doJSON(ctx context.Context, method, url string, headers map[string]string, payload any) ([]byte, error) {
	var body io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			// English text bug（English text），English text。
			return nil, Permanent(fmt.Errorf("English text: %w", err))
		}
		body = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		// URL English text——English text，English text。
		// English text err：url.Parse English text。
		return nil, Permanent(fmt.Errorf("English text: %s", redactRequestTarget(url)))
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		// English text、DNS English text、English text——English text，English text。
		//
		// English text。English text：http.Client.Do English text *url.Error，
		// English text Error() English text `Op "English textURL": English text`，English text**English text URL English text**
		// （English text access_token、English text key、English text hook id、Telegram /bot<token>/）。
		// English text，English text：notification_deliveries
		// English text last_error（English text）、English text（**English text**）、
		// English text、English text 502 English text。
		return nil, fmt.Errorf("English text: %s", redactTransportError(err))
	}
	defer resp.Body.Close()
	raw, readErr := io.ReadAll(io.LimitReader(resp.Body, respBodyLimit))
	if readErr != nil {
		return nil, fmt.Errorf("English text: %w", readErr)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return raw, nil
	}
	// 429（English text）English text 408（English text）English text；English text 4xx English text，English text。
	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusRequestTimeout {
		return nil, fmt.Errorf("English text (HTTP %d): %s", resp.StatusCode, snippet(raw))
	}
	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("English text (HTTP %d): %s", resp.StatusCode, snippet(raw))
	}
	return nil, Permanent(fmt.Errorf("English text (HTTP %d): %s", resp.StatusCode, snippet(raw)))
}

// snippet English text，English text。English text，
// English text last_error English text。
func snippet(raw []byte) string {
	return OneLine(string(raw), 200)
}

// redactRequestTarget English text「scheme://host/…」，English text。
//
// English text，English text**English text**：English text scheme English text host，
// English text。English text「English text」English text URL English text：
//
//	English text   English text query      /robot/send?access_token=xxx
//	English text   English text query      /cgi-bin/webhook/send?key=xxx
//	English text   English text**English text** /open-apis/bot/v2/hook/<hook_id>
//	Telegram English text**English text** /bot<token>/sendMessage
//
// English text「English text」English text，English text。
// English text host English text（DNS English text、English text、English text），
// English text。
//
// English text——English text。
func redactRequestTarget(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return "(English text)"
	}
	return u.Scheme + "://" + u.Host + "/…"
}

// redactTransportError English text，English text。
//
// *url.Error English text {Op, URL, Err}，Error() English text URL English text。
// English text Err English text，English text Error() ——English text，
// English text URL English text/English text，English text。
func redactTransportError(err error) string {
	var uerr *url.Error
	if errors.As(err, &uerr) {
		host := ""
		if u, parseErr := url.Parse(uerr.URL); parseErr == nil {
			host = u.Host
		}
		if uerr.Err != nil {
			return fmt.Sprintf("%s %s: %s", uerr.Op, host, uerr.Err)
		}
		return fmt.Sprintf("%s %s: English text", uerr.Op, host)
	}
	// English text *url.Error（English text）English text，English text。
	return redactURLsInText(err.Error())
}

// redactURLsInText English text http(s) English text。
//
// English text（English text、English text）。
// English text http/https English text，English text——English text。
func redactURLsInText(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); {
		rest := s[i:]
		if strings.HasPrefix(rest, "http://") || strings.HasPrefix(rest, "https://") {
			end := len(rest)
			if j := strings.IndexAny(rest, " \t\n\"'"); j >= 0 {
				end = j
			}
			b.WriteString(redactRequestTarget(rest[:end]))
			i += end
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}
