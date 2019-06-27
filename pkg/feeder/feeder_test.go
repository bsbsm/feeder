package feeder

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

func TestReadFeed(t *testing.T) {
	feedItem := &gofeed.Item{
		Title: "title 1",
		GUID:  "g-u-id-1",
	}

	var successRules = map[string]string{
		"Title": "new_title",
		"GUID":  "item_id",
	}

	var unreachRules = map[string]string{
		"Title1": "new_title",
		"GUID2":  "item_id",
	}

	payload := []byte("{\"new_title\":\"title 1\",\"item_id\":\"g-u-id-1\"}")
	var wantFields map[string]*json.RawMessage
	if err := json.Unmarshal(payload, &wantFields); err != nil {
		t.Error(err)
	}

	payload = []byte("{}")
	var wantEmptyFields map[string]*json.RawMessage
	if err := json.Unmarshal(payload, &wantEmptyFields); err != nil {
		t.Error(err)
	}

	type args struct {
		item  *gofeed.Item
		rules map[string]string
	}

	tests := []struct {
		name string
		args func() args
		// inspectResult func(res1 string, t *testing.T)
		want    map[string]*json.RawMessage
		wantErr error
	}{
		{
			name: "success parsing",
			args: func() args {
				return args{
					item:  feedItem,
					rules: successRules,
				}
			},
			want: wantFields,
		},
		{
			name: "parsing when fields not found",
			args: func() args {
				return args{
					item:  feedItem,
					rules: unreachRules,
				}
			},
			want: make(map[string]*json.RawMessage),
		},
		{
			name: "parsing with nil rules",
			args: func() args {
				return args{
					item:  feedItem,
					rules: nil,
				}
			},
			want: wantEmptyFields,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer time.Sleep(time.Millisecond)

			a := tt.args()
			res, err := parseFeedItem(a.item, a.rules)

			if tt.wantErr != nil && assert.Error(t, err) {
				assert.Equal(t, tt.wantErr, err, "feeder.parseFeedRecord returned unexpected error")
			} else {
				assert.NoError(t, err)

				var resFields map[string]*json.RawMessage
				if err := json.Unmarshal(res, &resFields); err != nil {
					t.Error(err)
				}

				assert.Equal(t, tt.want, resFields, "feeder.parseFeedRecord returned unexpected payload bytes")
			}
		})
	}
}
