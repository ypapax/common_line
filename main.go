package main

import (
	"flag"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/ypapax/logrus_conf"
	"io/ioutil"
	"strings"
)

func main(){
	if err := func() error {
		if err := logrus_conf.PrepareFromEnv("common_line"); err != nil {
			return errors.WithStack(err)
		}
		var filePath string
		const fileParamName = "file"
		flag.StringVar(&filePath, fileParamName, "", "file path of the log file")
		flag.Parse()
		if len(filePath) == 0 {
			flag.Usage()
			return errors.Errorf("missing -%+v path param", fileParamName)
		}
		logrus.Infof("starting reading file %+v	...", filePath)
		b, err := ioutil.ReadFile(filePath)
		if err != nil {
			return errors.WithStack(err)
		}
		lines := strings.Split(string(b), "\n")
		logrus.Infof("lines number: %+v", len(lines))
		return nil
	}(); err != nil {
		logrus.Errorf("%+v", err)
	}

}
