package homeassistant

import (
	"encoding/json"
	"testing"
)

func TestHomeAssistantStateStateWithUnit(t *testing.T) {
	state := &HomeAssistantState{
		State: "23.3",
		Attributes: map[string]json.RawMessage{
			"unit_of_measurement": json.RawMessage(`"°C"`),
		},
	}

	if got := state.StateWithUnit(); got != "23.3 °C" {
		t.Fatalf("StateWithUnit() = %q, want %q", got, "23.3 °C")
	}
}

func TestHomeAssistantStateStateWithUnitFallsBackToRawState(t *testing.T) {
	state := &HomeAssistantState{State: "on"}

	if got := state.StateWithUnit(); got != "on" {
		t.Fatalf("StateWithUnit() = %q, want %q", got, "on")
	}
}
