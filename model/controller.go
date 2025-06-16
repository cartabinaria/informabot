// SPDX-FileCopyrightText: 2023 - 2024 Omar Ayache <102006902+Jabbar03@users.noreply.github.com>
// SPDX-FileCopyrightText: 2023 - 2024 Samuele Musiani <samu@teapot.ovh>
// SPDX-FileCopyrightText: 2023 Angelo 'Flecart' Huang <xuanqiang.huang@studio.unibo.it>
// SPDX-FileCopyrightText: 2023 Santo Cariotti <santo@dcariotti.me>
// SPDX-FileCopyrightText: 2023 Stefano Volpe <foxy@teapot.ovh>
// SPDX-FileCopyrightText: 2023 Eyad Issa <eyadlorenzo@gmail.com>
// SPDX-FileCopyrightText: 2023 bogo8liuk <lucaborghi99@gmail.com>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package model

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strings"
	"time"

	tgbotapi "github.com/samuelemusiani/telegram-bot-api"
	"golang.org/x/exp/slices"

	"github.com/cartabinaria/informabot/utils"
)

func (data MessageData) HandleBotCommand(*tgbotapi.BotAPI, *tgbotapi.Message) CommandResponse {
	return makeResponseWithText(data.Text)
}

func buildHelpLine(builder *strings.Builder, name string, description string, slashes bool) {
	if slashes {
		builder.WriteString("/")
	}
	builder.WriteString(name + " - " + description + "\n")
}

func (data HelpData) HandleBotCommand(*tgbotapi.BotAPI, *tgbotapi.Message) CommandResponse {
	answer := strings.Builder{}
	for _, action := range Actions {
		description := action.Data.GetDescription()
		if description != "" && action.Type != "course" {
			buildHelpLine(&answer, action.Name, description, data.Slashes)
		}
	}
	for command, degree := range Degrees {
		buildHelpLine(&answer, command, "Menù "+degree.Name, data.Slashes)
	}

	return makeResponseWithText(answer.String())
}

func (data IssueData) HandleBotCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) CommandResponse {
	noMaintainerFound := true
	var answer strings.Builder
	var Ids []int64

	answer.WriteString(data.Response)

	for _, m := range Maintainers {
		Ids = append(Ids, int64(m.Id))
	}

	for i, participant := range utils.GetChatMembers(bot, message.Chat.ID, Ids) {
		if Ids[i] == participant.User.ID && participant.User.UserName != "???" {
			answer.WriteString("@" + participant.User.UserName + " ")
			noMaintainerFound = false
		}
	}
	if noMaintainerFound {
		return makeResponseWithText(data.Fallback)
	}
	return makeResponseWithText(answer.String())
}

func (data LookingForData) HandleBotCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) CommandResponse {
	chatTitle := strings.ToLower(message.Chat.Title)

	if (message.Chat.Type != "group" && message.Chat.Type != "supergroup") ||
		isAMainGroup(chatTitle) {
		log.Print("Error [LookingForData]: not a group or blacklisted")
		return makeResponseWithText(data.ChatError)
	}

	var chatId = message.Chat.ID
	var senderID = message.From.ID

	log.Printf("LookingForData: %d, %d", chatId, senderID)
	if chatArray, ok := ProjectsGroups[chatId]; ok {
		if !slices.Contains(chatArray, senderID) {
			ProjectsGroups[chatId] = append(chatArray, senderID)
		}
	} else {
		ProjectsGroups[chatId] = []int64{senderID}
	}
	err := SaveProjectsGroups(ProjectsGroups)
	if err != nil {
		log.Printf("Error [LookingForData]: %s\n", err)
	}

	chatMembers := utils.GetChatMembers(bot, message.Chat.ID, ProjectsGroups[chatId])

	var resultMsg string
	// Careful: additional arguments must be passed in the right order!
	if len(chatMembers) == 1 {
		resultMsg = fmt.Sprintf(data.SingularText, message.Chat.Title)
	} else {
		resultMsg = fmt.Sprintf(data.PluralText, len(chatMembers), message.Chat.Title)
	}

	for _, member := range chatMembers {
		userLastName := ""
		if member.User.LastName != "" {
			userLastName = " " + member.User.LastName
		}
		resultMsg += fmt.Sprintf("👤 <a href='tg://user?id=%d'>%s%s</a>\n",
			member.User.ID,
			member.User.FirstName,
			userLastName)
	}

	return makeResponseWithText(resultMsg)
}

func (data NotLookingForData) HandleBotCommand(_ *tgbotapi.BotAPI, message *tgbotapi.Message) CommandResponse {

	chatTitle := strings.ToLower(message.Chat.Title)

	if (message.Chat.Type != "group" && message.Chat.Type != "supergroup") ||
		isAMainGroup(chatTitle) {
		log.Print("Error [NotLookingForData]: not a group or yearly group")
		return makeResponseWithText(data.ChatError)
	} else if _, ok := ProjectsGroups[message.Chat.ID]; !ok {
		log.Print("Info [NotLookingForData]: group empty, user not found")
		return makeResponseWithText(fmt.Sprintf(data.NotFoundError, message.Chat.Title))
	}

	var chatId = message.Chat.ID
	var senderId = message.From.ID

	var msg string
	if idx := slices.Index(ProjectsGroups[chatId], senderId); idx == -1 {
		log.Print("Info [NotLookingForData]: user not found in group")
		msg = fmt.Sprintf(data.NotFoundError, chatTitle)
	} else {
		ProjectsGroups[chatId] = append(ProjectsGroups[chatId][:idx], ProjectsGroups[chatId][idx+1:]...)
		err := SaveProjectsGroups(ProjectsGroups)
		if err != nil {
			log.Printf("Error [NotLookingForData]: %s\n", err)
		}
		msg = fmt.Sprintf(data.Text, chatTitle)
	}

	return makeResponseWithText(msg)
}

func (data Lectures) HandleBotCommand(_ *tgbotapi.BotAPI, message *tgbotapi.Message) CommandResponse {
	rows := GetTimetableCoursesRows(&Timetables)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return makeResponseWithInlineKeyboard(keyboard)
}

func (data RepresentativesData) HandleBotCommand(_ *tgbotapi.BotAPI,
	message *tgbotapi.Message) CommandResponse {
	rows := make([][]tgbotapi.InlineKeyboardButton, len(Representatives))

	// get all keys in orderd to iterate on them sorted
	keys := make([]string, 0)
	for k := range Representatives {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, callback := range keys {
		repData := Representatives[callback]
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(repData.Course,
				fmt.Sprintf("representatives_%s", callback)))
		rows[i] = row
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return makeResponseWithInlineKeyboard(keyboard)
}

func (data ListData) HandleBotCommand(*tgbotapi.BotAPI, *tgbotapi.Message) CommandResponse {
	resultText := data.Header

	for _, item := range data.Items {
		itemInterface := make([]interface{}, len(item))
		for i, v := range item {
			itemInterface[i] = v
		}
		resultText += fmt.Sprintf(data.Template, itemInterface...)
	}

	return makeResponseWithText(resultText)
}

func (data LuckData) HandleBotCommand(_ *tgbotapi.BotAPI, message *tgbotapi.Message) CommandResponse {
	var emojis = []string{"🎲", "🎯", "🏀", "⚽", "🎳", "🎰"}

	// var canLuckGroup = ((message.Chat.Type == "group" || message.Chat.Type == "supergroup") && !isAMainGroup(message.Chat.Title))
	var canLuckGroup bool = true

	var msg string
	if canLuckGroup {
		rand.NewSource(time.Now().Unix())
		emoji := emojis[rand.Intn(len(emojis))]

		msg = emoji
	} else {
		msg = data.NoLuckGroupText
	}

	return makeResponseWithText(msg)
}

func (data InvalidData) HandleBotCommand(*tgbotapi.BotAPI, *tgbotapi.Message) CommandResponse {
	log.Printf("Probably a bug in the JSON action dictionary, got invalid data in command")
	return makeResponseWithText("Bot internal Error, contact developers")
}

func isAMainGroup(name string) bool {
	for _, i := range Settings.MainGroupsIdentifiers {
		if strings.Contains(name, i) {
			return true
		}
	}

	return false
}
