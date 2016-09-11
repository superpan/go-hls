package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownload(t *testing.T) {
	url := "http://videos-f.jwpsrv.com/content/conversions/zWLy8Jer/videos/21ETjILN-364765.mp4.m3u8?token=0_57d4d2eb_0x5273357cb28ea315190a086835a97c325c63202c"
	file := "/tmp/test.ts"

	err := Download(url, file)
	assert.Nil(t, err)
}
