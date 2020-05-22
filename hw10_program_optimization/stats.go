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
	return countDomains(u, domain), nil
}

type users [100000]User

func getUsers(r io.Reader) (*users, error) {
	var result users
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return &result, err
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		var user User
		if err = json.Unmarshal([]byte(line), &user); err != nil {
			return &result, err
		}
		result[i] = user
	}
	return &result, nil
}

func countDomains(u *users, domain string) DomainStat {
	result := make(DomainStat)

	for _, user := range u {
		if user.Email != "" && domain != "" && strings.HasSuffix(user.Email, domain) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result
}
