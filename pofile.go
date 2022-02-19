package pofile

import (
	"github.com/itchyny/timefmt-go"
	"github.com/pkg/errors"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"
)

const (
	MSGID = iota
	MSGIDPLURAL
	MSGSTR
	MSGCTXT
	COMPLATE
)

func Parse(path string) (p *Pofile, err error) {
	p, err = parse(path)
	return
}

func parseHeader(item Item) (h *Header, err error) {
	h = &Header{}
	// parse header
	v := reflect.ValueOf(h).Elem()
	// fmt.Println(item.MsgStr[0])
	for i := 0; i < v.NumField()-1; i++ {
		key := v.Type().Field(i).Tag.Get("key")
		regExp := regexp.MustCompile(key + `:[ ]+(.*?)\\n`)
		matchSlice := regExp.FindStringSubmatch(item.MsgStr[0])
		if len(matchSlice) < 1 {
			continue
		}
		match := strings.ReplaceAll(matchSlice[1], "\\n", "")
		if v.Type().Field(i).Type.String() == "string" {
			v.Field(i).Set(reflect.ValueOf(match))
		}

		if v.Type().Field(i).Type.String() == "*time.Time" {
			t, _ := timefmt.Parse(match, "%Y-%m-%d %H:%M%z")
			v.Field(i).Set(reflect.ValueOf(&t))
		}
	}
	h.rawText = item.MsgStr[0]
	// fmt.Println(h)
	return
}

func parse(path string) (p *Pofile, err error) {
	var bytes []byte
	bytes, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	text := string(bytes)
	lines := strings.Split(text, "\n")

	p = &Pofile{}
	lineLen := len(lines)
	type boundary struct {
		Start int
		End   int
	}
	var boundarySlice []boundary
	boundaryIndex := 0
	// determine the boundary of the entry
	for i := range lines {
		if strings.HasPrefix(lines[i], "msgid") &&
			!strings.HasPrefix(lines[i], "msgid_plural") {
			var start int
			for start = i - 1; start >= 0; start-- {
				if !strings.HasPrefix(lines[start], "msgctxt") &&
					strings.HasSuffix(strings.TrimSpace(lines[start]), "\"") {
					start++
					break
				}
			}
			boundarySlice = append(boundarySlice, boundary{
				Start: start,
			})
			if boundaryIndex > 0 {
				boundarySlice[boundaryIndex-1].End = start - 1
			}
			boundaryIndex++
		}
	}
	// fix header and bottom
	boundarySlice[0].Start = 0
	boundarySlice[boundaryIndex-1].End = lineLen - 1
	// parse each po entry
	for _, v := range boundarySlice {
		// fmt.Println(v.Start+1, v.End+1)
		item := Item{}
		state := COMPLATE
		msgStrIndex := -1

		for i := v.Start; i <= v.End; i++ {
			// fmt.Println(lines[i])
			line := lines[i]
			if strings.HasPrefix(line, "# ") {
				// Translator Comments
				line = strings.TrimPrefix(line, "# ")
				item.TranslatorComments = append(item.TranslatorComments, line)
			} else if strings.HasPrefix(line, "#. ") {
				// ExtractedComments
				line = strings.TrimPrefix(line, "#. ")
				line = strings.TrimRight(line, ",")
				item.ExtractedComments = append(item.ExtractedComments, line)
			} else if strings.HasPrefix(line, "#: ") {
				// Reference
				line = strings.TrimPrefix(line, "#: ")
				line = strings.TrimRight(line, ",")
				item.Reference = append(item.Reference, line)
			} else if strings.HasPrefix(line, "#, ") {
				// Flags
				line = strings.TrimPrefix(line, "#, ")
				item.Flags = strings.Split(line, ",")
			} else if strings.HasPrefix(line, "msgctxt") {
				regExp := regexp.MustCompile("^msgctxt\\s+\"(.*)\"")
				matchSlice := regExp.FindStringSubmatch(line)
				state = MSGCTXT
				if len(matchSlice) < 2 {
					continue
				}
				item.Msgctxt += matchSlice[1]
			} else if strings.HasPrefix(line, "msgid") &&
				!strings.HasPrefix(line, "msgid_plural") {
				regExp := regexp.MustCompile("^msgid\\s+\"(.*)\"")
				matchSlice := regExp.FindStringSubmatch(line)
				state = MSGID
				if len(matchSlice) < 2 {
					continue
				}
				item.MsgId += matchSlice[1]
			} else if strings.HasPrefix(line, "msgstr") {
				regExp := regexp.MustCompile("^msgstr[\\s\\S]*\"(.*)\"")
				matchSlice := regExp.FindStringSubmatch(line)
				state = MSGSTR
				if len(matchSlice) < 2 {
					continue
				}
				msgStrIndex++
				item.MsgStr = append(item.MsgStr, matchSlice[1])
			} else if strings.HasPrefix(line, "msgid_plural") {
				regExp := regexp.MustCompile("^msgid_plural\\s+\"(.*)\"")
				matchSlice := regExp.FindStringSubmatch(line)
				state = MSGIDPLURAL
				if len(matchSlice) < 2 {
					continue
				}
				item.MsgIdPlural += matchSlice[1]
			} else if strings.HasPrefix(line, "#~") {
				// ignore
				continue
			} else {
				strRegExp := regexp.MustCompile("\"(.*)\"")
				matchSlice := strRegExp.FindStringSubmatch(line)
				if len(matchSlice) < 2 {
					continue
				}
				switch state {
				case MSGCTXT:
					item.Msgctxt += matchSlice[1]
				case MSGID:
					item.MsgId += matchSlice[1]
				case MSGIDPLURAL:
					item.MsgIdPlural += matchSlice[1]
				case MSGSTR:
					item.MsgStr[msgStrIndex] += matchSlice[1]
				}
			}
		}
		// fmt.Println(item)
		// fmt.Println("=========")
		p.Items = append(p.Items, item)
	}

	// parse Header
	var h *Header
	h, err = parseHeader(p.Items[0])
	if err != nil {
		return nil, errors.Wrap(err, "error parse header")
	}
	p.Header = *h
	if len(p.Items) > 0 {
		p.Items = p.Items[1:]
	}
	return
}

func (h *Header) Get(key string) interface{} {
	regExp := regexp.MustCompile(key + ":[ ]+(.*?)\\\\n")
	matchSlice := regExp.FindStringSubmatch(h.rawText)
	if len(matchSlice) < 1 {
		return nil
	}
	return strings.ReplaceAll(matchSlice[1], "\\n", "")
}

type Dict map[string]interface{}

func (item *Item) isFuzzy() bool {
	for _, v := range item.Flags {
		if v == "fuzzy" {
			return true
		}
	}
	return false
}

func (p *Pofile) ToDict() (dict Dict) {
	dict = make(Dict)
	for _, item := range p.Items {
		if item.isFuzzy() {
			continue
		}

		//fmt.Println("Msgctxt", item.Msgctxt)
		//fmt.Println("MsgId", item.MsgId)
		//fmt.Println("MsgIdPlural", item.MsgIdPlural)
		//fmt.Println("MsgStr", item.MsgStr)
		var tmp interface{}
		if len(item.MsgStr) == 1 {
			tmp = item.MsgStr[0]
		} else if len(item.MsgStr) > 1 {
			var msgStrSlice []string
			for _, v := range item.MsgStr {
				msgStrSlice = append(msgStrSlice, v)
			}
			tmp = msgStrSlice
		}
		if item.Msgctxt != "" {
			dict[item.MsgId] = make(Dict)
			dict[item.MsgId].(Dict)[item.Msgctxt] = tmp
		} else {
			dict[item.MsgId] = tmp
		}
	}
	return
}
