package netatmo

type Repository interface {
	GetCurrent() (Current, error)
}

type Service interface {
	GetCurrent() (Current, error)
}

type service struct {
	repository Repository
}

func New(repository Repository) Service {
	return &service{repository: repository}
}

func (service *service) GetCurrent() (Current, error) {
	return service.repository.GetCurrent()
}
