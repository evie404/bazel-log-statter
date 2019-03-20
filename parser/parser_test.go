package parser

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name         string
		args         args
		wantTarget   string
		wantCached   bool
		wantStatus   Status
		wantDuration time.Duration
		wantErr      error
	}{
		{
			"cached passed line",
			args{
				"//admin/server:go_default_test                                  (cached) PASSED in 0.3s",
			},
			"//admin/server:go_default_test",
			true,
			StatusPassed,
			300 * time.Millisecond,
			nil,
		},
		{
			"no status line",
			args{
				"//summons/integration:go_default_test                                 NO STATUS",
			},
			"//summons/integration:go_default_test",
			false,
			StatusNoStatus,
			0,
			nil,
		},
		{
			"uncached line",
			args{
				"//social-graph/worker:go_default_test                                    PASSED in 53.8s",
			},
			"//social-graph/worker:go_default_test",
			false,
			StatusPassed,
			53800 * time.Millisecond,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTarget, gotCached, gotStatus, gotDuration, gotErr := ParseLine(tt.args.line)

			assert.Equal(t, tt.wantTarget, gotTarget)
			assert.Equal(t, tt.wantCached, gotCached)
			assert.Equal(t, string(tt.wantStatus), string(gotStatus))
			assert.Equal(t, tt.wantDuration, gotDuration)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}
