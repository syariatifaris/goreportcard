package exec

//New creates new runner object
func New(c *Config) Runner {
	if !c.UseFormat {
		return &standard{
			Config: c,
		}
	}
	return &formatted{
		Config: c,
	}
}

//Config as execution configuration
type Config struct {
	FailThres  float64
	Dir        string
	UseFormat  bool
	Verbose    bool
	FormatName string
}

type Grade struct {
	Rank       string `json:"rank"`
	Percentage string `json:"percentage"`
}

type ReportCard struct {
	Grade        Grade                  `json:"grade"`
	Files        int                    `json:"files"`
	Issues       int                    `json:"issues"`
	LinterScores map[string]interface{} `json:"linter_scores"`
}

//Runner as runner abstract
type Runner interface {
	Run() error
}
