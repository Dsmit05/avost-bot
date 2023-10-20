package models

import (
	"reflect"
	"testing"
)

func TestAnimeSpecs_UnIntersection(t *testing.T) {
	tests := []struct {
		name string
		a    AnimeSpecs
		arg  AnimeSpecs
		want AnimeSpecs
	}{
		{
			name: "case 1: well check 1 new anime",
			a: AnimeSpecs{
				AnimeSpec{Id: 1},
				AnimeSpec{Id: 2},
				AnimeSpec{Id: 3},
				AnimeSpec{Id: 4},
				AnimeSpec{Id: 5},
			},
			arg: AnimeSpecs{
				AnimeSpec{Id: 1},
				AnimeSpec{Id: 2},
				AnimeSpec{Id: 3},
				AnimeSpec{Id: 4},
				AnimeSpec{Id: 7},
			},
			want: AnimeSpecs{
				AnimeSpec{Id: 7},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.UnIntersection(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnIntersection() = %v, want %v", got, tt.want)
			}
		})
	}
}
