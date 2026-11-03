package service

type Service struct {
	// Dependencies will be injected here
}

func NewService() *Service {
	return &Service{}
}
