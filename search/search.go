package search

import (
	"context"
	"log"
	"sync"
	"time"
)

var matchers = make(map[string]Matcher)

type ResultPipe func(<-chan *Result) <-chan *Result

func Run(searchTerm string) {
	feeds, err := RetrieveFeeds()
	if err != nil {
		log.Fatal(err)
	}

	cs := make([]<-chan *Result, len(feeds))
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()

	// 여기서 여러개로 나누기
	for i, feed := range feeds {
		matcher, exist := matchers[feed.Type]
		if !exist {
			matcher = matchers["default"]
		}
		cs[i] = Match(ctx, matcher, feed, searchTerm)
	}

	Display(FanIn(cs...))

}

func FanIn(ins ...<-chan *Result) <-chan *Result {
	out := make(chan *Result)
	var wg sync.WaitGroup
	wg.Add(len(ins))
	for _, in := range ins {
		go func(in <-chan *Result) {
			defer wg.Done()
			for result := range in {
				out <- result
			}
		}(in)
	}

	go func() {
		defer close(out)
		wg.Wait()
	}()

	return out
}

func Register(feedType string, matcher Matcher) {
	if _, exists := matchers[feedType]; exists {
		log.Fatalln(feedType, " 검색기가 이미 등록되었습니")
	}

	log.Println("등록 완료: ", feedType, " 검색기")
	matchers[feedType] = matcher
}
