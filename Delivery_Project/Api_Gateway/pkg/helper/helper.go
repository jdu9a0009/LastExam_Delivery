package helper

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func ReplaceQueryParams(namedQuery string, params map[string]interface{}) (string, []interface{}) {
	var (
		i    int = 1
		args []interface{}
	)

	for k, v := range params {
		if k != "" {
			namedQuery = strings.ReplaceAll(namedQuery, ":"+k, "$"+strconv.Itoa(i))

			args = append(args, v)
			i++
		}
	}

	return namedQuery, args
}

func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

const otpChars = "1234567890"

func GenerateOTP(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}

func Difference(a, b []int32) []int32 {
	mb := make(map[int32]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []int32
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func ValMultipleQuery(query string, vals []int32) (string, []interface{}) {
	params := []interface{}{}

	for i, id := range vals {
		query += fmt.Sprintf("$%d,", i+1)
		params = append(params, id)
	}

	query = query[:len(query)-1] // remove trailing ","
	query += ")"

	return query, params
}

func InsertMultiple(queryInsert string, id string, vals []string) (string, []interface{}) {
	insertparams := []interface{}{}

	for i, d := range vals {
		p1 := i * 2 // starting position for insert params
		queryInsert += fmt.Sprintf("($%d, $%d),", p1+1, p1+2)
		insertparams = append(insertparams, id, d)
	}

	queryInsert = queryInsert[:len(queryInsert)-1] // remove trailing ","

	return queryInsert, insertparams
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NewNullBool(s bool) sql.NullBool {
	if !s {
		return sql.NullBool{}
	}
	return sql.NullBool{
		Bool:  s,
		Valid: true,
	}
}

func DoRequest(url string, method string, body interface{}) ([]byte, error) {
	data, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respByte, nil
}
func GenerateID(lastID string, word string) string {

	var (
		number int
		count  int
		ID     string
	)

	if len(lastID) < 1 {
		ID = word + "-0000001"
		return ID
	}

	for _, val := range lastID {
		c, ok := strconv.Atoi(string(val))
		if ok != nil || c == 0 {
			continue
		}
		number = number*10 + c
		count += 1
	}

	ID += string(lastID[0]) + "-"
	for i := 0; i < len(lastID)-2-count; i++ {
		ID += "0"
	}
	ID += strconv.Itoa(number + 1)

	return ID
}
func GenerateString(letter string, count int) string {
	letter = strings.ToUpper(letter)

	if count <= 0 {
		count = 1
	} else if count > 9999999 {
		count = 9999999
	}

	countStr := fmt.Sprintf("%07d", count)

	result := fmt.Sprintf("%s-%s", letter, countStr)
	return result
}
