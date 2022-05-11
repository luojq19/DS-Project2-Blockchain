package main

import (
	"errors"
	"bufio"
	"os"
)

func input() (*CLI, *os.File){
	if len(os.Args) < 3{
		panic(errors.New("Wrong number of arguments."))
	}
	cli := NewCLI(os.Args[1])
	workload := os.Args[2]

	file, err := os.Open(workload)
	if err != nil {
		panic(err)
	}

	return cli, file
}

func main() {

	cli, file := input()
	defer file.Close()
	
	scanner := bufio.NewScanner(file)

	for scanner.Scan(){
		cli.Atomic(scanner.Text())
	}
}