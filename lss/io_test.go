package lss;

import (
	"flag"
	"testing"
	"os"
	"log"
	"encoding/json"
)

var inputFileName *string= flag.String("if", "testdata/problem.json", "Input file to read a Problem from.");

func TestJSONInput(t *testing.T) {
	file, err := os.Open(*inputFileName)
	if err != nil {
		t.Errorf("Could not read JSON input: %v", err)
	}

	p, err := ReadLSProblem(file)
	if err != nil {
		t.Errorf("Could not read JSON input: %v", err)
	}

	p_str, err := json.Marshal(p)
	if err != nil {
		t.Errorf("Could not write JSON output: %v", err)
	}
	log.Printf(string(p_str))
}
