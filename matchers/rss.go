package matchers

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/seungrokoh/go-feedreader/search"
	"log"
	"net/http"
	"regexp"
	"time"
)

type (
	// item 구조체는 RSS 문서 내의 item 태그에
	// 정의된 필드들에 대응하는 필드들을 선언한다.
	item struct {
		XMLName     xml.Name `xml:"item"`
		PubDate     string   `xml:"pubDate"`
		Title       string   `xml:"title"`
		Description string   `xml:"description"`
		Link        string   `xml:"link"`
		GUID        string   `xml:"guid"`
		GeoRssPoint string   `xml:"georss:point"`
	}

	// image 구조체는 RSS 문서 내의 image 태그에
	// 정의된 필드들에 대응하는 필드들을 선언한다.
	image struct {
		XMLName xml.Name `xml:"image"`
		URL     string   `xml:"url"`
		Title   string   `xml:"title"`
		Link    string   `xml:"link"`
	}

	// channel 구조체는 RSS 문서 내의 channel 태그에
	// 정의된 필드들에 대응하는 필드들을 선언한다.
	channel struct {
		XMLName        xml.Name `xml:"channel"`
		Title          string   `xml:"title"`
		Description    string   `xml:"description"`
		Link           string   `xml:"link"`
		PubDate        string   `xml:"pubDate"`
		LastBuildDate  string   `xml:"lastBuildDate"`
		TTL            string   `xml:"ttl"`
		Language       string   `xml:"language"`
		ManagingEditor string   `xml:"managingEditor"`
		WebMaster      string   `xml:"webMaster"`
		Image          image    `xml:"image"`
		Item           []item   `xml:"item"`
	}

	// rssDocument 구조체는 RSS 문서에 정의된 필드들에 대응하는 필드들을 정의한다.
	rssDocument struct {
		XMLName xml.Name `xml:"rss"`
		Channel channel  `xml:"channel"`
	}
)

type rssMatcher struct{}

func init() {
	var matcher rssMatcher
	search.Register("rss", matcher)
}

func (m rssMatcher) Search(ctx context.Context, feed *search.Feed, searchTerm string) <-chan *search.Response {
	out := make(chan *search.Response)

	log.Printf("피드 종류[%s] 사이트[%s] 주소[%s]에서 검색을 수행합니다.\n", feed.Type, feed.Name, feed.URI)

	// 고루틴을 생성함으로써 장점은 없지만 연습용으로 구현
	go func() {
		document, err := m.retrieve(feed)
		if err != nil {
			out <- createResponse(nil, err)
			return
		}

		for _, channelItem := range document.Channel.Item {
			// context.Done이 들어오면 고루틴을 종료한다.
			if IsContextDone(ctx) {
				log.Println("Canceled Search Goroutine!!!")
				return
			}

			if response := containKeywordInItem(searchTerm, "Title", channelItem.Title); response != nil {
				out <- response
			}

			if response := containKeywordInItem(searchTerm, "Description", channelItem.Description); response != nil {
				out <- response
			}
		}
	}()
	return out
}

func containKeywordInItem(keyword, field, channelItem string) *search.Response {
	matched, err := regexp.MatchString(keyword, channelItem)
	// matchString 에러 발생시 err response 생성
	if err != nil {
		return createResponse(nil, err)
	}
	// keyword가 Item안에 존재 할 경우
	if matched {
		return createResponse(createResult(field, channelItem), nil)
	} else {
		return nil
	}
}

func createResult(field, content string) *search.Result {
	return &search.Result{
		Field:   field,
		Content: content,
	}
}

func createResponse(result *search.Result, err error) *search.Response {
	return &search.Response{
		Result: result,
		Error:  err,
	}
}

func IsContextDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func (m rssMatcher) retrieve(feed *search.Feed) (*rssDocument, error) {
	if feed.URI == "" {
		return nil, errors.New("검색할 RSS 피드가 정의되지 않았습니다")
	}
	client := http.Client{
		Timeout: 500 * time.Millisecond,
	}
	resp, err := client.Get(feed.URI)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP 응답 오류: %d\n", resp.StatusCode)
	}

	var document rssDocument
	err = xml.NewDecoder(resp.Body).Decode(&document)
	return &document, err
}
