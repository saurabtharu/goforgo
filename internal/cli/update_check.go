package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	updateCheckTimeout = 1200 * time.Millisecond
	updateInstallCmd   = "go install github.com/stonecharioteer/goforgo/cmd/goforgo@latest"
)

var defaultTagsURLForTests = "https://api.github.com/repos/stonecharioteer/goforgo/tags?per_page=1"

type githubTag struct {
	Name string `json:"name"`
}

type semVersion struct {
	Major int
	Minor int
	Patch int
}

func maybeNotifyUpdate(w io.Writer, currentVersion string) {
	maybeNotifyUpdateWithConfig(w, currentVersion, http.DefaultClient, defaultTagsURLForTests)
}

func maybeNotifyUpdateWithConfig(w io.Writer, currentVersion string, client *http.Client, tagsURL string) {
	if w == nil {
		return
	}

	latest, isNewer, err := checkForUpdate(currentVersion, tagsURL, client)
	if err != nil || !isNewer {
		return
	}

	fmt.Fprintf(w, "\n🔔 Update available: %s (current: %s)\n", latest, currentVersion)
	fmt.Fprintf(w, "   Update with: %s\n\n", updateInstallCmd)
}

func checkForUpdate(currentVersion, tagsURL string, client *http.Client) (latest string, isNewer bool, err error) {
	if client == nil {
		client = http.DefaultClient
	}

	current, ok := parseSemVersion(currentVersion)
	if !ok {
		return "", false, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), updateCheckTimeout)
	defer cancel()

	tag, err := fetchLatestTag(ctx, client, tagsURL)
	if err != nil {
		return "", false, err
	}

	latestParsed, ok := parseSemVersion(tag)
	if !ok {
		return "", false, nil
	}

	return tag, current.lessThan(latestParsed), nil
}

func fetchLatestTag(ctx context.Context, client *http.Client, tagsURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tagsURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "goforgo-update-check")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("tags request returned status %d", resp.StatusCode)
	}

	var tags []githubTag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return "", err
	}
	if len(tags) == 0 || tags[0].Name == "" {
		return "", fmt.Errorf("no tags returned")
	}

	return tags[0].Name, nil
}

func parseSemVersion(raw string) (semVersion, bool) {
	v := strings.TrimSpace(raw)
	v = strings.TrimPrefix(v, "v")
	if v == "" {
		return semVersion{}, false
	}

	// Drop any build metadata / prerelease suffix.
	if i := strings.IndexAny(v, "+-"); i >= 0 {
		v = v[:i]
	}

	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		return semVersion{}, false
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return semVersion{}, false
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return semVersion{}, false
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return semVersion{}, false
	}

	return semVersion{Major: major, Minor: minor, Patch: patch}, true
}

func (v semVersion) lessThan(other semVersion) bool {
	if v.Major != other.Major {
		return v.Major < other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor < other.Minor
	}
	return v.Patch < other.Patch
}
