package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownload(t *testing.T) {
	url := "http://devimages.apple.com/iphone/samples/bipbop/gear1/prog_index.m3u8"
	file := "/tmp/test.ts"

	err := Download(url, file)
	assert.Nil(t, err)
}
