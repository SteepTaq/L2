package main

import (
	"fmt"
	"time"
)

func main() {
	sig := func(after time.Duration) <-chan any {
		c := make(chan any)
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Minute),
		sig(1*time.Hour),
		sig(1*time.Second),
	)

	fmt.Printf("done after %v\n", time.Since(start))
}

func or(channels ...<-chan any) <-chan any {
	switch len(channels) {
	case 0:
		// Если в качестве аргумента не переданы каналы, то возвращается закрытый канал
		toClose := make(chan any)
		close(toClose)
		return toClose
	case 1:
		// Если передан один канал, то возвращается этот же самый канал
		return channels[0]
	default:
		return loop(channels)
	}
}

func loop(channels []<-chan any) <-chan any {
	// Крайний случай рекурсии, возвращаем единственный канал для ожидания
	if len(channels) == 1 {
		return channels[0]
	}

	// Создаем канал который будет закрыт при завершении работы горутины
	orDone := make(chan any)
	go func() {
		defer close(orDone)

		if len(channels) == 2 {
			// Если в слайсе осталось 2 канала, то ожидаем закрытия одного из них
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		} else {
			// Иначе делим слайс на 2 и рекурсивно вызываем функцию для каждой половины
			mid := len(channels) / 2
			select {
			case <-loop(channels[:mid]):
			case <-loop(channels[mid:]):
			}
		}
	}()

	return orDone
}