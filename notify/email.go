package notify

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"mime"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

// emailDialTimeout / emailSessionTimeout English text SMTP English text。
// net/smtp English text，English text，English text
// English text goroutine English text——English text dispatcher English text goroutine English text，
// English text。
const (
	emailDialTimeout    = 10 * time.Second
	emailSessionTimeout = 45 * time.Second
)

// emailChannel English text SMTP English text。
type emailChannel struct{}

func (emailChannel) Kind() string { return KindEmail }

// English text，English text；English text。
func (emailChannel) DefaultRatePerMin() int { return 60 }

// English text。SMTP English text、English text、English text，English text。
func (emailChannel) SecretKeys() []string { return []string{"password"} }

// host/port English text；tls English text。English text
// English text——English text「English text TLS」English text，English text。
func (emailChannel) DestinationKeys() []string { return []string{"host", "port", "tls"} }

func (emailChannel) Validate(cfg map[string]any) error {
	if cfgString(cfg, "host") == "" {
		return errors.New("English text SMTP English text")
	}
	port := cfgInt(cfg, "port")
	if port <= 0 || port > 65535 {
		return errors.New("SMTP English text（English text 1-65535）")
	}
	if cfgString(cfg, "from") == "" {
		return errors.New("English text")
	}
	if len(cfgStrings(cfg, "to")) == 0 {
		return errors.New("English text")
	}
	return nil
}

func (c emailChannel) Send(ctx context.Context, cfg map[string]any, m Message) (int, error) {
	if err := c.Validate(cfg); err != nil {
		return 0, Permanent(err)
	}
	host := cfgString(cfg, "host")
	port := cfgInt(cfg, "port")
	from := cfgString(cfg, "from")
	to := cfgStrings(cfg, "to")
	username := cfgString(cfg, "username")
	password := cfgString(cfg, "password")
	implicitTLS := cfgBool(cfg, "tls")

	msg, err := buildEmailMessage(from, to, m)
	if err != nil {
		return 0, Permanent(err)
	}

	addr := net.JoinHostPort(host, strconv.Itoa(port))
	client, err := emailDial(ctx, addr, host, implicitTLS)
	if err != nil {
		return 0, err
	}
	defer client.Close()

	// STARTTLS：English text。English text（English text auth English text）。
	if !implicitTLS {
		if ok, _ := client.Extension("STARTTLS"); ok {
			if err := client.StartTLS(&tls.Config{ServerName: host, MinVersion: tls.VersionTLS12}); err != nil {
				return 0, fmt.Errorf("STARTTLS English text: %w", err)
			}
		}
	}
	if username != "" {
		if err := client.Auth(smtp.PlainAuth("", username, password, host)); err != nil {
			// smtp.PlainAuth English text（English text localhost）。
			// English text**English text**English text，English text，English text——
			// English text「unencrypted connection」English text。
			if strings.Contains(err.Error(), "unencrypted connection") {
				return 0, Permanent(fmt.Errorf("English text：English text。English text TLS，English text 465 English text(English text TLS)，English text「English text TLS」English text (%w)", err))
			}
			return 0, Permanent(fmt.Errorf("SMTP English text: %w", err))
		}
	}
	if err := client.Mail(from); err != nil {
		return 0, smtpStageError(fmt.Sprintf("English text %s English text", from), err)
	}
	for _, rcpt := range to {
		if err := client.Rcpt(rcpt); err != nil {
			return 0, smtpStageError(fmt.Sprintf("English text %s English text", rcpt), err)
		}
	}
	w, err := client.Data()
	if err != nil {
		return 0, fmt.Errorf("SMTP DATA English text: %w", err)
	}
	if _, err := w.Write([]byte(msg)); err != nil {
		return 0, fmt.Errorf("English text: %w", err)
	}
	if err := w.Close(); err != nil {
		return 0, fmt.Errorf("English text: %w", err)
	}
	// Quit English text「English text」English text，English text。
	_ = client.Quit()
	// English text（HTML English text），English text。
	return len(m.Items), nil
}

// emailDial English text SMTP English text。
//
// implicitTLS=true English text 465 English text「English text TLS」English text；false English text 25/587 English text
// STARTTLS。English text：English text 465 English text greeting English text。
//
// English text**English text**English text（English text），English text net/smtp English text Client English text
// English text，English text；English text deadline
// English text。English text。
// Control English text blockInternalDial English text HTTP English text。English text SMTP
// English text SSRF English text：host English text 169.254.169.254 English text 127.0.0.1 English text，
// English text smtp.NewClient English text、English text last_error
// English text，English text；「English text vs English text」English text
// English text。English text，English text DNS English text。
func emailDial(ctx context.Context, addr, host string, implicitTLS bool) (*smtp.Client, error) {
	d := &net.Dialer{Timeout: emailDialTimeout, Control: blockInternalDial}
	var conn net.Conn
	var err error
	if implicitTLS {
		conn, err = tls.DialWithDialer(d, "tcp", addr, &tls.Config{ServerName: host, MinVersion: tls.VersionTLS12})
	} else {
		conn, err = d.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return nil, fmt.Errorf("English text SMTP English text: %w", err)
	}
	_ = conn.SetDeadline(time.Now().Add(emailSessionTimeout))
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("SMTP English text: %w", err)
	}
	return client, nil
}

// smtpStageError English text SMTP English text「English text」English text「English text」。
//
// English text：SMTP English text 4xx English text 5xx English text——
//   - 4xx（450 English text、451 English text、452 English text）English text**English text**English text，
//     English text；English text，English text。
//   - 5xx（550 English text、553 English text）English text，English text。
//
// English text，English text**English text**English text
// English text failed——English text。
// English text；English text（English text，
// English text）。
func smtpStageError(what string, err error) error {
	code := smtpReplyCode(err.Error())
	if code >= 500 && code < 600 {
		return Permanent(fmt.Errorf("%s: %w", what, err))
	}
	return fmt.Errorf("%s: %w", what, err)
}

// smtpReplyCode English text SMTP English text，English text 0。
// net/smtp English text，English text；English text「450 4.7.1 ...」。
func smtpReplyCode(text string) int {
	if len(text) < 3 {
		return 0
	}
	n, err := strconv.Atoi(text[:3])
	if err != nil {
		return 0
	}
	return n
}

// buildEmailMessage English text RFC 5322 English text。
//
// English text base64 English text：English text SMTP English text 1000 English text，English text HTML
// English text（English text）English text；English text base64 English text "." English text
// English text，English text SMTP English text。
func buildEmailMessage(from string, to []string, m Message) (string, error) {
	var b strings.Builder
	fmt.Fprintf(&b, "From: %s\r\n", from)
	fmt.Fprintf(&b, "To: %s\r\n", strings.Join(to, ", "))
	// English text RFC 2047 English text，English text。
	fmt.Fprintf(&b, "Subject: %s\r\n", mime.QEncoding.Encode("utf-8", htmlTitle(m)))
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	b.WriteString("Content-Transfer-Encoding: base64\r\n")
	// English text，English text。
	b.WriteString("\r\n")
	encoded := base64.StdEncoding.EncodeToString([]byte(htmlBody(m, 0)))
	// base64 English text 76 English text，English text RFC 2045。
	for len(encoded) > 76 {
		b.WriteString(encoded[:76] + "\r\n")
		encoded = encoded[76:]
	}
	b.WriteString(encoded + "\r\n")
	return b.String(), nil
}
