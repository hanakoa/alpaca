package services

import (
	"strconv"
	"net/http"
	"github.com/kevinmichaelchen/my-go-utils"
	"gopkg.in/guregu/null.v3"
)

const (
	defaultCursor = -1
	//200
	defaultPageSize = 5
	maxPageSize = 1000
)

type Page struct {
	Data          []interface{} `json:"data"`
	NextCursor    null.Int      `json:"next_cursor"`
	NextCursorStr string        `json:"next_cursor_str"`
	PrevCursor    int64         `json:"previous_cursor"`
	PrevCursorStr string        `json:"previous_cursor_str"`
}

func emptyPage() Page {
	data := make([]interface{}, 0)
	return Page{
		Data:          data,
		PrevCursor:    defaultCursor,
		PrevCursorStr: strconv.FormatInt(defaultCursor, 10),
		NextCursor:    null.IntFrom(defaultCursor),
		NextCursorStr: strconv.FormatInt(defaultCursor, 10)}
}

func makePage(count int, data []interface{}, cursor int64, lastId int64) Page {
	var nextCursor null.Int
	var nextCursorStr string
	if len(data) < count {
		nextCursor = null.NewInt(0, false)
	} else {
		nextCursor = null.IntFrom(lastId)
		nextCursorStr = strconv.FormatInt(lastId, 10)
	}

	return Page{
		Data:          data,
		PrevCursor:    cursor,
		PrevCursorStr: strconv.FormatInt(cursor, 10),
		NextCursor:    nextCursor,
		NextCursorStr: nextCursorStr}
}

func getCursor(r *http.Request) int64 {
	cursorString := r.FormValue("cursor")
	if cursorString == "" {
		return defaultCursor
	}
	if !utils.IsParseableAsInt64(cursorString) {
		return defaultCursor
	}
	cursor := utils.StringToInt64(cursorString)
	return cursor
}

func getSort(r *http.Request) string {
	// TODO actually parse it
	return "ASC"
}

func getCount(r *http.Request) int {
	countString := r.FormValue("count")
	if countString == "" {
		return defaultPageSize
	}
	count, _ := strconv.Atoi(countString)
	if count > maxPageSize {
		count = maxPageSize
	}
	if count < 1 {
		count = 1
	}
	return count
}