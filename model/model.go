// In this file we define all the structs used to parse JSON files into Go
// structs

package model

import (
	tgbotapi "github.com/musianisamuele/telegram-bot-api"
)

type DataInterface interface {
	HandleBotCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) CommandResponse
	GetDescription() string
}

// GetActionFromType returns an Action struct with the right DataInterface,
// inferred from the commandType string
func GetActionFromType(name string, commandType string) Action {
	var data DataInterface
	switch commandType {
	case "message":
		data = MessageData{}
	case "help":
		data = HelpData{}
	case "lookingFor":
		data = LookingForData{}
	case "notLookingFor":
		data = NotLookingForData{}
	case "yearly":
		data = YearlyData{}
	case "todayLectures":
		data = TodayLecturesData{}
	case "tomorrowLectures":
		data = TomorrowLecturesData{}
	case "list":
		data = ListData{}
	case "course":
		data = CourseData{}
	case "luck":
		data = LuckData{}
	default:
		data = InvalidData{}
	}

	return Action{
		Name: name,
		Type: commandType,
		Data: data,
	}
}

// SECTION GLOBAL JSON STRUCTS
type GroupsStruct = map[int64][]int64

type AutoReply struct {
	Text  string `json:"text"`
	Reply string `json:"reply"`
}

type SettingsStruct struct {
	LookingForBlackList []int64 `json:"lookingForBlackList"`
	GeneralGroups       []int64 `json:"generalGroups"`
}

type Meme struct {
	Name string
	Text string
}

type Action struct {
	Name string
	Type string        `json:"type"`
	Data DataInterface `json:"data"`
}

// SECTION ACTION STRUCTS DATA
type MessageData struct {
	Text        string `json:"text"`
	Description string `json:"description"`
}

type HelpData struct {
	Description string `json:"description"`
}

type LookingForData struct {
	Description  string `json:"description"`
	SingularText string `json:"singularText"`
	PluralText   string `json:"pluralText"`
	ChatError    string `json:"chatError"`
}

type NotLookingForData struct {
	Description   string `json:"description"`
	Text          string `json:"text"`
	ChatError     string `json:"chatError"`
	NotFoundError string `json:"notFoundError"`
}

type YearlyData struct {
	Description string `json:"description"`
	Command     string `json:"command"`
	NoYear      string `json:"noYear"`
}

type CourseId struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Year int    `json:"year"`
}

type TodayLecturesData struct {
	Description  string   `json:"description"`
	Course       CourseId `json:"course"`
	Title        string   `json:"title"`
	FallbackText string   `json:"fallbackText"`
}

type TomorrowLecturesData TodayLecturesData

type ListData struct {
	Description string     `json:"description"`
	Header      string     `json:"header"`
	Template    string     `json:"template"`
	Items       [][]string `json:"items"`
}

type CourseData struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Virtuale    string   `json:"virtuale"`
	Teams       string   `json:"teams"`
	Website     string   `json:"website"`
	Professors  []string `json:"professors"`
	Telegram    string   `json:"telegram"`
}

type LuckData struct {
	Description     string `json:"description"`
	NoLuckGroupText string `json:"noLuckGroupText"`
}

type InvalidData struct{}
