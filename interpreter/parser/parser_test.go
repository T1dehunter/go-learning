package parser

import (
	"go-learning/interpreter/parser/ast"
	"testing"
)

import (
	"go-learning/interpreter/lexer"
)

func TestAll(test *testing.T) {
	TestLetStatements(test)
	TestReturnStatements(test)
}

func TestLetStatements(test *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
	`
	l := lexer.New(input)
	parser := New(l)

	program := parser.ParseProgram()

	checkParserErrors(test, parser)

	if program == nil {
		test.Fatalf("ParseProgram() returned nil")
	}
	programStatementsLength := len(program.Statements)
	if programStatementsLength != 3 {
		test.Fatalf("program.Statements does not contain 3 statements. got=%d", programStatementsLength)
	}

	testData := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, dataItem := range testData {
		stmt := program.Statements[i]
		if !testLetStatement(test, stmt, dataItem.expectedIdentifier) {
			return
		}
	}

}

func checkParserErrors(test *testing.T, parser *Parser) {
	errors := parser.Errors()
	if len(errors) == 0 {
		return
	}
	test.Errorf("parser has %d errors", len(errors))
	for _, message := range errors {
		test.Errorf("parser error: %q", message)
	}
	test.FailNow()
}

func testLetStatement(test *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "let" {
		test.Errorf("s.TokenLiteral not 'let'. got=%q", statement.TokenLiteral())
		return false
	}

	letStmt, ok := statement.(*ast.LetStatement)
	if !ok {
		test.Errorf("s not *ast.LetStatement. got=%T", statement)
		return false
	}

	if letStmt.Name.Value != name {
		test.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		test.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatements(test *testing.T) {
	input := `
		return 5;
		return 10;
		return 999;
	`
	l := lexer.New(input)
	parser := New(l)

	program := parser.ParseProgram()

	checkParserErrors(test, parser)

	programStatementsLength := len(program.Statements)

	if programStatementsLength != 3 {
		test.Fatalf("program.Statements does not contain 3 statements. got=%d", programStatementsLength)
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			test.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			test.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
}
