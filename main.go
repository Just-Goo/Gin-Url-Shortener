package main

import "sync"

type mapper struct {
	mapping map[string]string
	sync.Mutex
}

var urlMapper Mapper

func main() {

}