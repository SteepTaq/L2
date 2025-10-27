package config

import "time"

type Config struct {
	MaxDepth      int           // Максимальная глубина рекурсии
	MaxConcurrent int           // Максимальное количество одновременных загрузок
	Timeout       time.Duration // Таймаут для HTTP запросов
	OutputDir     string        // Директория для сохранения файлов
	Verbose       bool          // Подробное логирование
}
