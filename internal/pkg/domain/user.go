package domain

type User struct {
	id          int64
	firstName   string
	lastName    string
	userName    string
	phoneNumber string
}

func (u *User) GetUserName() string {
	return u.userName
}
