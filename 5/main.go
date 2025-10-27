package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	return nil
}

func main() {
	// тип error - интерфейс, по умолчанию равен nil
	// соответственно, `tab` равен nil, `data` тоже nil
	var err error
	// `tab` указывает на `itab` с информацией о типе `*customError`
	err = test()

	// err != nil вернет true, так как `tab` не равен nil
	if err != nil {
		println("error") // "error"
		return
	}

	// Соответственно, дальше условия программа не дойдет
	println("ok")
}
