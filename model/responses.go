// SPDX-FileCopyrightText: 2023 Angelo 'Flecart' Huang <xuanqiang.huang@studio.unibo.it>
// SPDX-FileCopyrightText: 2023 Samuele Musiani <samu@teapot.ovh>
// SPDX-FileCopyrightText: 2023 Santo Cariotti <santo@dcariotti.me>
// SPDX-FileCopyrightText: 2023 Eyad Issa <eyadlorenzo@gmail.com>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package model contains the structure and functions used to pass information
// between the computations of the bot and the driver code, in bot.go
package model

import tgbotapi "github.com/samuelemusiani/telegram-bot-api"

// CommandResponse is returned by the command handler, it contains information
// about the command computation.
type CommandResponse struct {
	Text        string
	NextCommand string
	Rows        tgbotapi.InlineKeyboardMarkup
}

// makeResponse creates a CommandResponse with the given text and nextCommand
func makeResponse(text string, nextCommand string, rows tgbotapi.InlineKeyboardMarkup) CommandResponse {
	return CommandResponse{
		Text:        text,
		NextCommand: nextCommand,
		Rows:        rows,
	}
}

// makeResponseWithText creates a CommandResponse with the given text (and no nextCommand)
func makeResponseWithText(text string) CommandResponse {
	return makeResponse(text, "", tgbotapi.InlineKeyboardMarkup{})
}

// makeResponseWithNextCommand creates a CommandResponse with the given nextCommand (and no text)
func makeResponseWithNextCommand(nextCommand string) CommandResponse {
	return makeResponse("", nextCommand, tgbotapi.InlineKeyboardMarkup{})
}

// makeResponseWithInlineKeyboard creates a CommandResponse with the given array
// of elements for the keyboard array.
func makeResponseWithInlineKeyboard(rows tgbotapi.InlineKeyboardMarkup) CommandResponse {
	return makeResponse("", "", rows)
}

// IsEmpty returns true if the CommandResponse has no text and no nextCommand
func (r CommandResponse) IsEmpty() bool {
	return r.Text == "" && r.NextCommand == ""
}

// HasText returns true if the CommandResponse has some text
func (r CommandResponse) HasText() bool {
	return r.Text != ""
}

// HasNextCommand returns true if the CommandResponse has some nextCommand
func (r CommandResponse) HasNextCommand() bool {
	return r.NextCommand != ""
}

// HasButtonRows returns true if the CommandResponse has some rows
func (r CommandResponse) HasRows() bool {
	return len(r.Rows.InlineKeyboard) > 0
}
