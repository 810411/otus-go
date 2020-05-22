package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"fmt"
	"io"
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

	reader := bufio.NewReader(r)
	for i := 0; ; i++ {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return &result, err
			}
		}
		var user User
		if err = json.Unmarshal(line, &user); err != nil {
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
