package user

//go:generate mockgen -source=service.go -package=user -destination=mock_service.go
type Service interface {
	Get(userID uint) (User, error)
	GetByEmail(email string) (User, error)
	Save(user User) (uint, error)
	Delete(userID uint) error
}

type serviceImpl struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &serviceImpl{
		repo: repo,
	}
}

func (s *serviceImpl) Get(userID uint) (User, error) {
	user, err := s.repo.Get(userID)
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

func (s *serviceImpl) Save(user User) (uint, error) {
	userID, err := s.repo.Save(user)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (s *serviceImpl) Delete(userID uint) error {
	if err := s.repo.Delete(userID); err != nil {
		return err
	}
	return nil
}
