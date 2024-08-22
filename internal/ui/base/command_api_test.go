package base

import "testing"

func TestEmptyCommandHandler(t *testing.T) {
	handler := EmptyCommandHandler()
	if len(handler.GetCommands()) != 0 {
		t.Error("test failed")
	}
}
