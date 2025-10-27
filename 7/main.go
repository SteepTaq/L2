package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Возвращает канал, в который записываются числа из слайса vs
// и закрывается после записи всех чисел
func asChan(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v, ok := <-a:
				if ok {
					c <- v
				} else {
					// Чтобы select не попадал в бесконечный цикл,
					// когда канал a закрыт, нужно присвоить ему nil
					a = nil
				}
			case v, ok := <-b:
				if ok {
					c <- v
				} else {
					// Чтобы select не попадал в бесконечный цикл,
					// когда канал b закрыт, нужно присвоить ему nil
					b = nil
				}
			}
			// Если оба канала закрыты, то закрываем канал c и выходим из горутины
			if a == nil && b == nil {
				close(c)
				return
			}
		}
	}()

	return c
}

func main() {
	rand.Seed(time.Now().Unix())

	a := asChan(1, 3, 5, 7) // <-a, записываются числа 1, 3, 5, 7
	b := asChan(2, 4, 6, 8) // <-b, записываются числа 2, 4, 6, 8

	c := merge(a, b) // <-c, записываются числа с каналов a и b
	for v := range c {
		fmt.Print(v) // от 1 до 8 в случайном порядке
	}

	// Логика работы программы:
	// 1. Создаем 2 канала a и b, в которые записываются числа с задержкой
	// 2. Используем конвейер merge, который сливает числа из каналов a и b
	// 3. Читаем из канала c и выводим в консоль
}
