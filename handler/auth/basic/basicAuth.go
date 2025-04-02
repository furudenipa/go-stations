package basic

import "github.com/TechBowl-japan/go-stations/handler/auth"

func Authenticate(userID, password string, c *auth.Config) bool {
	if userID == c.UserID && password == c.Password {
		return true
	}
	return false
}
