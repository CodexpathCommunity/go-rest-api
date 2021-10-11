package entity

import "time"

// Users represents a user.
type Users struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"`
	Name      string    `json:"name"`
	Country   string    `json:"country"`
	Score     int       `json:"score"`
	IsAuth    bool      `json:"is_auth"`
	AuthCode    string      `json:"auth_code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u Users) GetID() string {
	return u.ID
}

// GetRole returns the user role.
func (u Users) GetRole() string {
	return u.Role
}

// GetName returns the user name.
func (u Users) GetName() string {
	return u.Name
}

// GetCountry returns the user country.
func (u Users) GetCountry() string {
	return u.Country
}

// GetScore returns the user score.
func (u Users) GetScore() int {
	return u.Score
}

// GetIsAuth returns the user isAuth.
func (u Users) GetIsAuth() bool {
	return u.IsAuth
}
