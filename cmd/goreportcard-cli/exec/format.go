package exec

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/syariatifaris/goreportcard/cmd/goreportcard-cli/exec/hook"

	"github.com/pkg/errors"

	"github.com/syariatifaris/goreportcard/check"
)

type formatted struct {
	*Config
	notifier hook.Notifier
}

func (f *formatted) Run(ctx context.Context) (bool, error) {
	result, err := check.Run(f.Dir)
	if err != nil {
		return false, errors.Wrap(err, fmt.Sprintf("unable to check dir %s", f.Dir))
	}
	avg := result.Average * 100
	pass := avg >= f.FailThres
	rc := &ReportCard{
		Grade: Grade{
			Rank:       fmt.Sprint(result.Grade),
			Percentage: fmt.Sprintf("%.1f%%", result.Average*100),
		},
		Files:     result.Files,
		Issues:    result.Issues,
		Passed:    pass,
		Threshold: f.FailThres,
	}
	ls := make(map[string]interface{}, 0)
	for _, c := range result.Checks {
		ls[c.Name] = fmt.Sprint(int64(c.Percentage*100), "%")
	}
	rc.LinterScores = ls
	bytes, err := json.MarshalIndent(rc, "", "\t")
	if err != nil {
		return false, errors.Wrap(err, "scoring failed")
	}
	fmt.Println(string(bytes))
	if f.notifier != nil {
		fmt.Println("contacting hook to send data ..")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(f.HookTimeout))
		defer cancel()
		responses, err := f.notifier.Do(ctx, string(bytes))
		if err != nil {
			return false, errors.Wrap(err, "unable to finish hook call(s)")
		}
		bytes, err := json.MarshalIndent(responses, "", "\t")
		fmt.Println("hook responses:")
		fmt.Println(string(bytes))
	}
	if !pass {
		return false, nil
	}
	return true, nil
}
