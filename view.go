package tigo

import (
	"html/template"
	"sync"
	"fmt"
	"os"
	"io"
	"path/filepath"
	"io/ioutil"
	"bytes"
)

type ViewRender struct {
	config   ViewRenderConfig
	tplMap   map[string]*template.Template
	tplMutex sync.RWMutex
}

type ViewRenderConfig struct {
	Root              string           //view root
	Extension         string           //template extension
	Master            string           //template master
	Partials          []string         //template partial, such as head, foot
	Funcs             template.FuncMap //template functions
	DisableCache      bool             //disable cache, debug mode
	DisableFilePartial bool             //enable render file use partial
}

func NewViewRender(config ViewRenderConfig) *ViewRender {
	return &ViewRender{
		config: config,
		tplMap: make(map[string]*template.Template),
		tplMutex: sync.RWMutex{},
	}
}

func (r *ViewRender) Init() error {
	//if r.config.Root != "" {
	//	if _, err := os.Stat(r.config.Root); os.IsNotExist(err) {
	//		return fmt.Errorf("ViewRender view root: %v not exists!", r.config.Root)
	//	}
	//}
	return nil
}


// Render a template to the screen
func (r *ViewRender) RenderFile(out io.Writer, name string, data interface{}) error {
	return r.execute(out, name, data, false)
}

// Render a template to the screen
func (r *ViewRender) Render(out io.Writer, name string, data interface{}) error {
	return r.execute(out, name, data, true)
}

func (r *ViewRender) execute(out io.Writer, name string, data interface{}, useMaster bool) error {
	var tpl *template.Template
	var err error
	var ok bool

	allFuncs := make(template.FuncMap, 0)
	allFuncs["include"] = func(layout string) (template.HTML, error) {
		buf := new(bytes.Buffer)
		err := r.execute(buf, layout, data, false)
		return template.HTML(buf.String()), err
	}

	// Get the plugin collection
	for k, v := range r.config.Funcs {
		allFuncs[k] = v
	}

	r.tplMutex.RLock()
	tpl, ok = r.tplMap[name]
	r.tplMutex.RUnlock()

	exeName := name
	if useMaster && r.config.Master != "" {
		exeName = r.config.Master
	}

	if !ok || r.config.DisableCache {

		tplList := make([]string, 0)
		if useMaster {
			//render()
			if r.config.Master != "" {
				tplList = append(tplList, r.config.Master)
			}
		}
		tplList = append(tplList, name)
		tplList = append(tplList, r.config.Partials...)

		// Loop through each template and test the full path
		tpl = template.New(name).Funcs(allFuncs)
		for _, v := range tplList {
			// Get the absolute path of the root template
			path, err := filepath.Abs(r.config.Root + string(os.PathSeparator) + v + r.config.Extension)
			if err != nil {
				return fmt.Errorf("TemplateEngine path:%v error: %v", path, err)
			}
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("TemplateEngine render read name:%v, path:%v, error: %v", v, path, err)
			}
			var tmpl *template.Template
			if v == name {
				tmpl = tpl
			} else {
				tmpl = tpl.New(v)
			}
			_, err = tmpl.Parse(string(data))
			if err != nil {
				return fmt.Errorf("TemplateEngine render parser name:%v, path:%v, error: %v", v, path, err)
			}
		}
		r.tplMutex.Lock()
		r.tplMap[name] = tpl
		r.tplMutex.Unlock()
	}

	// Display the content to the screen
	err = tpl.Funcs(allFuncs).ExecuteTemplate(out, exeName, data)
	if err != nil {
		return fmt.Errorf("ViewRender execute template error: %v", err)
	}

	return nil
}