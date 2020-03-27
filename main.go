package main

import (
	_ "github.com/seungrokoh/go-feedreader/matchers"
	"github.com/seungrokoh/go-feedreader/search"
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	search.Run("코로나")
}
