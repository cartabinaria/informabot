package bot

import (
	"fmt"
	"log"
	"log/slog"
	"strings"

	tgbotapi "github.com/samuelemusiani/telegram-bot-api"
	"golang.org/x/exp/slices"

	"github.com/cartabinaria/config-parser-go"
	"github.com/cartabinaria/informabot/model"
	"github.com/cartabinaria/informabot/utils"
)

func StartInformaBot(token string, debug bool) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Error creating new bot: %s", err)
	}
	bot.Debug = debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	model.InitGlobals()

	run(bot)
}

func run(bot *tgbotapi.BotAPI) {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.CallbackQuery != nil { // first process callback queries

			callbackText := update.CallbackQuery.Data

			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, callbackText)
			if _, err := bot.Request(callback); err != nil {
				log.Printf("Error [bot.Request() for the callback]: %s\n", err)
				continue
			}

			if strings.HasPrefix(callbackText, "lectures_") {
				handleCallback(bot, &update, "lezioni", callbackText)
			} else if strings.HasPrefix(callbackText, "representatives_") {
				handleCallback(bot, &update, "rappresentanti", callbackText)
			}

			continue
		} else if update.Message != nil {
			if filterMessage(bot, update.Message) {
				continue
			} else if update.Message.IsCommand() {
				handleCommand(bot, &update)
			} else {
				handleAutoreplies(bot, &update)
			}
		} else {
			slog.Debug("ignoring unknown update", "update", update)
		}

	}
}

type handlerBehavior = func(*tgbotapi.BotAPI, *tgbotapi.Update, string) bool
type handler = struct {
	handlerBehavior
	string
}

var handlers = []handler{
	{handleAction, "action"},
	{handleDegree, "degree"},
	{handleTeaching, "teaching"},
	{handleMeme, "meme"},
	{handleUnknown, "unknown"}}

func handleUnknown(bot *tgbotapi.BotAPI, update *tgbotapi.Update, _ string) bool {
	// If the bot is in a group and the command does NOT have the recipient bot
	// nothing is done
	if !update.Message.Chat.IsPrivate() {
		commandWithAt := update.Message.CommandWithAt()
		atIndex := strings.Index(commandWithAt, "@")
		if atIndex == -1 {
			return true
		}
	}

	handleAction(bot, update, "unknown")
	return true
}

func handleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	commandName := strings.ToLower(update.Message.Command())

	// Check if the command is for me
	commandWithAt := update.Message.CommandWithAt()
	atIndex := strings.Index(commandWithAt, "@")
	if atIndex != -1 {
		forName := commandWithAt[atIndex+1:]

		if bot.Self.UserName != forName {
			return
		}
	}

	for _, h := range handlers {
		if h.handlerBehavior(bot, update, commandName) {
			log.Printf("@%s: \t%s -> %s", update.Message.From.UserName, update.Message.Text, h.string)
			return
		}
	}
}

const DOMAIN = "@unibo.it\n"

func buildEmails(emails []string) string {
	return strings.Join(emails, DOMAIN) + DOMAIN
}

func teachingToString(teaching cparser.Teaching) string {
	var b strings.Builder
	if teaching.Name != "" {
		b.WriteString(fmt.Sprintf("<b>%s</b>\n", teaching.Name))
	}
	if teaching.Website != "" {
		currentAcademicYear := fmt.Sprint(utils.GetCurrentAcademicYear())
		b.WriteString(fmt.Sprintf("<a href='https://www.unibo.it/it/didattica/insegnamenti/insegnamento/%s/%s'>Sito</a>\n",
			currentAcademicYear, teaching.Website))
		b.WriteString(fmt.Sprintf("<a href='https://www.unibo.it/it/didattica/insegnamenti/insegnamento/%s/%s/orariolezioni'>Orario</a>\n",
			currentAcademicYear, teaching.Website))
	}
	if teaching.Professors != nil {
		b.WriteString(fmt.Sprintf("Professori:\n%s", buildEmails(teaching.Professors)))
	}
	if teaching.Name != "" {
		b.WriteString(fmt.Sprintf("<a href='https://risorse.students.cs.unibo.it/%s/'>📚 Risorse (istanza principale)</a>\n", teaching.Url))
		b.WriteString(fmt.Sprintf("<a href='https://dynamik.vercel.app/%s/'>📚 Risorse (istanza di riserva 1)</a>\n", teaching.Url))
		b.WriteString(fmt.Sprintf("<a href='https://github.com/cartabinaria/%s/'>📂 Repository GitHub delle risorse</a>\n", teaching.Url))
	}
	if teaching.Chat != "" {
		b.WriteString(fmt.Sprintf("<a href='%s'>👥 Gruppo Studenti</a>\n", teaching.Chat))
	}
	return b.String()
}

func handleTeaching(bot *tgbotapi.BotAPI, update *tgbotapi.Update, teachingName string) bool {
	teaching, ok := model.Teachings[teachingName]
	if !ok {
		return false
	}
	utils.SendHTML(bot, *update, teachingToString(teaching), false)
	return true
}

func degreeToTeaching(degree cparser.Degree) cparser.Teaching {
	return cparser.Teaching{
		Name: degree.Name,
		Url:  degree.Id,
		Chat: degree.Chat,
	}
}

func degreeToString(degree cparser.Degree) string {
	if len(degree.Teachings) == 0 {
		return teachingToString(degreeToTeaching(degree))
	}
	var b strings.Builder
	// header
	if degree.Icon != "" || degree.Name != "" || degree.Chat != "" {
		b.WriteString("<b>")
		elements := []string{}
		if degree.Icon != "" {
			elements = append(elements, degree.Icon)
		}
		if degree.Name != "" {
			elements = append(elements, degree.Name)
		}
		if degree.Chat != "" {
			elements = append(elements, fmt.Sprintf("(<a href='%s'>👥 Gruppo</a>)", degree.Chat))
		}
		b.WriteString(strings.Join(elements, " "))
		b.WriteString("</b>\n")
	}

	for y := range 5 {
		// For not classified teachings
		t1 := cparser.GetYearMandatoryTeachingsFromDegree(degree, int64(y))
		t2 := cparser.GetYearElectivesTeachingsFromDegree(degree, int64(y))
		if y != 0 && (len(t1) > 0 || len(t2) > 0) {
			// header
			b.WriteString(fmt.Sprintf("%d", y))

			var chat string
			for _, tmp := range degree.Years {
				if tmp.Year == int64(y) {
					chat = tmp.Chat
				}
			}

			if chat != "" {
				b.WriteString(fmt.Sprintf(" (<a href='%s'>👥 Gruppo</a>)", chat))
			}
			b.WriteString("\n")
		}
		for _, t := range append(t1, t2...) {
			b.WriteString(fmt.Sprintf("/%s\n", t.Name))
		}
	}
	return b.String()
}

func handleDegree(bot *tgbotapi.BotAPI, update *tgbotapi.Update, degreeId string) bool {
	degree, ok := model.Degrees[degreeId]
	if !ok {
		return false
	}
	utils.SendHTML(bot, *update, degreeToString(degree), false)
	return true
}

func handleMeme(bot *tgbotapi.BotAPI, update *tgbotapi.Update, memeName string) bool {
	memeIndex := slices.IndexFunc(model.MemeList, func(meme model.Meme) bool {
		return strings.ToLower(meme.Name) == memeName
	})

	if memeIndex != -1 {
		utils.SendHTML(bot, *update, model.MemeList[memeIndex].Text, false)
		return true
	}
	return false
}

// executes a given command in the command list, given its index
// if invalid index, does nothing
func executeCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update, commandIndex int) {
	if commandIndex >= 0 && commandIndex < len(model.Actions) {
		newCommand := model.Actions[commandIndex].Data.HandleBotCommand(bot, update.Message)

		if newCommand.HasText() {
			utils.SendHTML(bot, *update, newCommand.Text, false)
		}

		if newCommand.HasNextCommand() {
			handleAction(bot, update, newCommand.NextCommand)
		}

		if newCommand.HasRows() {

			var msg tgbotapi.MessageConfig
			if update.Message.IsTopicMessage {
				msg = tgbotapi.NewThreadMessage(update.Message.Chat.ID, update.Message.MessageThreadID, update.Message.Text)
			} else {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			}
			msg.ReplyMarkup = newCommand.Rows
			if _, err := bot.Send(msg); err != nil {
				errorMsg := "Error sending data"
				if update.Message.IsTopicMessage {
					msg = tgbotapi.NewThreadMessage(update.Message.Chat.ID, update.Message.MessageThreadID, errorMsg)
				} else {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, errorMsg)
				}
				bot.Send(msg)
			}
		}
	}
}

// executes a given command in the command list, given its name
// @return true if command was found, false otherwise
func handleAction(bot *tgbotapi.BotAPI, update *tgbotapi.Update, commandName string) bool {
	idx := slices.IndexFunc(model.Actions, func(action model.Action) bool {
		return action.Name == commandName
	})

	if idx != -1 {
		executeCommand(bot, update, idx)
		return true
	}

	return false
}

// Handle a callback searching a the good action
func handleCallback(bot *tgbotapi.BotAPI, update *tgbotapi.Update, commandName string, callback_text string) bool {
	idx := slices.IndexFunc(model.Actions, func(action model.Action) bool {
		return action.Name == commandName
	})

	if idx != -1 {
		model.Actions[idx].Data.HandleBotCallback(bot, update, callback_text)

		return true
	}

	return false
}

func filterMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) bool {
	if message.Dice != nil {
		// msg := tgbotapi.NewMessage(message.Chat.ID, "Found a dice")
		// bot.Send(msg)
		return true
	}
	return false
}

func handleAutoreplies(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	txt := strings.ToLower(update.Message.Text)
	for _, a := range model.Autoreplies {
		autoTxt := strings.ToLower(a.Text)
		if a.IsStrict {
			if strings.Contains(txt, autoTxt) {
				utils.SendHTML(bot, *update, a.Reply, true)
			}
		} else {
			sAutoTxt := strings.Split(autoTxt, " ")

			for i := range sAutoTxt {
				if !strings.Contains(txt, sAutoTxt[i]) {
					return
				}
			}
			utils.SendHTML(bot, *update, a.Reply, true)
		}
	}
}
