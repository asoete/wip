package embedder

import (
	"embed"
)

//go:embed images/*
var Images embed.FS
