package parsing

import (
	"testing"
	"zen/lang/parsing/expression"
)

func TestMapLiterals(t *testing.T) {
	print("parsing maps.zen\n")
	programNode := ParseTestFile(t, "maps.zen")
	if programNode == nil {
		return
	}

	// var person = { "name": "john", "age": 30 }
	varDecl := AssertVarDeclaration(t, programNode.Statements[0], "person", false, false)
	mapLit, ok := varDecl.Initializer.(*expression.MapLiteralExpression)
	if !ok {
		t.Errorf("Expected MapLiteralExpression, got %T", varDecl.Initializer)
		return
	}

	if len(mapLit.Entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(mapLit.Entries))
		return
	}

	// Check first entry
	AssertLiteralExpression(t, mapLit.Entries[0].Key, "name")
	AssertLiteralExpression(t, mapLit.Entries[0].Value, "john")

	// Check second entry
	AssertLiteralExpression(t, mapLit.Entries[1].Key, "age")
	AssertLiteralExpression(t, mapLit.Entries[1].Value, int64(30))

	// var settings : Map<string, float> = { "volume": 0.5, "brightness": 1.0 }
	varDecl = AssertVarDeclaration(t, programNode.Statements[1], "settings", false, false)

	// Check map literal
	mapLit, ok = varDecl.Initializer.(*expression.MapLiteralExpression)
	if !ok {
		t.Errorf("Expected MapLiteralExpression, got %T", varDecl.Initializer)
		return
	}

	if len(mapLit.Entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(mapLit.Entries))
		return
	}

	// Check entries
	AssertLiteralExpression(t, mapLit.Entries[0].Key, "volume")
	AssertLiteralExpression(t, mapLit.Entries[0].Value, float64(0.5))
	AssertLiteralExpression(t, mapLit.Entries[1].Key, "brightness")
	AssertLiteralExpression(t, mapLit.Entries[1].Value, float64(1.0))

	// var empty = { }
	varDecl = AssertVarDeclaration(t, programNode.Statements[2], "empty", false, false)
	mapLit, ok = varDecl.Initializer.(*expression.MapLiteralExpression)
	if !ok {
		t.Errorf("Expected MapLiteralExpression, got %T", varDecl.Initializer)
		return
	}

	if len(mapLit.Entries) != 0 {
		t.Errorf("Expected empty map, got %d entries", len(mapLit.Entries))
	}

	// var nested = { "user": { "name": "alice", "scores": [100, 95, 98] } }
	varDecl = AssertVarDeclaration(t, programNode.Statements[3], "nested", false, false)
	mapLit, ok = varDecl.Initializer.(*expression.MapLiteralExpression)
	if !ok {
		t.Errorf("Expected MapLiteralExpression, got %T", varDecl.Initializer)
		return
	}

	if len(mapLit.Entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(mapLit.Entries))
		return
	}

	// Check outer map
	AssertLiteralExpression(t, mapLit.Entries[0].Key, "user")
	innerMap, ok := mapLit.Entries[0].Value.(*expression.MapLiteralExpression)
	if !ok {
		t.Errorf("Expected MapLiteralExpression for nested map, got %T", mapLit.Entries[0].Value)
		return
	}

	// Check inner map
	if len(innerMap.Entries) != 2 {
		t.Errorf("Expected 2 entries in nested map, got %d", len(innerMap.Entries))
		return
	}

	AssertLiteralExpression(t, innerMap.Entries[0].Key, "name")
	AssertLiteralExpression(t, innerMap.Entries[0].Value, "alice")

	AssertLiteralExpression(t, innerMap.Entries[1].Key, "scores")
	scores, ok := innerMap.Entries[1].Value.(*expression.ArrayLiteralExpression)
	if !ok {
		t.Errorf("Expected ArrayLiteralExpression for scores, got %T", innerMap.Entries[1].Value)
		return
	}

	// Check scores array
	if len(scores.Elements) != 3 {
		t.Errorf("Expected 3 scores, got %d", len(scores.Elements))
		return
	}

	AssertLiteralExpression(t, scores.Elements[0], int64(100))
	AssertLiteralExpression(t, scores.Elements[1], int64(95))
	AssertLiteralExpression(t, scores.Elements[2], int64(98))
}
