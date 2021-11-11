package util

import "embed"

// content is our static web server content.
//go:embed content/*
var Content embed.FS
