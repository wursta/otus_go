package hw10programoptimization

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/mailru/easyjson"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (result DomainStat, err error) {
	result = DomainStat{}
	var regex *regexp.Regexp
	if regex, err = regexp.Compile("@.*\\." + domain); err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		var user User
		if err = easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}

		matched := regex.MatchString(user.Email)
		if matched {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
