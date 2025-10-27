package main

import (
	"fmt"
	"os"
)

func Foo() error {
	// err инициализируется как nil, но имеет тип *os.PathError
	var err *os.PathError = nil

	// err неявно имплементирует интерфейс error, поэтому может возвращаться в качестве error
	return err
}

func main() {
	// 1. Внутренняя стурктура интерфейсов:
	// type iface struct {
	//   tab  *itab          // 8 байт: таблица методов + информация о типе
	//   data unsafe.Pointer // 8 байт: указатель на данные
	// }
	//
	// type itab struct {
	//   inter *interfacetype // описание интерфейса
	//   _type *_type         // конкретный тип
	//   hash  uint32         // хеш типа для быстрого сравнения
	//   fun   [1]uintptr     // таблица методов
	// }

	// 2. Внутренняя структура пустого интерфейса:
	// type eface struct {
	//   _type *_type         // 8 байт: информация о типе
	//   data  unsafe.Pointer // 8 байт: указатель на данные
	// }

	// Ключевые отличия пустых интерфейсов от обычных интерфейсов, это то, что
	// `eface` содержит только информацию о типе, принимает любой тип, нет
	// предварительной вычисленной таблицы методов.

	// Foo() вернет *os.PathError как тип error. Создается структура `iface`,
	// поле `tab` указывает на `itab` с информацией о типе `*os.PathError`, `data`
	// равен nil.
	err := Foo()

	fmt.Println(err) // <nil>
	// интерфейс считается nil, если и `tab` и `data` равны nil
	// В нашем случае `tab` содержит информацию о типе `*os.PathError`,
	fmt.Println(err == nil)
}