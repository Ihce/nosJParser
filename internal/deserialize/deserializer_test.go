package nosj_parser

import "testing"

func TestDeserialize(t *testing.T) {
	test := []struct {
		input       string
		shouldPanic bool
	}{
		{"(<>)", false},
		{"(<>", true},
		{"<>", true},
		{"(<", true},
		{"", true},
		{"()", true},
		{"<>)", true},
		{"(<)", true},
	}

	for _, tt := range test {
		t.Run(tt.input, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.shouldPanic {
						t.Errorf("Deserialize() panicked for input %q, but it shouldn't have", tt.input)
					}
				} else {
					if tt.shouldPanic {
						t.Errorf("Deserialize() did not panic for input %q, but it should have", tt.input)
					}
				}
			}()
			deserialize(tt.input)
		})
	}
}
