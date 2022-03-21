package preference

//go:generate mockgen -source=service.go -package=preference -destination=mock_service.go
type Service interface {
	Save(defaultRepo DefaultRepo) error
	GetByUserID(userId uint) (DefaultRepo, error)
}

type serviceImpl struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &serviceImpl{
		repo: repo,
	}
}

func (s *serviceImpl) Save(defaultRepo DefaultRepo) error {
	if err := s.repo.Save(defaultRepo); err != nil {
		return err
	}
	return nil
}

func (s *serviceImpl) GetByUserID(userID uint) (DefaultRepo, error) {
	defaultRepo, err := s.repo.GetByUserID(userID)
	if err != nil {
		return defaultRepo, err
	}
	return defaultRepo, nil
}
