package model

type User struct {
	UserId       string      `json:"user_id,omitempty"`
	UserName     string      `json:"user_name,omitempty"`
	Credentials  Credentials `json:"credentials,omitempty"`
	RegisterDate string      `json:"register_date,omitempty"`
	Root         bool        `json:"-"`
}

type Credentials struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"-"`
}
