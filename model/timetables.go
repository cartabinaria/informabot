// SPDX-FileCopyrightText: 2023 - 2024 Samuele Musiani <samu@teapot.ovh>
// SPDX-FileCopyrightText: 2023 Santo Cariotti <santo@dcariotti.me>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package model

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/cartabinaria/config-parser-go"
	tgbotapi "github.com/samuelemusiani/telegram-bot-api"
)

type InlineKeyboardRows [][]tgbotapi.InlineKeyboardButton

// Returns a group of button rows for a selected groups on `timetables`
func GetTimetableCoursesRows(timetables *map[string]cparser.Timetable) InlineKeyboardRows {
	rows := make([][]tgbotapi.InlineKeyboardButton, len(*timetables))

	keys := make([]string, 0)

	for key := range *timetables {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for i, callback := range keys {
		timetable := (*timetables)[callback]
		row := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(timetable.Course, fmt.Sprintf("lectures_%s", callback)))
		rows[i] = row
		i++
	}

	return rows
}

// Returns buttons which permits to choose the day for the timetable
func ChooseTimetableDay(callback_text string) InlineKeyboardRows {
	rows := make([][]tgbotapi.InlineKeyboardButton, 7)
	var weekdays = [7]string{
		"Domenica", "Lunedì", "Martedì", "Mercoledì", "Giovedì", "Venerdì", "Sabato",
	}
	var months = [12]string{
		"Dicembre", "Gennaio", "Febbraio", "Marzo", "Aprile", "Maggio", "Giugno", "Luglio", "Agosto", "Settembre", "Ottobre", "Novembre",
	}

	dt := time.Now()

	for day := 0; day < 7; day++ {
		rows[day] = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s %d %s", weekdays[dt.Weekday()%7], dt.Day(), months[dt.Month()%12]), fmt.Sprintf("%s_day_%d", callback_text, dt.Unix())))
		dt = dt.AddDate(0, 0, 1)
	}

	return rows
}

// Returns a group of buttons rows for the available years of a `course`
func GetLectureYears(callback_text string, course string) InlineKeyboardRows {
	yearsNro := 3
	// Master degrees has a duration of only 2 years
	if strings.HasPrefix(callback_text, "lectures_lm") {
		yearsNro = 2
	}
	rows := make([][]tgbotapi.InlineKeyboardButton, yearsNro)

	i := 1
	for i <= yearsNro {
		buttonText := fmt.Sprintf("%d\u00b0 anno", i)
		buttonCallback := fmt.Sprintf("%s_y_%d", callback_text, i)
		row := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(buttonText, buttonCallback))
		rows[i-1] = row

		i++
	}

	return rows
}
