package users

type CreateUserReq struct {
	FirstName   string
	LastName    string
	NickName    string
	Password    string
	Email       string
	CountryCode string
}

type ListUsersReq struct {
	Page,
	Size int
	Country string
}
