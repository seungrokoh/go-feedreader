package search

import (
	"log"
	"sync"
)

var matchers = make(map[string]Matcher)

func Run(searchTerm string) {
	feeds, err := RetrieveFeeds()
	if err != nil {
		log.Fatal(err)
	}

	results := make(chan *Result)

	var waitGroup sync.WaitGroup

	waitGroup.Add(len(feeds))
	for _, feed := range feeds {
		matcher, exists := matchers[feed.Type]
		if !exists {
			matcher = matchers["default"]
		}

		go func(matcher Matcher, feed *Feed) {
			defer waitGroup.Done()
			Match(matcher, feed, searchTerm, results)
		}(matcher, feed)
	}

	go func() {
		waitGroup.Wait()
		close(results)
	}()

	Display(results)
}

func Register(feedType string, matcher Matcher) {
	if _, exists := matchers[feedType]; exists {
		log.Fatalln(feedType, " 검색기가 이미 등록되었습니")
	}

	log.Println("등록 완료: ", feedType, " 검색기")
	matchers[feedType] = matcher
}
