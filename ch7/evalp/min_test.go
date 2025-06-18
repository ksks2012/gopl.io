package eval

import (
	"testing"
)

func TestMin_Eval(t *testing.T) {
	tests := []struct {
		name string
		args []Expr
		env  Env
		want float64
	}{
		{
			name: "no args",
			args: []Expr{},
			env:  Env{},
			want: 0,
		},
		{
			name: "single literal arg",
			args: []Expr{literal(5)},
			env:  Env{},
			want: 5,
		},
		{
			name: "multiple literal args",
			args: []Expr{literal(10), literal(3), literal(7)},
			env:  Env{},
			want: 3,
		},
		{
			name: "with variables",
			args: []Expr{Var("a"), Var("b"), literal(8)},
			env:  Env{"a": 4, "b": 9},
			want: 4,
		},
		{
			name: "all variables",
			args: []Expr{Var("x"), Var("y")},
			env:  Env{"x": 2.5, "y": 2.7},
			want: 2.5,
		},
		{
			name: "negative values",
			args: []Expr{literal(-2), literal(-5), literal(0)},
			env:  Env{},
			want: -5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := min{args: tt.args}
			got := m.Eval(tt.env)
			if got != tt.want {
				t.Errorf("min.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}
