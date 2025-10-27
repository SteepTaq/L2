test: case1 case2 case3 case4 case5 case6

case1:
	go run . -f 1,3 data/tabs.txt && echo

case2:
	cat data/tabs.txt | go run . -f 2-4 && echo

case3:
	go run . -f 2-3 -d : data/colon.txt && echo

case4:
	cat data/colon.txt | go run . -f 2-3 -d : -s && echo

case5:
	go run . -f 10-12 data/tabs.txt && echo

case6:
	go run . -f 1 -d : data/colon.txt && echo