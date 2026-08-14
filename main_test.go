package main

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/timmo001/go-automate/homeassistant"
)

func TestPrintEntityStateIncludesUnitInDefaultBarText(t *testing.T) {
	state := &homeassistant.HomeAssistantState{
		State: "23.3",
		Attributes: map[string]json.RawMessage{
			"unit_of_measurement": json.RawMessage(`"°C"`),
		},
	}

	read, write, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	originalStdout := os.Stdout
	os.Stdout = write
	t.Cleanup(func() { os.Stdout = originalStdout })

	printEntityState(state, "Temperature", entityWatchOutputOptions{BarJSON: true})
	if err := write.Close(); err != nil {
		t.Fatal(err)
	}

	output, err := io.ReadAll(read)
	if err != nil {
		t.Fatal(err)
	}

	if got := strings.TrimSpace(string(output)); got != `{"class":"23.3","name":"Temperature","text":"23.3 °C","tooltip":"23.3 °C"}` {
		t.Fatalf("printEntityState() = %s", got)
	}
}

func TestParseCoverPosition(t *testing.T) {
	for _, value := range []string{"0", "10", "20", "30", "60", "100"} {
		if _, err := parseCoverPosition(value); err != nil {
			t.Fatalf("parseCoverPosition(%q) returned %v", value, err)
		}
	}

	for _, value := range []string{"", "10.5", "-1", "101"} {
		if _, err := parseCoverPosition(value); err == nil {
			t.Fatalf("parseCoverPosition(%q) succeeded", value)
		}
	}
}

func TestParseInputNumberValue(t *testing.T) {
	for _, value := range []string{"0", "23.8", "36", "-1.5"} {
		if _, err := parseInputNumberValue(value); err != nil {
			t.Fatalf("parseInputNumberValue(%q) returned %v", value, err)
		}
	}

	for _, value := range []string{"", "off", "NaN", "+Inf"} {
		if _, err := parseInputNumberValue(value); err == nil {
			t.Fatalf("parseInputNumberValue(%q) succeeded", value)
		}
	}
}

func TestAdjustedInputNumberValue(t *testing.T) {
	state := &homeassistant.HomeAssistantState{
		State: "23.8",
		Attributes: map[string]json.RawMessage{
			"min":  json.RawMessage(`16`),
			"max":  json.RawMessage(`36`),
			"step": json.RawMessage(`0.1`),
		},
	}

	if got, err := adjustedInputNumberValue(state, 1); err != nil || got != 23.9 {
		t.Fatalf("adjustedInputNumberValue(increment) = %v, %v", got, err)
	}
	if got, err := adjustedInputNumberValue(state, -1); err != nil || got != 23.7 {
		t.Fatalf("adjustedInputNumberValue(decrement) = %v, %v", got, err)
	}

	state.State = "36"
	if got, err := adjustedInputNumberValue(state, 1); err != nil || got != 36 {
		t.Fatalf("adjustedInputNumberValue(maximum) = %v, %v", got, err)
	}
}
