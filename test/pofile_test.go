package test

import (
	"encoding/json"
	"fmt"
	"github.com/0xJacky/pofile/profile"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
)

func TestPofile(t *testing.T) {
	p, err := profile.Parse("app.po")
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
	dict := make(profile.Dict)

	lang := []string{"de", "en", "fr", "ja", "ko", "zh_TW"}

	for _, v := range lang {
		p, err = profile.Parse(filepath.Join("locale", v, "LC_MESSAGES", "app.po"))
		if err != nil {
			log.Fatalln(err)
		}
		dict[p.Header.Language] = p.ToDict()
	}

	bytes, _ = json.Marshal(dict)
	_ = ioutil.WriteFile("translates.json", bytes, 0644)
}
