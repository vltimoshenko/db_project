package service

import (
	"github.com/db_project/app/forum"
	. "github.com/db_project/pkg/models"
)

type Service struct {
	Repository forum.RepositoryInterface
}

func (s Service) ClearDB() error {
	return s.Repository.ClearDB()
}

func (s Service) GetStatus() (Status, error) {
	return s.Repository.GetStatus()
}
