// SPDX-FileCopyrightText: 2023 - 2024 Omar Ayache <ayache.omar@gmail.com>
// SPDX-FileCopyrightText: 2023 - 2024 Samuele Musiani <samu@teapot.ovh>
// SPDX-FileCopyrightText: 2023 Angelo 'Flecart' Huang <xuanqiang.huang@studio.unibo.it>
// SPDX-FileCopyrightText: 2023 Gabriele Crestanello <gabriele.crestanello@studio.unibo.it>
// SPDX-FileCopyrightText: 2023 Santo Cariotti <santo@dcariotti.me>
// SPDX-FileCopyrightText: 2023 Stefano Volpe <foxy@teapot.ovh>
// SPDX-FileCopyrightText: 2023 Eyad Issa <eyadlorenzo@gmail.com>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

// In this file we define all the structs used to parse JSON files into Go
// structs

package model

import (
	tgbotapi "github.com/samuelemusiani/telegram-bot-api"
)

type DataInterface interface {
	HandleBotCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) CommandResponse
	HandleBotCallback(bot *tgbotapi.BotAPI, update *tgbotapi.Update, callback_text string)
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
	case "issue":
		data = IssueData{}
	case "lookingFor":
		data = LookingForData{}
	case "notLookingFor":
		data = NotLookingForData{}
	case "buttonsLecture":
		data = Lectures{}
	case "buttonsRepresentatives":
		data = RepresentativesData{}
	case "list":
		data = ListData{}
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
type ProjectsGroupsStruct = map[int64][]int64

type AutoReply struct {
	Text     string `json:"text"`
	IsStrict bool   `json:"strict"`
	Reply    string `json:"reply"`
}

type SettingsStruct struct {
	MainGroupsIdentifiers []string `json:"mainGroupsIdentifiers"`
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
	Slashes     bool   `json:"slashes"`
}

type IssueData struct {
	Description string `json:"description"`
	Response    string `json:"response"`
	Fallback    string `json:"fallback"`
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

type Lectures struct {
	Description  string `json:"description"`
	Title        string `json:"title"`
	FallbackText string `json:"fallbackText"`
}

type ListData struct {
	Description string     `json:"description"`
	Header      string     `json:"header"`
	Template    string     `json:"template"`
	Items       [][]string `json:"items"`
}

type LuckData struct {
	Description     string `json:"description"`
	NoLuckGroupText string `json:"noLuckGroupText"`
}

type RepresentativesData struct {
	Description  string `json:"description"`
	Title        string `json:"title"`
	FallbackText string `json:"fallbackText"`
}

type Representative struct {
	Course          string   `json:"course"`
	Representatives []string `json:"representatives"`
}

type InvalidData struct{}
