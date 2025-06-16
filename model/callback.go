// SPDX-FileCopyrightText: 2023 - 2024 Samuele Musiani <samu@teapot.ovh>
// SPDX-FileCopyrightText: 2023 Santo Cariotti <santo@dcariotti.me>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package model

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/samuelemusiani/telegram-bot-api"

	"github.com/cartabinaria/informabot/commands"
)

func (_ MessageData) HandleBotCallback(_bot *tgbotapi.BotAPI, _update *tgbotapi.Update, _callback_text string) {
	log.Printf("`HandleBotCallback` not defined for `MessageData`")
}

func (_ HelpData) HandleBotCallback(_bot *tgbotapi.BotAPI, _update *tgbotapi.Update, _callback_text string) {
	log.Printf("`HandleBotCallback` not defined for `HelpData`")
}

func (_ IssueData) HandleBotCallback(_bot *tgbotapi.BotAPI, _update *tgbotapi.Update, _callback_text string) {
	log.Printf("`HandleBotCallback` not defined for `IssueData`")
}

func (_ LookingForData) HandleBotCallback(_bot *tgbotapi.BotAPI, _update *tgbotapi.Update, _callback_text string) {
	log.Printf("`HandleBotCallback` not defined for `LookingForData`")
}

func (_ NotLookingForData) HandleBotCallback(_bot *tgbotapi.BotAPI, _update *tgbotapi.Update, _callback_text string) {
	log.Printf("`HandleBotCallback` not defined for `NotLookingForData`")
}

// Handle the callback for the lectures command (`/lezioni`)
// Parse the `callback_text` to check which operation it must to do:
// - If the string ends with "_today" or "_tomorrow" it returns the timetable
// - If the string just contains a "_y_<number>" it asks for today or tomorrow
// - Otherwise prints the course year of what timetable the user wants to see
func (data Lectures) HandleBotCallback(bot *tgbotapi.BotAPI, update *tgbotapi.Update, callback_text string) {
	var chatId = int64(update.CallbackQuery.Message.Chat.ID)
	var messageId = update.CallbackQuery.Message.MessageID

	if strings.Contains(callback_text, "_day_") {
		dayRegex, err := regexp.Compile(`_day_(\d+)`)
		if err != nil {
			log.Printf("Error [dayRegex]: %s\n", err)
			return
		}

		unixTime, err := strconv.ParseInt(dayRegex.FindString(callback_text)[5:], 10, 64)
		if err != nil {
			log.Printf("Error [unixTime]: %s\n", err)
			return
		}

		timeForLectures := time.Unix(unixTime, 0)

		yearRegex, err := regexp.Compile(`_y_(\d)_`)
		if err != nil {
			log.Printf("Error [yearRegex]: %s\n", err)
			return
		}

		year, err := strconv.Atoi(yearRegex.FindString(callback_text)[3:4])
		if err != nil {
			log.Printf("Error [convert to integer the year regex]: %s\n", err)
			return
		}

		timetableKey := callback_text[len("lectures_"):strings.Index(callback_text, "_y_")]

		timetable := Timetables[timetableKey]
		response, err := commands.GetTimeTable(timetable.Type, timetable.Name, timetable.Curriculum, year, timeForLectures)
		if err != nil {
			log.Printf("Error [GetTimeTable]: %s\n", err)
		}

		if response == "" {
			response = data.FallbackText
		} else {
			response = fmt.Sprintf(data.Title, timetable.Course, year, timeForLectures.Format("2006-01-02")) + "\n\n" + response
		}

		editConfig := tgbotapi.NewEditMessageText(chatId, messageId, response)
		editConfig.ParseMode = tgbotapi.ModeHTML

		_, err = bot.Send(editConfig)
		if err != nil {
			log.Printf("Error [bot.Send() for the NewEditMessageText]: %s\n", err)
		}
	} else if strings.Contains(callback_text, "_y_") {
		rows := ChooseTimetableDay(callback_text)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
		editConfig := tgbotapi.NewEditMessageReplyMarkup(chatId, messageId, keyboard)
		_, err := bot.Send(editConfig)
		if err != nil {
			log.Printf("Error [bot.Send() for the NewEditMessageReplyMarkup]: %s\n", err)
		}
	} else {
		timetableName := strings.TrimPrefix(callback_text, "lectures_")
		rows := GetLectureYears(callback_text, Timetables[timetableName].Course)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
		editConfig := tgbotapi.NewEditMessageReplyMarkup(chatId, messageId, keyboard)
		_, err := bot.Send(editConfig)
		if err != nil {
			log.Printf("Error [bot.Send() for the NewEditMessageReplyMarkup]: %s\n", err)
		}
	}
}

func (data RepresentativesData) HandleBotCallback(bot *tgbotapi.BotAPI, update *tgbotapi.Update, callback_text string) {
	var chatId = int64(update.CallbackQuery.Message.Chat.ID)
	var messageId = update.CallbackQuery.Message.MessageID

	degreeName := strings.TrimPrefix(callback_text, "representatives_")

	repData, _ := Representatives[degreeName]
	mails := repData.Representatives
	courseName := repData.Course

	var response string
	if len(mails) == 0 {
		response = fmt.Sprintf(data.FallbackText, courseName)
	} else {
		tmp := strings.Builder{}
		for _, m := range mails {
			tmp.WriteString(m)
			tmp.WriteString("\n")
		}
		response = fmt.Sprintf(data.Title, courseName) + "\n\n" + tmp.String()
	}

	editConfig := tgbotapi.NewEditMessageText(chatId, messageId, response)
	editConfig.ParseMode = tgbotapi.ModeHTML

	_, err := bot.Send(editConfig)
	if err != nil {
		log.Printf("Error [bot.Send() for the NewEditMessageText]: %s\n", err)
	}
}

func (_ ListData) HandleBotCallback(_bot *tgbotapi.BotAPI, _update *tgbotapi.Update, _callback_text string) {
	log.Printf("`HandleBotCallback` not defined for `ListData`")
}

func (_ LuckData) HandleBotCallback(_bot *tgbotapi.BotAPI, _update *tgbotapi.Update, _callback_text string) {
	log.Printf("`HandleBotCallback` not defined for `LuckData`")
}

func (_ InvalidData) HandleBotCallback(_bot *tgbotapi.BotAPI, _update *tgbotapi.Update, _callback_text string) {
	log.Printf("`HandleBotCallback` not defined for `InvalidData`")
}
