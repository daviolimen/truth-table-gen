package evaluator

import (
	"bytes"
	"fmt"
	"sort"
	"truth-table/ast"
)

type Environment struct {
	variableValue map[string]Object
}

type Object struct {
	value bool
	err   string
}

func NewEnvironment() *Environment {
	vv := make(map[string]Object)
	return &Environment{variableValue: vv}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.variableValue[name]
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.variableValue[name] = val
	return val
}

func CreateTruthTable(node ast.Node, e *Environment) string {
	var out bytes.Buffer

	findVariables(node, e)

	keys := make([]string, 0)
	for k, _ := range e.variableValue {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		out.WriteString(k + " ")
	}
	out.WriteString("ans")
	out.WriteString("\n")

	for mask := 0; mask < (1 << len(e.variableValue)); mask++ {
		for i, k := range keys {
			val := (mask & (1 << (len(e.variableValue) - i - 1))) != 0
			e.variableValue[k] = Object{value: val, err: ""}
			if val {
				out.WriteString("1")
			} else {
				out.WriteString("0")
			}
			out.WriteString(" ")
		}
		if eval(node, e).value {
			out.WriteString("1")
		} else {
			out.WriteString("0")
		}
		out.WriteString("\n")
	}

	return out.String()
}

func findVariables(node ast.Node, e *Environment) {
	switch node := node.(type) {
	case *ast.Variable:
		e.Set(node.Value, Object{value: false, err: ""})
	case *ast.PrefixExpression:
		findVariables(node.Right, e)
	case *ast.InfixExpression:
		findVariables(node.Left, e)
		findVariables(node.Right, e)
	}
}

func eval(node ast.Node, e *Environment) Object {
	switch node := node.(type) {
	case *ast.Literal:
		return Object{value: node.Value, err: ""}
	case *ast.Variable:
		obj, ok := e.Get(node.Value)
		if !ok {
			return Object{value: false, err: fmt.Sprintf("Variable named %s not found", node.Value)}
		}
		return obj
	case *ast.PrefixExpression:
		right := eval(node.Right, e)
		if right.err != "" {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := eval(node.Left, e)
		if left.err != "" {
			return left
		}
		right := eval(node.Right, e)
		if right.err != "" {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	}

	return Object{value: false, err: "Error when evaluating expression"}
}

func evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "~":
		return Object{value: !right.value, err: ""}
	default:
		return Object{value: false, err: fmt.Sprintf("The prefix operator %s does not exist", operator)}
	}
}

func evalInfixExpression(operator string, left Object, right Object) Object {
	switch operator {
	case "&":
		return Object{value: left.value && right.value, err: ""}
	case "|":
		return Object{value: left.value || right.value, err: ""}
	case "^":
		return Object{value: left.value != right.value, err: ""}
	case "=>":
		return Object{value: !left.value || right.value, err: ""}
	case "<>":
		return Object{value: left.value == right.value, err: ""}
	default:
		return Object{value: false, err: fmt.Sprintf("The infix operator %s does not exist", operator)}
	}
}
