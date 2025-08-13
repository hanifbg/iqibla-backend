package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hanifbg/landing_backend/internal/model/request"
	"github.com/hanifbg/landing_backend/internal/model/response"
)

func (p *TelegramAPI) SendMessage(ctx context.Context, ChatID int64, ThreadID int64, message string) (resp *response.Message, err error) {
	reqUrl := "https://api.telegram.org/bot" + p.cfg.TeleToken + "/sendMessage"
	if ChatID == 0 {
		return nil, fmt.Errorf("chat id is empty")
	}

	reqPayload := request.SendMessage{
		ChatID:          ChatID,
		MessageThreadID: ThreadID,
		Text:            message,
		ParseMode:       "Markdown",
	}

	reqJson, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, err
	}
	reqBody := bytes.NewBuffer(reqJson)

	req, err := http.NewRequest("POST", reqUrl, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	reqs, _ := p.client.Do(req)
	responseData, err := io.ReadAll(reqs.Body)
	defer reqs.Body.Close()

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseData, &resp)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%#v", responseData)
	return
}
