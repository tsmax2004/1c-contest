package app

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"gihub.com/bongerka/sberPDRIS/internal/repository"
	pdrisRepo "gihub.com/bongerka/sberPDRIS/internal/repository/pdris"
	"gihub.com/bongerka/sberPDRIS/internal/service"
	pdrisService "gihub.com/bongerka/sberPDRIS/internal/service/pdris"
)

type serviceProvider struct {
	conn *pgxpool.Pool

	repository repository.PdrisRepository
	service    service.PdrisService
}

func NewServiceProvider(conn *pgxpool.Pool) *serviceProvider {
	return &serviceProvider{
		conn: conn,
	}
}

func (s *serviceProvider) Repository() repository.PdrisRepository {
	if s.repository == nil {
		s.repository = pdrisRepo.NewRepository(s.conn)
	}

	return s.repository
}

func (s *serviceProvider) Service() service.PdrisService {
	if s.service == nil {
		s.service = pdrisService.NewService(s.Repository())
	}

	return s.service
}
