// Practice 7.13: Implement the String method for each of the expression types in the eval package, so that they return a string representation of the expression. The String method should be defined on the type itself, not as a method on the Expr interface. The String method should return a string that describes the expression in a way that is easy to read and understand.
package eval

import (
	"fmt"
	"strings"
)

func (v Var) String() string {
	return string(v)
}

func (l literal) String() string {
	return fmt.Sprintf("%g", l)
}

func (u unary) String() string {
	return fmt.Sprintf("%c(%s)", u.op, u.x.String())
}

func (b binary) String() string {
	return fmt.Sprintf("(%s %c %s)", b.x, b.op, b.y)
}

func (c call) String() string {
	argsStr := make([]string, len(c.args))
	for i, arg := range c.args {
		argsStr[i] = arg.String()
	}
	return fmt.Sprintf("%s(%s)", c.fn, strings.Join(argsStr, ", "))

}
