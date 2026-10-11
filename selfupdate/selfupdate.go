// Package selfupdate implements ARTEX English text：English text GitHub Release English text
// English text、English text、English text，English text。
//
// English text（English text start.sh / start.bat）：
//
//	English text  = English text，English text"English text"
//	English text      = English text（English text / SHA256 English text / English text / English text / English text）
//
// English text Go English text，English text sha256 English text sh English text bat English text
// English text（sha256sum / shasum / certutil），English text——English text
// English text，English text，English text。
//
// English text：
//
//	① English text server English text /api/update/apply → English text → English text artex.new → exit 75
//	② English text → Bootstrap English text artex.new → English text+English text → English text → exit 75
//	③ English text，English text → Bootstrap English text → English text
//
// English text：② English text；③ English text 3 English text
// English text（English text）English text artex.old English text。
package selfupdate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ExitRestart English text"English text"English text（EX_TEMPFAIL）。English text
// English text，English text。0 English text（English text），English text。
const ExitRestart = 75

// maxAttempts English text。English text +1，English text
// settleDelay English text；English text maxAttempts English text，English text。
const maxAttempts = 3

// Paths English text，English text**English text**English text。
// English text CWD：English text / English text，English text CWD English text
// English text，English text。
type Paths struct {
	Dir     string // English text
	Current string // English text        artex      / artex.exe
	New     string // English text            artex.new  / artex.new.exe
	Sum     string // English text sha256（hex）  artex.new.sha256 / artex.new.exe.sha256
	Old     string // English text      artex.old  / artex.old.exe
	Marker  string // English text            artex.upgrade.json
}

// ResolvePaths English text。
//
// Windows English text .new/.old English text .exe English text，English text，
// English text，English text。
func ResolvePaths() (Paths, error) {
	exe, err := os.Executable()
	if err != nil {
		return Paths{}, fmt.Errorf("English text: %w", err)
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		exe = resolved
	}
	dir := filepath.Dir(exe)
	name := filepath.Base(exe)
	ext := filepath.Ext(name) // Windows English text ".exe"，Unix English text
	stem := strings.TrimSuffix(name, ext)

	join := func(suffix string) string { return filepath.Join(dir, stem+suffix+ext) }
	return Paths{
		Dir:     dir,
		Current: exe,
		New:     join(".new"),
		Sum:     join(".new") + ".sha256",
		Old:     join(".old"),
		Marker:  filepath.Join(dir, stem+".upgrade.json"),
	}, nil
}

// marker English text，English text。
type marker struct {
	From     string `json:"from"`     // English text
	To       string `json:"to"`       // English text
	Attempts int    `json:"attempts"` // English text
	StagedAt int64  `json:"staged_at"`
}

func readMarker(path string) (marker, bool) {
	b, err := os.ReadFile(path)
	if err != nil {
		return marker{}, false
	}
	var m marker
	if json.Unmarshal(b, &m) != nil {
		return marker{}, false
	}
	return m, true
}

func writeMarker(path string, m marker) error {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

// cleanStaged English text。English text、English text、English text，English text
// artex.new English text。
func cleanStaged(p Paths) {
	_ = os.Remove(p.New)
	_ = os.Remove(p.Sum)
}

// CompareVersions English text，English text -1/0/1（a<b / a==b / a>b）。
// ok=false English text（English text "dev" English text
// git describe English text "0.3.7-2-gabc1234-dirty"），English text，
// English text"English text"English text、English text。
func CompareVersions(a, b string) (int, bool) {
	av, aok := parseVersion(a)
	bv, bok := parseVersion(b)
	if !aok || !bok {
		return 0, false
	}
	for i := range 3 {
		if av[i] != bv[i] {
			if av[i] < bv[i] {
				return -1, true
			}
			return 1, true
		}
	}
	return 0, true
}

// parseVersion English text "v0.3.7" / "0.3.7" English text [3]int。
//
// English text：build.sh English text tag English text git describe English text
// "0.3.7-2-gabc1234" English text，English text，English text
// 0.3.7 —— English text"English text"English text。
func parseVersion(s string) ([3]int, bool) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "v")
	if s == "" {
		return [3]int{}, false
	}
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return [3]int{}, false
	}
	var out [3]int
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return [3]int{}, false
		}
		out[i] = n
	}
	return out, true
}

// InDocker English text。Docker English text，
// `docker compose up -d` English text——English text
// （English text），English text。
func InDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	b, err := os.ReadFile("/proc/1/cgroup")
	if err != nil {
		return false
	}
	s := string(b)
	return strings.Contains(s, "docker") || strings.Contains(s, "containerd")
}
