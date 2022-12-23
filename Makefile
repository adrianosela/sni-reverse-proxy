NAME:=$(shell basename `git rev-parse --show-toplevel`)
HASH:=$(shell git rev-parse --verify --short HEAD)
EXAMPLE:=simple

# run default program
run: build
	./$(NAME)

# build default program (go source in *this* directory)
build:
	go build -o $(NAME)

# list examples (dirs in _examples_ dir)
list-examples:
	@ls _examples_/ | cat

# run an example by directory name (as set in EXAMPLE env var)
run-example:
	go run _examples_/$(EXAMPLE)/main.go
