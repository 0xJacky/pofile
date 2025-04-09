# Pofile

[![Go Report Card](https://goreportcard.com/badge/github.com/0xJacky/pofile)](https://goreportcard.com/report/github.com/0xJacky/pofile)
[![Go Reference](https://pkg.go.dev/badge/github.com/0xJacky/pofile.svg)](https://pkg.go.dev/github.com/0xJacky/pofile)
[![License](https://img.shields.io/github/license/0xJacky/pofile)](https://github.com/0xJacky/pofile/blob/main/LICENSE)

A Gettext `.po` file parser written in Go.

## Features

- Parse po file headers into a structured format with easy access methods
- Full support for parsing po file elements:
  - Translator comments
  - Extracted comments
  - References
  - Flags
  - Message context (msgctxt)
  - Message ID (msgid)
  - Plural message ID (msgid_plural)
  - Translations (msgstr)
- Automatic handling of fuzzy translations
- Convert po files to map/dictionary format (suitable for JSON conversion)
- Command-line interface for batch processing of po files

## Installation

### Library

```bash
go get github.com/0xJacky/pofile
```

### CLI Tool

Download the binary from [Releases](https://github.com/0xJacky/pofile/releases) or install from source:

```bash
go install github.com/0xJacky/pofile/cmd/pofile@latest
```

## Usage

### Library Usage

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/0xJacky/pofile"
	"log"
	"os"
)

func main() {
	// Parse a po file
	p, err := pofile.Parse("path/to/your.po")
	if err != nil {
		log.Fatalf("Failed to parse po file: %v", err)
	}
	
	// Access header information
	fmt.Printf("Language: %s\n", p.Header.Language)
	fmt.Printf("Custom header field: %v\n", p.Header.Get("X-Generator"))
	
	// Process translations
	for _, item := range p.Items {
		fmt.Printf("Original: %s\n", item.MsgId)
		fmt.Printf("Translation: %s\n", item.MsgStr[0])
	}
	
	// Convert to dictionary and save as JSON
	dict := p.ToDict()
	bytes, _ := json.MarshalIndent(dict, "", "  ")
	os.WriteFile("translations.json", bytes, 0644)
}
```

### Multiple Language Support

```go
package main

import (
	"encoding/json"
	"github.com/0xJacky/pofile"
	"log"
	"os"
	"path/filepath"
)

func main() {
	dict := make(pofile.Dict)
	languages := []string{"de", "en", "fr", "ja", "ko", "zh_TW"}
	
	for _, lang := range languages {
		poPath := filepath.Join("locale", lang, "LC_MESSAGES", "app.po")
		p, err := pofile.Parse(poPath)
		if err != nil {
			log.Printf("Failed to parse %s: %v", poPath, err)
			continue
		}
		dict[p.Header.Language] = p.ToDict()
	}
	
	bytes, _ := json.MarshalIndent(dict, "", "  ")
	os.WriteFile("all_translations.json", bytes, 0644)
}
```

### CLI Usage

Convert a single po file to JSON:

```bash
pofile build --file path/to/file.po
```

Convert all po files in a directory:

```bash
pofile build --file path/to/directory
```

## API Reference

### Core Types

#### Pofile
```go
type Pofile struct {
	Header Header
	Items  []Item
}

// Convert pofile to dictionary format
func (p *Pofile) ToDict() Dict
```

#### Header
```go
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

// Get custom header field
func (h *Header) Get(key string) interface{}
```

#### Item
```go
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

### Main Functions

```go
// Parse a .po file at the specified path
func Parse(path string) (p *Pofile, err error)
```

## License

MIT
