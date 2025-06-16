// SPDX-FileCopyrightText: 2023 Angelo 'Flecart' Huang <xuanqiang.huang@studio.unibo.it>
// SPDX-FileCopyrightText: 2023 Samuele Musiani <samu@teapot.ovh>
// SPDX-FileCopyrightText: 2023 Stefano Volpe <foxy@teapot.ovh>
// SPDX-FileCopyrightText: 2023 Eyad Issa <eyadlorenzo@gmail.com>
// SPDX-FileCopyrightText: 2023 bogo8liuk <lucaborghi99@gmail.com>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package utils

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode"

	tgbotapi "github.com/samuelemusiani/telegram-bot-api"
	"golang.org/x/text/unicode/norm"
)

/*
 * Has "bot" send "html" as a HTML-formatted text message in the same chat/topic
 * as "trigger". If "reply" is set to true, the new message is a reply to
 * "trigger".
 */
func SendHTML(bot *tgbotapi.BotAPI, trigger tgbotapi.Update, html string, reply bool) (tgbotapi.Message, error) {
	var msg tgbotapi.MessageConfig
	if trigger.Message.IsTopicMessage {
		msg = tgbotapi.NewThreadMessage(trigger.Message.Chat.ID,
			trigger.Message.MessageThreadID, html)
	} else {
		msg = tgbotapi.NewMessage(trigger.Message.Chat.ID,
			html)
	}
	if reply {
		msg.ReplyToMessageID = trigger.Message.MessageID
	}
	msg.ParseMode = tgbotapi.ModeHTML
	msg.DisableWebPagePreview = true
	return bot.Send(msg)
}

/*
ToSnakeCase converts a string into snake case. Useful for Telegram bots
commands.
*/
func ToSnakeCase(str string) string {
	return toLowerCaseConvention(str, '_')
}

func isSeparator(r rune) bool {
	return !unicode.IsLetter(r) && !unicode.IsDigit(r)
}

func toLowerCaseConvention(str string, delimiter rune) string {
	// normalize the string to NFD form
	normalizedStr := norm.NFD.String(strings.ToLower(strings.TrimSpace(str)))

	// remove diacritical marks from the string e.g. à -> a
	reg := regexp.MustCompile(`\p{M}`)
	normalizedStr = reg.ReplaceAllString(normalizedStr, "")

	splitted := strings.FieldsFunc(normalizedStr, isSeparator)

	// removing words before "'" character.
	for i := range splitted {
		apostropheSplit := strings.Split(splitted[i], "'")
		splitted[i] = apostropheSplit[len(apostropheSplit)-1]
	}

	return strings.Join(splitted, string(delimiter))
}

func WriteJSONFile(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	err = file.Close()
	return err
}

const BEGINNING_MONTH = time.September

/*
Returns the starting solar year of the current academic year (e.g. returns 2023
if the current academic year is 2023/24)
*/
func GetCurrentAcademicYear() int {
	now := time.Now()
	year := now.Year()
	if now.Month() >= BEGINNING_MONTH {
		return year
	} else {
		return year - 1
	}
}
