package temperaturesensor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTemperatureSensor(t *testing.T) {
	sensor := newTemperatureSensor("id", -10, 50, 20)

	assert.Equal(t, "id", sensor.ID)
	assert.Equal(t, -10.0, sensor.MinTemp)
	assert.Equal(t, 50.0, sensor.MaxTemp)
	assert.Equal(t, 20.0, sensor.Temperature)
}

func TestRead(t *testing.T) {
	sensor := newTemperatureSensor("Id", -10, 50, 20)

	// The Read() method should panic with "not implemented" message
	assert.PanicsWithValue(t, "not implemented", func() { sensor.Read() }, "The code did not panic")

	// or you can use require instead of assert to stop the test if it fails
	require.PanicsWithValue(t, "not implemented", func() { sensor.Read() }, "The code did not panic")
}
