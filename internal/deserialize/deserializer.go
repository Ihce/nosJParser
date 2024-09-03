package deserializer

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"unicode"
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

func (t TokenType) String() string {
	return [...]string{"START_MAP", "END_MAP", "KEY", "VALUE"}[t]
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

				tokens = append(tokens, Token{Type: START_MAP, Value: "(<"})
				inKey = true
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
				tokens = append(tokens, Token{Type: END_MAP, Value: ">)"})
				inKey = false
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

		case ',':
			if !inKey {
				tokens = append(tokens, Token{Type: VALUE, Value: currentToken.String()})
				currentToken.Reset()
				inKey = true
			} else {
				if currentToken.Len() > 0 {
					tokens = append(tokens, Token{Type: KEY, Value: currentToken.String()})
					currentToken.Reset()
				}
				tokens = append(tokens, Token{Type: VALUE, Value: ""})
			}

		default:
			currentToken.WriteByte(char)
		}
	}
	return tokens, nil
}

func parseTokens(tokens []Token) (map[string]interface{}, error) {
	stack := []map[string]interface{}{}
	var currentMap map[string]interface{}
	var currentKey string

	for _, token := range tokens {
		switch token.Type {
		case START_MAP:
			// Create a new map and push the current map onto the stack if it exists
			newMap := map[string]interface{}{}
			if currentMap != nil {
				stack = append(stack, currentMap)
				currentMap[currentKey] = newMap
			}
			currentMap = newMap
		case END_MAP:
			// Pop the map from the stack and set the current map as the value for the current key in the parent map
			if len(stack) == 0 {
				return currentMap, nil // Return the root map when the stack is empty
			}
			parentMap := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			currentMap = parentMap
		case KEY:
			// Set the current key to the token's value
			if strings.Contains(token.Value, " ") {
				return nil, fmt.Errorf("ERROR -- Invalid key -- Spaces Found In Key: %s", strings.Trim(token.Value, " "))
			}
			currentKey = token.Value
		case VALUE:
			// Set the value for the current key in the current map
			currentMap[currentKey] = token.Value
			currentKey = ""
		}
	}
	return currentMap, nil
}

func printMap(dMap map[string]interface{}) (string, error) {
	var result strings.Builder
	result.WriteString("begin-map\n")
	for key, value := range dMap {
		switch v := value.(type) {
		case map[string]interface{}:
			result.WriteString(fmt.Sprintf("%s -- map -- \n", key))
			nestedMapValue, err := printMap(v)
			if err != nil {
				return "", err
			}
			result.WriteString(nestedMapValue)
		case string:
			parsedValueType, parsedValue, err := parseValues(Token{Value: v})
			if err != nil {
				return "", err
			}
			result.WriteString(fmt.Sprintf("%s -- %s -- %s\n", key, parsedValueType, parsedValue))
		default:
			return "", fmt.Errorf("ERROR -- Unsupported value type for key %s", key)
		}
	}
	result.WriteString("end-map\n")
	return result.String(), nil
}

// This is awful. I'm sorry.
func parseValues(token Token) (string, string, error) {
	if valid, err := isValidComplexString(token.Value); err != nil {
		return "", "", err
	} else if valid {
		resultString, err := url.QueryUnescape(token.Value)
		if err != nil {
			return "", "", fmt.Errorf("ERROR -- Invalid value -- Could not escape -- %s", token.Value)
		}
		return "string", resultString, nil
	}

	// Check if the token value is a valid binary string
	if valid, err := isBinaryString(token.Value); err != nil {
		return "", "", err
	} else if valid {
		value, err := strconv.ParseInt(token.Value, 2, 64)
		if err != nil {
			return "", "", fmt.Errorf("ERROR -- Invalid value %s", token.Value)
		}

		bitLength := len(token.Value)

		if token.Value[0] == '1' {
			value = value - (1 << bitLength)
		}

		return "num", strconv.FormatInt(value, 10), nil
	}

	// Check if the token value is a valid simple string
	if valid, err := isValidSimpleString(token.Value); err != nil {
		return "", "", err
	} else if valid {
		resultString := token.Value[:len(token.Value)-1]
		return "string", resultString, nil
	}

	// If none of the above, return an error
	return "Unknown", "", fmt.Errorf("ERROR -- Invalid value %s", token.Value)
}

func isBinaryString(input string) (bool, error) {
	for _, char := range input {
		if char != '0' && char != '1' {
			return false, nil
		}
	}
	return true, nil
}

//	func isValidMap(input string) bool {
//		return strings.HasPrefix(input, "(<") && strings.HasSuffix(input, ">)")
//	}
func isValidSimpleString(input string) (bool, error) {
	if len(input) == 0 {
		return false, fmt.Errorf("ERROR -- Invalid Value -- Simple strings can not be empty")
	}
	trimmedString := strings.Trim(input, " ")
	if input[0] == 0x20 || input[0] == 0x09 {
		return false, fmt.Errorf("ERROR -- Invalid Value -- Simple strings can not start with whitespace")
	}
	if trimmedString[len(trimmedString)-1:] == "s" && input[len(input)-1:] == " " {
		return false, fmt.Errorf("ERROR -- Invalid Value -- Simple strings can not end with whitespace")
	}
	if input[len(input)-1:] == "s" {
		for _, char := range []byte(input) {
			if !(unicode.IsLetter(rune(char)) || char == 0x20 || char == 0x09) {
				return false, fmt.Errorf("ERROR -- Invalid Value -- Simple strings has an illegal character(s)")
			}
		}
		return true, nil
	}
	return false, nil
}

func isValidComplexString(input string) (bool, error) {
	if strings.Contains(input, "%") {
		return true, nil
	}
	return false, nil
}

func Deserialize(input string) (string, error) {
	if len(input) < 4 {
		return "", fmt.Errorf("ERROR -- Input is too short to contain a root map")
	}
	input = strings.Trim(input, " ")
	tokens, err := tokenize(input)
	if err != nil {
		return "", err
	}

	dMap, err := parseTokens(tokens)
	if err != nil {
		return "", err
	}

	result, err := printMap(dMap)

	if err != nil {
		return "", err
	}

	return result, nil
}
