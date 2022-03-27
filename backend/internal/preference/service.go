package preference

// Service represents a preference service.
// It provides methods to manager user preferences.
//go:generate mockgen -source=service.go -package=preference -destination=mock_service.go
type Service interface {
	Save(defaultRepo DefaultRepo) error
	GetByUserID(userID uint) (DefaultRepo, error)
}

type serviceImpl struct {
	repo Repo
}

// NewService creates and return a new preference service.
func NewService(repo Repo) Service {
	return &serviceImpl{
		repo: repo,
	}
}

// Save stores the user default repository.
// It returns any error occurred while storing it.
func (s *serviceImpl) Save(defaultRepo DefaultRepo) error {
	if err := s.repo.Save(defaultRepo); err != nil {
		return err
	}
	return nil
}

// GetByUserID retrieves user's default repository.
// It returns user's default repository along with any error occurred while retrieving it.
func (s *serviceImpl) GetByUserID(userID uint) (DefaultRepo, error) {
	defaultRepo, err := s.repo.GetByUserID(userID)
	if err != nil {
		return defaultRepo, err
	}
	return defaultRepo, nil
}
