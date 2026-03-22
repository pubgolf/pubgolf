package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInfraError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "address already in use",
			err:  fmt.Errorf("listen tcp :5432: bind: address already in use: %w", errSimulated),
			want: true,
		},
		{
			name: "docker daemon unavailable",
			err:  fmt.Errorf("Cannot connect to the Docker daemon at unix:///var/run/docker.sock: %w", errSimulated),
			want: true,
		},
		{
			name: "connection refused",
			err:  fmt.Errorf("dial tcp 127.0.0.1:5432: connection refused: %w", errSimulated),
			want: true,
		},
		{
			name: "permission denied",
			err:  fmt.Errorf("open /etc/secrets: permission denied: %w", errSimulated),
			want: true,
		},
		{
			name: "no such host",
			err:  fmt.Errorf("dial tcp: lookup db.example.com: no such host: %w", errSimulated),
			want: true,
		},
		{
			name: "OOM killed container",
			err:  fmt.Errorf("container exited with code 137: %w", errSimulated),
			want: true,
		},
		{
			name: "case insensitive match",
			err:  fmt.Errorf("ADDRESS ALREADY IN USE: %w", errSimulated),
			want: true,
		},
		{
			name: "test failure is not infra",
			err:  fmt.Errorf("exit status 1: %w", errSimulated),
			want: false,
		},
		{
			name: "lint failure is not infra",
			err:  fmt.Errorf("run golangci-lint cmd: exit status 1: %w", errSimulated),
			want: false,
		},
		{
			name: "compilation error is not infra",
			err:  fmt.Errorf("cannot find package foo: %w", errSimulated),
			want: false,
		},
		{
			name: "generic error is not infra",
			err:  fmt.Errorf("something went wrong: %w", errSimulated),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, isInfraError(tt.err))
		})
	}
}
