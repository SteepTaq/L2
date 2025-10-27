package grepper

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli/v3"
)

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

func ParseConfig(cmd *cli.Command) (*Config, error) {
	after := cmd.Int("A")
	before := cmd.Int("B")

	context := cmd.Int("C")
	if context > 0 {
		after = context
		before = context
	}

	pattern := cmd.Args().Get(0)
	if pattern == "" {
		return nil, fmt.Errorf("pattern is required")
	}

	cfg := &Config{
		After:       after,
		Before:      before,
		Context:     context,
		Count:       cmd.Bool("c"),
		IgnoreCase:  cmd.Bool("i"),
		Invert:      cmd.Bool("v"),
		Fixed:       cmd.Bool("F"),
		LineNumbers: cmd.Bool("n"),
		Pattern:     pattern,
		InputFile:   cmd.Args().Get(1),
	}

	return cfg, nil
}

type Grepper struct {
	cfg     *Config
	matcher func(string) bool
	writer  *bufio.Writer
}

func NewGrepper(cfg *Config) (*Grepper, error) {
	g := &Grepper{
		cfg:    cfg,
		writer: bufio.NewWriter(os.Stdout),
	}

	var baseMatcher func(string) bool

	pattern := cfg.Pattern
	if cfg.Fixed {
		if cfg.IgnoreCase {
			pattern = strings.ToLower(pattern)
			baseMatcher = func(s string) bool {
				return strings.Contains(strings.ToLower(s), pattern)
			}
		} else {
			baseMatcher = func(s string) bool {
				return strings.Contains(s, pattern)
			}
		}
	} else {
		if cfg.IgnoreCase {
			pattern = "(?i)" + pattern
		}

		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern: %w", err)
		}
		baseMatcher = re.MatchString
	}

	if cfg.Invert {
		g.matcher = func(s string) bool {
			return !baseMatcher(s)
		}
	} else {
		g.matcher = baseMatcher
	}

	return g, nil
}

type lineInfo struct {
	number int
	text   string
}

func (g *Grepper) Grep(reader io.Reader) error {
	defer g.writer.Flush()

	scanner := bufio.NewScanner(reader)

	if g.cfg.Count {
		count := 0
		for scanner.Scan() {
			if g.matcher(scanner.Text()) {
				count++
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}

		fmt.Println(count)
		return nil
	}

	beforeBuffer := make([]lineInfo, 0, g.cfg.Before)
	afterCounter := 0
	lineNumber := 0
	lastPrintedLineNumber := 0

	for scanner.Scan() {
		lineNumber++
		currentLine := scanner.Text()

		if g.matcher(currentLine) {
			toSeparate := lastPrintedLineNumber > 0 &&
				lineNumber > lastPrintedLineNumber+1 &&
				g.cfg.Before > 0

			if toSeparate {
				fmt.Println("-----")
			}

			for _, li := range beforeBuffer {
				if li.number > lastPrintedLineNumber {
					g.print(li.number, li.text)
				}
			}
			beforeBuffer = make([]lineInfo, 0, g.cfg.Before)

			g.print(lineNumber, currentLine)
			lastPrintedLineNumber = lineNumber

			afterCounter = g.cfg.After
		} else if afterCounter > 0 {
			g.print(lineNumber, currentLine)
			afterCounter--
			lastPrintedLineNumber = lineNumber
		} else if g.cfg.Before > 0 {
			beforeBuffer = append(beforeBuffer, lineInfo{number: lineNumber, text: currentLine})
			if len(beforeBuffer) > g.cfg.Before {
				beforeBuffer = beforeBuffer[1:]
			}
		}
	}

	return scanner.Err()
}

func (g *Grepper) print(num int, line string) {
	if g.cfg.LineNumbers {
		fmt.Printf("%d: %s\n", num, line)
	} else {
		fmt.Println(line)
	}
}