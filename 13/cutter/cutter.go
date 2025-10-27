package cutter

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/urfave/cli/v3"
)

type Config struct {
	FieldsSpec    string
	Delimiter     string
	OnlySeparated bool
	InputFile     string
}

func ParseConfig(cmd *cli.Command) (*Config, error) {
	fields := cmd.String("f")
	if strings.TrimSpace(fields) == "" {
		return nil, fmt.Errorf("-f is required")
	}

	cfg := &Config{
		FieldsSpec:    fields,
		Delimiter:     cmd.String("d"),
		OnlySeparated: cmd.Bool("s"),
		InputFile:     cmd.Args().Get(0),
	}

	r, size := utf8.DecodeRuneInString(cfg.Delimiter)
	if r == utf8.RuneError || size != len(cfg.Delimiter) {
		return nil, fmt.Errorf("-d must be a single character")
	}

	if cfg.Delimiter == "" {
		cfg.Delimiter = "\t"
	}
	return cfg, nil
}

type Cutter struct {
	cfg    *Config
	writer *bufio.Writer
}

func NewCutter(cfg *Config) *Cutter {
	return &Cutter{
		cfg:    cfg,
		writer: bufio.NewWriter(os.Stdout),
	}
}

func (c *Cutter) Cut(reader io.Reader) error {
	defer c.writer.Flush()

	ranges, err := parseFieldsSpec(c.cfg.FieldsSpec)
	if err != nil {
		return err
	}

	delimBytes := []byte(c.cfg.Delimiter)

	scanner := bufio.NewScanner(reader)
	const maxCapacity = 10 * 1024 * 1024
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		line := scanner.Bytes()
		hasDelim := bytes.Contains(line, delimBytes)
		if !hasDelim {
			if c.cfg.OnlySeparated || len(line) == 0 {
				continue
			}

			if _, err := c.writer.Write(line); err != nil {
				return err
			}

			if err := c.writer.WriteByte('\n'); err != nil {
				return err
			}

			continue
		}

		parts := bytes.Split(line, delimBytes)
		selected := make([][]byte, 0, len(parts))
		// ????
		for _, rg := range ranges {
			for i := rg.start; i <= rg.end; i++ {
				idx := i - 1
				if idx >= 0 && idx < len(parts) {
					selected = append(selected, parts[idx])
				}
			}
		}

		if len(selected) == 0 {
			continue
		}

		var out []byte
		if len(selected) == 1 {
			out = selected[0]
		} else {
			out = bytes.Join(selected, delimBytes)
		}

		if len(out) == 0 {
			continue
		}

		if _, err := c.writer.Write(out); err != nil {
			return err
		}

		if err := c.writer.WriteByte('\n'); err != nil {
			return err
		}
	}

	return scanner.Err()
}

type fieldRange struct {
	start int
	end   int
}

func parseFieldsSpec(spec string) ([]fieldRange, error) {
	trimmed := strings.TrimSpace(spec)
	if trimmed == "" {
		return nil, errors.New("fields specification is empty")
	}

	tokens := strings.Split(trimmed, ",")
	ranges := make([]fieldRange, 0, len(tokens))
	for _, tok := range tokens {
		t := strings.TrimSpace(tok)
		if t == "" {
			continue
		}

		if strings.Contains(t, "-") {
			parts := strings.SplitN(t, "-", 2)
			if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
				continue
			}

			var s, e int
			s, errS := strconv.Atoi(parts[0])
			e, errE := strconv.Atoi(parts[1])
			if errS != nil || errE != nil || s <= 0 || e <= 0 || s > e {
				continue
			}

			ranges = append(ranges, fieldRange{start: s, end: e})
			continue
		}

		var n int
		if _, err := fmt.Sscanf(t, "%d", &n); err == nil && n > 0 {
			ranges = append(ranges, fieldRange{start: n, end: n})
		}
	}

	if len(ranges) == 0 {
		return nil, errors.New("no valid fields specified")
	}

	return ranges, nil
}