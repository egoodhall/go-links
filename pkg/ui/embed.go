package ui

import "embed"

//go:embed templates/*.go.html
var templates embed.FS

//go:embed static/*
var StaticFiles embed.FS
