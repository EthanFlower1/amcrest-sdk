package amcrest

import "testing"

func TestParseKV(t *testing.T) {
	input := "table.General.MachineName=TestCam\ntable.General.LocalNo=1\n"
	result := parseKV(input)
	if result["table.General.MachineName"] != "TestCam" {
		t.Errorf("expected TestCam, got %s", result["table.General.MachineName"])
	}
	if result["table.General.LocalNo"] != "1" {
		t.Errorf("expected 1, got %s", result["table.General.LocalNo"])
	}
}

func TestParseKVSingleValue(t *testing.T) {
	input := "result=1\n"
	result := parseKV(input)
	if result["result"] != "1" {
		t.Errorf("expected 1, got %s", result["result"])
	}
}

func TestParseKVWithSpaces(t *testing.T) {
	input := "result = 2011-7-3 21:02:32\n"
	result := parseKV(input)
	if result["result"] != "2011-7-3 21:02:32" {
		t.Errorf("expected time string, got %s", result["result"])
	}
}

func TestStripTablePrefix(t *testing.T) {
	input := "table.General.MachineName=TestCam\ntable.General.LocalNo=1\n"
	result := parseKVWithPrefix(input, "table.General.")
	if result["MachineName"] != "TestCam" {
		t.Errorf("expected TestCam, got %s", result["MachineName"])
	}
	if result["LocalNo"] != "1" {
		t.Errorf("expected 1, got %s", result["LocalNo"])
	}
}
