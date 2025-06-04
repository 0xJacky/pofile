package pofile

import (
	"os"
	"reflect"
	"regexp"
	"strings"

	"slices"

	"github.com/itchyny/timefmt-go"
	"github.com/pkg/errors"
)

const (
	MSGID = iota
	MSGIDPLURAL
	MSGSTR
	MSGCTXT
	COMPLATE
)

func Parse(path string) (p *Pofile, err error) {
	var bytes []byte
	bytes, err = os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "parse read file error")
	}
	text := string(bytes)

	p, err = ParseText(text)
	return
}

func parseHeader(item Item) (h *Header, err error) {
	h = &Header{}
	// parse header
	v := reflect.ValueOf(h).Elem()
	// fmt.Println(item.MsgStr[0])
	for i := 0; i < v.NumField()-1; i++ {
		key := v.Type().Field(i).Tag.Get("key")
		// 使用与 Header.Get 相同的正则表达式
		regExp := regexp.MustCompile(key + ":[ ]+(.*?)(?:\\n|$)")
		matchSlice := regExp.FindStringSubmatch(item.MsgStr[0])
		if len(matchSlice) < 2 {
			continue
		}
		match := strings.TrimSpace(matchSlice[1])

		if v.Type().Field(i).Type.String() == "string" {
			v.Field(i).Set(reflect.ValueOf(match))
		}

		if v.Type().Field(i).Type.String() == "*time.Time" {
			if match != "" {
				t, parseErr := timefmt.Parse(match, "%Y-%m-%d %H:%M%z")
				if parseErr == nil {
					v.Field(i).Set(reflect.ValueOf(&t))
				}
			}
		}
	}
	h.rawText = item.MsgStr[0]
	return
}

// extractQuotedString 提取引号内的字符串内容并正确处理转义字符
func extractQuotedString(line string) string {
	// 寻找第一个未转义的双引号
	start := -1
	for i, r := range line {
		if r == '"' {
			// 检查是否被转义
			backslashCount := 0
			for j := i - 1; j >= 0 && line[j] == '\\'; j-- {
				backslashCount++
			}
			// 如果反斜杠数量为偶数（包括0），则双引号未被转义
			if backslashCount%2 == 0 {
				start = i
				break
			}
		}
	}

	if start == -1 {
		return ""
	}

	// 寻找匹配的结束双引号
	end := -1
	for i := start + 1; i < len(line); i++ {
		if line[i] == '"' {
			// 检查是否被转义
			backslashCount := 0
			for j := i - 1; j >= 0 && line[j] == '\\'; j-- {
				backslashCount++
			}
			// 如果反斜杠数量为偶数，则双引号未被转义
			if backslashCount%2 == 0 {
				end = i
				break
			}
		}
	}

	if end == -1 {
		return ""
	}

	// 提取并处理转义字符
	result := line[start+1 : end]
	// 处理转义序列
	result = strings.ReplaceAll(result, "\\\"", "\"")
	result = strings.ReplaceAll(result, "\\\\", "\\")
	result = strings.ReplaceAll(result, "\\n", "\n")
	result = strings.ReplaceAll(result, "\\t", "\t")
	result = strings.ReplaceAll(result, "\\r", "\r")

	return result
}

func ParseText(text string) (p *Pofile, err error) {
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
	// fixed header and bottom boundary
	boundarySlice[0].Start = 0
	boundarySlice[boundaryIndex-1].End = lineLen - 1
	// parse each po entry
	init := true
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
				state = MSGCTXT
				extracted := extractQuotedString(line)
				item.Msgctxt += extracted
			} else if strings.HasPrefix(line, "msgid") &&
				!strings.HasPrefix(line, "msgid_plural") {
				state = MSGID
				extracted := extractQuotedString(line)
				item.MsgId += extracted
			} else if strings.HasPrefix(line, "msgstr") {
				state = MSGSTR
				extracted := extractQuotedString(line)
				msgStrIndex++
				item.MsgStr = append(item.MsgStr, extracted)
			} else if strings.HasPrefix(line, "msgid_plural") {
				state = MSGIDPLURAL
				extracted := extractQuotedString(line)
				item.MsgIdPlural += extracted
			} else if strings.HasPrefix(line, "#~") {
				// ignore
				break
			} else {
				extracted := extractQuotedString(line)
				if extracted != "" {
					switch state {
					case MSGCTXT:
						item.Msgctxt += extracted
					case MSGID:
						item.MsgId += extracted
					case MSGIDPLURAL:
						item.MsgIdPlural += extracted
					case MSGSTR:
						item.MsgStr[msgStrIndex] += extracted
					}
				}
			}
		}
		// fmt.Println(item)
		// fmt.Println("=========")
		if init || item.MsgId != "" {
			p.Items = append(p.Items, item)
			init = false
		}
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
	regExp := regexp.MustCompile(key + ":[ ]+(.*?)(?:\\n|$)")
	matchSlice := regExp.FindStringSubmatch(h.rawText)
	if len(matchSlice) < 2 {
		return nil
	}
	return strings.TrimSpace(matchSlice[1])
}

func (item *Item) isFuzzy() bool {
	return slices.Contains(item.Flags, "fuzzy")
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
			msgStrSlice = append(msgStrSlice, item.MsgStr...)
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
