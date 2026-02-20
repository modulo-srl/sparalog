package writers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/modulo-srl/sparalog/env"
	"github.com/modulo-srl/sparalog/logs"
)

type TelegramWriter struct {
	Writer

	apiKey    string
	channelID int
}

// NewTelegramWriter returns a telegramWriter.
func NewTelegramWriter(botAPIKey string, channelID int) *TelegramWriter {
	w := TelegramWriter{
		apiKey:    botAPIKey,
		channelID: channelID,
	}

	return &w
}

func (w *TelegramWriter) Start() error {
	w.StartQueue(100, w.onQueueItem)

	return nil
}

// Write enqueue an item and returns immediately,
// or blocks while the internal queue is full.
func (w *TelegramWriter) Write(item *logs.Item) {
	w.Enqueue(item)
}

func (w *TelegramWriter) Stop() {
	w.StopQueue(3)
}

func (w *TelegramWriter) onQueueItem(i *logs.Item) error {
	var s string

	prog, host := env.Device()

	var msg string
	if i.Prefix != "" {
		msg = "[" + i.Prefix + "] " + i.Message
	} else {
		msg = i.Message
	}

	s = logs.LevelsIcons[i.Level] + " " +
		prog + "<i>[ " + host + " ]</i>" + "\n\n" +
		"<b>" + msg + "</b>\n"

	stack := i.StackTrace
	if stack != "" {
		s += "\n<pre>" + stack + "</pre>\n"
	}

	s += "\n" + env.Runtime()

	return w.sendMessage(s)
}

type telegramReq struct {
	ChatID         int    `json:"chat_id"`
	Text           string `json:"text"`
	ParseMode      string `json:"parse_mode"`
	DisablePreview bool   `json:"disable_web_page_preview"`
}

type telegramResp struct {
	Result    bool        `json:"ok"`
	Data      interface{} `json:"result,omitempty"`
	ErrorDesc string      `json:"description,omitempty"`
}

type telegramErrorResp struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func (w *TelegramWriter) sendMessage(s string) error {
	url := "https://api.telegram.org/bot" + w.apiKey + "/sendMessage"

	reqData := telegramReq{
		ChatID:         w.channelID,
		ParseMode:      "HTML",
		DisablePreview: true,
		Text:           s,
	}

	requestBody, _ := json.Marshal(&reqData)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	header := http.Header{}
	header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header = header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read response.
	if resp.StatusCode != 200 {
		// Prova a decodificare l'errore
		b, err := io.ReadAll(resp.Body)
		if err == nil {
			var respError telegramErrorResp
			err = json.Unmarshal(b, &respError)
			if err == nil {
				return fmt.Errorf("[%d] %s", respError.ErrorCode, respError.Description)
			}
		}

		return errors.New("http status: " + resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var respData telegramResp
	err = json.Unmarshal(b, &respData)
	if err != nil {
		return err
	}

	if !respData.Result {
		return errors.New("result not ok: " + respData.ErrorDesc)
	}

	return nil
}
