package writers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/modulo-srl/sparalog"
	"github.com/modulo-srl/sparalog/env"
	"github.com/modulo-srl/sparalog/writers/templates"
)

type telegramWriter struct {
	templates.Writer

	apiKey    string
	channelID int

	worker *templates.Worker
}

// NewTelegramWriter returns a telegramWriter.
func NewTelegramWriter(botAPIKey string, channelID int) (sparalog.Writer, error) {
	w := telegramWriter{
		apiKey:    botAPIKey,
		channelID: channelID,
	}

	w.worker = templates.NewWorker(&w, 100)

	return &w, nil
}

// Write enqueue an item and returns immediately,
// or blocks while the internal queue is full.
func (w *telegramWriter) Write(item sparalog.Item) sparalog.WriterError {
	w.worker.Enqueue(item)
	return nil
}

func (w *telegramWriter) Close() {
	w.worker.Close(3)
}

func (w *telegramWriter) ProcessQueueItem(item sparalog.Item) sparalog.WriterError {
	var s string

	prog, host := env.Device()

	s = sparalog.LevelsIcons[item.Level] + " " +
		prog + "<i>[ " + host + " ]</i>" + "\n\n" +
		"<b>" + item.Line + "</b>\n"

	if item.StackTrace != "" {
		s += "\n<pre>" + item.StackTrace + "</pre>\n"
	}

	s += "\n" + env.Runtime()

	err := w.sendMessage(s)
	if err != nil {
		return w.ErrorItem(err)
	}

	return nil
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

func (w *telegramWriter) sendMessage(s string) error {
	url := "https://api.telegram.org/bot" + w.apiKey + "/sendMessage"

	reqData := telegramReq{
		ChatID:         w.channelID,
		ParseMode:      "HTML",
		DisablePreview: true,
		Text:           s,
	}

	requestBody, err := json.Marshal(&reqData)

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
		return errors.New("http status: " + resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
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
