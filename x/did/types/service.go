package types

func NewService(sid, name, description string, status ServiceStatus, issuer string) Service {

	return Service{
		Sid:         sid,
		Name:        name,
		Description: description,
		Issuer:      issuer,
		Status:      status,
	}
}
