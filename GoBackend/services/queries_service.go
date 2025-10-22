package services

import (
	"frontdesk/models"
	"frontdesk/repositories"
)

type QueriesService interface {
	CreateQuery(query *models.Query) error
	GetQueries() ([]models.Query, error)
	ResolveQuery(queryStatus *models.QueryStatus, id int) error
	GetFAQs() ([]models.FAQ, error)
}

type queriesService struct {
	queriesRepository repositories.QueriesRepository
}

func NewQueriesService(queriesRepository repositories.QueriesRepository) *queriesService {
	return &queriesService{queriesRepository: queriesRepository}
}

func (s *queriesService) CreateQuery(query *models.Query) error {
	return s.queriesRepository.SaveQuery(query)
}

func (s *queriesService) GetQueries() ([]models.Query, error) {
	return s.queriesRepository.FetchQueries()
}

func (s *queriesService) ResolveQuery(queryStatus *models.QueryStatus, id int) error {
	return s.queriesRepository.UpdateQueryStatus(queryStatus, id)
}

func (s *queriesService) GetFAQs() ([]models.FAQ, error) {
	faqs, err := s.queriesRepository.FetchFAQs()

	if err != nil {
		return nil, err
	}

	return faqs, nil
}
