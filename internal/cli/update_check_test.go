package cli

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
	// Swap URL to local server for deterministic output.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `[{"name":"v9.9.9"}]`)
	}))
	defer ts.Close()

	oldURL := defaultTagsURLForTests
	defaultTagsURLForTests = ts.URL
	defer func() { defaultTagsURLForTests = oldURL }()

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
