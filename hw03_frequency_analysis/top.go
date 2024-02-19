package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(str string) []string {
	if str == "" {
		return nil
	}

	words := strings.Fields(str)
	if len(words) == 0 {
		return nil
	}

	// Наполняем мапу количеством упоминаний, обрабатывая слова и обрезая некоторые символы
	wordsMap := make(map[string]int, len(words))
	for _, v := range words {
		if v == "-" {
			continue
		}

		word := trimAndLowerWord(v)
		if _, ok := wordsMap[word]; !ok {
			wordsMap[word] = 0
		}
		wordsMap[word]++
	}

	// Создаём слайс со всеми словами
	top := make([]string, len(wordsMap))
	i := 0
	for word := range wordsMap {
		top[i] = word
		i++
	}

	// Сортируем слайс по значениям упоминаний из мапы.
	// Если количесто упоминаний одинаковое, то сравниваем по алфавиту.
	sort.Slice(top, func(i, j int) bool {
		wi := top[i]
		wj := top[j]
		if wordsMap[wi] == wordsMap[wj] {
			return strings.Compare(wi, wj) <= 0
		}

		return wordsMap[wi] > wordsMap[wj]
	})

	if len(top) > 10 {
		return top[:10]
	}
	return top
}

// Приводит к нижнему регистру и обрезает символы !(воскл. знак), '(кавычка), .(точка), ,(запятая).
func trimAndLowerWord(v string) string {
	v = strings.ToLower(v)
	v = strings.Trim(v, "!',.")
	return v
}
