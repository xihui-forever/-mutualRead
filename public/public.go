package public

import (
	"embed"
)

//go:embed build/*
var Public embed.FS
