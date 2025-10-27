test: case1 case2 case3 case4 case5 case6 case7 case8 case9 case10

case1:
	go run . -B 2 "eternal summer" data/poetry.txt && echo ""

case2:
	cat data/poetry.txt | go run . -C 1 "So long" && echo "\n"

case3:
	go run . -v -n "fair" data/poetry.txt && echo "\n"

case4:
	go run . -i "warning" data/logs.txt && echo "\n"

case5:
	cat data/logs.txt | go run . "User '[^']+'" && echo "\n"

case6:
	go run . "service." data/logs.txt && echo "\n"

case7:
	go run . -B 2 "START" data/edge.txt && echo "\n"

case8:
	cat data/edge.txt | go run . -A 2 "END" && echo "\n"

case9:
	go run . -c "anything" data/empty.txt && echo "\n"

case10:
	go run . -B 1 "line" data/edge.txt