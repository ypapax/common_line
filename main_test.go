package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseLine(t *testing.T){
	type testCase struct {
		inputLine string
		expected string
	}
	cases := []testCase{{
		inputLine: "INFO[2022-02-14 15:59:17.906586 +0300]/Users/some/Dropbox/golang/src/github.com/file/path/some-file.go:95 github.com/ypapax/log_conf/packages/backend/internal/log_conf.Files() log file /Users/user/tmp/project__cmd-2022-02-14--15-59-14/debug.project.log for levels: [debug]",
		expected: "/Users/some/Dropbox/golang/src/github.com/file/path/some-file.go:95",
	}}
	for _, c := range cases {
		t.Run(c.inputLine, func(t *testing.T) {
			as := assert.New(t)
			actual, err := ParseLine(c.inputLine)
			if !as.NoError(err) {
				return
			}
			if !as.Equal(c.expected, actual) {
				return
			}
		})
	}

}
