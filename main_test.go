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
