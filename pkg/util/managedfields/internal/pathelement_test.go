package internal

import "testing"

func TestPathElementRoundTrip(t *testing.T) {
	tests := []string{
		`i:0`,
		`i:1234`,
		`f:`,
		`f:spec`,
		`f:more-complicated-string`,
		`k:{"name":"my-container"}`,
		`k:{"port":"8080","protocol":"TCP"}`,
		`k:{"optionalField":null}`,
		`k:{"jsonField":{"A":1,"B":null,"C":"D","E":{"F":"G"}}}`,
		`k:{"listField":["1","2","3"]}`,
		`v:null`,
		`v:"some-string"`,
		`v:1234`,
		`v:{"some":"json"}`,
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			pe, err := NewPathElement(test)
			if err != nil {
				t.Fatalf("Failed to create path element: %v", err)
			}
			output, err := PathElementString(pe)
			if err != nil {
				t.Fatalf("Failed to create string from path element: %v", err)
			}
			if test != output {
				t.Fatalf("Expected round-trip:\ninput: %v\noutput: %v", test, output)
			}
		})
	}
}

func TestPathElementIgnoreUnknown(t *testing.T) {
	_, err := NewPathElement("r:Hello")
	if err != nil {
		t.Fatalf("Unknown qualifiers should be ignored")
	}
}

func TestNewPathElementError(t *testing.T) {
	tests := []string{
		``,
		`no-colon`,
		`i:index is not a number`,
		`i:1.23`,
		`i:`,
		`v:invalid json`,
		`v:`,
		`k:invalid json`,
		`k:{"name":invalid}`,
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			_, err := NewPathElement(test)
			if err == nil {
				t.Fatalf("Expected error, no error found")
			}
		})
	}
}
