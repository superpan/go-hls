package main

import (
	"log"
	"os"

	"github.com/superpan/go-hls/hls"
	"github.com/urfave/cli"
)

// DownloadCli is the cli wrapper for download functionality
func DownloadCli(c *cli.Context) error {
	if c.String("url") == "" {
		log.Fatal("url required")
	}

	if c.String("output") == "" {
		log.Fatal("output required")
	}
	return hls.Download(c.String("url"), c.String("output"))
}

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url",
			Usage: "m3u8 media url",
		},
		cli.StringFlag{
			Name:  "output",
			Usage: "path to output file",
		},
	}
	app.Action = DownloadCli

	app.Run(os.Args)
}
