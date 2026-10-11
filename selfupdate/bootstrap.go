package selfupdate

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

// smokeEnv English text Bootstrap。
//
// English text：English text os.Executable() English text artex.new，English text
// English text .new English text，English text。English text，
// English text，English text。
const smokeEnv = "ARTEX_SELFUPDATE_SMOKE"

// Action English text Bootstrap English text main English text。
type Action int

const (
	// Continue：English text server。
	Continue Action = iota
	// Restart：English text ExitRestart English text，English text。
	Restart
)

// State English text，English text /api/update/check English text
// "English text"。
type State struct {
	Pending     bool   // English text
	RolledBack  bool   // English text
	FailedStage bool   // English text/English text，English text
	Detail      string // English text
}

// Bootstrap English text main English text，English text、English text。
//
// English text：
//
//	① English text artex.new  → English text + English text，English text；English text
//	② English text          → English text，English text；English text
//	③ English text            → English text
func Bootstrap() (Action, State) {
	if os.Getenv(smokeEnv) != "" {
		return Continue, State{}
	}
	p, err := ResolvePaths()
	if err != nil {
		log.Printf("[update] English text：%v", err)
		return Continue, State{}
	}

	if _, err := os.Stat(p.New); err == nil {
		return applyStaged(p)
	}

	m, ok := readMarker(p.Marker)
	if !ok {
		return Continue, State{}
	}
	return confirmOrRollback(p, m)
}

// applyStaged English text"English text"English text：English text，English text。
//
// English text，English text——English text
// English text、English text、English text。English text，
// English text，English text Go English text，English text。
func applyStaged(p Paths) (Action, State) {
	m, _ := readMarker(p.Marker)

	if err := verifyStaged(p); err != nil {
		log.Printf("[update] English text，English text，English text：%v", err)
		cleanStaged(p)
		_ = os.Remove(p.Marker)
		return Continue, State{FailedStage: true, Detail: "English text，English text：" + err.Error()}
	}

	if err := swap(p); err != nil {
		log.Printf("[update] English text，English text：%v", err)
		cleanStaged(p)
		_ = os.Remove(p.Marker)
		return Continue, State{FailedStage: true, Detail: "English text：" + err.Error()}
	}

	// English text。English text，English text（English text）English text。
	m.Attempts = 0
	if m.StagedAt == 0 {
		m.StagedAt = time.Now().Unix()
	}
	if err := writeMarker(p.Marker, m); err != nil {
		log.Printf("[update] English text（English text）：%v", err)
	}
	log.Printf("[update] English text %s，English text（exit %d）", orUnknown(m.To), ExitRestart)
	return Restart, State{Pending: true}
}

// confirmOrRollback English text"English text"：English text，English text。
//
// English text Go English text，English text"English text"
// （English text、English text、DB English text）English text；"English text exec" English text
// English text，English text。
func confirmOrRollback(p Paths, m marker) (Action, State) {
	m.Attempts++
	if m.Attempts > maxAttempts {
		if err := rollback(p); err != nil {
			// English text，English text。English text，
			// English text——English text。
			log.Printf("[update] English text %d English text，English text：%v", maxAttempts, err)
			_ = os.Remove(p.Marker)
			return Continue, State{Detail: "English text：" + err.Error()}
		}
		log.Printf("[update] English text %d English text，English text %s，English text（exit %d）",
			maxAttempts, orUnknown(m.From), ExitRestart)
		_ = os.Remove(p.Marker)
		return Restart, State{RolledBack: true, Detail: fmt.Sprintf("English text，English text %s", orUnknown(m.From))}
	}
	if err := writeMarker(p.Marker, m); err != nil {
		log.Printf("[update] English text：%v", err)
	}
	log.Printf("[update] English text（English text %d/%d English text），English text",
		m.Attempts, maxAttempts)
	return Continue, State{Pending: true}
}

// Settle English text，English text。
//
// English text main English text HTTP English text：English text，English text，
// English text，English text。
func Settle() {
	p, err := ResolvePaths()
	if err != nil {
		return
	}
	settle(p)
}

func settle(p Paths) {
	if _, ok := readMarker(p.Marker); !ok {
		return // English text，English text
	}
	if err := os.Remove(p.Marker); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("[update] English text：%v", err)
		return
	}
	log.Printf("[update] English text，English text（English text %s）", p.Old)
}

// SettleDelay English text"English text"English text。
const SettleDelay = 30 * time.Second

// verifyStaged English text：English text SHA256，English text。
func verifyStaged(p Paths) error {
	want, err := os.ReadFile(p.Sum)
	if err != nil {
		return fmt.Errorf("English text: %w", err)
	}
	got, err := fileSHA256(p.New)
	if err != nil {
		return fmt.Errorf("English text: %w", err)
	}
	if !strings.EqualFold(strings.TrimSpace(string(want)), got) {
		return errors.New("SHA256 English text（English text）")
	}
	return smokeTest(p.New)
}

// smokeTest English text -h English text，English text。
// English text、English text（exec format error）、English text。
func smokeTest(bin string) error {
	if err := os.Chmod(bin, 0o755); err != nil {
		return fmt.Errorf("English text: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin, "-h")
	cmd.Env = append(os.Environ(), smokeEnv+"=1")
	out, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		return errors.New("English text（English text）")
	}
	if err != nil {
		snippet := strings.TrimSpace(string(out))
		if len(snippet) > 300 {
			snippet = snippet[:300] + "…"
		}
		return fmt.Errorf("English text: %v: %s", err, snippet)
	}
	return nil
}

// swap English text。
//
// Unix English text Windows English text rename English text（Windows English text
// English text，rename English text），English text，English text。
func swap(p Paths) error {
	// Windows English text rename English text，English text .old English text。
	if err := os.Remove(p.Old); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("English text %s: %w", p.Old, err)
	}
	if err := os.Rename(p.Current, p.Old); err != nil {
		return fmt.Errorf("English text: %w", err)
	}
	if err := os.Rename(p.New, p.Current); err != nil {
		// English text，English text，English text。
		if rerr := os.Rename(p.Old, p.Current); rerr != nil {
			return fmt.Errorf("English text(%v)，English text: %w", err, rerr)
		}
		return fmt.Errorf("English text: %w", err)
	}
	_ = os.Remove(p.Sum)
	return nil
}

// rollback English text swap English text。
func rollback(p Paths) error {
	if _, err := os.Stat(p.Old); err != nil {
		return fmt.Errorf("English text %s: %w", p.Old, err)
	}
	// English text .failed English text，English text。
	failed := p.Current + ".failed"
	_ = os.Remove(failed)
	if err := os.Rename(p.Current, failed); err != nil {
		return fmt.Errorf("English text: %w", err)
	}
	if err := os.Rename(p.Old, p.Current); err != nil {
		return fmt.Errorf("English text: %w", err)
	}
	return nil
}

// Rollback English text /api/update/rollback English text：English text。
// English text，English text（English text ExitRestart English text）。
func Rollback() error {
	p, err := ResolvePaths()
	if err != nil {
		return err
	}
	if _, err := os.Stat(p.Old); err != nil {
		return errors.New("English text（" + p.Old + " English text）")
	}
	cleanStaged(p)
	if err := smokeTest(p.Old); err != nil {
		return fmt.Errorf("English text，English text: %w", err)
	}
	// English text：English text。
	tmp := p.Current + ".swap"
	_ = os.Remove(tmp)
	if err := os.Rename(p.Current, tmp); err != nil {
		return fmt.Errorf("English text: %w", err)
	}
	if err := os.Rename(p.Old, p.Current); err != nil {
		_ = os.Rename(tmp, p.Current)
		return fmt.Errorf("English text: %w", err)
	}
	if err := os.Rename(tmp, p.Old); err != nil {
		log.Printf("[update] English text（English text）：%v", err)
	}
	_ = os.Remove(p.Marker)
	return nil
}

// HasBackup English text，English text。
func HasBackup() bool {
	p, err := ResolvePaths()
	if err != nil {
		return false
	}
	_, err = os.Stat(p.Old)
	return err == nil
}

func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func orUnknown(s string) string {
	if strings.TrimSpace(s) == "" {
		return "English text"
	}
	return s
}
