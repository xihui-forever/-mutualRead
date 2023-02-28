package cmd

type Cmd struct {
	Path  string
	Role  int
	Logic interface{}
}

var CmdList = []Cmd{}
