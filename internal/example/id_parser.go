package example

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func IDParser(route string) func(r *http.Request) (uint, error) {
	index := -1
	for i, p := range strings.Split(route, "/") {
		if p == ":id" {
			index = i
			break
		}
	}
	if index == -1 {
		panic("invalid route")
	}
	return func(r *http.Request) (uint, error) {
		splited := strings.Split(r.URL.Path, "/")
		if len(splited) <= index {
			return 0, errors.New("invalid route")
		}
		id, err := strconv.ParseUint(splited[index], 10, 64)
		if err != nil {
			return 0, err
		}
		return uint(id), nil
	}
}
