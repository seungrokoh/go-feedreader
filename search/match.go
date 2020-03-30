package search

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Result struct {
	Field   string
	Content string
}

// 검색 타입의 필요한 동작을 정의
type Matcher interface {
	Search(feed *Feed, searchTerm string) ([]*Result, error)
}

// 고루틴으로써 호출되며 개별 피드 타입에 대한 검색을 동시에 수행
func Match(ctx context.Context, matcher Matcher, feed *Feed, searchTerm string) <-chan *Result {
	// 지정된 검색기를 이용해 검색을 수행
	out := make(chan *Result)
	// 검색 결과를 채널에 기록
	go func() {
		defer close(out)
		searchResults, err := matcher.Search(feed, searchTerm)
		if err != nil {
			// 시간초과로 err 발생할 시 고루틴 종료
			log.Println("Search error timeout!!!!! : ", err)
			return
		}

		for _, result := range searchResults {
			select {
			case <-ctx.Done():
				fmt.Println("Canceled goroutine!!!!")
				return
			case out <- result:
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()
	return out
}

// 개별 고루틴이 전달한 검색 결과를 콘솔에 출력
func Display(results <-chan *Result) {
	for result := range results {
		log.Printf("%s:\n%s\n\n", result.Field, result.Content)
	}
}
