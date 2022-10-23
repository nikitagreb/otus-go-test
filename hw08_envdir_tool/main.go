package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	flag.Parse()

	envDirPath := flag.Arg(0)

	if envDirPath == "" {
		log.Fatalf("env dir not passed")
	}

	env, err := ReadDir(envDirPath)
	if err != nil {
		log.Fatalf(err.Error())
	}

	cmd := flag.Args()[1:]

	code := RunCmd(cmd, env)

	os.Exit(code)
}
