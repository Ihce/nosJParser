package deserializer

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type TokenType int

const (
	START_MAP TokenType = iota
	END_MAP
	KEY
	VALUE
)

type Token struct {
	Type  TokenType
	Value string
}

func tokenize(input string) ([]Token, error) {
	var tokens []Token
	var currentToken strings.Builder
	inKey := true

	for i := 0; i < len(input); i++ {
		char := input[i]
		switch char {
		case '(':
			if input[i+1] == '<' && i+1 < len(input) {
				tokens = append(tokens, Token{Type: START_MAP})
				i++
			} else {
				return nil, fmt.Errorf("ERROR -- Invalid start of a map at index %d", i)
			}
		case '>':
			if input[i+1] == ')' && i+1 < len(input) {
				if currentToken.Len() > 0 {
					tokens = append(tokens, Token{Type: VALUE, Value: currentToken.String()})
					currentToken.Reset()
				}
				tokens = append(tokens, Token{Type: END_MAP})
				i++
			} else {
				return nil, fmt.Errorf("ERROR -- Invalid end of a map at index %d", i)
			}
		case ':':
			if inKey {
				tokens = append(tokens, Token{Type: KEY, Value: currentToken.String()})
				currentToken.Reset()
				inKey = false
			} else {
				currentToken.WriteByte(char)
			}
		default:
			currentToken.WriteByte(char)
		}
	}
	return tokens, nil
}

func parseTokens(tokens []Token) (string, error) {
	var result strings.Builder
	var currentKey string
	var stack []TokenType
	for _, token := range tokens {
		switch token.Type {
		case START_MAP:
			stack = append(stack, START_MAP)
			result.WriteString("begin-map\n")
		case END_MAP:
			if len(stack) == 0 || stack[len(stack)-1] != START_MAP {
				return "", fmt.Errorf("ERROR -- Unmatched end of map")
			}
			result.WriteString("end-map\n")
			stack = stack[:len(stack)-1]
		case KEY:
			currentKey = token.Value
		case VALUE:
			if currentKey == "" {
				return "", fmt.Errorf("ERROR -- Value without a key")
			}
			valueType, value, err := parseValues(token)
			if err != nil {
				return "", err
			}
			result.WriteString(fmt.Sprintf("%s -- %s -- %s\n", currentKey, valueType, value))
			currentKey = ""
		}

	}
	if len(stack) != 0 {
		return "", fmt.Errorf("ERROR -- Unmatched begin-map")
	}
	return result.String(), nil
}

func parseValues(token Token) (string, string, error) {
	switch {
	case isValidMap(token.Value):
		// Tokenize the nested map value
		nestedTokens, err := tokenize(token.Value)
		if err != nil {
			return "", "", fmt.Errorf("ERROR -- Invalid nested map %s", token.Value)
		}
		// Parse the nested tokens
		nestedResult, err := parseTokens(nestedTokens)
		if err != nil {
			return "", "", fmt.Errorf("ERROR -- Failed to parse nested map %s", token.Value)
		}
		return "map", nestedResult, nil
	case isValidComplexString(token.Value):
		resultString, err := url.QueryUnescape(token.Value)
		if err != nil {
			return "", "", fmt.Errorf("ERROR -- Invalid value -- Could not escape -- %s", token.Value)
		}
		return "string", resultString, nil

	case isBinaryString(token.Value):
		value, err := strconv.ParseInt(token.Value, 2, 64)
		if err != nil {
			return "", "", fmt.Errorf("ERROR -- Invalid value %s", token.Value)
		}

		bitLength := len(token.Value)

		if token.Value[0] == '1' {
			value = value - (1 << bitLength)
		}

		return "num", strconv.FormatInt(value, 10), nil

	case isValidSimpleString(token.Value):
		resultString := token.Value[:len(token.Value)-1]
		return "string", resultString, nil

	default:
		return "Unknown", "", fmt.Errorf("ERROR -- Invalid value %s", token.Value)
	}
}

func isBinaryString(input string) bool {
	for _, char := range input {
		if char != '0' && char != '1' {
			return false
		}
	}
	return true
}
func isValidMap(input string) bool {
	return strings.HasPrefix(input, "(<") && strings.HasSuffix(input, ">)")
}
func isValidSimpleString(input string) bool {
	if input[len(input)-1:] == "s" {
		for _, char := range []byte(input) {
			if char < 32 || char > 126 {
				return false
			}
		}
	}
	return true
}

func isValidComplexString(input string) bool {
	return strings.Contains(input, "%")
}

//	func tokenTypeToString(tokenType TokenType) string {
//		switch tokenType {
//		case START_MAP:
//			return "START_MAP"
//		case END_MAP:
//			return "END_MAP"
//		case KEY:
//			return "KEY"
//		case VALUE:
//			return "VALUE"
//		default:
//			return "UNKNOWN"
//		}
//	}
func Deserialize(input string) (string, error) {
	if len(input) < 4 {
		return "", fmt.Errorf("ERROR -- Input is too short to contain a root map")
	}
	tokens, err := tokenize(input)
	if err != nil {
		return "", err
	}

	result, err := parseTokens(tokens)
	if err != nil {
		return "", err
	}
	return result, err
}
