package login

/*
/   This is just a small file to define the interfaces that
/   should be used if a particular login is required.
*/

type Loginer interface {
	Login(user, password string) (string, error)
}

type SocialLoginer interface {
	Login(token string) string
}
