package feeder

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

var ErrEmptyRule = errors.New("Parsing rule is empty")

func NewFeeder(s FeedStorage) (*Feeder, error) {
	if s == nil || reflect.ValueOf(s).IsNil() {
		return nil, errors.New("FeedStorage is nil")
	}

	return &Feeder{storage: s}, nil
}

type FeedStorage interface {
	CreateNews(sourceID int, title string, payloadJSON []byte) error
	GetFeedSources() ([]*FeedSource, error)
}

type FeedSource struct {
	Rule map[string]string
	URL  string
	ID   int
}

func ImplementRule(s *FeedSource, rule string) error {
	if rule == "" {
		return ErrEmptyRule
	}

	pairs := strings.Split(rule, ",")
	if len(pairs) == 0 {
		return ErrEmptyRule
	}

	s.Rule = make(map[string]string)

	for _, p := range pairs {
		nameAndValue := strings.SplitN(p, "=", 2)

		switch len(nameAndValue) {
		case 2:
			s.Rule[nameAndValue[0]] = nameAndValue[1]
		case 1:
			s.Rule[nameAndValue[0]] = nameAndValue[0]
		}
	}

	return nil
}

type Feeder struct {
	storage FeedStorage
	sources []*FeedSource
}

func (f *Feeder) Reading(period time.Duration) {
	for {
		for _, s := range f.sources {
			f.readFeed(s)
		}

		newSources, err := f.storage.GetFeedSources()
		if err != nil {
			fmt.Println("Feed reading: error while get feed sources. Use cache")
		}

		if len(newSources) > 0 {
			f.sources = newSources
		}

		time.Sleep(period)
	}
}

var fp = gofeed.NewParser()

func (f *Feeder) readFeed(s *FeedSource) {
	if len(s.Rule) == 0 {
		return
	}

	feed, _ := fp.ParseURL(s.URL)

	for _, item := range feed.Items {
		payload, err := json.Marshal(item)

		if err != nil {
			fmt.Printf("Error while feed reading: %s\n", err)
		}

		var fields map[string]*json.RawMessage
		if err := json.Unmarshal(payload, &fields); err != nil {
			fmt.Printf("Error while feed reading: %s\n", err)
		}

		mapToSave := make(map[string]interface{})

		for k, newK := range s.Rule {
			if val, exist := fields[strings.ToLower(k)]; exist && val != nil {
				mapToSave[newK] = val
			}
		}

		payloadToSave, err := json.Marshal(mapToSave)
		if err != nil {
			fmt.Printf("Error while feed reading: %s\n", err)
		}

		if err := f.storage.CreateNews(s.ID, item.Title, payloadToSave); err != nil &&
			!strings.HasPrefix(err.Error(), "UNIQUE") {
			fmt.Printf("Error while feed reading: %s\n", err)
		}
	}
}
