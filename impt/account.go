package impt

type Any interface{}

type Account struct {
	Raw      string
	Email    string
	Password string
	Name string
}

func (a *Account) String() string {
	return "{raw: " + a.Raw + ", email: " + a.Email + ", password: " + a.Password + ", user name: " + a.Name + "}\n"
}
