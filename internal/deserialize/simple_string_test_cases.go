package deserializer

type SimpleStringTestCase struct {
	Input    string
	Expected bool
	ErrMsg   string
}

type EdgeTestCase struct {
	Input    string
	Expected string
	ErrMsg   interface{}
}

// Generated loosely by LLM
var SimpleStringTestCases = []SimpleStringTestCase{
	// Valid cases
	{"hellos", true, ""},
	{"hello worlds", true, ""},
	{"hello\tworlds", true, ""},
	{"hellos", true, ""},

	// Invalid cases
	{" hellos", false, "ERROR -- Invalid Value -- Simple strings can not start with whitespace"},
	{"hellos ", false, "ERROR -- Invalid Value -- Simple strings can not end with whitespace"},
	{"hello!", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello123s", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello@worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"", false, "ERROR -- Invalid Value -- Simple strings can not be empty"},
	{"\thellos", false, "ERROR -- Invalid Value -- Simple strings can not start with whitespace"},
	{"hellos\t", false, "ERROR -- Invalid Value -- Simple strings can not end with whitespace"},
	{"hello\nworlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello\rworlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello\vworlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello\fworlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello\bworlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello\\worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello/worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello.worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello,worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello:worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello;worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello'worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello\"worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello<worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello>worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello[worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello]worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello{worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello}worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello|worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello~worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello`worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello^worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello&worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello*worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello(worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello)worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello-worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello+worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello=worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello_worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello#worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello$worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello%worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello?worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello!worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},
	{"hello\u00A0worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Non-breaking space
	{"hello\u200Bworlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Zero-width space
	{"hello\u202Eworlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Right-to-left override
	{"hello\u0000worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Null character
	{"hello\u0007worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Bell character
	{"hello\u001Bworlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Escape character
	{"hello\u007Fworlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Delete character
	{"hello\xF0\x9F\x98\x80worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"}, // Emoji
	{"hello\uFFFDworlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Replacement character
	{"hello\u2028worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Line separator
	{"hello\u2029worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Paragraph separator
	{"hello\u2060worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Word joiner
	{"hello\uFEFFworlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Byte order mark
	{"hello\u200Dworlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Zero-width joiner
	{"hello\u200Cworlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Zero-width non-joiner
	{"hello\u200Eworlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Left-to-right mark
	{"hello\u200Fworlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Right-to-left mark
	{"hello\u202Aworlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Left-to-right embedding
	{"hello\u202Bworlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Right-to-left embedding
	{"hello\u202Cworlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Pop directional formatting
	{"hello\u202Dworlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Left-to-right override
	{"hello\u202Eworlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Right-to-left override
	{"hello\u2066worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Left-to-right isolate
	{"hello\u2067worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Right-to-left isolate
	{"hello\u2068worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // First strong isolate
	{"hello\u2069worlds", false, "ERROR -- Invalid Value -- Simple strings has an illegal character(s)"},           // Pop directional isolate
}

var EdgeTestCases = []EdgeTestCase{
	{"", "", "ERROR -- Input is empty"},
	{"test", "", "ERROR -- Input does not contain a map"},
	{"(<a:abcds,b:1001,c:efghis,d:1011,e:jklmops>)", "begin-map\na -- string -- abcd\nb -- num -- -7\nc -- string -- efghi\nd -- num -- -5\ne -- string -- jklmop\nend-map\n", nil},
	{"(<a:s>)", "begin-map\na -- string -- \nend-map\n", nil},
	{"(<a:1010>)", "begin-map\na -- num -- -6\nend-map\n", nil},
	{"(<a:1010,a:1010>)", "begin-map\na -- num -- -6\nend-map\n", nil},
	{"(<a:s>)", "begin-map\na -- string -- \nend-map\n", nil},
	{"(<!:s>)", "", "ERROR -- Invalid Key"},
}
