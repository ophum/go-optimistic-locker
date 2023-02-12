package example

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func MakePetsKey(id uint) string {
	return fmt.Sprintf("pets-%d", id)
}

func KeyParser(route string) func(r *http.Request) (string, error) {
	index := -1
	for i, v := range strings.Split(route, "/") {
		if v == ":id" {
			index = i
		}
	}
	if index == -1 {
		panic("keyParser: invalid route")
	}

	return func(r *http.Request) (string, error) {
		splited := strings.Split(r.URL.Path, "/")
		if len(splited) <= index {
			return "", errors.New("keyParser: invalid route")
		}

		id, err := strconv.ParseUint(splited[index], 10, 64)
		if err != nil {
			return "", errors.Wrap(err, "keyParser")
		}
		return MakePetsKey(uint(id)), nil
	}
}
