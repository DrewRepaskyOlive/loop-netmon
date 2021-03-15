package main

import "bitbucket.org/crosschx/loop-netmon/loop"

func main() {
	if err := loop.Serve(); err != nil {
		panic(err)
	}
}
