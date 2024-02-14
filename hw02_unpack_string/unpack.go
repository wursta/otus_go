package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	source := []rune(str)

	b := strings.Builder{}

	for i, v := range source {
		repeatCount, err := strconv.Atoi(string(v))
		if err != nil {
			b.WriteRune(v)
			continue
		}

		// Если мы тут, то это число, но если это первый элемент, то отдаём ошибку
		if i == 0 {
			return "", ErrInvalidString
		}

		// Числа не могут идти друг за другом.
		// Если предыдущее значение - это число, то отдаём ошибку.
		_, err = strconv.Atoi(string(source[i-1]))
		if err == nil {
			return "", ErrInvalidString
		}

		if repeatCount > 0 {
			b.WriteString(strings.Repeat(string(source[i-1]), repeatCount-1))
		} else {
			// Если пришёл 0, то нужно пересобрать строку, убрав последний добавленный символ
			tmpStr := []rune(b.String())
			tmpStr = tmpStr[:len(tmpStr)-1]
			b.Reset()
			for _, tmpV := range tmpStr {
				b.WriteRune(tmpV)
			}
		}
	}
	return b.String(), nil
}
