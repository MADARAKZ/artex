package selfupdate

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Repo English text。English text：English text
// English text，English text。
const Repo = "Autumn-27/artex"

// latestURL English text GitHub English text"English text"English text。English text prerelease English text draft。
const latestURL = "https://api.github.com/repos/" + Repo + "/releases/latest"

// allowedHosts English text。English text checkRedirect，
// English text——English text DNS English text / English text
// English text，English text SHA256SUMS English text。
var allowedHosts = map[string]bool{
	"api.github.com":                       true,
	"github.com":                           true,
	"objects.githubusercontent.com":        true, // release English text
	"release-assets.githubusercontent.com": true,
	"raw.githubusercontent.com":            true,
}

// Release English text GitHub Release English text。
type Release struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	Draft       bool      `json:"draft"`
	Prerelease  bool      `json:"prerelease"`
	PublishedAt time.Time `json:"published_at"`
	HTMLURL     string    `json:"html_url"`
	Assets      []Asset   `json:"assets"`
}

// Asset English text Release English text。
type Asset struct {
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
	Size int64  `json:"size"`
}

// NewClient English text GitHub English text HTTP English text。proxy English text。
//
// English text Transport：English text TLS English text，English text
// English text InsecureSkipVerify English text。
func NewClient(proxy string) *http.Client {
	tr := &http.Transport{
		ForceAttemptHTTP2:   true,
		TLSHandshakeTimeout: 15 * time.Second,
	}
	if p := strings.TrimSpace(proxy); p != "" {
		if pu, err := url.Parse(p); err == nil {
			tr.Proxy = http.ProxyURL(pu)
		}
	}
	return &http.Client{
		Transport: tr,
		Timeout:   30 * time.Minute, // English text，English text
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("English text")
			}
			return checkURL(req.URL)
		},
	}
}

// checkURL English text https + English text。
func checkURL(u *url.URL) error {
	if u.Scheme != "https" {
		return fmt.Errorf("English text HTTPS English text: %s", u.Scheme+"://"+u.Host)
	}
	if !allowedHosts[strings.ToLower(u.Hostname())] {
		return fmt.Errorf("English text GitHub English text: %s", u.Hostname())
	}
	return nil
}

// FetchLatest English text。
func FetchLatest(ctx context.Context, c *http.Client) (*Release, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, latestURL, nil)
	if err != nil {
		return nil, err
	}
	if err := checkURL(req.URL); err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "artex-selfupdate")

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("English text GitHub English text（English text）: %w", err)
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusForbidden, resp.StatusCode == http.StatusTooManyRequests:
		// English text GitHub API English text IP English text 60 English text，English text IP English text。
		return nil, fmt.Errorf("GitHub English text（English text 60 English text），English text")
	case resp.StatusCode == http.StatusNotFound:
		return nil, fmt.Errorf("English text %s English text", Repo)
	case resp.StatusCode != http.StatusOK:
		return nil, fmt.Errorf("GitHub English text %d", resp.StatusCode)
	}

	var rel Release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, fmt.Errorf("English text Release English text: %w", err)
	}
	if strings.TrimSpace(rel.TagName) == "" {
		return nil, fmt.Errorf("Release English text tag")
	}
	return &rel, nil
}

// AssetName English text，English text build.sh English text package_binary English text：
// artex-<English text>-<os>-<arch>.zip（English text v English text）。
func AssetName(tag, goos, goarch string) string {
	return fmt.Sprintf("artex-%s-%s-%s.zip", strings.TrimPrefix(tag, "v"), goos, goarch)
}

// FindAsset English text Release English text。
func (r *Release) FindAsset(name string) (Asset, bool) {
	for _, a := range r.Assets {
		if strings.EqualFold(a.Name, name) {
			return a, true
		}
	}
	return Asset{}, false
}
