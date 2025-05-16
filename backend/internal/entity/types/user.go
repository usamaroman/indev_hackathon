package types

type UserType string

const (
	AdminUser UserType = "admin"
	Customer  UserType = "customer"
)

func (t UserType) String() string {
	return string(t)
}

var UserTypeToString = map[UserType]string{
	AdminUser: AdminUser.String(),
	Customer:  Customer.String(),
}

var StringToUserType = map[string]UserType{
	"admin":    AdminUser,
	"customer": Customer,
}


