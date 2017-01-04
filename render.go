package tigo

import (
	"io"
)
//Map for data with map[string]interface{}
type M map[string]interface{}

type Render interface {
	//Init
	Init() error
	//Render layout and name
	Render(out io.Writer, name string, data interface{}) error
	//Render only file.
	RenderFile(out io.Writer, name string, data interface{}) error
}

