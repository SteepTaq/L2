package main

import (
	"fmt"
)

func main() {
	var s = []string{"1", "2", "3"} // [1,2,3], len = 3, cap = 3
	modifySlice(s)                  // передаем слайс по значению
	fmt.Println(s)                  // [3,2,3]
}

func modifySlice(i []string) {
	// структура слайса хранит ссылку на массив
	i[0] = "3" // [3,2,3], len = 3, cap = 3

	// добавление 4 превысит capacity, поэтому будет создан новый массив
	// и скопированны значения, с увеличенным вдвое capacity
	i = append(i, "4") // [3,2,3,4],   len = 4, cap = 6
	i[1] = "5"         // [3,5,3,4],   len = 4, cap = 6
	i = append(i, "6") // [3,5,3,4,5], len = 5, cap = 6
}
