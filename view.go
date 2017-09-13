package main

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
)

//  @see https://gohugo.io/templates/partials/
//  @see https://gohugo.io/templates/introduction/
//  @see https://elithrar.github.io/article/approximating-html-template-inheritance/
//  @see https://hackernoon.com/golang-template-2-template-composition-and-how-to-organize-template-files-4cb40bcdf8f6
//
type View struct {
	RootDir string
	Layout  string
	Writer  io.Writer
}

func NewView(w io.Writer) View {
	return View{
		RootDir: "views",
		Layout:  "views/layouts/base.html",
		Writer:  w,
	}
}

func (this View) Rend(viewFile string, data interface{}) {
	/**
	  	data example:

	         map[string]interface{}{
	          "user": GetUser(r),
	          "someotherdata": someStructWithData,
	          }
	  **/
	viewFile = this.RootDir + "/" + viewFile

	if _, err := os.Open(viewFile); err != nil {
		viewFilePattern := viewFile + ".[0-9A-Za-z]*"

		// fmt.Println(viewFilePattern)
		// os.Exit(1)

		matches, err2 := filepath.Glob(viewFilePattern)
		if err2 != nil {
			// 未找到匹配的视图文件
			panic("view file not exists: " + viewFile)
		}
		viewFile = matches[0]
	}
	//  layout 和 视图文件拼接
	viewFiles := []string{this.Layout, viewFile}

	// log.Println(viewFiles)

	t := template.Must(template.ParseFiles(viewFiles...))
	t.Execute(this.Writer, data)
}
