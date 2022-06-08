package utils

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestParseLogLevel(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		level        string
		want         zerolog.Level
		wantErr      bool
		errExcpected error
	}{
		"debug": {
			level:        "debug",
			want:         zerolog.DebugLevel,
			wantErr:      false,
			errExcpected: nil,
		},
		"info": {
			level:        "info",
			want:         zerolog.InfoLevel,
			wantErr:      false,
			errExcpected: nil,
		},
		"unknown": {
			level:        "unknown",
			want:         zerolog.InfoLevel,
			wantErr:      true,
			errExcpected: ErrLogLevel,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ParseLogLevel(tc.level)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Equal(t, err, tc.errExcpected)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
