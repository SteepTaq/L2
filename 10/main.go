package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"wb_l2/10/sorter"

	"github.com/urfave/cli/v3"
)

func main() {
	// Убираем флаг help
	cli.HelpFlag = nil

	cmd := &cli.Command{
		Name:                   "sort",
		Description:            "simpler analogue of UNIX-utility sort",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			// Флаги для сортировки
			&cli.IntFlag{Name: "k", Usage: "key column"},
			&cli.BoolFlag{Name: "n", Usage: "numeric sort"},
			&cli.BoolFlag{Name: "r", Usage: "reverse sort"},
			&cli.BoolFlag{Name: "u", Usage: "unique sort"},
			&cli.BoolFlag{Name: "M", Usage: "month sort"},
			&cli.BoolFlag{Name: "b", Usage: "ignore trailing blanks"},
			&cli.BoolFlag{Name: "c", Usage: "check if sorted"},
			&cli.BoolFlag{Name: "h", Usage: "human readable numbers sort"},
		},
		Action: action,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(ctx context.Context, cmd *cli.Command) error {
	// На основе флагов формируем структуру конфигурации
	cfg, err := sorter.ParseConfig(cmd)
	if err != nil {
		return err
	}

	// Считываем входные данные
	lines, err := readInput(cfg.InputFile)
	if err != nil {
		return err
	}

	// Создаем структуру Sorter
	sorter := sorter.NewSorter(lines, cfg)

	// Если передан флаг -c, проверяем, отсортированы ли данные
	if cfg.CheckSorted {
		if err := sorter.CheckSorted(); err != nil {
			return err
		}
	}

	// Сортируем данные
	sorter.Sort()

	// Если передан флаг -u, фильтруем уникальные строки
	sorted := sorter.Lines
	if cfg.Uniquer {
		sorted = sorter.FilterUnique()
	}

	// Выводим отсортированные данные
	for _, line := range sorted {
		fmt.Println(line)
	}

	return nil
}

func readInput(fileName string) ([]string, error) {
	if fileName == "" {
		scanner := bufio.NewScanner(os.Stdin)

		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		return lines, nil
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}
