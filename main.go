// SPDX-FileCopyrightText: 2023 Angelo 'Flecart' Huang <xuanqiang.huang@studio.unibo.it>
// SPDX-FileCopyrightText: 2023 Gabriele Crestanello <gabriele.crestanello@studio.unibo.it>
// SPDX-FileCopyrightText: 2023 Stefano Volpe <foxy@teapot.ovh>
// SPDX-FileCopyrightText: 2023 Eyad Issa <eyadlorenzo@gmail.com>
// SPDX-FileCopyrightText: 2024 Samuele Musiani <samu@teapot.ovh>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"log"
	"os"

	"github.com/cartabinaria/informabot/bot"
)

const tokenKey = "TOKEN"

func main() {
	token, found := os.LookupEnv(tokenKey)
	if !found {
		log.Fatalf("token not found. please set the %s environment variable",
			tokenKey)
	}

	bot.StartInformaBot(token, false)
}
