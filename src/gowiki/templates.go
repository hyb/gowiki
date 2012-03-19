package main 

import (
	"bytes"
	"log"
	"os"
	"template"
)

var wikiTmplSet *template.Set

func init() {
	// parse the templates
	wikiTmplSet = template.SetMust(template.ParseSetGlob("templates/*.tmpl"))
	// wikiTmplSet.Funcs(template.FuncMap{

	// })
}

type TmplInput struct {
	Common *CommonData
	Contents interface{}
}

type CommonData struct {
	Flash string
}

func RenderPage(tmplName string, common *CommonData, contents interface{}) ([]byte, os.Error) {
	buf := bytes.NewBuffer(nil)
	err := wikiTmplSet.Execute(buf, tmplName, &TmplInput{Common: common, Contents: contents})
	if err != nil {
		log.Println("Error: template failure: ", err.String())
	}
	return buf.Bytes(), err
}