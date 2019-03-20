package parser

import (
	"testing"
	"time"
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
	}{
		{
			"cached passed line",
			args{
				"//admin/server:go_default_test                                  (cached) PASSED in 0.3s",
			},
			"//admin/server:go_default_test",
			true,
			Status_Passed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTarget, gotCached, gotStatus, gotDuration := ParseLine(tt.args.line)
			if gotTarget != tt.wantTarget {
				t.Errorf("ParseLine() gotTarget = %v, want %v", gotTarget, tt.wantTarget)
			}
			if gotCached != tt.wantCached {
				t.Errorf("ParseLine() gotCached = %v, want %v", gotCached, tt.wantCached)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("ParseLine() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
			if gotDuration != tt.wantDuration {
				t.Errorf("ParseLine() gotDuration = %v, want %v", gotDuration, tt.wantDuration)
			}
		})
	}
}
