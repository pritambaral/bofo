package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

const (
	MEM_SIZE       = 30000
	MAX_LOOP_DEPTH = 1024
)

type Script struct {
	program []byte
	max     int
	ptr     int
}

func initScript(filepath string) (script *Script, err error) {
	script = &Script{}
	script.program, err = ioutil.ReadFile(filepath)
	if err != nil {
		return
	}
	script.ptr = 0
	script.max = len(script.program) - 2 //One EOF and one due to zero-indexed ptr
	return
}

func (script *Script) nextOp() bool {
	if script.ptr >= script.max {
		return false
	}
	script.ptr++
	return true
}

type Memory struct {
	bank []byte
	ptr  int16
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: %s /path/to/script.bf", os.Args[0])
		os.Exit(2)
	}

	script, err := initScript(os.Args[1])
	if err != nil {
		fmt.Println("Couldn't read file", err)
		os.Exit(2)
	}
	script.ptr = -1

	data := &Memory{make([]byte, MEM_SIZE), 0}
	loop_jumps := make([]int, MAX_LOOP_DEPTH)
	loop_depth := -1

	for script.nextOp() {
		switch script.program[script.ptr] {
		case '>':
			data.ptr++
		case '<':
			data.ptr--
		case '+':
			data.bank[data.ptr]++
		case '-':
			data.bank[data.ptr]--
		case '.':
			fmt.Print(string(data.bank[data.ptr]))
		case ',':
			buf := make([]byte, 1)
			os.Stdin.Read(buf)
			data.bank[data.ptr] = buf[0]
		case '[':
			//To jump or not to jump
			if data.bank[data.ptr] == 0 {
				for depth := 1; depth != 0 && script.nextOp(); {
					if script.program[script.ptr] == '[' {
						depth++
					} else if script.program[script.ptr] == ']' {
						depth--
					}
				}
			} else {
				//Remember where we enter loop, so we can quickly reiterate
				loop_depth += 1
				loop_jumps[loop_depth] = script.ptr
			}
		case ']':
			if data.bank[data.ptr] == 0 {
				loop_jumps[loop_depth] = -1
				loop_depth--
			} else {
				script.ptr = loop_jumps[loop_depth]
			}
		}
	}
}
