package telegram

import (
	"WarnNotify/types"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"
)

func (t *Notify) WarnMessage(text string) error {
	for _, u := range t.users {
		if err := t.warnMessage(u, text); err != nil {
			return err
		}
	}
	return nil
}

func (t *Notify) warnMessage(uid int64, msg string) error {
	return t.pushRequest(&request{
		User: uid,
		Msg:  msg,
	})
}

func (t *Notify) WarnStructMessage(message types.StructMessage) error {
	switch message.Tp {
	case types.StructMessageTypeText:
		temp, err := template.New("").Funcs(map[string]any{
			"sub": func(a, b int64) int64 { return a - b },
			"add": func(a, b int) int { return a + b },
			"escapeMarkdown": func(v string) string {
				replacer := strings.NewReplacer(
					"_", "\\_",
					"*", "\\*",
					"[", "\\[",
					"]", "\\]",
					"(", "\\(",
					")", "\\)",
					"~", "\\~",
					"`", "\\`",
					">", "\\>",
					"#", "\\#",
					"+", "\\+",
					"-", "\\-",
					"=", "\\=",
					"|", "\\|",
					"{", "\\{",
					"}", "\\}",
					".", "\\.",
					"!", "\\!",
				)
				return replacer.Replace(v)
			},
			"escapeURL": func(v string) string {
				return url.PathEscape(v)
			},
		}).Parse(message.Template)
		if err != nil {
			return err
		}
		buf := bytes.NewBuffer(nil)
		if err := temp.Execute(buf, message.Meta); err != nil {
			return err
		}
		return t.WarnMessage(buf.String())
	default:
		return errors.New("unsupported message type")
	}
}

func (t *Notify) Watch(getMeta func() string, duration time.Duration) (context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	go func(ctx context.Context) {
		for range time.Tick(duration) {
			text := getMeta()
			if err := t.WarnMessage(text); err != nil {
				log.Printf("【WarnNotify】 send watch info failed;err(%+v)", err)
			}
		}
		log.Println("【WarnNotify】 finish watch")
	}(ctx)
	return cancel, nil
}

func (t *Notify) pushRequest(req *request) error {
	select {
	case t.queue <- req:
	default:
		return errors.New("call frequently")
	}
	return nil
}

func (t *Notify) asyncRequest() {
	for range time.Tick(time.Second) {
		select {
		case req := <-t.queue:
			if err := t.request(req); err != nil {
				log.Printf("【WarnNotify】 request failed;err(%+v)", err)
			}
		default:
		}
	}
}

const (
	// tgBotAPIEndpoint is the endpoint for all API methods,
	// with formatting for Sprintf.
	tgBotAPIEndpoint = "https://api.telegram.org/bot%s/%s"
	// tgBotFileEndpoint is the endpoint for downloading a file from PlatformTypeTelegram.
	tgBotFileEndpoint = "https://api.telegram.org/file/bot%s/%s"
)

type tgBotRsp struct {
	Ok          bool            `json:"ok"`
	Result      json.RawMessage `json:"result,omitempty"`
	ErrorCode   int             `json:"error_code,omitempty"`
	Description string          `json:"description,omitempty"`
}

func (t *Notify) request(req *request) error {
	body, err := json.Marshal(req)
	if err != nil {
		return nil
	}
	ctx, cancel := context.WithDeadline(context.TODO(), time.Now().Add(time.Second*3))
	defer cancel()
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf(tgBotAPIEndpoint, t.token, "sendMessage"), bytes.NewReader(body))
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/json")
	rsp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	data, err := io.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	rspMsg := tgBotRsp{}
	if err := json.Unmarshal(data, &rspMsg); err != nil {
		return err
	}
	if !rspMsg.Ok {
		log.Printf("【WarnNotify】 request failed;rsp(%s)", string(data))
	}
	return nil
}

func (t *Notify) close() {
	t.cancel()
}
