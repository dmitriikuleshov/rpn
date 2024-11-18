package rpn

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

func Calc(expression string) (float64, error) {
	tokens, err := tokenize(expression)
	if err != nil {
		return 0, err
	}

	postfixTokens, err := infixToPostfix(tokens)
	if err != nil {
		return 0, err
	}

	return evaluatePostfix(postfixTokens)
}

func tokenize(expression string) ([]string, error) {
	var tokens []string
	var num strings.Builder

	for i, ch := range expression {
		if unicode.IsDigit(ch) || ch == '.' {
			num.WriteRune(ch)
		} else {
			if num.Len() > 0 {
				tokens = append(tokens, num.String())
				num.Reset()
			}
			if ch == ' ' {
				continue
			}
			if strings.ContainsRune("+-*/()", ch) {
				if ch == '-' && (i == 0 || tokens[len(tokens)-1] == "(" || strings.ContainsRune("+-*/", rune(expression[i-1]))) {
					tokens = append(tokens, "_")
				} else {
					tokens = append(tokens, string(ch))
				}
			} else {
				return nil, errors.New("invalid character in expression")
			}
		}
	}
	if num.Len() > 0 {
		tokens = append(tokens, num.String())
	}
	return tokens, nil
}

func infixToPostfix(tokens []string) ([]string, error) {
	var postfix []string
	var ops []string

	precedence := map[string]int{
		"+": 1, "-": 1,
		"*": 2, "/": 2,
		"_": 3, // Unary minus
	}

	for _, token := range tokens {
		switch token {
		case "(":
			ops = append(ops, token)
		case ")":
			for len(ops) > 0 && ops[len(ops)-1] != "(" {
				postfix = append(postfix, ops[len(ops)-1])
				ops = ops[:len(ops)-1]
			}
			if len(ops) == 0 || ops[len(ops)-1] != "(" {
				return nil, errors.New("mismatched parentheses")
			}
			ops = ops[:len(ops)-1]
		case "+", "-", "*", "/", "_":
			for len(ops) > 0 && precedence[ops[len(ops)-1]] >= precedence[token] {
				postfix = append(postfix, ops[len(ops)-1])
				ops = ops[:len(ops)-1]
			}
			ops = append(ops, token)
		default:
			postfix = append(postfix, token)
		}
	}

	for len(ops) > 0 {
		if ops[len(ops)-1] == "(" {
			return nil, errors.New("mismatched parentheses")
		}
		postfix = append(postfix, ops[len(ops)-1])
		ops = ops[:len(ops)-1]
	}

	return postfix, nil
}

func evaluatePostfix(postfix []string) (float64, error) {
	var stack []float64

	for _, token := range postfix {
		switch token {
		case "+", "-", "*", "/":
			if len(stack) < 2 {
				return 0, errors.New("invalid expression")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var result float64
			switch token {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				if b == 0 {
					return 0, errors.New("division by zero")
				}
				result = a / b
			}
			stack = append(stack, result)
		case "_":
			if len(stack) < 1 {
				return 0, errors.New("invalid expression")
			}
			stack[len(stack)-1] *= -1
		default:
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, errors.New("invalid number")
			}
			stack = append(stack, num)
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("invalid expression")
	}

	return stack[0], nil
}

