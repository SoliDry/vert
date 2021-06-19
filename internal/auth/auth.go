package auth

import (
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/arthurkushman/pgo"
)

const alphabetForRandom = "abcdefghijklmnopqrstuvwxyz0123456789"

type User struct {
	login      string
	password   string
	isLoggedIn bool
}

func NewUser(login, password string, isLoggedIn bool) *User {
	return &User{login: login, password: password, isLoggedIn: isLoggedIn}
}

type Authenticator interface {
	Authenticate() (ok bool, err error)
}

//Authenticate authenticates user by login/password
func (user *User) Authenticate() (ok bool, err error) {

	return ok, err
}

func calcPassword(password, salt string) (string, error) {
	sha1Pwd := pgo.Sha1(password)
	sha1Mix := pgo.Sha1(salt + pgo.Sha1(pgo.Sha1(password)))

	binSha1Pwd, err := hex.DecodeString(sha1Pwd)
	if err != nil {
		return "", err
	}

	binSha1Mix, err := hex.DecodeString(sha1Mix)
	if err != nil {
		return "", err
	}

	var xorBytes []byte
	for i, binNum := range binSha1Pwd {
		xorBytes = append(xorBytes, binNum^binSha1Mix[i])
	}

	return hex.EncodeToString(xorBytes), nil
}

func genSalt() string {
	var symbols = []rune(alphabetForRandom)
	rand.Seed(time.Now().UnixNano())

	var s string
	for i := 0; i < 37; i++ {
		randInt := rand.Intn(37-0) + 0
		s += string(symbols[randInt])
	}

	return s
}
