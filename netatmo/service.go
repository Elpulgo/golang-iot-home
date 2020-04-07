package netatmo

type Repository interface {
	GetCurrent(chan CurrentResult)
}

type Service interface {
	GetCurrent() (Current, error)
}

type service struct {
	repository Repository
}

type CurrentResult struct {
	Current Current
	Error   error
}

func New(repository Repository) Service {
	return &service{repository: repository}
}

func (service *service) GetCurrent() (Current, error) {
	channel := make(chan CurrentResult)

	go service.repository.GetCurrent(channel)

	response := <-channel

	return response.Current, response.Error

}
