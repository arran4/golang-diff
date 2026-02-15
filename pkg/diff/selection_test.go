package diff

import (
	"reflect"
	"testing"
)

func TestLineSelectionCompile(t *testing.T) {
	tests := []struct {
		input    string
		expected *ShortCodeLineSelectionLang
		wantErr  bool
	}{
		{
			input: "a:1-8,b:2-9",
			expected: &ShortCodeLineSelectionLang{
				RangesA: []Range{{1, 8}},
				RangesB: []Range{{2, 9}},
			},
		},
		{
			input: "1-8,2-9",
			expected: &ShortCodeLineSelectionLang{
				RangesA: []Range{{1, 8}},
				RangesB: []Range{{2, 9}},
			},
		},
		{
			input: "1-8",
			expected: &ShortCodeLineSelectionLang{
				RangesA: []Range{{1, 8}},
				RangesB: []Range{{1, 8}},
			},
		},
		{
			input: "a:10",
			expected: &ShortCodeLineSelectionLang{
				RangesA: []Range{{10, 10}},
				RangesB: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			s := WithLineSelectionShortCode(tt.input)
			got, err := s.Compile()
			if (err != nil) != tt.wantErr {
				t.Errorf("Compile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Compile() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestWidthSelectionCompile(t *testing.T) {
	tests := []struct {
		input    string
		expected *ShortCodeWidthSelectionLang
		wantErr  bool
	}{
		{
			input: "a:1-8,b:2-9",
			expected: &ShortCodeWidthSelectionLang{
				RangesA: []Range{{1, 8}},
				RangesB: []Range{{2, 9}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			s := WithWidthSelectionShortCode(tt.input)
			got, err := s.Compile()
			if (err != nil) != tt.wantErr {
				t.Errorf("Compile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Compile() = %v, want %v", got, tt.expected)
			}
		})
	}
}
