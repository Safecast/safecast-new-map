package database

import "testing"

func TestRedactDSN(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "password hidden",
			in:   "postgres://readonly_user:secret@127.0.0.1:5433/safecast?sslmode=disable",
			want: "postgres://readonly_user:redacted@127.0.0.1:5433/safecast?sslmode=disable",
		},
		{
			name: "no password",
			in:   "postgres://readonly_user@127.0.0.1:5433/safecast?sslmode=disable",
			want: "postgres://readonly_user@127.0.0.1:5433/safecast?sslmode=disable",
		},
		{
			name: "file path",
			in:   "database-8765.sqlite",
			want: "database-8765.sqlite",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := redactDSN(tt.in); got != tt.want {
				t.Errorf("redactDSN(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
