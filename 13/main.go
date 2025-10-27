package main

import (
	"context"
	"io"
	"log"
	"os"

	"wb_l2/13/cutter"

	"github.com/urfave/cli/v3"
)

func main() {
	cli.HelpFlag = nil

	cmd := &cli.Command{
		Name:                   "cut",
		Description:            "simpler analogue of UNIX-utility cut",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "f",
				Usage:    "fields to output",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "d",
				Usage: "delimiter character",
				Value: "\t",
			},
			&cli.BoolFlag{
				Name:  "s",
				Usage: "only print lines that contain the delimiter",
			},
		},
		Action: action,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(ctx context.Context, cmd *cli.Command) error {
	cfg, err := cutter.ParseConfig(cmd)
	if err != nil {
		return err
	}

	var reader io.Reader
	if cfg.InputFile != "" {
		file, err := os.Open(cfg.InputFile)
		if err != nil {
			return err
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}

	c := cutter.NewCutter(cfg)
	return c.Cut(reader)
}