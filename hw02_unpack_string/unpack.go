package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	source := []rune(str)
	sLen := len(source)

	if len(source) == 0 {
		return "", nil
	}

	// Если первый символ это число, то отдаём ошибку
	_, isNum := getNum(source[0])
	if isNum {
		return "", ErrInvalidString
	}

	b := strings.Builder{}

	for i, v := range source {
		repeatCount, isNum := getNum(v)

		if !isNum {
			isLast := i == sLen-1
			if isLast {
				b.WriteRune(v)
			} else {
				// Если следующий символ - число и он равен 0, то пропускаем символ
				nextNum, isNum := getNum(source[i+1])
				if (isNum && nextNum > 0) || !isNum {
					b.WriteRune(v)
				}
			}
			continue
		}

		// Числа не могут идти друг за другом.
		// Если предыдущее значение - это число, то отдаём ошибку.
		_, isNum = getNum(source[i-1])
		if isNum {
			return "", ErrInvalidString
		}

		// Нули обрабатываются в блоке считывания букв
		if repeatCount == 0 {
			continue
		}

		b.WriteString(strings.Repeat(string(source[i-1]), repeatCount-1))
	}
	return b.String(), nil
}

func getNum(letter rune) (int, bool) {
	num, err := strconv.Atoi(string(letter))
	if err != nil {
		return 0, false
	}
	return num, err == nil
}
