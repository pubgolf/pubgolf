package cmd

import (
	"hash/fnv"
	"testing"
)

func TestNormalizeSlug(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "simple", in: "fix-auth", want: "fix-auth"},
		{name: "uppercase", in: "UPPER-Case", want: "upper-case"},
		{name: "special chars", in: "special!@#chars", want: "special-chars"},
		{name: "consecutive hyphens", in: "a--b--c", want: "a-b-c"},
		{name: "leading trailing hyphens", in: "-leading-trailing-", want: "leading-trailing"},
		{name: "mixed special", in: "hello___world...test", want: "hello-world-test"},
		{name: "single char", in: "x", want: "x"},
		{name: "all special", in: "!@#$%", want: ""},
		{name: "exactly 20 chars", in: "abcdefghijklmnopqrst", want: "abcdefghijklmnopqrst"},
		{
			name: "truncation with hash",
			in:   "issue-1234-some-very-long-description-of-the-bug",
			want: "issue-1234-some-very-f0c825",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := normalizeSlug(tt.in)
			if got != tt.want {
				t.Errorf("normalizeSlug(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestNormalizeSlug_TruncationMaxLength(t *testing.T) {
	t.Parallel()

	got := normalizeSlug("issue-1234-some-very-long-description-of-the-bug")

	// 20 chars + "-" + 6 hex chars = 27 max
	if len(got) > 27 {
		t.Errorf("normalizeSlug produced %d chars, want <= 27: %q", len(got), got)
	}

	if len(got) != 27 {
		t.Errorf("normalizeSlug produced %d chars, want exactly 27 for truncated input: %q", len(got), got)
	}
}

func TestNormalizeSlug_TruncationDeterministic(t *testing.T) {
	t.Parallel()

	input := "a-very-long-worktree-name-that-exceeds-the-limit"
	first := normalizeSlug(input)

	second := normalizeSlug(input)
	if first != second {
		t.Errorf("normalizeSlug is not deterministic: %q != %q", first, second)
	}
}

func TestPortOffsetForSlug_MainTree(t *testing.T) {
	t.Parallel()

	got := portOffsetForSlug("")
	if got != 0 {
		t.Errorf("portOffsetForSlug(\"\") = %d, want 0", got)
	}
}

func TestPortOffsetForSlug_Deterministic(t *testing.T) {
	t.Parallel()

	slug := "fix-auth"
	first := portOffsetForSlug(slug)

	second := portOffsetForSlug(slug)
	if first != second {
		t.Errorf("portOffsetForSlug not deterministic: %d != %d", first, second)
	}
}

func TestPortOffsetForSlug_Range(t *testing.T) {
	t.Parallel()

	slugs := []string{"fix-auth", "add-leaderboard", "refactor-db", "test-worktree", "my-feature"}
	for _, slug := range slugs {
		got := portOffsetForSlug(slug)
		if got < 1 || got > 500 {
			t.Errorf("portOffsetForSlug(%q) = %d, want in [1, 500]", slug, got)
		}
	}
}

func TestPortOffsetForSlug_EnvOverride(t *testing.T) {
	t.Setenv("PUBGOLF_PORT_OFFSET", "42")

	got := portOffsetForSlug("any-slug")
	if got != 42 {
		t.Errorf("portOffsetForSlug with PUBGOLF_PORT_OFFSET=42 = %d, want 42", got)
	}
}

func TestPortOffsetForSlug_EnvOverrideInvalid(t *testing.T) {
	slug := "fix-auth"

	// Get the expected hash-based offset.
	expected := func() int {
		h := fnv.New32a()
		h.Write([]byte(slug))

		return int(h.Sum32()%500) + 1
	}()

	tests := []struct {
		name string
		val  string
	}{
		{name: "zero", val: "0"},
		{name: "negative", val: "-1"},
		{name: "too high", val: "500"},
		{name: "way too high", val: "1000"},
		{name: "not a number", val: "abc"},
		{name: "empty", val: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.val != "" {
				t.Setenv("PUBGOLF_PORT_OFFSET", tt.val)
			}

			got := portOffsetForSlug(slug)
			if got != expected {
				t.Errorf("portOffsetForSlug(%q) with PUBGOLF_PORT_OFFSET=%q = %d, want %d", slug, tt.val, got, expected)
			}
		})
	}
}

func TestDockerProjectForSlug(t *testing.T) {
	t.Parallel()

	tests := []struct {
		slug string
		want string
	}{
		{slug: "", want: "pubgolf"},
		{slug: "fix-auth", want: "pubgolf-fix-auth"},
		{slug: "add-leaderboard", want: "pubgolf-add-leaderboard"},
	}

	for _, tt := range tests {
		t.Run(tt.slug, func(t *testing.T) {
			t.Parallel()

			got := dockerProjectForSlug(tt.slug)
			if got != tt.want {
				t.Errorf("dockerProjectForSlug(%q) = %q, want %q", tt.slug, got, tt.want)
			}
		})
	}
}

func TestDataDirForSlug(t *testing.T) {
	t.Parallel()

	tests := []struct {
		base string
		slug string
		want string
	}{
		{base: "data/postgres", slug: "", want: "data/postgres"},
		{base: "data/postgres", slug: "fix-auth", want: "data/postgres-fix-auth"},
		{base: "data/go-test-coverage", slug: "fix-auth", want: "data/go-test-coverage-fix-auth"},
	}

	for _, tt := range tests {
		t.Run(tt.base+"/"+tt.slug, func(t *testing.T) {
			t.Parallel()

			got := dataDirForSlug(tt.base, tt.slug)
			if got != tt.want {
				t.Errorf("dataDirForSlug(%q, %q) = %q, want %q", tt.base, tt.slug, got, tt.want)
			}
		})
	}
}
