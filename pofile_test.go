package pofile

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestPofile(t *testing.T) {
	p, err := Parse("test_data/app.po")
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
	_ = os.WriteFile("test_data/output_test.json", bytes, 0644)

	fmt.Println("Test Pofile ToDict")
	fmt.Println(p.ToDict())
	dict := make(Dict)

	lang := []string{"de", "en", "fr", "ja", "ko", "zh_TW"}

	for _, v := range lang {
		p, err = Parse(filepath.Join("test_data/locale", v, "LC_MESSAGES", "app.po"))
		if err != nil {
			log.Fatalln(err)
		}
		dict[p.Header.Language] = p.ToDict()
	}

	bytes, _ = json.Marshal(dict)
	_ = os.WriteFile("test_data/translates.json", bytes, 0644)
}

func TestEscapeCharacters(t *testing.T) {
	p, err := Parse("test_data/escape_test.po")
	if err != nil {
		t.Fatalf("解析 PO 文件失败: %v", err)
	}

	// 验证转义字符正确处理
	expectedResults := map[string]string{
		"To confirm revocation, please type \"Revoke\" in the field below:": "要确认撤销，请在下面的字段中输入 \"撤销\"：",
		"The \"Settings\" menu":         "\"设置\" 菜单",
		"Click \"OK\" to continue":      "点击 \"确定\" 继续",
		"Use \"double quotes\" in text": "在文本中使用 \"双引号\"",
	}

	dict := p.ToDict()

	for expectedKey, expectedValue := range expectedResults {
		actualValue, exists := dict[expectedKey]
		if !exists {
			t.Errorf("字典中缺少键: %s", expectedKey)
			continue
		}

		if actualValue != expectedValue {
			t.Errorf("翻译不匹配:\n键: %s\n期望: %s\n实际: %s", expectedKey, expectedValue, actualValue)
		}
	}

	fmt.Println("转义字符处理测试通过")
}

func TestHeaderParsing(t *testing.T) {
	p, err := Parse("test_data/app.po")
	if err != nil {
		t.Fatalf("解析 PO 文件失败: %v", err)
	}

	// 验证头部字段正确解析
	if p.Header.Language != "zh_CN" {
		t.Errorf("Language 字段解析错误，期望: zh_CN，实际: %s", p.Header.Language)
	}

	if p.Header.LastTranslator != "Automatically generated" {
		t.Errorf("LastTranslator 字段解析错误，期望: Automatically generated，实际: %s", p.Header.LastTranslator)
	}

	if p.Header.ContentType != "text/plain; charset=UTF-8" {
		t.Errorf("ContentType 字段解析错误，期望: text/plain; charset=UTF-8，实际: %s", p.Header.ContentType)
	}

	// 验证 Header.Get 方法
	xGenerator := p.Header.Get("X-Generator")
	if xGenerator != "Poedit 3.0.1" {
		t.Errorf("X-Generator 获取错误，期望: Poedit 3.0.1，实际: %v", xGenerator)
	}

	fmt.Println("头部解析测试通过")
}
