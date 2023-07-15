package simulation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnomaly_String(t *testing.T) {
	tests := []struct {
		name string
		a    Anomaly
		want string
	}{
		{
			name: "None",
			a:    None,
			want: "None",
		},
		{
			name: "Spike",
			a:    Spike,
			want: "Spike",
		},
		{
			name: "Drift",
			a:    Drift,
			want: "Drift",
		},
		{
			name: "Noise",
			a:    Noise,
			want: "Noise",
		},
		{
			name: "Flatline",
			a:    Flatline,
			want: "Flatline",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.a.String())
		})
	}
}

func TestParseAnomaly(t *testing.T) {
	tests := []struct {
		name string
		a    interface{}
		want Anomaly
	}{
		{
			name: "None",
			a:    "None",
			want: None,
		},
		{
			name: "Spike",
			a:    "Spike",
			want: Spike,
		},
		{
			name: "Drift",
			a:    "Drift",
			want: Drift,
		},
		{
			name: "Noise",
			a:    "Noise",
			want: Noise,
		},
		{
			name: "Flatline",
			a:    "Flatline",
			want: Flatline,
		},
		{
			name: "Invalid String",
			a:    "Invalid",
			want: None,
		},
		{
			name: "Anomaly Type",
			a:    Spike,
			want: Spike,
		},
		{
			name: "Invalid Type",
			a:    123,
			want: None,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ParseAnomaly(tt.a))
		})
	}
}
