package notify

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

// English text：
//   ① English text / English text（SSRF）
//   ② English text
//
// English text：English text 127.0.0.1 English text httptest English text，English text
// English text。English text TestMain English text AllowLocalTargetsEnv，English text SSRF English text
// English text，English text**English text**English text。

func TestMain(m *testing.M) {
	// English text；SSRF English text。
	_ = os.Setenv(AllowLocalTargetsEnv, "1")
	os.Exit(m.Run())
}

// TestDialGuardRejectsLoopbackByDefault English text SSRF English text：
// English text，English text**English text**English text。
func TestDialGuardRejectsLoopbackByDefault(t *testing.T) {
	var hit bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hit = true
		_, _ = io.WriteString(w, `{"errcode":0}`)
	}))
	defer srv.Close()

	t.Setenv(AllowLocalTargetsEnv, "") // English text = English text
	_, err := (dingTalkChannel{}).Send(context.Background(),
		map[string]any{"webhook": srv.URL + "/robot/send"}, Message{Items: []Item{{Severity: "high"}}})
	if err == nil {
		t.Fatal("English text")
	}
	if hit {
		t.Fatal("English text——English text")
	}
	// English text（English text SMTP English text）。
	if !strings.Contains(err.Error(), AllowLocalTargetsEnv) {
		t.Errorf("English text: %v", err)
	}
}

// TestDialGuardAllowsLoopbackWhenOptedIn English text：English text，
// English text postfix / English text。
func TestDialGuardAllowsLoopbackWhenOptedIn(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"errcode":0,"errmsg":"ok"}`)
	}))
	defer srv.Close()

	t.Setenv(AllowLocalTargetsEnv, "1")
	if _, err := (dingTalkChannel{}).Send(context.Background(),
		map[string]any{"webhook": srv.URL + "/robot/send"}, Message{Items: []Item{{Severity: "high"}}}); err != nil {
		t.Fatalf("English text: %v", err)
	}
}

func TestIsBlockedDialIP(t *testing.T) {
	blocked := []string{
		"127.0.0.1", "127.1.2.3", "::1",
		"169.254.169.254", // English text——English text
		"169.254.1.1", "fe80::1",
		"0.0.0.0", "::",
		"224.0.0.1", "ff02::1",
		"::ffff:127.0.0.1", // IPv4-mapped English text，English text
		"",
	}
	for _, s := range blocked {
		if !isBlockedDialIP(net.ParseIP(s)) {
			t.Errorf("%s English text", s)
		}
	}
	// RFC1918 English text**English text**：English text Mattermost / SMTP English text。
	// English text——English text，English text，
	// English text（English text）。
	allowed := []string{"10.0.0.5", "172.16.3.4", "192.168.1.10", "8.8.8.8", "2606:4700::1111"}
	for _, s := range allowed {
		if isBlockedDialIP(net.ParseIP(s)) {
			t.Errorf("%s English text（English text）", s)
		}
	}
}

// TestValidateHTTPURLRejectsLiteralPrivateTargets English text：
// English text IP English text，English text。
func TestValidateHTTPURLRejectsLiteralPrivateTargets(t *testing.T) {
	t.Setenv(AllowLocalTargetsEnv, "")
	for _, raw := range []string{
		"http://127.0.0.1:8080/hook",
		"http://169.254.169.254/latest/meta-data/",
		"http://[::1]:8080/hook",
	} {
		if err := validateHTTPURL(raw); err == nil {
			t.Errorf("%s English text", raw)
		}
	}
	// English text（English text，English text）。
	for _, raw := range []string{"https://oapi.dingtalk.com/robot/send", "http://10.0.0.9/hook"} {
		if err := validateHTTPURL(raw); err != nil {
			t.Errorf("%s English text: %v", raw, err)
		}
	}
}

// TestValidateHTTPURLErrorNeverLeaksCredentials English text。
//
// url.Parse **English text**English text *url.Error，English text Error() English text。English text
// English text http.Client.Do English text，English text；English text「English text」English text
// （file://、gopher://、ftp://）English text url.Parse English text、English text scheme English text，
// English text——English text。
func TestValidateHTTPURLErrorNeverLeaksCredentials(t *testing.T) {
	cases := []string{
		"http://127.0.0.1/%zz?access_token=" + leakProbeToken,         // English text
		"https://a.example.com:port/x?access_token=" + leakProbeToken, // English text
		"http://[::1?access_token=" + leakProbeToken,                  // English text
	}
	for _, raw := range cases {
		// English text**English text**English text url.Parse English text。English text，English text
		// English text（English text）。
		if _, err := url.Parse(raw); err == nil {
			t.Errorf("%q English text，English text", raw)
			continue
		}
		err := validateHTTPURL(raw)
		if err == nil {
			t.Errorf("%q English text", raw)
			continue
		}
		assertNoSecret(t, err.Error(), leakProbeToken)
	}
	// English text。
	t.Setenv(AllowLocalTargetsEnv, "")
	err := (dingTalkChannel{}).Validate(map[string]any{"webhook": cases[0]})
	if err == nil {
		t.Fatal("English text")
	}
	assertNoSecret(t, err.Error(), leakProbeToken)
}

// TestEmailDialGuardRejectsLoopbackByDefault English text SMTP English text。
//
// English text net.Dialer，English text SSRF English text：host English text
// 169.254.169.254 English text 127.0.0.1 English text，English text smtp.NewClient English text
// English text、English text last_error English text——English text
// English text；「English text vs English text」English text。
//
// English text TestMain English text AllowLocalTargetsEnv（English text 127.0.0.1 English text
// English text），English text——English text，
// English text。
func TestEmailDialGuardRejectsLoopbackByDefault(t *testing.T) {
	f := newFakeSMTP(t)
	cfg := emailCfg(t, f, nil)

	t.Setenv(AllowLocalTargetsEnv, "") // English text = English text
	_, err := (emailChannel{}).Send(context.Background(), cfg, singleMsg())
	if err == nil {
		t.Fatal("English text")
	}
	// English text：English text Control English text，EHLO English text。
	if f.sawCommand("EHLO") || f.sawCommand("HELO") {
		t.Fatal("SMTP English text——English text")
	}
	// English text（English text postfix English text）。
	if !strings.Contains(err.Error(), AllowLocalTargetsEnv) {
		t.Errorf("English text: %v", err)
	}
}

// TestEmailDialGuardAllowsLoopbackWhenOptedIn English text：English text
// English text。English text SMTP / English text，English text。
func TestEmailDialGuardAllowsLoopbackWhenOptedIn(t *testing.T) {
	f := newFakeSMTP(t)
	cfg := emailCfg(t, f, nil)

	t.Setenv(AllowLocalTargetsEnv, "1")
	if _, err := (emailChannel{}).Send(context.Background(), cfg, singleMsg()); err != nil {
		t.Fatalf("English text SMTP English text: %v", err)
	}
	if !f.sawCommand("EHLO") {
		t.Fatal("English text EHLO——English text")
	}
}
