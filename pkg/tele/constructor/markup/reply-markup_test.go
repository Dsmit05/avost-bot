package markup

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v3"
)

func TestCreateInlineKeyboardMany(t *testing.T) {
	type args struct {
		maxCountWidth int
		maxCountBtn   int
		btns          []tele.InlineButton
	}
	tests := []struct {
		name string
		args args
		want []tele.ReplyMarkup
	}{
		{
			name: "case 1: 2 msg for 2 btn [1, 2] [3, 4]",
			args: args{
				maxCountWidth: 2,
				maxCountBtn:   2,
				btns: []tele.InlineButton{
					{
						Unique: "1",
						Text:   "1",
					}, {
						Unique: "2",
						Text:   "2",
					},
					{
						Unique: "3",
						Text:   "3",
					}, {
						Unique: "4",
						Text:   "4",
					},
				},
			},

			want: []tele.ReplyMarkup{
				{
					InlineKeyboard: [][]tele.InlineButton{
						{
							{
								Unique: "1",
								Text:   "1",
							},
							{
								Unique: "2",
								Text:   "2",
							},
						},
					},
				}, {
					InlineKeyboard: [][]tele.InlineButton{
						{
							{
								Unique: "3",
								Text:   "3",
							},
							{
								Unique: "4",
								Text:   "4",
							},
						},
					},
				},
			},
		},
		{
			name: "case 2: 2 msg [1, 2, 3] [4]",
			args: args{
				maxCountWidth: 3,
				maxCountBtn:   2,
				btns: []tele.InlineButton{
					{
						Unique: "1",
						Text:   "1",
					}, {
						Unique: "2",
						Text:   "2",
					},
					{
						Unique: "3",
						Text:   "3",
					}, {
						Unique: "4",
						Text:   "4",
					},
				},
			},

			want: []tele.ReplyMarkup{
				{
					InlineKeyboard: [][]tele.InlineButton{
						{
							{
								Unique: "1",
								Text:   "1",
							},
							{
								Unique: "2",
								Text:   "2",
							},
						},
					},
				}, {
					InlineKeyboard: [][]tele.InlineButton{
						{
							{
								Unique: "3",
								Text:   "3",
							},
							{
								Unique: "4",
								Text:   "4",
							},
						},
					},
				},
			},
		},
		{
			name: "case 3: 4 msg with 1 btn [1] [2] [3] [4]",
			args: args{
				maxCountWidth: 2,
				maxCountBtn:   1,
				btns: []tele.InlineButton{
					{
						Unique: "1",
						Text:   "1",
					}, {
						Unique: "2",
						Text:   "2",
					},
					{
						Unique: "3",
						Text:   "3",
					}, {
						Unique: "4",
						Text:   "4",
					},
				},
			},

			want: []tele.ReplyMarkup{
				{
					InlineKeyboard: [][]tele.InlineButton{
						{
							{
								Unique: "1",
								Text:   "1",
							},
						},
					},
				}, {
					InlineKeyboard: [][]tele.InlineButton{
						{
							{
								Unique: "2",
								Text:   "2",
							},
						},
					},
				},
				{
					InlineKeyboard: [][]tele.InlineButton{
						{
							{
								Unique: "3",
								Text:   "3",
							},
						},
					},
				},
				{
					InlineKeyboard: [][]tele.InlineButton{
						{
							{
								Unique: "4",
								Text:   "4",
							},
						},
					},
				},
			},
		},
		{
			name: "case 4: 2 msg for 2 btn [[1, 2], [3]] [[4, 5]]",
			args: args{
				maxCountWidth: 2,
				maxCountBtn:   3,
				btns: []tele.InlineButton{
					{
						Unique: "1",
						Text:   "1",
					}, {
						Unique: "2",
						Text:   "2",
					},
					{
						Unique: "3",
						Text:   "3",
					}, {
						Unique: "4",
						Text:   "4",
					}, {
						Unique: "5",
						Text:   "5",
					},
				},
			},

			want: []tele.ReplyMarkup{
				{
					InlineKeyboard: [][]tele.InlineButton{
						{
							{
								Unique: "1",
								Text:   "1",
							},
							{
								Unique: "2",
								Text:   "2",
							},
						},
						{
							{
								Unique: "3",
								Text:   "3",
							},
						},
					},
				}, {
					InlineKeyboard: [][]tele.InlineButton{
						{
							{
								Unique: "4",
								Text:   "4",
							},
							{
								Unique: "5",
								Text:   "5",
							},
						},
					},
				},
			},
		},
		{
			name: "case 5: all zero Count",
			args: args{
				maxCountWidth: 0,
				maxCountBtn:   0,
				btns: []tele.InlineButton{
					{
						Unique: "1",
						Text:   "1",
					}, {
						Unique: "2",
						Text:   "2",
					},
					{
						Unique: "3",
						Text:   "3",
					}, {
						Unique: "4",
						Text:   "4",
					},
				},
			},
			want: []tele.ReplyMarkup{},
		},
		{
			name: "case 6: empty btn",
			args: args{
				maxCountWidth: 2,
				maxCountBtn:   2,
				btns:          []tele.InlineButton{},
			},
			want: []tele.ReplyMarkup{},
		},
		{
			name: "case 7: btn is nil",
			args: args{
				maxCountWidth: 9999,
				maxCountBtn:   9999,
				btns:          nil,
			},
			want: []tele.ReplyMarkup{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateInlineKeyboardMany(tt.args.maxCountWidth, tt.args.maxCountBtn, tt.args.btns...)
			assert.Equal(t, got, tt.want)
		})
	}
}
