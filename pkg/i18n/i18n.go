/*
Package i18n is a simple language manager, use INI format file.

Source code and other details for the project are available at GitHub:

	https://github.com/gookit/i18n

Language files:

	// structs on mode is FileMode(default)
	lang/
		en.ini
		zh-CN.ini

	// structs on mode is DirMode
	lang/
		en/
			default.ini
			other.ini
		zh-CN/
			default.ini
			other.ini

Load:

	defaultLang = "en"
	languages := map[string]string{
	    "en": "English",
	    "zh-CN": "简体中文",
	    "zh-TW": "繁体中文",
	}

	i18n.Init("conf/lang", defaultLang, languages)

Usage:

	// translate from default language
	val := i18n.Dtr("key")

	// translate from special language
	val := i18n.Tr("en", "key")

这里改造成加载yaml
*/
package i18n

import (
	"errors"
	"fmt"
	"bailu/utils"
	"os"
	"strings"
)

const (
	// ---- I18n.LoadMode language file load mode

	// FileMode language name is file name. "en" -> "lang/en.ini"
	FileMode uint8 = 0
	// DirMode language name is dir name, will load all file in the dir. "en" -> "lang/en/*.ini"
	DirMode uint8 = 1

	// ---- I18n.TransMode translate message mode.

	// SprintfMode render message arguments by fmt.Sprintf
	SprintfMode uint8 = 0
	// ReplaceMode render message arguments by string replace
	ReplaceMode uint8 = 1
)

// I18n language manager
type I18n struct {
	// languages data
	data map[string]map[string]string

	// language files directory
	langDir string
	// language list {en:English, zh-CN:简体中文}
	languages map[string]string
	// loaded lang files
	// loadedFiles []string

	// ------------ config for i18n ------------

	// LoadMode mode for the load language files.
	//  0 single language file
	//  1 multi-language directory
	LoadMode uint8
	// TransMode translate mode.
	//  0 sprintf
	//  1 replace
	TransMode uint8
	// DefaultLang default language name. eg. "en"
	DefaultLang string
	// FallbackLang spare(fallback) language name. eg. "en"
	FallbackLang string
}

/************************************************************
 * create instance
 ************************************************************/

// New an i18n instance
func New(langDir, defLang string, languages map[string]string) *I18n {
	return &I18n{
		data: make(map[string]map[string]string, 0),
		// language data config
		langDir:   langDir,
		languages: languages,
		// set default lang
		DefaultLang: defLang,
	}
}

// NewEmpty new an empty i18n instance
func NewEmpty() *I18n {
	return &I18n{
		data: make(map[string]map[string]string, 0),
		// init languages
		languages: make(map[string]string, 0),
	}
}

// NewWithInit an i18n instance and call init
func NewWithInit(langDir, defLang string, languages map[string]string) *I18n {
	return New(langDir, defLang, languages).Init()
}

// Config the manager instance
func (l *I18n) Config(fn func(l *I18n)) {
	fn(l)
}

// Init load add language files
func (l *I18n) Init() *I18n {
	if l.LoadMode == FileMode {
		l.loadSingleFiles()
	} else if l.LoadMode == DirMode {
		l.loadDirFiles()
	} else {
		panic("invalid load mode setting. only allow 0, 1")
	}

	return l
}

/************************************************************
 * translate message
 ************************************************************/

// Dt translate from default lang
func (l *I18n) Dt(key string, args ...interface{}) string {
	return l.Tr(l.DefaultLang, key, args...)
}

// Dtr translate from default lang
func (l *I18n) Dtr(key string, args ...interface{}) string {
	return l.Tr(l.DefaultLang, key, args...)
}

// DefTr translate from default lang
func (l *I18n) DefTr(key string, args ...interface{}) string {
	return l.Tr(l.DefaultLang, key, args...)
}

// T translate from a lang by key
func (l *I18n) T(lang, key string, args ...interface{}) string {
	return l.Tr(lang, key, args...)
}

// Tr translate from a lang by key
//
// Config:
//
//	[site]
//	name = my blog
//
// Read:
//
//	site.name => "my blog"
func (l *I18n) Tr(lang, key string, args ...interface{}) string {
	if !l.HasLang(lang) {
		// find from fallback lang
		msg := l.transFromFallback(key)
		if msg == "" {
			return key
		}

		// has args for the message
		if len(args) > 0 {
			msg = l.renderMessage(msg, args...)
		}
		return msg
	}

	// find message by key
	msg, ok := l.data[lang][key]
	if !ok {
		// key not exists, find from fallback lang
		msg = l.transFromFallback(key)
	}

	// message is empty
	if msg == "" {
		return key
	}

	// has args for the message
	if len(args) > 0 {
		msg = l.renderMessage(msg, args...)
	}
	return msg
}

// HasKey in the language data
func (l *I18n) HasKey(lang, key string) (ok bool) {
	if !l.HasLang(lang) {
		return
	}
	_, ok = l.data[lang][key]
	return
}

// translate from fallback language
func (l *I18n) transFromFallback(key string) string {
	fl := l.FallbackLang
	if !l.HasLang(fl) {
		return ""
	}
	return l.data[fl][key]
}

const errMsg = "CANNOT-TO-STRING"

func (l *I18n) renderMessage(msg string, args ...interface{}) string {
	if l.TransMode == SprintfMode {
		return fmt.Sprintf(msg, args...)
	}

	// if args[0] is []string
	if ss, ok := args[0].([]string); ok {
		return strings.NewReplacer(ss...).Replace(msg)
	}

	var ss []string

	// if args is map[string]interface{}
	if mp, ok := args[0].(map[string]interface{}); ok {
		for k, v := range mp {
			str, is := v.(string)
			if !is {
				str = errMsg
			}

			ss = append(ss, "{"+k+"}")
			ss = append(ss, str)
		}
	} else {
		// if args is: {field1, value1, field2, value2, ...}, try convert all element to string.
		for i, val := range args {
			str, is := val.(string)
			if !is {
				str = errMsg
			}

			if i%2 == 0 {
				str = "{" + str + "}"
			}
			ss = append(ss, str)
		}
	}

	return strings.NewReplacer(ss...).Replace(msg)
}

/************************************************************
 * data manage
 ************************************************************/

// load language files when LoadMode is 0
func (l *I18n) loadSingleFiles() {
	pathSep := string(os.PathSeparator)

	for lang := range l.languages {
		//data, err := utils.LoadYaml(l.langDir + pathSep + lang + ".yaml")
		//if err != nil {
		//	panic("fail to load language: " + lang + ", error " + err.Error())
		//}
		data, err := utils.LoadFlatYaml[string](l.langDir + pathSep + lang + ".yaml")
		if err != nil {
			panic("fail to load language: " + lang + ", error " + err.Error())
		}
		l.data[lang] = data
	}
}

// load language files when LoadMode is 1
func (l *I18n) loadDirFiles() {
	pathSep := string(os.PathSeparator)

	for lang := range l.languages {
		dirPath := l.langDir + pathSep + lang
		files, err := os.ReadDir(dirPath)
		if err != nil {
			panic("read dir fail: " + dirPath + ", error " + err.Error())
		}

		for _, fi := range files {
			// skip dir and filter the specified format
			if fi.IsDir() || !strings.HasSuffix(fi.Name(), ".yaml") {
				continue
			}

			d, err := utils.LoadFlatYaml[string](dirPath + pathSep + fi.Name())
			if err != nil {
				panic("fail to load language file: " + lang + ", error " + err.Error())
			}
			l.data[lang] = d
		}
	}
}

// Add register and init new language. alias of NewLang()
func (l *I18n) Add(lang string, name string) {
	l.NewLang(lang, name)
}

// AddLang register and init new language. alias of NewLang()
func (l *I18n) AddLang(lang string, name string) {
	l.NewLang(lang, name)
}

// WithLang register and init new language. alias of NewLang()
func (l *I18n) WithLang(lang string, name string) *I18n {
	l.NewLang(lang, name)
	return l
}

// NewLang create/add a new language
//
// Usage:
//
//	i18n.NewLang("zh-CN", "简体中文")
func (l *I18n) NewLang(lang string, name string) {
	// language exist
	if _, ok := l.languages[lang]; ok {
		return
	}

	if name == "" {
		name = utils.UpperFirst(lang)
	}
	l.data[lang] = make(map[string]string)
	l.languages[lang] = name
}

// LoadFile append data to a exist language
//
// Usage:
//
//	i18n.LoadFile("zh-CN", "path/to/zh-CN.ini")
func (l *I18n) LoadFile(lang string, file string) error {
	// append data
	if _, ok := l.data[lang]; ok {
		data, err := utils.LoadFlatYaml[string](file)
		if err != nil {
			return err
		}
		l.data[lang] = data
		return nil
	} else {
		return errors.New("language" + lang + " not exist, please create it before load data")
	}
}

// LoadString load language data form a string
//
// Usage:
//
//	i18n.LoadString("zh-CN", `
//	name = blog
//	age = 233
//	`)
func (l *I18n) LoadString(lang string, jsonStr string) error {
	_, ok := l.data[lang]
	if !ok {
		return fmt.Errorf("language '%s' is not registered", lang)
	}
	// append data
	ld, err := utils.JsonStr2FlatMap[string](jsonStr)
	l.data[lang] = ld
	return err
}

// Lang get language data instance
func (l *I18n) Lang(lang string) map[string]string {
	if _, ok := l.languages[lang]; ok {
		return l.data[lang]
	}
	return nil
}

// HasLang in the manager
func (l *I18n) HasLang(lang string) bool {
	_, ok := l.languages[lang]
	return ok
}

// DelLang from the i18n manager
func (l *I18n) DelLang(lang string) bool {
	_, ok := l.languages[lang]
	if ok {
		delete(l.data, lang)
		delete(l.languages, lang)
	}

	return ok
}

// Languages get all languages
func (l *I18n) Languages() map[string]string {
	return l.languages
}
