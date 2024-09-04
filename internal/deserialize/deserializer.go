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

type OrderedMap struct {
	keys []string
	data map[string]interface{}
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		keys: []string{},
		data: make(map[string]interface{}),
	}
}

func (om *OrderedMap) Set(key string, value interface{}) {
	if _, exists := om.data[key]; !exists {
		om.keys = append(om.keys, key)
	}
	om.data[key] = value
}

func (om *OrderedMap) Get(key string) (interface{}, bool) {
	value, exists := om.data[key]
	return value, exists
}

func (om *OrderedMap) Keys() []string {
	return om.keys
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
				currentToken.Reset()
				tokens = append(tokens, Token{Type: START_MAP, Value: "(<"})
				inKey = true
				i++
			} else {
				return nil, fmt.Errorf("ERROR -- Invalid start of a map at index %d", i)
			}
		case '>':
			if input[i+1] == ')' && i+1 < len(input) {
				if currentToken.Len() > 0 {
					if len(strings.Trim(currentToken.String(), " ")) == 0 {
						currentToken.Reset()
					} else {
						tokens = append(tokens, Token{Type: VALUE, Value: currentToken.String()})
						currentToken.Reset()
					}
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

func parseTokens(tokens []Token) (*OrderedMap, error) {
	stack := []*OrderedMap{}
	var currentMap *OrderedMap
	var currentKey string

	for _, token := range tokens {
		switch token.Type {
		case START_MAP:
			newMap := NewOrderedMap()
			if currentMap != nil {
				stack = append(stack, currentMap)
				currentMap.Set(currentKey, newMap)
			}
			currentMap = newMap
		case END_MAP:
			if len(stack) == 0 {
				return currentMap, nil
			}
			parentMap := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			currentMap = parentMap
		case KEY:
			if valid, err := isValidKey(token.Value); err != nil {
				return nil, err
			} else if !valid {
				return nil, fmt.Errorf("ERROR -- Invalid Key")
			}
			currentKey = token.Value
		case VALUE:
			currentMap.Set(currentKey, token.Value)
			currentKey = ""
		}
	}
	return currentMap, nil
}

func printMap(dMap *OrderedMap) (string, error) {
	var result strings.Builder
	result.WriteString("begin-map\n")
	for _, key := range dMap.Keys() {
		value := dMap.data[key]
		switch v := value.(type) {
		case *OrderedMap:
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

	if valid, err := isValidSimpleString(token.Value); err != nil {
		return "", "", err
	} else if valid {
		resultString := token.Value[:len(token.Value)-1]
		return "string", resultString, nil
	}

	return "Unknown", "", fmt.Errorf("ERROR -- Invalid value %s", token.Value)
}

func isBinaryString(input string) (bool, error) {
	for _, char := range input {
		if char != 0x30 && char != 0x31 {
			return false, nil
		}
	}
	return true, nil
}

func isValidKey(input string) (bool, error) {
	for _, char := range []byte(input) {
		if !(0x61 <= char && char <= 0x7A) {
			return false, nil
		}

	}
	return true, nil
}

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
	if len(input) == 0 {
		return "", fmt.Errorf("ERROR -- Input is empty")
	}
	if !strings.Contains(input, "(<") && !strings.Contains(input, ">)") {
		return "", fmt.Errorf("ERROR -- Input does not contain a map")
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
