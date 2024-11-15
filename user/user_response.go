package user

type UserData struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	Login       string   `json:"login"`
	Enabled     bool     `json:"enabled"`
}

func (UserData) IsEntity() {}
