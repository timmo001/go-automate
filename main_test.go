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
