package main

import (
	"flag"
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
	defer fmt.Println("num of goroutine : ", runtime.NumGoroutine())
	keyword := flag.String("keyword", "", "search keyword")
	flag.Parse()
	search.Run(*keyword)
}

// 1. go run main.go -keyword=corona flag
// 2. http.Get에 타임아웃 걸기 100 * time.Milliseconds
// 3. cancel goroutine
