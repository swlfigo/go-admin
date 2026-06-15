package service

import (
	"go-admin/internal/model"
	"go-admin/internal/repository"
)

type LogService struct{ repo *repository.LogRepo }

func NewLogService(r *repository.LogRepo) *LogService { return &LogService{repo: r} }

func (s *LogService) List(keyword string, page, size int) ([]model.OperationLog, int64, error) {
	return s.repo.List(keyword, page, size)
}

func (s *LogService) Clear() error { return s.repo.Clear() }
