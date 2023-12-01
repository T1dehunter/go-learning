package repl

import (
	"bufio"
	"fmt"
	"go-learning/interpreter/lexer"
	"go-learning/interpreter/token"
	"io"
)

const PROMPT = ">>> "

func Start(input io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(input)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for currentToken := l.NextToken(); currentToken.Type != token.EOF; currentToken = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", currentToken)
		}
	}
}
