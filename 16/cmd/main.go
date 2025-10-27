package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"wb_l2/16/src/config"
	"wb_l2/16/src/downloader"
)

var (
	maxDepth  = flag.Int("depth", 1, "Максимальная глубина рекурсии")
	outputDir = flag.String("output", "./downloaded", "Выходная директория")
	maxConc   = flag.Int("concurrent", 5, "Максимальное количество одновременных загрузок")
	timeout   = flag.Duration("timeout", 30*time.Second, "Таймаут для HTTP запросов")
	verbose   = flag.Bool("verbose", false, "Подробное логирование")
)

func init() {
	flag.Parse()
}

func main() {
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	urlStr := flag.Arg(0)

	cfg := config.Config{
		MaxDepth:      *maxDepth,
		MaxConcurrent: *maxConc,
		Timeout:       *timeout,
		OutputDir:     *outputDir,
		Verbose:       *verbose,
	}

	dl := downloader.NewDownloader(cfg)

	fmt.Printf("Loading: %s\n", urlStr)
	if *verbose {
		fmt.Printf("Output directory: %s\n", *outputDir)
		fmt.Printf("Depth: %d, Concurrent: %d\n", *maxDepth, *maxConc)
	}

	if err := dl.Download(urlStr); err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Done! Downloaded %d files\n", dl.GetDownloadedCount())
}
