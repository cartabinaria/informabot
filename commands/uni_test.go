// SPDX-FileCopyrightText: 2023 Angelo 'Flecart' Huang <xuanqiang.huang@studio.unibo.it>
// SPDX-FileCopyrightText: 2023 Gabriele Crestanello <gabriele.crestanello@studio.unibo.it>
// SPDX-FileCopyrightText: 2023 Santo Cariotti <santo@dcariotti.me>
// SPDX-FileCopyrightText: 2023 Eyad Issa <eyadlorenzo@gmail.com>
// SPDX-FileCopyrightText: 2024 Samuele Musiani <samu@teapot.ovh>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package commands

import (
	"strings"
	"testing"
	"time"
)

func TestLezioniTime_Format(t *testing.T) {
	a := time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC)
	ti := LezioniTime(a)

	type args struct {
		format string
	}
	tests := []struct {
		name string
		tr   *LezioniTime
		args args
		want string
	}{
		{
			name: "Date formatting",
			tr:   &ti,
			args: args{
				format: "2006-01-02T15:04:05",
			},
			want: "2023-01-01T12:00:00",
		},
		{
			name: "Empty date",
			tr:   &LezioniTime{},
			args: args{
				format: "2006-01-02T15:04:05",
			},
			want: "0001-01-01T00:00:00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := tt.tr
			if got := tr.Format(tt.args.format); got != tt.want {
				t.Errorf("LezioniTime.Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLezioniTime_UnmarshalJSON(t *testing.T) {
	var lezioniTime LezioniTime
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		tr      *LezioniTime
		args    args
		wantErr bool
	}{
		{
			name: "Not an error",
			tr:   &lezioniTime,
			args: args{
				data: []byte{34, 50, 48, 50, 51, 45, 48, 51, 45, 49, 51, 84, 49, 50, 58, 48, 48, 58, 48, 48, 34},
			},
			wantErr: false,
		},
		{
			name: "Error #1",
			tr:   &lezioniTime,
			args: args{
				data: []byte("A"),
			},
			wantErr: true,
		},
		{
			name: "Error #2",
			tr:   &lezioniTime,
			args: args{
				data: []byte{},
			},
			wantErr: true,
		},
		{
			name: "Error #2",
			tr:   &lezioniTime,
			args: args{
				data: []byte{34},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := tt.tr
			if err := tr.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("LezioniTime.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetTimeTable(t *testing.T) {
	type args struct {
		courseType string
		courseName string
		year       int
		day        time.Time
		curriculum string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		error bool
	}{
		{
			name: "Weekend",
			args: args{
				courseType: "laurea",
				courseName: "informatica",
				year:       1,
				day:        time.Date(2023, 3, 11, 0, 0, 0, 0, time.UTC),
			},
			want: "",
		},
		{
			name: "Weekday",
			args: args{
				courseType: "laurea",
				courseName: "informatica",
				year:       1,
				day:        time.Date(2023, 10, 31, 0, 0, 0, 0, time.UTC),
			},
			want: `🕘`,
		},
		{
			name: "Not a valid url",
			args: args{
				courseType: "test",
				courseName: "test",
			},
			want:  "",
			error: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTimeTable(tt.args.courseType, tt.args.courseName, tt.args.curriculum, tt.args.year, tt.args.day)
			if err != nil && !tt.error {
				t.Errorf("GetTimeTable() error = %v", err)
				return
			} else if got == "" {
			} else {
				got = strings.ReplaceAll(got, " ", "")
				want := strings.ReplaceAll(tt.want, " ", "")
				want = strings.ReplaceAll(want, "\t", "")
				if !strings.Contains(got, want) {
					t.Errorf("GetTimeTable() = %v, want %v", got, want)
				}
			}
		})
	}
}
func TestWeekend(t *testing.T) {

	date := time.Date(2023, 3, 11, 0, 0, 0, 0, time.UTC)
	result, err := GetTimeTable("laurea", "informatica", "", 1, date)
	if err != nil {
		t.Fatalf("Error while getting timetable: %s", err)
	}
	if result != "" {
		t.Errorf("Expected empty string in weekend, got %s", result)
	}
}
