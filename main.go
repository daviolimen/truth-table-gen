package main

import (
	"fmt"
	"os"
	"os/user"
	"truth-table/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the truth table solver!\n", user.Username)
	fmt.Printf("Feel free to type in expressions\n")
	repl.Start(os.Stdin, os.Stdout)
}
