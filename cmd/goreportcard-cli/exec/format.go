package exec

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/syariatifaris/goreportcard/check"
)

type formatted struct {
	*Config
}

func (f *formatted) Run() error {
	result, err := check.Run(f.Dir)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to check dir %s", f.Dir))
	}
	rc := &ReportCard{
		Grade: Grade{
			Rank:       fmt.Sprint(result.Grade),
			Percentage: fmt.Sprintf("%.1f%%", result.Average*100),
		},
		Files:  result.Files,
		Issues: result.Issues,
	}
	ls := make(map[string]interface{}, 0)
	for _, c := range result.Checks {
		ls[c.Name] = fmt.Sprint(int64(c.Percentage*100), "%")
	}
	rc.LinterScores = ls
	bytes, err := json.Marshal(rc)
	if err != nil {
		return errors.Wrap(err, "scoring failed")
	}
	fmt.Println(string(bytes))
	return nil
}
