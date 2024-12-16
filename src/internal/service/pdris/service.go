package pdris

import (
	"context"

	"go.uber.org/zap"

	"gihub.com/bongerka/sberPDRIS/internal/repository"
	def "gihub.com/bongerka/sberPDRIS/internal/service"
)

type service struct {
	repo repository.PdrisRepository
}

func NewService(repo repository.PdrisRepository) def.PdrisService {
	return &service{
		repo: repo,
	}
}

func (s *service) UpdateValue(ctx context.Context, value int) error {
	err := s.repo.UpdateValue(ctx, value)
	if err != nil {
		zap.L().Error("unable to update value", zap.Error(err))
	}

	return nil
}

func (s *service) GetValue(ctx context.Context) (int, error) {
	val, err := s.repo.GetValue(ctx)
	if err != nil {
		zap.L().Error("unable to get value", zap.Error(err), zap.Int("val", val))
	}

	return val, nil
}
