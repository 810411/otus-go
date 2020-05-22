package hw10_program_optimization //nolint:golint,stylecheck

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %s", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		var user User
		if err = json.Unmarshal([]byte(line), &user); err != nil {
			return
		}
		result[i] = user
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	var b strings.Builder
	b.WriteRune('.')
	b.WriteString(domain)
	suf := b.String()

	result := make(DomainStat)
	for _, user := range u {
		if strings.HasSuffix(user.Email, suf) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
