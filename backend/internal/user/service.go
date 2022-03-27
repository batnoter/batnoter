package user

// Service represents a user service.
// It provides different methods to manage app user.
//go:generate mockgen -source=service.go -package=user -destination=mock_service.go
type Service interface {
	Get(userID uint) (User, error)
	GetByEmail(email string) (User, error)
	Save(user User) (uint, error)
	Delete(userID uint) error
}

type service struct {
	repo Repo
}

// NewService creates and return a new user service.
func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

// Get retrieves a user with given user id.
// It returns a user along with any error occurred while retrieving it.
func (s *service) Get(userID uint) (User, error) {
	user, err := s.repo.Get(userID)
	if err != nil {
		return user, err
	}
	return user, nil
}

// GetByEmail retrieves a user with given email.
// It returns a user along with any error occurred while retrieving it.
func (s *service) GetByEmail(email string) (User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return user, err
	}
	return user, nil
}

// Save stores the user.
// It returns the user id of the user along with any error occurred while storing the user.
func (s *service) Save(user User) (uint, error) {
	userID, err := s.repo.Save(user)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// Delete deletes the user with given user id.
// It returns any error occurred while deleting the user.
func (s *service) Delete(userID uint) error {
	if err := s.repo.Delete(userID); err != nil {
		return err
	}
	return nil
}
