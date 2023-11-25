package model

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	tgbotapi "github.com/musianisamuele/telegram-bot-api"
	"golang.org/x/exp/slices"

	"github.com/csunibo/informabot/commands"
	"github.com/csunibo/informabot/utils"
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
	noMantainerFound := true
	var answer strings.Builder
	var Ids []int64

	answer.WriteString(data.Response)

	for _, m := range Mantainers {
		Ids = append(Ids, int64(m.Id))
	}

	for i, participant := range utils.GetChatMembers(bot, message.Chat.ID, Ids) {
		if Ids[i] == participant.User.ID && participant.User.UserName != "???" {
			answer.WriteString("@" + participant.User.UserName + " ")
			noMantainerFound = false
		}
	}
	if noMantainerFound {
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
	if chatArray, ok := Groups[chatId]; ok {
		if !slices.Contains(chatArray, senderID) {
			Groups[chatId] = append(chatArray, senderID)
		}
	} else {
		Groups[chatId] = []int64{senderID}
	}
	err := SaveGroups(Groups)
	if err != nil {
		log.Printf("Error [LookingForData]: %s\n", err)
	}

	chatMembers := utils.GetChatMembers(bot, message.Chat.ID, Groups[chatId])

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
	} else if _, ok := Groups[message.Chat.ID]; !ok {
		log.Print("Info [NotLookingForData]: group empty, user not found")
		return makeResponseWithText(fmt.Sprintf(data.NotFoundError, message.Chat.Title))
	}

	var chatId = message.Chat.ID
	var senderId = message.From.ID

	var msg string
	if idx := slices.Index(Groups[chatId], senderId); idx == -1 {
		log.Print("Info [NotLookingForData]: user not found in group")
		msg = fmt.Sprintf(data.NotFoundError, chatTitle)
	} else {
		Groups[chatId] = append(Groups[chatId][:idx], Groups[chatId][idx+1:]...)
		err := SaveGroups(Groups)
		if err != nil {
			log.Printf("Error [NotLookingForData]: %s\n", err)
		}
		msg = fmt.Sprintf(data.Text, chatTitle)
	}

	return makeResponseWithText(msg)
}

func (data YearlyData) HandleBotCommand(_ *tgbotapi.BotAPI, message *tgbotapi.Message) CommandResponse {
	chatTitle := strings.ToLower(message.Chat.Title)

	// check if string contains the year number
	if strings.Contains(chatTitle, "primo") ||
		strings.Contains(chatTitle, "first") {
		return makeResponseWithNextCommand(data.Command + "1")
	} else if strings.Contains(chatTitle, "secondo") ||
		strings.Contains(chatTitle, "second") {
		return makeResponseWithNextCommand(data.Command + "2")
	} else if strings.Contains(chatTitle, "terzo") ||
		strings.Contains(chatTitle, "third") {
		return makeResponseWithNextCommand(data.Command + "3")
	} else {
		return makeResponseWithText(data.NoYear)
	}
}

func (data TodayLecturesData) HandleBotCommand(*tgbotapi.BotAPI, *tgbotapi.Message) CommandResponse {

	response, err := commands.GetTimeTable(data.Course.Type, data.Course.Name, data.Course.Year, time.Now())
	if err != nil {
		log.Printf("Error [TodayLecturesData]: %s\n", err)
		return makeResponseWithText("Bot internal Error, contact developers")
	}

	var msg string
	if response != "" {
		msg = data.Title + response
	} else {
		msg = data.FallbackText
	}

	return makeResponseWithText(msg)
}

func (data TomorrowLecturesData) HandleBotCommand(*tgbotapi.BotAPI, *tgbotapi.Message) CommandResponse {
	tomorrowTime := time.Now().AddDate(0, 0, 1)

	response, err := commands.GetTimeTable(data.Course.Type, data.Course.Name, data.Course.Year, tomorrowTime)
	if err != nil {
		log.Printf("Error [TomorrowLecturesData]: %s\n", err)
		return makeResponseWithText("Bot internal Error, contact developers")
	}

	var msg string
	if response != "" {
		msg = data.Title + response
	} else {
		msg = data.FallbackText
	}

	return makeResponseWithText(msg)
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
	var noLuckGroups = []int64{-1563447632} // NOTE: better way to handle this?

	var canLuckGroup = true

	if slices.Index(noLuckGroups, message.Chat.ID) != -1 {
		canLuckGroup = false
	}

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
