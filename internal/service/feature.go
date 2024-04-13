package service

import "github.com/Cr4z1k/Avito-test-task/internal/repository"

type FeatureService struct {
	r repository.Feature
}

func NewFeatureService(r repository.Feature) *FeatureService {
	return &FeatureService{r: r}
}

func (s *FeatureService) CreateFeatures(featureIDs []int) error {
	if err := s.r.CreateFeatures(featureIDs); err != nil {
		return err
	}

	return nil
}
