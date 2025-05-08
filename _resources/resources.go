package _resources

import (
	_ "embed"
)

//go:embed fonts/Roboto.ttf
var FontRoboto []byte

var BuildHash string
var BuildPipelineLink string
