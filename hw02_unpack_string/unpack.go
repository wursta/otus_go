package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

//func main() {
//	str, err := Unpack("aaa10b")
//	fmt.Println(str, err)
//}

func Unpack(str string) (string, error) {
	//fmt.Println(str)
	source := []rune(str)

	b := strings.Builder{}

	for i, v := range source {
		repeatCount, err := strconv.Atoi(string(v))

		if err != nil {
			b.WriteRune(v)
		} else {
			// Если предыдущего элемента нет, то это первый элемент и он число -> то отдаём ошибку
			if i == 0 {
				return "", ErrInvalidString
			}

			// Числа не могут идти друг за другом.
			// Если предыдущее значение - это число отдаём ошибку.
			_, err := strconv.Atoi(string(source[i-1]))
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

		//fmt.Println(b.String())
	}

	return b.String(), nil
}
