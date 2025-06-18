// Practice 7.14: Implement the min function in the eval package that takes a variable number of arguments and returns the minimum value among them. The function should handle both literal values and variables, and it should return 0 if no arguments are provided. The function should also implement a String method that returns a string representation of the expression in a way that is easy to read and understand.
package eval

import "strings"

type min struct {
	args []Expr
}

func (m min) Eval(env Env) float64 {
	if len(m.args) == 0 {
		return 0 // or handle as an error
	}
	minValue := m.args[0].Eval(env)
	for _, arg := range m.args[1:] {
		value := arg.Eval(env)
		if value < minValue {
			minValue = value
		}
	}
	return minValue
}

func (m min) String() string {
	argsStr := make([]string, len(m.args))
	for i, arg := range m.args {
		argsStr[i] = arg.String()
	}
	return "min(" + strings.Join(argsStr, ", ") + ")"
}
