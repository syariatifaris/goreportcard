[![Go Report Card](https://goreportcard.com/badge/gojp/goreportcard)](https://goreportcard.com/report/gojp/goreportcard) [![Build Status](https://travis-ci.org/gojp/goreportcard.svg?branch=master)](https://travis-ci.org/gojp/goreportcard) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/gojp/goreportcard/blob/master/LICENSE)

# Go Report Card

This repository is a fork version from `github.com/gojp/goreportcard`. All inherited functionalities such as the web dashboard and goreportcard-cli persists in this repository. This forked version provide more robust presentation and hook communication between the cli tools and user's REST API.

`goreportcard` is a web application that generates a report on the quality of an open source go project. It uses several measures, including `gofmt`, `go vet`, `go lint` and `gocyclo`. To get a report on your own project, try using the hosted version of this code running at [goreportcard.com](https://goreportcard.com).

### Requirement

For custom `goreportcard` installation, you need to install `Go Dep`. 

### Custom Installation

Run

```$xslt
sh install.sh
```

This will run dep ensure, and installing `goreportcard-cli` to your PATH

### Standard Installation

This standard installation will pull the latest version of `github.com/gojp/goreportcard`. 
Assuming you already have a recent version of Go installed, pull down the code with `go get`:

```
go get github.com/gojp/goreportcard
```

Go into the source directory and pull down the project dependencies:

```
cd $GOPATH/src/github.com/gojp/goreportcard
make install
```

Now run

```
make start-dev
```

and you should see

```
Running on 127.0.0.1:8000...
```

Navigate to that URL in your browser and check that you can see the front page.

When running the site in a production environment, instead of `make start-dev`, run:

```
make start
```

### Command Line Interface

There is also a CLI available for grading applications on your local machine.

#### JSON Format

To format the response as a JSON, run:
```$xslt
goreportcard-cli -f
```

It will show result such as:
```$xslt
{
	"grade": {
		"rank": "A+",
		"percentage": "99.8%"
	},
	"files": 433,
	"issues": 5,
	"linter_scores": {
		"go_vet": "99%",
		"gocyclo": "99%",
		"gofmt": "100%",
		"golint": "99%",
		"ineffassign": "99%",
		"license": "100%",
		"misspell": "100%"
	}
}
```

#### Simple Format

Example usage:
```
go get github.com/gojp/goreportcard/cmd/goreportcard-cli
cd $GOPATH/src/github.com/gojp/goreportcard
goreportcard-cli
```

```
Grade: A+ (99.9%)
Files: 362
Issues: 2
gofmt: 100%
go_vet: 99%
gocyclo: 99%
golint: 100%
ineffassign: 100%
license: 100%
misspell: 100%
```

Verbose output is also available:
```
goreportcard-cli -v
```

```
Grade: A+ (99.9%)
Files: 332
Issues: 2
gofmt: 100%
go_vet: 99%
go_vet  vendor/github.com/prometheus/client_golang/prometheus/desc.go:25
        error: cannot find package "github.com/prometheus/client_model/go" in any of: (vet)

gocyclo: 99%
gocyclo download/download.go:22
        warning: cyclomatic complexity 17 of function download() is high (> 15) (gocyclo)

golint: 100%
ineffassign: 100%
license: 100%
misspell: 100%
```

### Hook

This fork version `github.com/syariatifaris/goreportcard` enables POST to enlisted hooks to send the report data for various purpose. 
You need to prepare the accessible endpoint so the goreportcard able to send the data using REST. 

1. Create hook config json file such as:
```aidl
{
    "hooks": [
        {
            "url": "http://hook1.service.consul/",
            "headers":[
                {
                    "key": "Auth",
                    "value": "abcd"
                },
                {
                    "key": "Content-Type",
                    "value": "application/json"
                }
            ]
        }
    ]
}
```

2. Save the file as a json file i.e: `hook.json`
3. Run:
```aidl
goreportcard-cli -f -hook "pathto/hook.json" -timeout 10
```
4. You will see the output such as: 
```aidl
{
	"grade": {
		"rank": "A+",
		"percentage": "99.8%"
	},
	"files": 433,
	"issues": 5,
	"linter_scores": {
		"go_vet": "99%",
		"gocyclo": "99%",
		"gofmt": "100%",
		"golint": "99%",
		"ineffassign": "99%",
		"license": "100%",
		"misspell": "100%"
	}
}
contacting hook to send data ..
hook responses:
[
	{
		"url": "http://hook.service.consul",
		"status_code": 200,
		"message": "some message from hook"
	}
]
```
### Contributing

Go Report Card is an open source project run by volunteers, and contributions are welcome! Check out the [Issues](https://github.com/gojp/goreportcard/issues) page to see if your idea for a contribution has already been mentioned, and feel free to raise an issue or submit a pull request.

### License

The code is licensed under the permissive Apache v2.0 licence. This means you can do what you like with the software, as long as you include the required notices. [Read this](https://tldrlegal.com/license/apache-license-2.0-(apache-2.0)) for a summary.

### Notes

We don't support this on Windows since we have no way to test it on Windows.
