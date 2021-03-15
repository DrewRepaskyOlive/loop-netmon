package main

import "github.com/DrewRepaskyOlive/loop-netmon/loop"

func main() {
	if err := loop.Serve(); err != nil {
		panic(err)
	}
}
