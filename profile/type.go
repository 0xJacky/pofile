package profile

import (
	"time"
)

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
	rawText                 string
}

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

type Pofile struct {
	Header Header
	Items  []Item
}
