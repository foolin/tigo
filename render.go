package tigo

import (
	"io"
	"path/filepath"
	"strings"
	"os"
	"fmt"
	"io/ioutil"
	"html/template"
	"bytes"
	"log"
)

type Render interface {
	//init
	Init() error
	//Render layout and name
	Render(out io.Writer, name string, data interface{}) error
	//Render only file
	RenderFile(out io.Writer, name string, data interface{}) error
}

// HtmlRender implements Render interface, but based on golang templates.
type HtmlRender struct {
	viewRoot string
	master   string
	ext   string
	template *template.Template
}

type HtmlRenderConfig struct {
	ViewRoot   string
	MasterPage string
	Extension     string
}

var defaultFuncs = template.FuncMap{
	"content": func() (string, error) {
		return ">>>Error:content tag not support!!!", nil
	},
	"render": func(partialName string) (template.HTML, error) {
		return template.HTML(fmt.Sprintf(">>>Error:render %v tag not support !!!", partialName)), nil
	},
}

//NewSimpleView returns a SimpleView with templates loaded from viewDir
func NewHtmlRender(config ...HtmlRenderConfig) Render {
	var c HtmlRenderConfig
	if len(config) > 0{
		c = config[0]
	}else {
		c = HtmlRenderConfig{
			ViewRoot: "views",
			MasterPage: "master",
			Extension: ".html",
		}
	}
	return &HtmlRender{
		viewRoot: c.ViewRoot,
		master: c.MasterPage,
		ext: c.Extension,
		template:    template.New(filepath.Base(c.ViewRoot)),
	}
}

func (s *HtmlRender) Init() error {
	info, err := os.Stat(s.viewRoot)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("tigo: view root:%s is not a directory", s.viewRoot)
	}
	werr := filepath.Walk(s.viewRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		extension := filepath.Ext(path)
		if s.ext != extension{
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		// We remove the directory name from the path
		// this means if we have directory foo, with file bar.tpl
		// full path for bar file foo/bar.tpl
		// we trim the foo part and remain with /bar.tpl
		//
		// NOTE we don't account for the opening slash, when dir ends with /.
		name := path[len(s.viewRoot):]

		name = filepath.ToSlash(name)

		name = strings.TrimPrefix(name, "/") // case  we missed the opening slash

		name = strings.TrimSuffix(name, extension) // remove extension

		t := s.template.New(name)
		_, err = t.Funcs(defaultFuncs).Parse(string(data))
		if err != nil {
			return err
		}
		return nil
	})

	if werr != nil {
		return werr
	}

	////master
	//if s.master == "" {
	//	for _, ext := range supported {
	//		if tpl := s.template.Lookup(fmt.Sprintf("master%v", ext)); tpl != nil {
	//			s.master = tpl.Name()
	//			break
	//		}
	//	}
	//}

	return nil
}

// Render executes template named name, passing data as context, the output is written to out.
func (s *HtmlRender) Render(out io.Writer, name string, data interface{}) error {
	if s.master == "" {
		return fmt.Errorf("master page not exist in views root: %v.", s.viewRoot)
	}
	renderPartialNames := make(map[string]bool, 0)
	return s.executeTemplate(out, s.master, data, template.FuncMap{
		"content": func() (template.HTML, error) {
			buff, err := s.executeTemplateBuf(name, data)
			if err != nil {
				return "", err
			}
			return template.HTML(buff.String()), nil
		},
		"render": func(partialName string) (template.HTML, error) {
			if renderPartialNames[partialName] == true{
				return "", fmt.Errorf("render cycle error, render \"%v\"", partialName)
			}
			if s.template.Lookup(partialName) != nil {
				renderPartialNames[partialName] = true
				buf, err := s.executeTemplateBuf(partialName, data)
				return template.HTML(buf.String()), err
			}
			return "", nil
		},
	})
}

// Render executes template named name, passing data as context, the output is written to out.
func (s *HtmlRender) RenderFile(out io.Writer, name string, data interface{}) error {
	log.Printf("renderFile: %v", name)
	return s.executeTemplate(out, name, data, template.FuncMap{
		"content": defaultFuncs["content"],
		"render": func(partialName string) (template.HTML, error) {
			if s.template.Lookup(partialName) != nil {
				buf, err := s.executeTemplateBuf(partialName, data)
				return template.HTML(buf.String()), err
			}
			return "", nil
		},
	})
}

func (s *HtmlRender) executeTemplate(out io.Writer, name string, data interface{}, funcs template.FuncMap) error {
	log.Printf("executeTemplate: %v", name)
	if tpl := s.template.Lookup(name); tpl != nil {
		return tpl.Funcs(funcs).Execute(out, data)
	}
	return fmt.Errorf("template:%v not found!", name)
}

func (s *HtmlRender) executeTemplateBuf(name string, data interface{}) (*bytes.Buffer, error) {
	log.Printf("executeTemplateBuf: %v", name)
	buf := new(bytes.Buffer)
	err := s.template.Funcs(template.FuncMap{
		"content": defaultFuncs["content"],
	}).ExecuteTemplate(buf, name, data)
	return buf, err
}