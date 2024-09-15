package util

import (
	"errors"
	"strings"
	"unicode"
)

var (
	ErrHeadersNotFound  = errors.New("can't find headers")
	ErrTableParseFailed = errors.New("unable to parse table")
)

func GetHeadsIDs(header []string, heads []string) (map[string]int, error) {
	res := make(map[string]int, len(heads))
	for i, head := range header {
		for _, neededHhead := range heads {
			if head == neededHhead {
				if _, ok := res[head]; ok {
					return nil, ErrHeadersNotFound
				}
				res[head] = i
			}
		}
	}
	if len(res) != len(heads) {
		return nil, ErrHeadersNotFound
	}
	return res, nil
}

func ParseTable(table string, lineSep string, skipLines, maxLines int) ([][]string, error) {
	lines := strings.Split(table, lineSep)
	if len(lines) <= skipLines {
		return nil, ErrTableParseFailed
	}
	trimmed := make([]string, 0, len(lines)-skipLines)
	for _, l := range lines[skipLines:] {
		trimmed = append(trimmed, strings.TrimSpace(l))
	}
	return splitTable(trimmed, delimetersFinder(trimmed, maxLines)), nil
}

// private section

func delimetersFinder(lines []string, maxLines int) []int {
	var colDelims []int
	var isPossibleDelimeter bool

	if maxLines < 2 {
		maxLines = 2
	}
	if len(lines) > maxLines {
		lines = lines[:maxLines]
	}
	pos := 0
	// Движемся по всем строкам для анализа с начала в конец.
	// Определение разделителя:
	//	- на всех строках на данной позиции стоит пробельный симовол
	//	- хотя бы на одной из строк следующий символ - не пробельный
	// Если хотя бы в одной строке дошли до конца - завершаем работу
	//! Не подходит для таблиц, где могут быть пустые "последние" поля до концв

	for {
		charFound := false
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			if len(line) <= pos {
				return colDelims
			}
			if !unicode.IsSpace([]rune(line)[pos]) {
				charFound = true
				if isPossibleDelimeter {
					colDelims = append(colDelims, pos)
				}
				break
			}
		}
		isPossibleDelimeter = !charFound
		pos++
	}
}

func splitTable(lines []string, colDelims []int) [][]string {
	res := make([][]string, 0, len(lines))
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}
		res = append(res, make([]string, 0, len(colDelims)+1))
		beg := 0
		for _, end := range colDelims {
			res[i] = append(res[i], strings.TrimSpace(line[beg:end]))
			beg = end
		}
		res[i] = append(res[i], strings.TrimSpace(line[beg:]))
	}
	return res
}
