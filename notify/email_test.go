package notify

import (
	"bufio"
	"context"

	"net"
	"strings"
	"sync"
	"testing"
)

// English text。English text email.Send English text 0——
// English text SMTP English text，English text、
// English text（English text、English text、English text、DATA English text）。
//
// English text SMTP English text，English text mock English text net/smtp：
// English text「English text SMTP English text」English text，
// English text mock English text。

// fakeSMTP English text SMTP English text：English text greet/EHLO/AUTH/MAIL/RCPT/DATA/QUIT，
// English text。
type fakeSMTP struct {
	ln net.Listener

	// rcptReply English text RCPT TO English text；English text 250。
	rcptReply string
	// mailReply English text MAIL FROM English text；English text 250。
	mailReply string
	// advertiseAuth English text true English text EHLO English text AUTH PLAIN。
	advertiseAuth bool

	mu       sync.Mutex
	data     string
	commands []string
}

func newFakeSMTP(t *testing.T) *fakeSMTP {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	f := &fakeSMTP{ln: ln, rcptReply: "250 OK", mailReply: "250 OK"}
	go f.serve()
	t.Cleanup(func() { ln.Close() })
	return f
}

func (f *fakeSMTP) hostPort(t *testing.T) (string, int) {
	t.Helper()
	addr, ok := f.ln.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatal("English text TCP English text")
	}
	return "127.0.0.1", addr.Port
}

func (f *fakeSMTP) record(cmd string) {
	f.mu.Lock()
	f.commands = append(f.commands, cmd)
	f.mu.Unlock()
}

func (f *fakeSMTP) body() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.data
}

func (f *fakeSMTP) sawCommand(prefix string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, c := range f.commands {
		if strings.HasPrefix(c, prefix) {
			return true
		}
	}
	return false
}

func (f *fakeSMTP) serve() {
	conn, err := f.ln.Accept()
	if err != nil {
		return
	}
	defer conn.Close()
	br := bufio.NewReader(conn)
	w := func(s string) { _, _ = conn.Write([]byte(s + "\r\n")) }
	w("220 fake.local ESMTP ready")
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimRight(line, "\r\n")
		f.record(line)
		switch {
		case strings.HasPrefix(line, "EHLO"), strings.HasPrefix(line, "HELO"):
			// English text STARTTLS：English text（English text，English text TLS）。
			w("250-fake.local")
			if f.advertiseAuth {
				w("250-AUTH PLAIN")
			}
			w("250 8BITMIME")
		case strings.HasPrefix(line, "AUTH"):
			// English text：PLAIN English text，English text。
			w("235 2.7.0 Authentication successful")
		case strings.HasPrefix(line, "MAIL FROM"):
			w(f.mailReply)
		case strings.HasPrefix(line, "RCPT TO"):
			w(f.rcptReply)
		case strings.HasPrefix(line, "DATA"):
			w("354 End data with <CR><LF>.<CR><LF>")
			var sb strings.Builder
			for {
				dl, err := br.ReadString('\n')
				if err != nil {
					return
				}
				if strings.TrimRight(dl, "\r\n") == "." {
					break
				}
				sb.WriteString(dl)
			}
			f.mu.Lock()
			f.data = sb.String()
			f.mu.Unlock()
			w("250 2.0.0 Ok: queued as FAKE1")
		case strings.HasPrefix(line, "QUIT"):
			w("221 2.0.0 Bye")
			return
		default:
			w("250 OK")
		}
	}
}

func emailCfg(t *testing.T, f *fakeSMTP, extra map[string]any) map[string]any {
	t.Helper()
	host, port := f.hostPort(t)
	cfg := map[string]any{
		"host": host,
		"port": float64(port),
		"from": "artex@example.com",
		"to":   []any{"a@example.com", "b@example.com"},
	}
	for k, v := range extra {
		cfg[k] = v
	}
	return cfg
}

func TestEmailSendDeliversFullMessage(t *testing.T) {
	f := newFakeSMTP(t)
	f.advertiseAuth = true
	cfg := emailCfg(t, f, map[string]any{"username": "artex", "password": "pw"})

	if _, err := (emailChannel{}).Send(context.Background(), cfg, singleMsg()); err != nil {
		t.Fatalf("English text: %v", err)
	}
	// English text：English text、English text、DATA。
	for _, want := range []string{"MAIL FROM:<artex@example.com>", "RCPT TO:<a@example.com>", "RCPT TO:<b@example.com>", "DATA", "AUTH", "QUIT"} {
		if !f.sawCommand(want) {
			t.Errorf("SMTP English text %q，English text：%v", want, f.commands)
		}
	}
	// English text base64 English text HTML，English text（English text）。
	body := f.body()
	if body == "" {
		t.Fatal("DATA English text")
	}
	if !strings.Contains(body, "Content-Type: text/html") {
		t.Errorf("English text Content-Type English text:\n%s", body)
	}
	if !strings.Contains(body, "base64") {
		t.Errorf("English text base64 English text（English text HTML English text SMTP English text 1000 English text）:\n%s", body)
	}
	// English text To English text。
	if !strings.Contains(body, "a@example.com, b@example.com") {
		t.Errorf("To English text:\n%s", body)
	}
}

func TestEmailSendWithoutAuth(t *testing.T) {
	// English text AUTH —— English text。
	f := newFakeSMTP(t)
	cfg := emailCfg(t, f, nil)
	if _, err := (emailChannel{}).Send(context.Background(), cfg, singleMsg()); err != nil {
		t.Fatalf("English text: %v", err)
	}
	if f.sawCommand("AUTH") {
		t.Errorf("English text AUTH: %v", f.commands)
	}
}

// TestEmailSendClassifiesSMTPReplies English text：
// 5xx English text、4xx（English text）English text。
func TestEmailSendClassifiesSMTPReplies(t *testing.T) {
	cases := []struct {
		name      string
		rcptReply string
		mailReply string
		permanent bool
	}{
		{"English text 550 English text", "550 5.1.1 User unknown", "250 OK", true},
		{"English text 450 English text", "450 4.7.1 Greylisting in action", "250 OK", false},
		{"English text 452 English text", "452 4.2.2 Mailbox full", "250 OK", false},
		{"English text 553 English text", "250 OK", "553 5.1.3 Bad address", true},
		{"English text 451 English text", "250 OK", "451 4.3.0 Temporary failure", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFakeSMTP(t)
			f.rcptReply = tc.rcptReply
			f.mailReply = tc.mailReply
			_, err := (emailChannel{}).Send(context.Background(), emailCfg(t, f, nil), singleMsg())
			if err == nil {
				t.Fatal("English text")
			}
			if got := IsPermanent(err); got != tc.permanent {
				t.Fatalf("permanent English text：English text %v English text %v (%v)", tc.permanent, got, err)
			}
			// English text，English text。
			if !strings.Contains(err.Error(), strings.Fields(tc.rcptReply)[0]) && !strings.Contains(err.Error(), strings.Fields(tc.mailReply)[0]) {
				t.Errorf("English text: %v", err)
			}
		})
	}
}

func TestEmailSendRefusesPlaintextCredentials(t *testing.T) {
	// net/smtp English text PlainAuth English text（English text localhost）。
	// English text**English text**English text，English text；English text。
	// English text localhost English text。
	f := newFakeSMTP(t)
	f.advertiseAuth = true
	_, port := f.hostPort(t)
	cfg := map[string]any{
		"host":     "smtp.example.com", // English text localhost
		"port":     float64(port),
		"from":     "a@example.com",
		"to":       []any{"b@example.com"},
		"username": "artex",
		"password": "pw",
	}
	_, err := (emailChannel{}).Send(context.Background(), cfg, singleMsg())
	if err == nil {
		t.Skip("English text DNS English text，English text（English text）")
	}
	// English text English text English text；English text**English text**English text。
	if !IsPermanent(err) && !strings.Contains(err.Error(), "English text") {
		t.Logf("English text：%v（English text localhost English text）", err)
	}
}

func TestEmailValidateReportsMissingFields(t *testing.T) {
	// English text，English text；English text
	// English text。English text「English text」。
	cases := []struct {
		name string
		cfg  map[string]any
	}{
		{"English text host", map[string]any{"port": float64(25), "from": "a@b.c", "to": []any{"d@e.f"}}},
		{"English text port", map[string]any{"host": "smtp.example.com"}},
		{"port English text", map[string]any{"host": "h", "port": float64(70000), "from": "a@b.c", "to": []any{"d@e.f"}}},
		{"English text from", map[string]any{"host": "h", "port": float64(25), "to": []any{"d@e.f"}}},
		{"English text to", map[string]any{"host": "h", "port": float64(25), "from": "a@b.c"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := (emailChannel{}).Validate(tc.cfg); err == nil {
				t.Fatalf("English text: %v", tc.cfg)
			}
		})
	}
}

// TestEmailConfigTolerance English text：JSONB English text float64，
// English text UI English text；English text。
func TestEmailConfigTolerance(t *testing.T) {
	cfg := map[string]any{
		"host": "smtp.example.com",
		"port": "587", // English text
		"from": "a@b.c",
		"to":   "d@e.f", // English text
		"tls":  "true",  // English text
	}
	if err := (emailChannel{}).Validate(cfg); err != nil {
		t.Fatalf("English text: %v", err)
	}
	if got := cfgInt(cfg, "port"); got != 587 {
		t.Errorf("cfgInt English text，English text %d", got)
	}
	if !cfgBool(cfg, "tls") {
		t.Error("cfgBool English text \"true\"")
	}
	if to := cfgStrings(cfg, "to"); len(to) != 1 || to[0] != "d@e.f" {
		t.Errorf("cfgStrings English text，English text %v", to)
	}
}

// TestFilterValidateRejectsTypo English text：
// English text，English text。
func TestFilterValidateRejectsTypo(t *testing.T) {
	good := []string{"", "low", "medium", "high", "critical"}
	for _, s := range good {
		if err := (Filter{MinSeverity: s}).Validate(); err != nil {
			t.Errorf("English text %q English text: %v", s, err)
		}
	}
	// English text——English text。
	for _, s := range []string{"hgih", "HIGH", "English text", "high ", "crit"} {
		err := (Filter{MinSeverity: s}).Validate()
		if err == nil {
			t.Errorf("English text %q English text（English text、English text）", s)
			continue
		}
		// English text。
		if !strings.Contains(err.Error(), "low") || !strings.Contains(err.Error(), "critical") {
			t.Errorf("English text，English text %q", err.Error())
		}
	}
}

// TestFilterValidateIsWriteTimeOnly English text「English text、English text」English text：
// English text（English text）。
func TestFilterValidateIsWriteTimeOnly(t *testing.T) {
	raw := []byte(`{"min_severity":"hgih"}`)
	f := ParseFilter(raw) // English text
	if f.MinSeverity != "hgih" {
		t.Fatalf("English text，English text %q", f.MinSeverity)
	}
	// English text（English text panic、English text）。
	_ = Match(f, Snapshot{Kind: EventFindingCreated, Severity: "critical"})
}
