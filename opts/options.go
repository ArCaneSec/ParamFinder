package opts

import (
	"errors"

	"github.com/projectdiscovery/goflags"
)

type Options struct {
	Url           string
	Headless      bool
	DirectoryPath string
	Silent        bool
	Crawl         bool
	OutputPath    string
	Headers       goflags.StringSlice
}

func (o *Options) Validate() error {
	if o.Url == "" && o.DirectoryPath == "" {
		return errors.New("options: you must provide 1 url or a path to directory containing request/response files")
	}

	if len(o.Headers)%2 != 0 {
		return errors.New("options: each header must must have a value")
	}

	return nil
}
