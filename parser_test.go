package main

import "testing"

func verifyNode(t *testing.T, node ASTNode, name string) {
	if node.GetName() != name {
		t.Errorf("Expected node of type %s, got %s", name, node.GetName())
	}
}

func TestParseIntConstAssign(t *testing.T) {
	prog, err := Parse(Tokenize("x", "a = 10"))
	if err != nil {
		t.Fatal(err)
	}
	verifyNode(t, prog, "Program")
	if len(prog.Init) != 1 {
		t.Errorf("Expected 1 Init, got %d", len(prog.Init))
	}
}
