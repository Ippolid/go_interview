package generics

import (
	"reflect"
	"strconv"
	"testing"
)

func TestDouble(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  int
	}{
		{
			name:  "positive value",
			input: 21,
			want:  42,
		},
		{
			name:  "negative value",
			input: -4,
			want:  -8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Double(tt.input)
			if got != tt.want {
				t.Fatalf("Double(%d) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestDoubleFloat(t *testing.T) {
	got := Double(2.5)
	if got != 5.0 {
		t.Fatalf("Double(2.5) = %v, want 5.0", got)
	}
}

func TestDotProduct(t *testing.T) {
	got := DotProduct([]int{1, 2, 3}, []int{4, 5, 6})
	if got != 32 {
		t.Fatalf("DotProduct() = %d, want 32", got)
	}
}

func TestDotProductPanicsOnUnequalLength(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for slices of unequal length")
		}
	}()

	_ = DotProduct([]int{1, 2}, []int{1})
}

func TestSum(t *testing.T) {
	input := map[string]int64{
		"a": 5,
		"b": 7,
		"c": 10,
	}

	got := Sum(input)
	if got != 22 {
		t.Fatalf("Sum() = %d, want 22", got)
	}
}

func TestMapSlice(t *testing.T) {
	got := MapSlice([]int{1, 2, 3}, func(v int) string {
		return "n=" + strconv.Itoa(v)
	})
	want := []string{"n=1", "n=2", "n=3"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("MapSlice() = %v, want %v", got, want)
	}
}

func TestFilterSlice(t *testing.T) {
	got := FilterSlice([]int{1, 2, 3, 4, 5, 6}, func(v int) bool {
		return v%2 == 0
	})
	want := []int{2, 4, 6}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FilterSlice() = %v, want %v", got, want)
	}
}

func TestIdentity(t *testing.T) {
	if got := Identity("go"); got != "go" {
		t.Fatalf("Identity(\"go\") = %q, want %q", got, "go")
	}

	if got := Identity(123); got != 123 {
		t.Fatalf("Identity(123) = %d, want 123", got)
	}
}

func TestContains(t *testing.T) {
	if !Contains([]int{1, 2, 3}, 2) {
		t.Fatal("Contains should return true for existing value")
	}

	if Contains([]string{"go", "gen"}, "rust") {
		t.Fatal("Contains should return false for missing value")
	}
}
