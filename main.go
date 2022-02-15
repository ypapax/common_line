package main

import (
	"flag"
	"github.com/pkg/errors"
	logrus "github.com/sirupsen/logrus"
	"github.com/ypapax/logrus_conf"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
	"sync"
)

const maxPrint = 10

type Parsed struct {
	FilePath string
	Copies   []string
}

type Lines struct {
	Parsed []Parsed
	ByParsedLine map[string]Parsed
	//ByCount map[int][]Parsed
	mtx sync.Mutex
}

var ll Lines

func addLine(parsed, full string) {
	ll.mtx.Lock()
	defer ll.mtx.Unlock()

	if ll.ByParsedLine == nil {
		ll.ByParsedLine = make(map[string]Parsed)
	}
	if v, ok := ll.ByParsedLine[parsed]; ok {
		v.Copies = append(ll.ByParsedLine[parsed].Copies, full)
		ll.ByParsedLine[parsed] = v
		return
	}
	p := Parsed{FilePath: parsed, Copies: []string{full}}
	ll.Parsed = append(ll.Parsed, p)
	ll.ByParsedLine[parsed] = p
}

type ByCopiesCount []Parsed

func (a ByCopiesCount) Len() int           { return len(a) }
func (a ByCopiesCount) Less(i, j int) bool { return len(a[i].Copies) > len(a[j].Copies) }
func (a ByCopiesCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }


func sortByCount() {
	ll.mtx.Lock()
	defer ll.mtx.Unlock()
	sort.Sort(ByCopiesCount(ll.Parsed))
}

func printMostUsed() {
	ll.mtx.Lock()
	defer ll.mtx.Unlock()
	for i, l := range ll.Parsed {
		lc := logrus.WithField("len(l.Copies)", len(l.Copies))
		lc.Infof("filePath: %+v", l.FilePath)
		if i > maxPrint {
			break
		}
	}
}


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
		for _, l := range lines {
			parsed, errP := ParseLine(l)
			if errP != nil {
				return errors.WithStack(errP)
			}
			if len(parsed) == 0 {
				logrus.Tracef("skip line %+v because it gives empty parsed result", l)
				continue
			}
			addLine(parsed, l)
		}
		sortByCount()
		printMostUsed()
		return nil
	}(); err != nil {
		logrus.Errorf("%+v", err)
	}

}

var reg= regexp.MustCompile(`\](.+\:\d+) `)

func ParseLine(inputLine string) (string, error) {
	sm := reg.FindAllStringSubmatch(inputLine, -1)
	for i, sm1 := range sm {
		logrus.Tracef("%+v sm: %+v", i, strings.Join(sm1, "; "))
	}
	if len(sm) == 0 {
		//return "", errors.Errorf("not enough parts for inputLine '%+v'", inputLine)
		return "", nil
	}
	if len(sm[0]) <= 1 {
		return "", errors.Errorf("not enough parts")
	}
	return sm[0][1], nil
}
