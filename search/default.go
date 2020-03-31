package search

import "context"

// 기본 검색기를 구현할 defaultMatcher 타입
type defaultMatcher struct{}

// 기본 검색기를 프로그램에 등록
func init() {
	var matcher defaultMatcher
	Register("default", matcher)
}

// 기본 검색기의 동작을 구현
func (m defaultMatcher) Search(ctx context.Context, feed *Feed, searchTerm string) <-chan *Response {
	return make(chan *Response)
}
