package repl

import (
	"bufio"
	"fmt"
	"io"
	"truth-table/evaluator"
	"truth-table/lexer"
	"truth-table/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		expr, errorList := p.Parse()

		if len(errorList) != 0 {
			printParserErrors(out, errorList)
			continue
		}

		env := evaluator.NewEnvironment()
		truthTable := evaluator.CreateTruthTable(expr, env)
		io.WriteString(out, truthTable)
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
