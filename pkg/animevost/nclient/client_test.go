//go:build test_api
// +build test_api

package nclient

import (
	"fmt"
	"testing"
)

var cli, _ = NewClient("https://api.animevost.org/v1")

func TestClient_GetPage(t *testing.T) {
	type args struct {
		page     int
		quantity int
	}

	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "case 1: Well get 1 anime",
			args: args{
				page:     1444,
				quantity: 1,
			},
			want:    187,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.GetPage(tt.args.page, tt.args.quantity)
			fmt.Printf("%+v", got)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != 1 {
				t.Errorf("GetPage() len models.AnimeSpecs not 1, got = %v, want %v", got, tt.want)
			}

			if got[0].Id != tt.want {
				t.Errorf("GetPage() got = %v, want id: %v", got, tt.want)
			}
		})
	}
}

func TestClient_SearchForName(t *testing.T) {
	type args struct {
		name string
	}

	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "case 1: Well check, search first anime",
			args: args{
				name: "Блич",
			},
			want:    1292,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.SearchForName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchForName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) == 0 {
				t.Errorf("SearchForName() len models.AnimeSpecs not 1, got = %v, want %v", got, tt.want)
			}

			if got[0].Id != tt.want {
				t.Errorf("SearchForName() got = %v, want id: %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetPlayList(t *testing.T) {
	type args struct {
		id int
	}

	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "case 1: Well check, with len",
			args: args{
				id: 55,
			},
			want:    13,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.GetPlayList(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPlayList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != tt.want {
				t.Errorf("SearchForName() len models.Playlists not want %v, got = %v", tt.want, got)
			}
		})
	}
}

func TestClient_GetInfo(t *testing.T) {
	type args struct {
		id int
	}

	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "case 1: Well check, correct id",
			args: args{
				id: 2829,
			},
			want:    2829,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.GetInfo(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.Id != tt.want {
				t.Errorf("GetInfo() got = %v, want id: %v", got, tt.want)
			}
		})
	}
}
