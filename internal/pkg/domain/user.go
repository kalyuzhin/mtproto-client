package domain

type User struct {
	id          int64
	firstName   string
	lastName    string
	userName    string
	phoneNumber string
}

func (u *User) GetUserID() int64 {
	return u.id
}

func (u *User) GetFirstName() string {
	return u.firstName
}

func (u *User) GetLastName() string {
	return u.lastName
}

func (u *User) GetUserName() string {
	return u.userName
}

func (u *User) GetPhoneNumber() string {
	return u.phoneNumber
}
