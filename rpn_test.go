package rpn

import (
	"math"
	"testing"
)


func TestCalc(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
		expectErr  bool
	}{
		{"3 + 5", 8, false},
		{"10 - 2 * 3", 4, false},
		{"10 / 2 + 3", 8, false},
		{"(10 + 5) * 2", 30, false},
		{"3 + 5 * (2 - 8)", -27, false},
		{"3.5 + 2.5", 6, false},
		{"5 / 0", 0, true},                         // zero division
		{"2 * (3 + 5))", 0, true},                  // extra parenthesis
		{"-3 * 2", -6, false},
		{"2 * -3", 0, true},                        // invalid unary minus
		{"2 * (-3)", -6, false},
		{"(1 + 2) * (3 / (4 - 4))", 0, true},       // zero division
		{"10 + -5", 0, true},                       // invalid unary minus
		{"(2 + 2) * 2", 8, false},                 
		{"", 0, true},                              // empty line
		{"10.5 + 5.5 - 2", 14, false},        
		{"sqrt(4)", math.NaN(), true},              // invalid syntax
		{"5 - - 10", 0, true},                      // invalid unary minus
		{"((((((5*(5*5)*5))))))", 625, false},
		{"-5*5*5*5", -625, false},
		{"-(5*5*5*5)", -625, false},
		{"-((-(5*5*5*5)))", 625, false},
		{"1 + 1 + 1 +     1", 4, false},
		{"1 + 1 + 1 + + 1", 0, true},                // invalid syntax
		{"(2 + 2) * 5 - (1 - 1) * 100", 20, false},
	}

	for _, test := range tests {
		result, err := Calc(test.expression)
		if test.expectErr {
			if err == nil {
				t.Errorf("Expected error for expression %q, but got nil", test.expression)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for expression %q: %v", test.expression, err)
			}
			if math.Abs(result-test.expected) > 1e-9 {
				t.Errorf("For expression %q, expected %f, but got %f", test.expression, test.expected, result)
			}
		}
	}
}
