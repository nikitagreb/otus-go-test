package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"

	json "github.com/json-iterator/go"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	res := make(DomainStat)
	if domain == "" {
		return res, errors.New("searched domain is empty")
	}
	sc := bufio.NewScanner(r)
	sc.Split(bufio.ScanLines)
	user := new(User)
	suffix := "." + domain
	for sc.Scan() {
		if err := json.Unmarshal(sc.Bytes(), user); err != nil {
			return res, err
		}
		if hasSuffix := strings.HasSuffix(user.Email, suffix); hasSuffix {
			d := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			res[d]++
		}
		user = &User{}
	}

	return res, nil
}
