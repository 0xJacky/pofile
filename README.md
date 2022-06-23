# Pofile
Gettext po file parse written in Go

## Install
```
go get github.com/0xJacky/pofile
```

or download binary from [Release](https://github.com/0xJacky/pofile/releases) for using cli mode.

## Feature
1. Parse po file headers to a struct type, provide Get method for getting custom field.
2. Parse po file items, support translator comments, extracted comments, reference, flags, msgctxt, msgid, msgid_plural, msgstr types.
3. Automatically ignore fuzzy items.
4. Parse po file to map[string]interface{} which can be further converted to JSON.
5. Provide cli mode for parsing pofile(s).

## Type
1. Pofile
```
type Pofile struct {
	Header Header
	Items  []Item
}
```
2. Convert pofile to dict
```
func (p *Pofile) ToDict() (dict Dict)
```
3. Entry
```
func Parse(path string) (p *Pofile, err error)
```
4. Pofile Header Struct
```
type Header struct {
	ProjectIdVersion        string     `key:"Project-Id-Version"`
	ReportMsgBugsTo         string     `key:"Report-Msgid-Bugs-To"`
	POTCreationDate         *time.Time `key:"POT-Creation-Date"`
	PORevisionDate          *time.Time `key:"PO-Revision-Date"`
	LastTranslator          string     `key:"Last-Translator"`
	Language                string     `key:"Language"`
	LanguageTeam            string     `key:"Language-Team"`
	ContentType             string     `key:"Content-Type"`
	ContentTransferEncoding string     `key:"Content-Transfer-Encoding"`
	PluralForms             string     `key:"Plural-Forms"`
}
```
5. Visit custom header field
```
func (h *Header) Get(key string) interface{}
```
6. Profile item struct
```
type Item struct {
	TranslatorComments []string
	ExtractedComments  []string
	Reference          []string
	Flags              []string
	Msgctxt            string
	MsgId              string
	MsgIdPlural        string
	MsgStr             []string
}
```

## CLI Mode
1. Convert a single pofile to JSON.
```
./pofile build --file <path-to-pofile>
```

2. Convert all pofiles from a directory to JSON.
```
./pofile build --file <path-to-dir>
```

## Example
```
package test

import (
	"encoding/json"
	"fmt"
	"github.com/0xJacky/pofile"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
)

func TestPofile(t *testing.T) {
	p, err := pofile.Parse("en.po")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Test Header")
	header := p.Header
	fmt.Println("ProjectIdVersion", header.ProjectIdVersion)
	fmt.Println("ReportMsgBugsTo", header.ReportMsgBugsTo)
	fmt.Println("POTCreationDate", header.POTCreationDate)
	fmt.Println("PORevisionDate", header.PORevisionDate)
	fmt.Println("LastTranslator", header.LastTranslator)
	fmt.Println("Language", header.Language)
	fmt.Println("LanguageTeam", header.LanguageTeam)
	fmt.Println("ContentType", header.ContentType)
	fmt.Println("ContentTransferEncoding", header.ContentTransferEncoding)
	fmt.Println("PluralForms", header.PluralForms)

	fmt.Println("==========")
	fmt.Println("Test Header.Get")
	fmt.Println(header.Get("X-Generator"))
	fmt.Println("==========")

	fmt.Println("Test Items")
	for _, item := range p.Items {
		fmt.Println("TranslatorComments", item.TranslatorComments)
		fmt.Println("ExtractedComments", item.ExtractedComments)
		fmt.Println("Reference", item.Reference)
		fmt.Println("Flags", item.Flags)
		fmt.Println("Msgctxt", item.Msgctxt)
		fmt.Println("MsgId", item.MsgId)
		fmt.Println("MsgIdPlural", item.MsgIdPlural)
		fmt.Println("MsgStr", item.MsgStr)
		fmt.Println("==========")
	}

	// Test Pofile ToDict
	bytes, _ := json.Marshal(p.ToDict())
	_ = ioutil.WriteFile("output_test.json", bytes, 0644)

	fmt.Println("Test Pofile ToDict")
	fmt.Println(p.ToDict())
	dict := make(pofile.Dict)

	lang := []string{"de", "en", "fr", "ja", "ko", "zh_TW"}

	for _, v := range lang {
		p, err = pofile.Parse(filepath.Join("locale", v, "LC_MESSAGES", "app.po"))
		if err != nil {
			log.Fatalln(err)
		}
		dict[p.Header.Language] = p.ToDict()
	}

	bytes, _ = json.Marshal(dict)
	_ = ioutil.WriteFile("translates.json", bytes, 0644)
}

```