package db

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func prepareDatabase() *sql.DB {
	databaseFilePath = "./test.db"
	os.Remove(databaseFilePath)
	databaseInstance = nil
	return getDb()
}

func prepareDbForRead(t *testing.T) {
	db := prepareDatabase()

	query := `
	INSERT INTO news(
		Title,
		SourceID,
		PayloadJSON
	) values('NewTitle1', 1, "json");
	INSERT INTO news(
		Title,
		SourceID,
		PayloadJSON
	) values('NewTitle2', 2, "json");
	INSERT INTO news(
		Title,
		SourceID,
		PayloadJSON
	) values('NewTitle3', 3, "json");
	INSERT INTO sources(
		URL,
		Rule
	) values('uselessurl1', 'title:newtitle');
	INSERT INTO sources(
		URL,
		Rule
	) values('uselessurl2', 'title:newtitle');
	INSERT INTO sources(
		URL,
		Rule
	) values('uselessurl3', 'title:newtitle');
			`
	_, err := db.Exec(query)
	if err != nil {
		t.Error(err)
	}
}

func TestGetNews(t *testing.T) {
	sqlite := SQLiteDatabase{}

	type args struct {
		offset, count int
	}

	tests := []struct {
		name       string
		in         args
		prepare    func(t *testing.T)
		inspect    func(t *testing.T) //inspects database after execution of GetNews
		wantLen    int
		wantErr    bool
		inspectErr func(err error, t *testing.T)
	}{
		{
			name:    "success getting",
			in:      args{0, 20},
			prepare: prepareDbForRead,
			wantLen: 3,
			wantErr: false,
			inspect: func(t *testing.T) {
				db := getDb()
				query := `SELECT count(*) FROM news`

				row := db.QueryRow(query)
				if row == nil {
					t.Error("database return nil")
				}

				var c int
				err := row.Scan(&c)
				if err != nil || c != 3 {
					t.Errorf("Expect 3 but got %v. Err: %s", c, err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer time.Sleep(time.Millisecond)

			tt.prepare(t)
			res, err := sqlite.GetNews(tt.in.offset, tt.in.count)

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantLen, len(res), "SQLiteDatabase.GetNews returned unexpected array's len")

			if tt.inspect != nil {
				tt.inspect(t)
			}
		})
	}
}

func TestGetNewsDetail(t *testing.T) {
	sqlite := SQLiteDatabase{}

	tests := []struct {
		name       string
		in         int
		prepare    func(t *testing.T)
		inspect    func(t *testing.T) //inspects database after execution of GetNewsDetail
		want       *NewsDetail
		wantErr    bool
		inspectErr func(err error, t *testing.T)
	}{
		{
			name:    "finded",
			in:      1,
			prepare: prepareDbForRead,
			want:    &NewsDetail{"NewTitle1", "json", "uselessurl1"},
			wantErr: false,
			inspect: func(t *testing.T) {
				db := getDb()
				query := `SELECT count(*) FROM news`

				row := db.QueryRow(query)
				if row == nil {
					t.Error("database return nil")
				}

				var c int
				err := row.Scan(&c)
				if err != nil || c != 3 {
					t.Errorf("Expect 3 but return %v. Err: %s", c, err)
				}
			},
		},
		{
			name:    "not found",
			in:      5,
			prepare: prepareDbForRead,
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, err, ErrNotFound, "SQLiteDatabase.GetNewsDetail returned unexpected error")
			},
			inspect: func(t *testing.T) {
				db := getDb()
				query := `SELECT count(*) FROM news`

				row := db.QueryRow(query)
				if row == nil {
					t.Error("database return nil")
				}

				var c int
				err := row.Scan(&c)
				if err != nil || c != 3 {
					t.Errorf("Expect 3 but return %v OR err: %s", c, err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer time.Sleep(time.Millisecond)

			tt.prepare(t)
			res, err := sqlite.GetNewsDetail(tt.in)

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err, "SQLiteDatabase.GetNewsDetail returned error but expected not error")
				assert.NotNil(t, res, "SQLiteDatabase.GetNewsDetail returned nil but expected not nil")
			}

			assert.Equal(t, res, tt.want, "SQLiteDatabase.GetNewsDetail. Wanted %v but got %v", tt.want, res)

			if tt.inspect != nil {
				tt.inspect(t)
			}
		})
	}
}

func TestGetNewsWithTitle(t *testing.T) {
	sqlite := SQLiteDatabase{}

	type args struct {
		offset, count int
		title         string
	}

	tests := []struct {
		name       string
		in         args
		prepare    func(t *testing.T)
		inspect    func(t *testing.T) //inspects database after execution of GetNewsWithTitle
		wantLen    int
		wantErr    bool
		inspectErr func(err error, t *testing.T)
	}{
		{
			name:    "success getting",
			in:      args{0, 20, "New"},
			prepare: prepareDbForRead,
			wantLen: 3,
			wantErr: false,
			inspect: func(t *testing.T) {
				db := getDb()
				query := `SELECT count(*) FROM news`

				row := db.QueryRow(query)
				if row == nil {
					t.Error("database return nil")
				}

				var c int
				err := row.Scan(&c)
				if err != nil || c != 3 {
					t.Errorf("Expect 3 but return %v. Err: %s", c, err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer time.Sleep(time.Millisecond)

			tt.prepare(t)
			res, err := sqlite.GetNewsWithTitle(tt.in.title, tt.in.offset, tt.in.count)

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}

			assert.Equal(t, tt.wantLen, len(res), "SQLiteDatabase.GetNewsWithTitle returned unexpected array's len")

			if tt.inspect != nil {
				tt.inspect(t)
			}
		})
	}
}

func TestGetFeedSource(t *testing.T) {
	sqlite := SQLiteDatabase{}

	tests := []struct {
		name       string
		prepare    func(t *testing.T)
		wantLen    int
		inspect    func(t *testing.T) //inspects database after execution of GetFeedSource
		wantErr    bool
		inspectErr func(err error, t *testing.T)
	}{
		{
			name:    "success getting",
			prepare: prepareDbForRead,
			wantLen: 3,
			wantErr: false,
			inspect: func(t *testing.T) {
				db := getDb()
				query := `SELECT count(*) FROM sources`

				row := db.QueryRow(query)
				if row == nil {
					t.Error("database return nil")
				}

				var c int
				err := row.Scan(&c)
				if err != nil || c != 3 {
					t.Errorf("Expect 3 but return %v. Err: %s", c, err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer time.Sleep(time.Millisecond)

			tt.prepare(t)
			res, err := sqlite.GetFeedSources()

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}

			assert.Equal(t, tt.wantLen, len(res), "SQLiteDatabase.GetFeedSources returned unexpected array's len")

			if tt.inspect != nil {
				tt.inspect(t)
			}
		})
	}
}

func TestCreateNews(t *testing.T) {
	sqlite := SQLiteDatabase{}

	type args struct {
		sourceID    int
		title       string
		payloadJSON []byte
	}

	tests := []struct {
		name       string
		in         args
		inspect    func(t *testing.T) //inspects database after execution of CreateNews
		wantErr    bool
		inspectErr func(err error, t *testing.T)
	}{
		{
			name:    "success adding",
			in:      args{1, "News title", []byte("{\"id\":11,\"title\":\"News title\",\"body\":\"News body\"}")},
			wantErr: false,
			inspect: func(t *testing.T) {
				db := getDb()
				query := `SELECT count(*) FROM news`

				row := db.QueryRow(query)
				if row == nil {
					t.Error("database return nil")
				}

				var c int
				err := row.Scan(&c)
				if err != nil || c < 1 {
					t.Errorf("Expect 1 but return 0. Err: %s", err)
				}
			},
		},
		{
			name:    "failed adding",
			in:      args{1, "", nil},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, err, ErrIncorrectArgs, "SQLiteDatabase.CreateNews returned unexpected error")
			},
			inspect: func(t *testing.T) {
				db := getDb()
				query := `SELECT count(*) FROM news`

				row := db.QueryRow(query)
				if row == nil {
					t.Error("database return nil")
				}

				var c int
				err := row.Scan(&c)
				if err != nil || c > 0 {
					t.Errorf("Expect 0 but return %v OR err: %s", c, err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer time.Sleep(time.Millisecond)

			prepareDatabase()

			err := sqlite.CreateNews(tt.in.sourceID, tt.in.title, tt.in.payloadJSON)

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

			if tt.inspect != nil {
				tt.inspect(t)
			}
		})
	}
}

func TestCreateFeedSource(t *testing.T) {
	sqlite := SQLiteDatabase{}

	type args struct {
		url  string
		rule string
	}

	tests := []struct {
		name       string
		in         args
		inspect    func(t *testing.T) //inspects database after execution of CreateFeedSource
		wantErr    bool
		inspectErr func(err error, t *testing.T)
	}{
		{
			name:    "success adding",
			in:      args{url: "https://www.netroby.com/rss", rule: "Title=NewTitle"},
			wantErr: false,
			inspect: func(t *testing.T) {
				db := getDb()
				query := `SELECT count(*) FROM sources`

				row := db.QueryRow(query)
				if row == nil {
					t.Error("database return nil")
				}

				var c int
				err := row.Scan(&c)
				if err != nil || c < 1 {
					t.Errorf("Expect 1 but return 0. Err: %s", err)
				}
			},
		},
		{
			name:    "failed adding",
			in:      args{url: "https", rule: ""},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, err, ErrIncorrectArgs, "SQLiteDatabase.CreateFeedSource returned unexpected error")
			},
			inspect: func(t *testing.T) {
				db := getDb()
				query := `SELECT count(*) FROM sources`

				row := db.QueryRow(query)
				if row == nil {
					t.Error("database return nil")
				}

				var c int
				err := row.Scan(&c)
				if err != nil || c > 0 {
					t.Errorf("Expect 0 but return %v OR err: %s", c, err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer time.Sleep(time.Millisecond)

			prepareDatabase()

			err := sqlite.CreateFeedSource(tt.in.url, tt.in.rule)

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

			if tt.inspect != nil {
				tt.inspect(t)
			}
		})
	}
}
