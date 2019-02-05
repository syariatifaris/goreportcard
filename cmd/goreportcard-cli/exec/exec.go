package exec

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/syariatifaris/goreportcard/cmd/goreportcard-cli/exec/hook"
)

//New creates new runner object
func New(c *Config) (Runner, error) {
	if !c.UseFormat {
		return &standard{
			Config: c,
		}, nil
	}
	var n hook.Notifier
	if c.HookFile != "" {
		notifier, err := buildNotifier(c.HookFile)
		if err != nil {
			return nil, err
		}
		n = notifier
	}
	return &formatted{
		Config:   c,
		notifier: n,
	}, nil
}

//Config as execution configuration
type Config struct {
	FailThres   float64
	Dir         string
	UseFormat   bool
	Verbose     bool
	FormatName  string
	HookFile    string
	HookTimeout int
}

type Grade struct {
	Rank       string `json:"rank"`
	Percentage string `json:"percentage"`
}

type ReportCard struct {
	Grade        Grade                  `json:"grade"`
	Files        int                    `json:"files"`
	Issues       int                    `json:"issues"`
	Passed       bool                   `json:"passed"`
	Threshold    float64                `json:"threshold"`
	LinterScores map[string]interface{} `json:"linter_scores"`
}

type HookConfig struct {
	Hooks []*Hook `json:"hooks"`
}

type Hook struct {
	URL     string        `json:"url"`
	Headers []*HookHeader `json:"headers"`
}

type HookHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//Runner as runner abstract
type Runner interface {
	Run(ctx context.Context) (bool, error)
}

func buildNotifier(path string) (hook.Notifier, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var config *HookConfig
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	reqs := make(map[string]*hook.RequestBody, len(config.Hooks))
	for _, h := range config.Hooks {
		headers := make(map[string]string, len(h.Headers))
		for _, header := range h.Headers {
			headers[header.Key] = header.Value
		}
		reqs[h.URL] = &hook.RequestBody{
			Headers: headers,
		}
	}
	return hook.NewNotifier(reqs), nil
}
