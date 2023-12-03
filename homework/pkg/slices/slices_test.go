package slices

import (
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	type args struct {
		s       []any
		portion int
	}
	type testCase struct {
		name string
		args args
		want [][]any
	}
	tests := []testCase{
		{"zero-length slice", args{[]any{}, 10}, [][]any{}},
		{"nil slice", args{nil, 10}, [][]any{}},
		{"small slice", args{[]any{1, 2, 3, 4, 5, 6, 7, 8}, 10}, [][]any{{1, 2, 3, 4, 5, 6, 7, 8}}},
		{"multiple-length slice", args{[]any{1, 2, 3, 4, 5, 6, 7, 8, 9}, 3}, [][]any{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}},
		{"non-multiple-length slice", args{[]any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 4}, [][]any{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Split(tt.args.s, tt.args.portion); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Split() = %v, want %v", got, tt.want)
			}
		})
	}
}
