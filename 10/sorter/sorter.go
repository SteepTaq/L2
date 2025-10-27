package sorter

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/urfave/cli/v3"
)

// Маппинг месяцев на числа
var months = map[string]int{
	"JAN": 1, "FEB": 2, "MAR": 3, "APR": 4,
	"MAY": 5, "JUN": 6, "JUL": 7, "AUG": 8,
	"SEP": 9, "OCT": 10, "NOV": 11, "DEC": 12,
}

// Маппинг суффиксов для человекочитаемых размеров
var humanSuffixes = map[rune]int{
	'K': 1024,
	'M': 1024 * 1024,
	'G': 1024 * 1024 * 1024,
	'T': 1024 * 1024 * 1024 * 1024,
}

// Парсинг чисел с суффиксами
func parseHumanReadable(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}

	lastChar := rune(s[len(s)-1])
	multiplier, exists := humanSuffixes[unicode.ToUpper(lastChar)]

	if exists {
		numStr := s[:len(s)-1]
		num, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return 0, err
		}

		return int64(num * float64(multiplier)), nil
	}

	// Если суффикса нет, пытаемся парсить как обычное число
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return int64(num), nil
}

type Config struct {
	// Сортировка по столбцу
	KeyColumn int
	// Сортировка по числам
	Numeric bool
	// Сортировка в обратном порядке
	Reverse bool
	// Сортировка уникальных строк
	Uniquer bool
	// Сортировка по месяцам
	MonthSort bool
	// Игнорирование пробелов
	IgnoreTrailingBlanks bool
	// Проверка, отсортированы ли данные
	CheckSorted bool
	// Сортировка по числам с суффиксами
	HumanReadable bool
	// Имя входного файла
	InputFile string
}

func ParseConfig(cmd *cli.Command) (*Config, error) {
	k := cmd.Int("k")
	if k < 0 {
		return nil, fmt.Errorf("key column must be positive or zero")
	}

	cfg := &Config{
		KeyColumn:            k,
		Numeric:              cmd.Bool("n"),
		Reverse:              cmd.Bool("r"),
		Uniquer:              cmd.Bool("u"),
		MonthSort:            cmd.Bool("M"),
		IgnoreTrailingBlanks: cmd.Bool("b"),
		CheckSorted:          cmd.Bool("c"),
		HumanReadable:        cmd.Bool("h"),
		InputFile:            cmd.Args().Get(0),
	}

	return cfg, nil
}

type Sorter struct {
	Lines  []string
	Config *Config
}

func NewSorter(lines []string, config *Config) *Sorter {
	return &Sorter{
		Lines:  lines,
		Config: config,
	}
}

func (s *Sorter) Sort() {
	sort.Slice(s.Lines, s.less)
}

// Сравнение строк
func (s *Sorter) less(i, j int) bool {
	lineI := s.Lines[i]
	lineJ := s.Lines[j]

	// Получаем нужные части строк для сравнения
	strI := s.getTarget(lineI)
	strJ := s.getTarget(lineJ)

	isLess := false

	// Сортировка в зависимости от флагов
	switch {
	case s.Config.HumanReadable:
		numI, errI := parseHumanReadable(strI)
		numJ, errJ := parseHumanReadable(strJ)

		if errI != nil {
			isLess = true
			break
		}
		if errJ != nil {
			isLess = false
			break
		}

		isLess = numI < numJ
	case s.Config.MonthSort:
		monthI, okI := months[strings.ToUpper(strI)]
		monthJ, okJ := months[strings.ToUpper(strJ)]

		if !okI {
			isLess = true
			break
		}
		if !okJ {
			isLess = false
			break
		}

		isLess = monthI < monthJ
	case s.Config.Numeric:
		numI, errI := strconv.ParseFloat(strI, 64)
		numJ, errJ := strconv.ParseFloat(strJ, 64)

		if errI != nil && errJ != nil {
			isLess = strI < strJ
			break
		} else if errI != nil {
			isLess = true
			break
		} else if errJ != nil {
			isLess = false
			break
		}

		isLess = numI < numJ
	default:
		isLess = strI < strJ
	}

	// Сортировка в обратном порядке
	if s.Config.Reverse {
		return !isLess
	}

	return isLess
}

func (s *Sorter) getTarget(line string) string {
	var key string

	if s.Config.KeyColumn == 0 {
		key = line
	} else {
		columns := strings.Split(line, "\t")
		i := s.Config.KeyColumn - 1

		if i >= 0 && i < len(columns) {
			key = columns[i]
		} else {
			key = ""
		}
	}

	if s.Config.IgnoreTrailingBlanks {
		key = strings.TrimSpace(key)
	}

	return key
}

func (s *Sorter) CheckSorted() error {
	fileName := s.Config.InputFile
	if fileName == "" {
		fileName = "stdin"
	}

	for i := 1; i < len(s.Lines); i++ {
		// Проверяем, что текущая строка больше или равна предыдущей
		if s.less(i, i-1) {
			return fmt.Errorf("sort: %s:%d: disorder: %s", fileName, i+1, s.Lines[i])
		}
	}

	return nil
}

func (s *Sorter) FilterUnique() []string {
	if len(s.Lines) == 0 {
		return s.Lines
	}

	result := make([]string, 0, len(s.Lines))
	result = append(result, s.Lines[0])
	lastStr := s.getTarget(s.Lines[0])

	for i := 1; i < len(s.Lines); i++ {
		currentStr := s.getTarget(s.Lines[i])

		if currentStr != lastStr {
			result = append(result, s.Lines[i])
			lastStr = currentStr
		}
	}

	return result
}
