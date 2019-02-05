package hook

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/parnurzeal/gorequest"
)

type Response struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	HookErr    string `json:"hook_err,omitempty"`
}

type RequestBody struct {
	Headers map[string]string `json:"headers"`
}

func NewNotifier(requests map[string]*RequestBody) Notifier {
	return &notifier{
		rbs: requests,
	}
}

type Notifier interface {
	Do(ctx context.Context, report string) ([]*Response, error)
}

type notifier struct {
	rbs map[string]*RequestBody
	to  time.Duration
	mux sync.Mutex
}

func (n *notifier) Do(ctx context.Context, report string) ([]*Response, error) {
	if len(n.rbs) == 0 {
		return nil, errors.New("empty webhook")
	}
	n.mux.Lock()
	defer n.mux.Unlock()
	done := make(chan bool, 1)
	rchan := make(chan *Response, len(n.rbs))
	go func() {
		var wg sync.WaitGroup
		for url, rb := range n.rbs {
			wg.Add(1)
			go func(u string, r *RequestBody) {
				defer wg.Done()
				req := gorequest.New()
				sa := req.Post(u)
				for key, val := range r.Headers {
					sa.Set(key, val)
				}
				resp, body, errs := sa.Send(report).End()
				response := &Response{
					Message: body,
					URL:     u,
				}
				if resp != nil {
					response.StatusCode = resp.StatusCode
				}
				if response.StatusCode == 0 {
					response.StatusCode = http.StatusNotFound
				}
				if len(errs) > 0 {
					response.HookErr = errs[0].Error()
				}
				rchan <- response
			}(url, rb)
		}
		wg.Wait()
		done <- true
		close(done)
		close(rchan)
	}()
	select {
	case <-ctx.Done():
		return nil, errors.New("process timeout")
	case <-done:
		var rs []*Response
		for res := range rchan {
			rs = append(rs, res)
		}
		return rs, nil
	}
}
