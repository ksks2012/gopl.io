package eval

import (
	"testing"
)

func TestVar_String(t *testing.T) {
	v := Var("foo")
	if got := v.String(); got != "foo" {
		t.Errorf("Var.String() = %q, want %q", got, "foo")
	}
}

func TestLiteral_String(t *testing.T) {
	l := literal(42.5)
	if got := l.String(); got != "42.5" {
		t.Errorf("literal.String() = %q, want %q", got, "42.5")
	}
}

func TestUnary_String(t *testing.T) {
	u := unary{op: '+', x: literal(7)}
	want := "+(7)"
	if got := u.String(); got != want {
		t.Errorf("unary.String() = %q, want %q", got, want)
	}
}

func TestBinary_String(t *testing.T) {
	b := binary{op: '-', x: literal(10), y: Var("z")}
	want := "(10 - z)"
	if got := b.String(); got != want {
		t.Errorf("binary.String() = %q, want %q", got, want)
	}
}

func TestCall_String(t *testing.T) {
	c := call{fn: "max", args: []Expr{literal(1), Var("b")}}
	want := "max(1, b)"
	if got := c.String(); got != want {
		t.Errorf("call.String() = %q, want %q", got, want)
	}
}

func TestCall_String_NoArgs(t *testing.T) {
	c := call{fn: "rand", args: []Expr{}}
	want := "rand()"
	if got := c.String(); got != want {
		t.Errorf("call.String() = %q, want %q", got, want)
	}
}

func TestExprStringRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		expr Expr
	}{
		{
			name: "Var",
			expr: Var("X"),
		},
		{
			name: "Literal",
			expr: literal(42.0),
		},
		{
			name: "Unary",
			expr: unary{op: '-', x: literal(5)},
		},
		{
			name: "Binary",
			expr: binary{op: '*', x: literal(2), y: Var("Y")},
		},
		{
			name: "Call with args",
			expr: call{fn: "min", args: []Expr{literal(1), literal(2)}},
		},
		{
			name: "Call no args",
			expr: call{fn: "rand", args: []Expr{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.expr.String()
			parsed, err := Parse(s)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", s, err)
			}
			if parsed.String() != s {
				t.Errorf("String() round-trip failed: got %q, want %q", parsed.String(), s)
			}
		})
	}
}
