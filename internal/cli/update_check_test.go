package cli

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestParseSemVersion(t *testing.T) {
	tests := []struct {
		in   string
		ok   bool
		want semVersion
	}{
		{in: "0.9.1", ok: true, want: semVersion{0, 9, 1}},
		{in: "v0.9.1", ok: true, want: semVersion{0, 9, 1}},
		{in: "v1.2.3-beta", ok: true, want: semVersion{1, 2, 3}},
		{in: "dev", ok: false},
		{in: "1.2", ok: false},
	}

	for _, tt := range tests {
		got, ok := parseSemVersion(tt.in)
		if ok != tt.ok {
			t.Fatalf("parseSemVersion(%q) ok=%v want %v", tt.in, ok, tt.ok)
		}
		if ok && got != tt.want {
			t.Fatalf("parseSemVersion(%q) got=%+v want %+v", tt.in, got, tt.want)
		}
	}
}

func TestCheckForUpdate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `[{"name":"v0.9.2"}]`)
	}))
	defer ts.Close()

	tag, newer, err := checkForUpdate("0.9.1", ts.URL, ts.Client())
	if err != nil {
		t.Fatalf("checkForUpdate error: %v", err)
	}
	if !newer {
		t.Fatalf("expected newer=true")
	}
	if tag != "v0.9.2" {
		t.Fatalf("expected tag v0.9.2 got %s", tag)
	}
}

func TestMaybeNotifyUpdate(t *testing.T) {
	resetUpdateCheckGlobalsForTests(t)

	// Swap URL to local server for deterministic output.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `[{"name":"v9.9.9"}]`)
	}))
	defer ts.Close()

	oldURL := defaultTagsURLForTests
	defaultTagsURLForTests = ts.URL
	defer func() { defaultTagsURLForTests = oldURL }()

	updateCachePathForTests = t.TempDir() + "/update-check.json"

	var b strings.Builder
	maybeNotifyUpdateWithConfig(&b, "0.9.1", ts.Client(), defaultTagsURLForTests)

	out := b.String()
	if !strings.Contains(out, "Update available") {
		t.Fatalf("expected update message, got: %q", out)
	}
	if !strings.Contains(out, "go install github.com/stonecharioteer/goforgo/cmd/goforgo@latest") {
		t.Fatalf("expected install command in message, got: %q", out)
	}
}

func TestResolveUpdateStatus_UsesFreshCache(t *testing.T) {
	resetUpdateCheckGlobalsForTests(t)
	updateCachePathForTests = t.TempDir() + "/update-check.json"
	now := time.Date(2026, 4, 7, 12, 0, 0, 0, time.UTC)
	updateCheckNow = func() time.Time { return now }

	if err := saveUpdateCache(updateCheckCache{
		LastChecked: now,
		Current:     "0.9.3",
		Latest:      "v9.9.9",
		IsNewer:     true,
	}); err != nil {
		t.Fatalf("save cache: %v", err)
	}

	hit := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hit++
		_, _ = io.WriteString(w, `[{"name":"v0.9.4"}]`)
	}))
	defer ts.Close()

	latest, newer, err := resolveUpdateStatus("0.9.3", ts.URL, ts.Client())
	if err != nil {
		t.Fatalf("resolveUpdateStatus error: %v", err)
	}
	if !newer || latest != "v9.9.9" {
		t.Fatalf("expected cached newer tag, got latest=%q newer=%v", latest, newer)
	}
	if hit != 0 {
		t.Fatalf("expected no network call when cache is fresh, got %d", hit)
	}
}

func TestResolveUpdateStatus_RefreshesExpiredCache(t *testing.T) {
	resetUpdateCheckGlobalsForTests(t)
	updateCachePathForTests = t.TempDir() + "/update-check.json"
	now := time.Date(2026, 4, 7, 12, 0, 0, 0, time.UTC)
	updateCheckNow = func() time.Time { return now }

	if err := saveUpdateCache(updateCheckCache{
		LastChecked: now.Add(-48 * time.Hour),
		Current:     "0.9.3",
		Latest:      "v9.9.9",
		IsNewer:     true,
	}); err != nil {
		t.Fatalf("save cache: %v", err)
	}

	hit := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hit++
		_, _ = io.WriteString(w, `[{"name":"v0.9.4"}]`)
	}))
	defer ts.Close()

	latest, newer, err := resolveUpdateStatus("0.9.3", ts.URL, ts.Client())
	if err != nil {
		t.Fatalf("resolveUpdateStatus error: %v", err)
	}
	if !newer || latest != "v0.9.4" {
		t.Fatalf("expected refreshed tag, got latest=%q newer=%v", latest, newer)
	}
	if hit != 1 {
		t.Fatalf("expected one network call when cache is stale, got %d", hit)
	}
}

func resetUpdateCheckGlobalsForTests(t *testing.T) {
	t.Helper()
	updateCheckNow = time.Now
	updateCheckCacheTTL = 24 * time.Hour
	updateCachePathForTests = ""
	setUpdateNotice("")
}
