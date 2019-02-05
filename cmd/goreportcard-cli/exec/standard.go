package exec

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/syariatifaris/goreportcard/check"
)

type standard struct {
	*Config
}

func (s *standard) Run(ctx context.Context) error {
	result, err := check.Run(s.Dir)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to check dir %s", s.Dir))
	}
	showReport(&result)
	for _, c := range result.Checks {
		fmt.Printf("%s: %d%%\n", c.Name, int64(c.Percentage*100))
		if s.Verbose && len(c.FileSummaries) > 0 {
			for _, f := range c.FileSummaries {
				fmt.Printf("\t%s\n", f.Filename)
				for _, e := range f.Errors {
					fmt.Printf("\t\tLine %d: %s\n", e.LineNumber, e.ErrorString)
				}
			}
		}
	}
	if avg := result.Average * 100; avg < s.FailThres {
		return fmt.Errorf("examination not pass at %f, try again", avg)
	}
	return nil
}

func showReport(result *check.ChecksResult) {
	fmt.Printf("Grade: %s (%.1f%%)\n", result.Grade, result.Average*100)
	fmt.Printf("Files: %d\n", result.Files)
	fmt.Printf("Issues: %d\n", result.Issues)
}
