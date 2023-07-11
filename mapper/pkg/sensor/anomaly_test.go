package temperaturesensor

import (
	"testing"

	assertT "github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert := assertT.New(t)

	testCases := []struct {
		name string
		args interface{}
		want Anomaly
	}{
		{
			name: "Anomaly_None",
			args: "None",
			want: None,
		},
		{
			name: "Anomaly_Spike",
			args: "Spike",
			want: Spike,
		},
		{
			name: "Anomaly_Drift",
			args: "Drift",
			want: Drift,
		},
		{
			name: "Anomaly_Noise",
			args: "Noise",
			want: Noise,
		},
		{
			name: "Anomaly_Flatline",
			args: "Flatline",
			want: Flatline,
		},
		{
			name: "Anomaly_Spike_As_Anomaly_Type",
			args: Spike,
			want: Spike,
		},
		{
			name: "Invalid_String",
			args: "Invalid",
			want: None,
		},
		{
			name: "Integer_Value",
			args: 123,
			want: None,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseAnomaly(tt.args)
			assert.Equal(tt.want, got, "they should be equal")
		})
	}
}
