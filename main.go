package main

import (
	"github.com/seungrokoh/go-feedreader/search"
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	search.Run("Sherlock")
}
