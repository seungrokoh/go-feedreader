package main

import (
	"fmt"
	_ "github.com/seungrokoh/go-feedreader/matchers"
	"github.com/seungrokoh/go-feedreader/search"
	"log"
	"os"
	"runtime"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	search.Run("코로나")
	fmt.Println("num of goroutine: ", runtime.NumGoroutine())
}
