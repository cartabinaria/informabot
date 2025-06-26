// SPDX-FileCopyrightText: 2023 Angelo 'Flecart' Huang <xuanqiang.huang@studio.unibo.it>
// SPDX-FileCopyrightText: 2023 Samuele Musiani <samu@teapot.ovh>
// SPDX-FileCopyrightText: 2023 Stefano Volpe <foxy@teapot.ovh>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package utils

import (
	tgbotapi "github.com/samuelemusiani/telegram-bot-api"
)

func makeUnknownMember(chatConfigWithUser tgbotapi.ChatConfigWithUser) tgbotapi.ChatMember {
	return tgbotapi.ChatMember{
		User: &tgbotapi.User{
			ID:        chatConfigWithUser.UserID,
			FirstName: "???",
			LastName:  "???",
			UserName:  "???",
		},
	}
}

func GetChatMembers(bot *tgbotapi.BotAPI, chatID int64, memberIds []int64) []tgbotapi.ChatMember {
	var members []tgbotapi.ChatMember

	for _, id := range memberIds {
		chatConfigWithUser := tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: id,
		}

		getChatMemberConfig := tgbotapi.GetChatMemberConfig{
			ChatConfigWithUser: chatConfigWithUser,
		}

		member, err := bot.GetChatMember(getChatMemberConfig)
		if err != nil {
			member = makeUnknownMember(chatConfigWithUser)
		}
		members = append(members, member)
	}

	return members
}
