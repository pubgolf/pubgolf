package cmd

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

			assert.Equal(t, tt.want, normalizeSlug(tt.in))
		})
	}
}

func TestNormalizeSlug_TruncationMaxLength(t *testing.T) {
	t.Parallel()

	got := normalizeSlug("issue-1234-some-very-long-description-of-the-bug")

	// 20 chars + "-" + 6 hex chars = 27 max
	assert.LessOrEqual(t, len(got), 27, "slug should be at most 27 chars: %q", got)
	assert.Len(t, got, 27, "truncated slug should be exactly 27 chars: %q", got)
}

func TestNormalizeSlug_TruncationDeterministic(t *testing.T) {
	t.Parallel()

	input := "a-very-long-worktree-name-that-exceeds-the-limit"
	first := normalizeSlug(input)
	second := normalizeSlug(input)

	assert.Equal(t, first, second, "normalizeSlug should be deterministic")
}

func TestPortOffsetForSlug_MainTree(t *testing.T) {
	t.Parallel()

	got, err := portOffsetForSlug("")
	require.NoError(t, err)
	assert.Equal(t, 0, got)
}

func TestPortOffsetForSlug_Deterministic(t *testing.T) {
	t.Parallel()

	slug := "fix-auth"
	first, err := portOffsetForSlug(slug)
	require.NoError(t, err)

	second, err := portOffsetForSlug(slug)
	require.NoError(t, err)

	assert.Equal(t, first, second, "portOffsetForSlug should be deterministic")
}

func FuzzPortOffsetForSlug_Range(f *testing.F) {
	f.Add("fix-auth")
	f.Add("add-leaderboard")
	f.Add("x")
	f.Add("a-very-long-worktree-name")

	f.Fuzz(func(t *testing.T, slug string) {
		if slug == "" {
			t.Skip("empty slug returns 0, not in [1, 500]")
		}

		got, err := portOffsetForSlug(slug)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, got, 1, "portOffsetForSlug(%q) should be >= 1", slug)
		assert.LessOrEqual(t, got, 500, "portOffsetForSlug(%q) should be <= 500", slug)
	})
}

func FuzzPortOffsetForSlug_EnvOverride(f *testing.F) {
	f.Add(1)
	f.Add(42)
	f.Add(499)
	f.Add(250)

	f.Fuzz(func(t *testing.T, offset int) {
		if offset < 1 || offset > 499 {
			t.Skip("out of valid range")
		}

		t.Setenv("PUBGOLF_PORT_OFFSET", strconv.Itoa(offset))

		got, err := portOffsetForSlug("any-slug")
		require.NoError(t, err)
		assert.Equal(t, offset, got)
	})
}

func TestPortOffsetForSlug_EnvOverrideInvalid(t *testing.T) {
	tests := []struct {
		name string
		val  string
	}{
		{name: "zero", val: "0"},
		{name: "negative", val: "-1"},
		{name: "too high", val: "500"},
		{name: "way too high", val: "1000"},
		{name: "not a number", val: "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PUBGOLF_PORT_OFFSET", tt.val)

			_, err := portOffsetForSlug("fix-auth")
			require.ErrorIs(t, err, errInvalidPortOffset)
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

			assert.Equal(t, tt.want, dockerProjectForSlug(tt.slug))
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

			assert.Equal(t, tt.want, dataDirForSlug(tt.base, tt.slug))
		})
	}
}
