package models

import (
	"math/rand"

	"github.com/samsung-cnct/technical-on-boarding/onboarding/app/jobs/onboarding"
)

// User model object to manage user authentication
type User struct {
	ID       int
	Username string
	AuthEnv  *onboarding.AuthEnvironment
}

var db = make(map[int]*User)

// GetUser returns a User by id
func GetUser(id int) *User {
	return db[id]
}

// NewUser creates a new user
func NewUser() *User {
	user := &User{ID: rand.Intn(10000)}
	db[user.ID] = user
	return user
}
