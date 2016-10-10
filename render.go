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
)
//Map for data with map[string]interface{}
type M map[string]interface{}

type Render interface {
	//Init
	Init() error
	//Render layout and name
	Render(out io.Writer, name string, data interface{}) error
	//Render only file without layout
	RenderFile(out io.Writer, name string, data interface{}) error
}

// HtmlRender implements Render interface, but based on golang templates.
type HtmlRender struct {
	viewRoot string
	master   string
	ext      string
	template *template.Template
}

type HtmlRenderConfig struct {
	ViewRoot   string
	MasterPage string
	Extension  string
}

var defaultFuncs = template.FuncMap{
	"content": func() (string, error) {
		return ">>>Error:content tag not support!!!<<<", nil
	},
	"render": func(partialName string) (template.HTML, error) {
		return template.HTML(fmt.Sprintf(">>>Error:render %v tag not support !!!", partialName)), nil
	},
}

const (
	maxRenderFileNum = 20        //max render one file times, to prevent the infinite loop.
)

//NewHtmlRender returns a default render with templates loaded from viewRoot
func NewHtmlRender(config HtmlRenderConfig) Render {
	return &HtmlRender{
		viewRoot: config.ViewRoot,
		master: config.MasterPage,
		ext: config.Extension,
		template:    template.New(filepath.Base(config.ViewRoot)),
	}
}

// Init for initialize, when running, this method is executed.
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
		if s.ext != extension {
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

	return nil
}

// Render executes template named name, passing data as context, the output is written to out.
// This method render use layout with master page.
func (s *HtmlRender) Render(out io.Writer, name string, data interface{}) error {
	if s.master == "" {
		return fmt.Errorf("master page not exist in views root: %v.", s.viewRoot)
	}
	renderPartialNames := make(map[string]int, 0)
	return s.executeTemplate(out, s.master, data, template.FuncMap{
		"content": func() (template.HTML, error) {
			buff, err := s.executeTemplateBuf(name, data)
			if err != nil {
				return "", err
			}
			return template.HTML(buff.String()), nil
		},
		"render": func(partialName string) (template.HTML, error) {
			renderPartialNames[partialName] = renderPartialNames[partialName] + 1
			if renderPartialNames[partialName] > maxRenderFileNum {
				return "", fmt.Errorf("render cycle error, render \"%v\" max allow %v times.", partialName, maxRenderFileNum)
			}
			if s.template.Lookup(partialName) != nil {
				buf, err := s.executeTemplateBuf(partialName, data)
				return template.HTML(buf.String()), err
			}
			return "", nil
		},
	})
}

// Render executes template named name, passing data as context, the output is written to out.
// This method render with no layout.
func (s *HtmlRender) RenderFile(out io.Writer, name string, data interface{}) error {
	renderPartialNames := make(map[string]int, 0)
	return s.executeTemplate(out, name, data, template.FuncMap{
		"content": defaultFuncs["content"],
		"render": func(partialName string) (template.HTML, error) {
			renderPartialNames[partialName] = renderPartialNames[partialName] + 1
			if renderPartialNames[partialName] > maxRenderFileNum {
				return "", fmt.Errorf("render cycle error, render \"%v\" max allow %v times.", partialName, maxRenderFileNum)
			}
			if s.template.Lookup(partialName) != nil {
				buf, err := s.executeTemplateBuf(partialName, data)
				return template.HTML(buf.String()), err
			}
			return "", nil
		},
	})
}

func (s *HtmlRender) executeTemplate(out io.Writer, name string, data interface{}, funcs template.FuncMap) error {
	if tpl := s.template.Lookup(name); tpl != nil {
		return tpl.Funcs(funcs).Execute(out, data)
	}
	return fmt.Errorf("template:%v not found!", name)
}

func (s *HtmlRender) executeTemplateBuf(name string, data interface{}) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := s.template.Funcs(template.FuncMap{
		"content": defaultFuncs["content"],
	}).ExecuteTemplate(buf, name, data)
	return buf, err
}