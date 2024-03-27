package evaluator

import "testing"

import (
	"go-learning/interpreter/lexer"
	"go-learning/interpreter/object"
	"go-learning/interpreter/parser"
)

func TestAll(t *testing.T) {
	TestEvalIntegerExpression(t)
	TestEvalBooleanExpression(t)
	TestBangOperator(t)
}

func TestEvalIntegerExpression(test *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, testData := range tests {
		evaluated := testEval(testData.input)
		testIntegerObject(test, evaluated, testData.expected)
	}
}

func testEval(input string) object.Object {
	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()

	return Eval(program)
}

func testIntegerObject(test *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		test.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		test.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func TestBangOperator(test *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, testData := range tests {
		evaluated := testEval(testData.input)
		testBooleanObject(test, evaluated, testData.expected)
	}
}
