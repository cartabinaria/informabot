// SPDX-FileCopyrightText: 2026 Informabot contributors
// SPDX-License-Identifier: AGPL-3.0-or-later

package bot

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	cparser "github.com/cartabinaria/config-parser-go"
	"github.com/cartabinaria/informabot/model"
	tgbotapi "github.com/samuelemusiani/telegram-bot-api"
)

const syntheticPrivate = "SYNTHETIC_TOKEN synthetic_user 987654321 https://synthetic.invalid/private"
const deliveryAPIError = `{"ok":false,"error_code":403,"description":"` + syntheticPrivate + `"}`

var deliveryTransportError = &url.Error{Op: "Post", URL: "https://synthetic.invalid/botSYNTHETIC_TOKEN/sendMessage", Err: errors.New(syntheticPrivate)}

type deliveryResponse struct {
	body string
	err  error
}

type deliveryRequest struct {
	method string
	params url.Values
}

// Every request is answered in memory, never delegated to a network transport.
type deliveryClient struct {
	t         *testing.T
	responses []deliveryResponse
	requests  []deliveryRequest
}

func (c *deliveryClient) Do(req *http.Request) (*http.Response, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		c.t.Fatal(err)
	}
	params, err := url.ParseQuery(string(body))
	if err != nil {
		c.t.Fatal(err)
	}
	index := len(c.requests)
	c.requests = append(c.requests, deliveryRequest{method: req.URL.Path, params: params})
	response := deliveryResponse{}
	if index < len(c.responses) {
		response = c.responses[index]
	}
	if response.err != nil {
		return nil, response.err
	}
	if response.body == "" {
		response.body = `{"ok":true,"result":{"message_id":1,"chat":{"id":1,"type":"private"}}}`
	}
	return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(response.body)), Header: make(http.Header)}, nil
}

type deliveryData struct {
	response  model.CommandResponse
	callbacks *int
}

func (d deliveryData) HandleBotCommand(*tgbotapi.BotAPI, *tgbotapi.Message) model.CommandResponse {
	return d.response
}

func (d deliveryData) HandleBotCallback(*tgbotapi.BotAPI, *tgbotapi.Update, string) {
	*d.callbacks++
}

func (d deliveryData) GetDescription() string {
	return "synthetic"
}

func deliveryFixture(t *testing.T, responses ...deliveryResponse) (*tgbotapi.BotAPI, *deliveryClient, *bytes.Buffer) {
	t.Helper()
	actions, degrees, teachings := model.Actions, model.Degrees, model.Teachings
	memes, autoreplies := model.MemeList, model.Autoreplies
	model.Actions = nil
	model.Degrees = nil
	model.Teachings = nil
	model.MemeList = nil
	model.Autoreplies = nil
	t.Cleanup(func() {
		model.Actions, model.Degrees, model.Teachings = actions, degrees, teachings
		model.MemeList, model.Autoreplies = memes, autoreplies
	})
	output := new(bytes.Buffer)
	writer, flags, prefix := log.Writer(), log.Flags(), log.Prefix()
	log.SetOutput(output)
	log.SetFlags(0)
	log.SetPrefix("")
	t.Cleanup(func() {
		log.SetOutput(writer)
		log.SetFlags(flags)
		log.SetPrefix(prefix)
	})
	client := &deliveryClient{t: t, responses: responses}
	api := &tgbotapi.BotAPI{Client: client, Self: tgbotapi.User{UserName: "synthetic_bot"}}
	api.SetAPIEndpoint("https://offline.invalid/bot%s/%s")
	return api, client, output
}

func deliveryUpdate(text string) tgbotapi.Update {
	message := &tgbotapi.Message{
		MessageID: 987654321,
		Chat:      &tgbotapi.Chat{ID: 987654321, Type: "private"},
		From:      &tgbotapi.User{ID: 987654321, UserName: "synthetic_user"},
		Text:      text,
	}
	if strings.HasPrefix(text, "/") {
		message.Entities = []tgbotapi.MessageEntity{{Type: "bot_command", Length: len(strings.Fields(text)[0])}}
	}
	return tgbotapi.Update{Message: message}
}

func assertDeliveryLog(t *testing.T, output *bytes.Buffer, want string) {
	t.Helper()
	var lines []string
	for _, line := range strings.Split(output.String(), "\n") {
		if strings.Contains(line, "Error") {
			lines = append(lines, line)
		}
	}
	got := strings.Join(lines, "\n")
	if got != want {
		t.Errorf("delivery log = %q, want %q", got, want)
	}
	// The existing command audit log is out of scope; only new error lines
	// must exclude all synthetic identity, URL, token and message canaries.
	for _, private := range []string{"SYNTHETIC_TOKEN", "synthetic_user", "987654321", "https://", "/fixture", "synthetic reply", "lectures_private"} {
		if strings.Contains(got, private) {
			t.Errorf("delivery log leaked %q", private)
		}
	}
}

func TestDeliveryCommands(t *testing.T) {
	for _, route := range []string{"action", "degree", "teaching", "meme", "unknown"} {
		for _, failure := range []string{"success", "api", "transport", "decode"} {
			t.Run(route+"/"+failure, func(t *testing.T) {
				response := deliveryResponse{}
				wantLog := ""
				if failure != "success" {
					class := "request"
					switch failure {
					case "api":
						response.body, class = deliveryAPIError, "telegram_api"
					case "transport":
						response.err = deliveryTransportError
					case "decode":
						response.body = `{"ok":true,"result":"` + syntheticPrivate + `"}`
					}
					wantLog = "Error delivering command: " + class
				}
				api, client, output := deliveryFixture(t, response)
				command := route
				model.Actions = []model.Action{{Name: "unknown", Data: model.MessageData{Text: "synthetic unknown"}}}
				want := "synthetic reply"
				switch route {
				case "action":
					command = "fixture"
					model.Actions = append(model.Actions, model.Action{Name: command, Data: model.MessageData{Text: want}})
					// A matched action must not fall through to another handler on failure.
					model.Degrees = map[string]cparser.Degree{command: {Chat: "should not send"}}
				case "degree":
					model.Degrees = map[string]cparser.Degree{command: {Chat: "https://synthetic.invalid/degree"}}
					model.Teachings = map[string]cparser.Teaching{command: {Chat: "should not send"}}
					want = "<a href='https://synthetic.invalid/degree'>👥 Gruppo Studenti</a>\n"
				case "teaching":
					model.Teachings = map[string]cparser.Teaching{command: {Chat: "https://synthetic.invalid/teaching"}}
					model.MemeList = []model.Meme{{Name: command, Text: "should not send"}}
					want = "<a href='https://synthetic.invalid/teaching'>👥 Gruppo Studenti</a>\n"
				case "meme":
					model.MemeList = []model.Meme{{Name: "MeMe", Text: want}}
				case "unknown":
					command, want = "missing", "synthetic unknown"
				}
				update := deliveryUpdate("/" + command)
				handleUpdate(api, &update)
				if len(client.requests) != 1 {
					t.Fatalf("requests = %d, want 1", len(client.requests))
				}
				params := client.requests[0].params
				if params.Get("text") != want || params.Get("parse_mode") != "HTML" || params.Get("disable_web_page_preview") != "true" || params.Get("reply_to_message_id") != "" {
					t.Errorf("unexpected message: %v", params)
				}
				assertDeliveryLog(t, output, wantLog)
			})
		}
	}
}

func TestDeliveryKeyboardAndChain(t *testing.T) {
	for _, thread := range []bool{false, true} {
		for _, failure := range []string{"success", "keyboard", "fallback", "all"} {
			t.Run(failure+map[bool]string{false: "/chat", true: "/thread"}[thread], func(t *testing.T) {
				responses := []deliveryResponse{{}, {}, {}}
				wantLog := ""
				wantTexts := []string{"first", "next", "/fixture"}
				if failure != "success" {
					responses[2].body = deliveryAPIError
					wantTexts = append(wantTexts, "Error sending data")
					wantLog = "Error delivering command: telegram_api"
				}
				if failure == "fallback" || failure == "all" {
					responses = append(responses, deliveryResponse{err: deliveryTransportError})
					wantLog += ",request"
				}
				if failure == "all" {
					responses[0].err = deliveryTransportError
					responses[1].body = deliveryAPIError
					wantLog = "Error delivering command: request,telegram_api,telegram_api,request"
				}
				api, client, output := deliveryFixture(t, responses...)
				rows := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("synthetic", "synthetic_callback")))
				model.Actions = []model.Action{
					{Name: "fixture", Data: deliveryData{response: model.CommandResponse{Text: "first", NextCommand: "next", Rows: rows}}},
					{Name: "next", Data: model.MessageData{Text: "next"}},
				}
				update := deliveryUpdate("/fixture")
				update.Message.IsTopicMessage, update.Message.MessageThreadID = thread, 42
				handleUpdate(api, &update)
				var gotTexts []string
				for i, request := range client.requests {
					gotTexts = append(gotTexts, request.params.Get("text"))
					wantThread := ""
					if thread {
						wantThread = "42"
					}
					if request.params.Get("message_thread_id") != wantThread {
						t.Errorf("request %d lost thread: %v", i, request.params)
					}
					wantMarkup := ""
					if i == 2 {
						wantMarkup = `{"inline_keyboard":[[{"text":"synthetic","callback_data":"synthetic_callback"}]]}`
					}
					if request.params.Get("reply_markup") != wantMarkup {
						t.Errorf("request %d markup = %q, want %q", i, request.params.Get("reply_markup"), wantMarkup)
					}
				}
				if !reflect.DeepEqual(gotTexts, wantTexts) {
					t.Errorf("send order = %v, want %v", gotTexts, wantTexts)
				}
				assertDeliveryLog(t, output, wantLog)
			})
		}
	}
}

func TestDeliveryAutoreplies(t *testing.T) {
	for _, failure := range []bool{false, true} {
		for _, earlyReturn := range []bool{false, true} {
			t.Run(map[bool]string{false: "success", true: "failure"}[failure]+map[bool]string{false: "/all", true: "/early_return"}[earlyReturn], func(t *testing.T) {
				responses := []deliveryResponse{{}, {}}
				wantLog := ""
				if failure {
					responses = []deliveryResponse{{body: deliveryAPIError}, {err: deliveryTransportError}}
					wantLog = "Error delivering autoreply: telegram_api"
					if !earlyReturn {
						wantLog += ",request"
					}
				}
				api, client, output := deliveryFixture(t, responses...)
				model.Autoreplies = []model.AutoReply{
					{Text: "HELLO", IsStrict: true, Reply: "first"},
					{Text: "world hello", Reply: "second"},
					{Text: "world", IsStrict: true, Reply: "third"},
					{Text: "absent", IsStrict: true, Reply: "not sent"},
				}
				if earlyReturn {
					model.Autoreplies[1].Text = "missing word"
					model.Autoreplies = append(model.Autoreplies, model.AutoReply{Text: "hello", IsStrict: true, Reply: "must stay skipped"})
				}
				update := deliveryUpdate("Hello World")
				update.Message.IsTopicMessage, update.Message.MessageThreadID = true, 42
				handleUpdate(api, &update)
				wantTexts := []string{"first", "second", "third"}
				if earlyReturn {
					wantTexts = wantTexts[:1]
				}
				var gotTexts []string
				for _, request := range client.requests {
					gotTexts = append(gotTexts, request.params.Get("text"))
					if request.params.Get("reply_to_message_id") != "987654321" || request.params.Get("message_thread_id") != "42" || request.params.Get("parse_mode") != "HTML" {
						t.Errorf("unexpected autoreply: %v", request.params)
					}
				}
				if !reflect.DeepEqual(gotTexts, wantTexts) {
					t.Errorf("autoreplies = %v, want %v", gotTexts, wantTexts)
				}
				assertDeliveryLog(t, output, wantLog)
			})
		}
	}
}

func TestDeliveryRouting(t *testing.T) {
	for _, route := range []string{"other_bot", "unknown_group", "missing_unknown", "noncommand", "dice", "unknown_update"} {
		t.Run(route, func(t *testing.T) {
			api, client, output := deliveryFixture(t)
			model.Actions = []model.Action{{Name: "fixture", Data: model.MessageData{Text: "not sent"}}}
			update := deliveryUpdate("/fixture")
			switch route {
			case "other_bot":
				update = deliveryUpdate("/fixture@other_bot")
			case "unknown_group":
				update = deliveryUpdate("/missing")
				update.Message.Chat.Type = "group"
			case "missing_unknown":
				update = deliveryUpdate("/missing")
			case "noncommand":
				update = deliveryUpdate("fixture")
				model.Autoreplies = []model.AutoReply{{Text: "absent", IsStrict: true, Reply: "not sent"}}
			case "dice":
				update.Message.Dice = &tgbotapi.Dice{Emoji: "🎲", Value: 1}
			case "unknown_update":
				update = tgbotapi.Update{}
			}
			handleUpdate(api, &update)
			if len(client.requests) != 0 {
				t.Errorf("unexpected requests: %v", client.requests)
			}
			assertDeliveryLog(t, output, "")
		})
	}
}

func TestDeliveryErrorContract(t *testing.T) {
	failures := []error{errors.New("synthetic first"), errors.New("synthetic next"), errors.New("synthetic keyboard"), errors.New("synthetic fallback")}
	api, client, output := deliveryFixture(t,
		deliveryResponse{err: failures[0]}, deliveryResponse{err: failures[1]},
		deliveryResponse{err: failures[2]}, deliveryResponse{err: failures[3]})
	rows := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("synthetic", "synthetic")))
	model.Actions = []model.Action{
		{Name: "fixture", Data: deliveryData{response: model.CommandResponse{Text: "first", NextCommand: "next", Rows: rows}}},
		{Name: "next", Data: model.MessageData{Text: "next"}},
	}
	update := deliveryUpdate("/fixture")
	matched, err := handleAction(api, &update, "fixture")
	if !matched {
		t.Fatal("failed delivery must still match the command")
	}
	for i, failure := range failures {
		if !errors.Is(err, failure) {
			t.Errorf("lost delivery error %d", i)
		}
	}
	if len(client.requests) != len(failures) || output.Len() != 0 {
		t.Errorf("inner handler retried or logged: requests = %d, log = %q", len(client.requests), output.String())
	}
	if matched, err := handleAction(api, &update, "missing"); matched || err != nil {
		t.Errorf("missing action: matched = %v, error = %v", matched, err)
	}
	for _, index := range []int{-1, len(model.Actions)} {
		if err := executeCommand(api, &update, index); err != nil {
			t.Errorf("invalid action index: %v", err)
		}
	}
	if len(client.requests) != len(failures) {
		t.Fatal("missing or invalid action sent a message")
	}
}

func TestDeliveryCallbackAck(t *testing.T) {
	for _, failure := range []string{"success", "api", "transport"} {
		t.Run(failure, func(t *testing.T) {
			response := deliveryResponse{}
			wantLog := ""
			if failure == "api" {
				response.body, wantLog = deliveryAPIError, "Error delivering callback acknowledgement: telegram_api"
			} else if failure == "transport" {
				response.err, wantLog = deliveryTransportError, "Error delivering callback acknowledgement: request"
			}
			api, client, output := deliveryFixture(t, response)
			callbacks := 0
			model.Actions = []model.Action{{Name: "lezioni", Data: deliveryData{callbacks: &callbacks}}}
			update := deliveryUpdate("/fixture")
			update.CallbackQuery = &tgbotapi.CallbackQuery{ID: "987654321", Data: "lectures_private"}
			handleUpdate(api, &update)
			wantCallbacks := 0
			if failure == "success" {
				wantCallbacks = 1
			}
			if callbacks != wantCallbacks || len(client.requests) != 1 {
				t.Fatalf("callbacks = %d, requests = %d", callbacks, len(client.requests))
			}
			request := client.requests[0]
			if request.method != "/bot/answerCallbackQuery" || request.params.Get("callback_query_id") != "987654321" || request.params.Get("text") != "lectures_private" {
				t.Errorf("unexpected acknowledgement: %v", request)
			}
			assertDeliveryLog(t, output, wantLog)
		})
	}
}
