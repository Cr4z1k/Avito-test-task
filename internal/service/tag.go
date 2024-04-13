package service

import "github.com/Cr4z1k/Avito-test-task/internal/repository"

type TagService struct {
	r repository.Tag
}

func NewTagService(r repository.Tag) *TagService {
	return &TagService{r: r}
}

func (s *TagService) CreateTags(tagIDs []int) error {
	if err := s.r.CreateTags(tagIDs); err != nil {
		return err
	}

	return nil
}
