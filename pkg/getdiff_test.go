package pkg

import (
	"reflect"
	"testing"
)

func TestGetDiff(test *testing.T) {
	tests := []struct {
		name          string
		oldStr        string
		newStr        string
		wantAdditions []Coordinate
		wantDeletions []Coordinate
	}{
		{
			name:          "Identical strings",
			oldStr:        "abc",
			newStr:        "abc",
			wantAdditions: []Coordinate{},
			wantDeletions: []Coordinate{},
		},
		{
			name:   "Single character addition",
			oldStr: "abc",
			newStr: "abdc",
			wantAdditions: []Coordinate{
				{StartX: 2, StartY: 2, DestX: 3, DestY: 2},
			},
			wantDeletions: []Coordinate{},
		},
		{
			name:          "Single character deletion",
			oldStr:        "abdc",
			newStr:        "abc",
			wantAdditions: []Coordinate{},
			wantDeletions: []Coordinate{
				{StartX: 2, StartY: 2, DestX: 2, DestY: 3},
			},
		},
		{
			name:   "One replace (both addition and deletion)",
			oldStr: "abc",
			newStr: "adc",
			wantAdditions: []Coordinate{
				{StartX: 1, StartY: 1, DestX: 2, DestY: 1},
			},
			wantDeletions: []Coordinate{
				{StartX: 1, StartY: 1, DestX: 1, DestY: 2},
			},
		},
		{
			name:   "Empty old string",
			oldStr: "",
			newStr: "abc",
			wantAdditions: []Coordinate{
				{StartX: 0, StartY: 0, DestX: 1, DestY: 0},
				{StartX: 1, StartY: 0, DestX: 2, DestY: 0},
				{StartX: 2, StartY: 0, DestX: 3, DestY: 0},
			},
			wantDeletions: []Coordinate{},
		},
		{
			name:          "Empty new string",
			oldStr:        "abc",
			newStr:        "",
			wantAdditions: []Coordinate{},
			wantDeletions: []Coordinate{
				{StartX: 0, StartY: 0, DestX: 0, DestY: 1},
				{StartX: 0, StartY: 1, DestX: 0, DestY: 2},
				{StartX: 0, StartY: 2, DestX: 0, DestY: 3},
			},
		},
		{
			name:   "Multiple adds and deletes",
			oldStr: "abcdef",
			newStr: "abXYef",
			wantAdditions: []Coordinate{
				{StartX: 2, StartY: 2, DestX: 3, DestY: 2},
				{StartX: 3, StartY: 2, DestX: 4, DestY: 2},
			},
			wantDeletions: []Coordinate{
				{StartX: 2, StartY: 2, DestX: 2, DestY: 3},
				{StartX: 2, StartY: 3, DestX: 2, DestY: 4},
			},
		},
	}

	for _, tt := range tests {
		test.Run(tt.name, func(t *testing.T) {
			gotAdditions, gotDeletions := GetDiff(tt.oldStr, tt.newStr)
			if !reflect.DeepEqual(gotAdditions, tt.wantAdditions) {
				t.Errorf("GetDiff() additions = %v, want %v", gotAdditions, tt.wantAdditions)
			}
			if !reflect.DeepEqual(gotDeletions, tt.wantDeletions) {
				t.Errorf("GetDiff() deletions = %v, want %v", gotDeletions, tt.wantDeletions)
			}
		})
	}
}
