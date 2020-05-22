// +build !bench

package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	data = `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`
	shortData = `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}`
	user      = User{1, "Howard Mendoza", "0Oliver", "aliquid_qui_ea@Browsedrive.gov", "6-866-899-36-79", "InAQJvsq", "Blackbird Place 25"}
)

func TestGetDomainStat(t *testing.T) {
	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}

func Test_countDomains(t *testing.T) {
	uArr := users{user, user, user, user, user, user}
	t.Run("counting", func(t *testing.T) {
		result, err := countDomains(uArr, "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 6}, result)
	})
	t.Run("nil users", func(t *testing.T) {
		var u users
		result, err := countDomains(u, "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
	t.Run("empty domain", func(t *testing.T) {
		result, err := countDomains(uArr, "")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}

func Test_getUsers(t *testing.T) {
	t.Run("unmarshal", func(t *testing.T) {
		result, err := getUsers(bytes.NewBufferString(shortData))
		require.NoError(t, err)
		require.Equal(t, users{user}, result)
	})
}

func BenchmarkGetDomainStat(b *testing.B) {
	doms := []string{"com", "gov", "net", "unknown"}
	for _, d := range doms {
		if _, err := GetDomainStat(bytes.NewBufferString(data), d); err != nil {
			log.Fatal(err)
		}
	}
}
