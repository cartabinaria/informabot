// SPDX-FileCopyrightText: 2023 - 2024 Omar Ayache <ayache.omar@gmail.com>
// SPDX-FileCopyrightText: 2023 - 2024 Samuele Musiani <samu@teapot.ovh>
// SPDX-FileCopyrightText: 2023 Angelo 'Flecart' Huang <xuanqiang.huang@studio.unibo.it>
// SPDX-FileCopyrightText: 2023 Santo Cariotti <santo@dcariotti.me>
// SPDX-FileCopyrightText: 2023 Stefano Volpe <foxy@teapot.ovh>
// SPDX-FileCopyrightText: 2023 Eyad Issa <eyadlorenzo@gmail.com>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

// This file contains all the global variables of the bot that are initialized
// with the start of the bot.
//
// This file should be here because it had circular imports with the Model (bot
// imported model, which imported bot in order to access the global variables,
// especially for the settings)

package model

import (
	"github.com/cartabinaria/config-parser-go"
	"log"
)

var (
	Autoreplies     []AutoReply
	Actions         []Action
	Degrees         map[string]cparser.Degree
	MemeList        []Meme
	Settings        SettingsStruct
	Teachings       map[string]cparser.Teaching
	ProjectsGroups  ProjectsGroupsStruct
	Timetables      map[string]cparser.Timetable
	Maintainers     []cparser.Maintainer
	Representatives map[string]cparser.Representative
)

func InitGlobals() {
	var err error
	Autoreplies, err = ParseAutoReplies()
	if err != nil {
		log.Fatalf("Error reading autoreply.json file: %s", err.Error())
	}

	Teachings, err = ParseTeachings()
	if err != nil {
		log.Fatalf(err.Error())
	}

	Degrees, err = ParseDegrees()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// This should be executed AFTER ParseDegrees()
	Actions, err = ParseActions()
	if err != nil {
		log.Fatalf("Error reading actions.json file: %s", err.Error())
	}

	Settings, err = ParseSettings()
	if err != nil {
		log.Fatalf("Error reading settings.json file: %s", err.Error())
	}

	MemeList, err = ParseMemeList()
	if err != nil {
		log.Fatalf("Error reading memes.json file: %s", err.Error())
	}

	ProjectsGroups, err = ParseOrCreateProjectsGroups()
	if err != nil {
		log.Fatalf("Error reading or creating groups.json file: %s", err.Error())
	}

	Timetables, err = ParseTimetables()
	if err != nil {
		log.Fatalf(err.Error())
	}

	Maintainers, err = ParseMaintainers()
	if err != nil {
		log.Fatalf("Error parsing maintainers.json file: %s", err.Error())
	}

	Representatives, err = ParseRepresentatives()
	if err != nil {
		log.Fatalf("Error parsing representatives.json file: %s", err.Error())
	}
}
