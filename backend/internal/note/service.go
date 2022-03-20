package note

//go:generate mockgen -source=service.go -package=note -destination=mock_service.go
type Service interface {
	GetAll(email string) ([]Note, error)
	Get(noteId int) (Note, error)
	Save(note Note) error
	Delete(noteId int) error
}

type serviceImpl struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &serviceImpl{
		repo: repo,
	}
}

func (s *serviceImpl) GetAll(email string) ([]Note, error) {
	notes, err := s.repo.GetAll(email)
	if err != nil {
		return notes, err
	}
	return notes, nil
}

func (s *serviceImpl) Get(noteId int) (Note, error) {
	note, err := s.repo.Get(noteId)
	if err != nil {
		return note, err
	}
	return note, nil
}

func (s *serviceImpl) Save(note Note) error {
	if err := s.repo.Save(note); err != nil {
		return err
	}
	return nil
}

func (s *serviceImpl) Delete(noteId int) error {
	if err := s.repo.Delete(noteId); err != nil {
		return err
	}
	return nil
}
