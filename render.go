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
	"path"
)
//Map for data with map[string]interface{}
type M map[string]interface{}

type Render interface {
	//Init
	Init() error
	//Render layout and name
	Render(out io.Writer, name string, data interface{}) error
}

// HtmlRender implements Render interface, but based on golang templates.
// {{layout "layout/master"}} use master page.
// {{render "layout/header"}} render sub page
// {{content}} only use layout page, this tag for render content.
type HtmlRender struct {
	viewRoot string
	ext      string
	template *template.Template
	funcs    template.FuncMap
}

type HtmlRenderConfig struct {
	ViewRoot  string
	Extension string
	Funcs     template.FuncMap
}

var emptyFuncs = template.FuncMap{
	"content": func() (template.HTML, error) {
		return ">>>Error:{content} tag only support at layout page!!!<<<", nil
	},
	"layout": func(layoutName string) (template.HTML, error) {
		return template.HTML(fmt.Sprintf(">>>Error:{layout %v} tag not support !!!<<<", layoutName)), nil
	},
	"render": func(partialName string) (template.HTML, error) {
		return template.HTML(fmt.Sprintf(">>>Error:{render %v} tag not support !!!<<<", partialName)), nil
	},
}

const (
	maxRenderFileNum = 20        //max render one file times, to prevent the infinite loop.
)

//NewHtmlRender returns a default render with templates loaded from viewRoot
func NewHtmlRender(config HtmlRenderConfig) Render {
	return &HtmlRender{
		viewRoot: config.ViewRoot,
		ext: config.Extension,
		template: nil,
		funcs: config.Funcs,
	}
}

//interface
func (r *HtmlRender) Init() error {
	info, err := os.Stat(r.viewRoot)
	if err != nil {
		//return fmt.Errorf("tigo: view root:%s is not a directory.", r.viewRoot)
		return nil
	}
	if !info.IsDir() {
		//return fmt.Errorf("tigo: view root:%s is not a directory", r.viewRoot)
		return nil
	}
	allFuncs := template.FuncMap{}
	for k, v := range emptyFuncs {
		allFuncs[k] = v
	}
	for k, v := range r.funcs {
		allFuncs[k] = v
	}
	r.template = template.New(filepath.Base(r.viewRoot))
	werr := filepath.Walk(r.viewRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("tigo: view root:%s, error: %v", r.viewRoot, err)
		}
		if info.IsDir() {
			return nil
		}

		extension := filepath.Ext(path)
		if r.ext != extension {
			return nil
		}

		// We remove the directory name from the path
		// this means if we have directory foo, with file bar.tpl
		// full path for bar file foo/bar.tpl
		// we trim the foo part and remain with /bar.tpl
		//
		// NOTE we don't account for the opening slash, when dir ends with /.
		name := path[len(r.viewRoot):]
		name = filepath.ToSlash(name)
		name = strings.TrimPrefix(name, "/") // case  we missed the opening slash
		name = strings.TrimSuffix(name, extension) // remove extension

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("tigo: view root:%s, error: %v", r.viewRoot, err)
		}
		t := r.template.New(name)
		content := fmt.Sprintf("%s", data)
		_, err = t.Funcs(allFuncs).Parse(content)
		if err != nil {
			return fmt.Errorf("tigo: view root:%s, error: %v", r.viewRoot, err)
		}
		return nil
	})

	if werr != nil {
		return fmt.Errorf("tigo: view root:%s, error: %v", r.viewRoot, werr)
	}

	return nil
}

func (r *HtmlRender) Render(out io.Writer, name string, data interface{}) error {
	err := r.executeRender(out, name, data)
	if err != nil {
		return fmt.Errorf("HtmlRender render \"%v\" happen error, %v", path.Join(r.viewRoot, name + r.ext), err)
	}
	return nil
}

// Render executes template named name, passing data as context, the output is written to out.
func (r *HtmlRender) executeRender(out io.Writer, name string, data interface{}) error {
	var masterName string
	var renderTimes map[string]int
	var allFuncs = template.FuncMap{
		"content": emptyFuncs["content"],
		"layout": func(layoutName string) (template.HTML, error) {
			masterName = layoutName
			return "", nil
		},
		"render": func(partialName string) (template.HTML, error) {
			renderTimes[partialName] = renderTimes[partialName] + 1
			if renderTimes[partialName] > maxRenderFileNum {
				return "", fmt.Errorf("render cycle error, render \"%v\" max allow %v times.", partialName, maxRenderFileNum)
			}
			if r.template.Lookup(partialName) != nil {
				buf, err := r.executeTemplateBuf(partialName, data, nil)
				return template.HTML(buf.String()), err
			}
			return "", nil
		},
	}
	//执行页面
	renderTimes = make(map[string]int, 0)
	buf, err := r.executeTemplateBuf(name, data, allFuncs)
	if err != nil {
		return err
	}
	if masterName == "" {
		//直接输出
		_, err = out.Write(buf.Bytes())
		return err
	}

	//执行母版页
	allFuncs["content"] = func() (template.HTML, error) {
		return template.HTML(buf.Bytes()), nil
	}
	//如果含有layout，则执行
	renderTimes = make(map[string]int, 0)
	return r.executeTemplateRaw(out, masterName, data, allFuncs)
}

func (r *HtmlRender) executeTemplateRaw(out io.Writer, name string, data interface{}, funcs template.FuncMap) error {
	allFuncs := template.FuncMap{}
	for k, v := range funcs {
		allFuncs[k] = v
	}
	for k, v := range r.funcs {
		allFuncs[k] = v
	}
	return r.template.Funcs(allFuncs).ExecuteTemplate(out, name, data)
}

func (r *HtmlRender) executeTemplateBuf(name string, data interface{}, funcs template.FuncMap) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := r.executeTemplateRaw(buf, name, data, funcs)
	return buf, err
}