package user

//go:generate mockgen -source=service.go -package=user -destination=mock_service.go
type Service interface {
	Get(userId uint) (User, error)
	GetByEmail(email string) (User, error)
	Save(user User) error
	Delete(userId uint) error
}

type serviceImpl struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &serviceImpl{
		repo: repo,
	}
}

func (s *serviceImpl) Get(userId uint) (User, error) {
	user, err := s.repo.Get(userId)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *serviceImpl) GetByEmail(email string) (User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *serviceImpl) Save(user User) error {
	if err := s.repo.Save(user); err != nil {
		return err
	}
	return nil
}

func (s *serviceImpl) Delete(userId uint) error {
	if err := s.repo.Delete(userId); err != nil {
		return err
	}
	return nil
}
