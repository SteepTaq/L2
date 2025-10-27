package main

import (
	"context"
	"io"
	"log"
	"os"

	"wb_l2/12/grepper"

	"github.com/urfave/cli/v3"
)

func main() {
	cli.HelpFlag = nil

	cmd := &cli.Command{
		Name:                   "grep",
		Description:            "simpler analogue of UNIX-utility grep",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.IntFlag{Name: "A", Usage: "print N lines after match"},
			&cli.IntFlag{Name: "B", Usage: "print N lines before match"},
			&cli.IntFlag{Name: "C", Usage: "print N lines around match"},
			&cli.BoolFlag{Name: "c", Usage: "print only count of matches"},
			&cli.BoolFlag{Name: "i", Usage: "ignore case"},
			&cli.BoolFlag{Name: "v", Usage: "invert matching"},
			&cli.BoolFlag{Name: "F", Usage: "pattern as a fixed string"},
			&cli.BoolFlag{Name: "n", Usage: "show line numbers"},
		},
		Action: action,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	After       int
	Before      int
	Context     int
	Count       bool
	IgnoreCase  bool
	Invert      bool
	Fixed       bool
	LineNumbers bool
	Pattern     string
	InputFile   string
}

func action(ctx context.Context, cmd *cli.Command) error {
	cfg, err := grepper.ParseConfig(cmd)
	if err != nil {
		return err
	}

	grepper, err := grepper.NewGrepper(cfg)
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

	return grepper.Grep(reader)
}
