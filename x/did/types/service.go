package types

func NewService(sid, name, description string, status ServiceStatus, issuers []string) Service {

	return Service{
		Sid:         sid,
		Name:        name,
		Description: description,
		Issuers:     issuers,
		Status:      status,
	}
}
