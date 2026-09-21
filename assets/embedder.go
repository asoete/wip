package embedder

import (
	"embed"
)

//go:embed fonts/*
var Fonts embed.FS

//go:embed images/*
var Images embed.FS
