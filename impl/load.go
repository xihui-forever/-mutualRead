package impl

type Cmd struct {
	Path  string
	Role  int
	Logic interface{} // func(ctx, req) (resp, err)
}

var CmdList = []Cmd{}
