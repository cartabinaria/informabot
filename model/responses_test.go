// SPDX-FileCopyrightText: 2023 Gabriele Crestanello <gabriele.crestanello@studio.unibo.it>
// SPDX-FileCopyrightText: 2023 Samuele Musiani <samu@teapot.ovh>
// SPDX-FileCopyrightText: 2023 Santo Cariotti <santo@dcariotti.me>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package model

import (
	"reflect"
	"testing"

	tgbotapi "github.com/samuelemusiani/telegram-bot-api"
)

func TestCommandResponse_IsEmpty(t *testing.T) {
	type fields struct {
		Text        string
		NextCommand string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Empty string fields",
			fields: fields{
				Text:        "",
				NextCommand: "",
			},
			want: true,
		},
		{
			name: "Text string field",
			fields: fields{
				Text:        "randomtext",
				NextCommand: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respose := CommandResponse{
				Text:        tt.fields.Text,
				NextCommand: tt.fields.NextCommand,
			}
			if got := respose.IsEmpty(); got != tt.want {
				t.Errorf("CommandResponse.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandResponse_HasText(t *testing.T) {
	type fields struct {
		Text        string
		NextCommand string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Empty string fields",
			fields: fields{
				Text:        "",
				NextCommand: "",
			},
			want: false,
		},
		{
			name: "Text string field",
			fields: fields{
				Text:        "something",
				NextCommand: "",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := CommandResponse{
				Text:        tt.fields.Text,
				NextCommand: tt.fields.NextCommand,
			}
			if got := response.HasText(); got != tt.want {
				t.Errorf("CommandResponse.HasText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandResponse_HasNextCommand(t *testing.T) {
	type fields struct {
		Text        string
		NextCommand string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Empty string fields",
			fields: fields{
				Text:        "",
				NextCommand: "",
			},
			want: false,
		},
		{
			name: "Nextcommand string field",
			fields: fields{
				Text:        "",
				NextCommand: "next",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := CommandResponse{
				Text:        tt.fields.Text,
				NextCommand: tt.fields.NextCommand,
			}
			if got := response.HasNextCommand(); got != tt.want {
				t.Errorf("CommandResponse.HasNextCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeResponse(t *testing.T) {
	type args struct {
		text        string
		nextCommand string
	}
	tests := []struct {
		name string
		args args
		want CommandResponse
	}{
		{
			name: "Full fields",
			args: args{
				text:        "testo",
				nextCommand: "next",
			},
			want: CommandResponse{
				Text:        "testo",
				NextCommand: "next",
			},
		},
		{
			name: "Empty fields",
			args: args{
				text:        "",
				nextCommand: "",
			},
			want: CommandResponse{
				Text:        "",
				NextCommand: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeResponse(tt.args.text, tt.args.nextCommand, tgbotapi.InlineKeyboardMarkup{}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeResponseWithText(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want CommandResponse
	}{
		{
			name: "Full text",
			args: args{
				text: "testo",
			},
			want: CommandResponse{
				Text:        "testo",
				NextCommand: "",
			},
		},
		{
			name: "Empty text",
			args: args{
				text: "",
			},
			want: CommandResponse{
				Text:        "",
				NextCommand: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeResponseWithText(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeResponseWithText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeResponseWithNextCommand(t *testing.T) {
	type args struct {
		nextCommand string
	}
	tests := []struct {
		name string
		args args
		want CommandResponse
	}{
		{
			name: "Full next",
			args: args{
				nextCommand: "next",
			},
			want: CommandResponse{
				Text:        "",
				NextCommand: "next",
			},
		},
		{
			name: "Empty next",
			args: args{
				nextCommand: "",
			},
			want: CommandResponse{
				Text:        "",
				NextCommand: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeResponseWithNextCommand(tt.args.nextCommand); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeResponseWithNextCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
