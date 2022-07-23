all:
	cat .base.md > README.md
	go run main.go | csv2md /dev/stdin >> README.md