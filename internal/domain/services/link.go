package services

import (
	"context"
	"url-shortener-api/internal/domain/entities"
)

type Storage interface {
	GetOne(ctx context.Context, id string) (*entities.Link, error)
	Create(ctx context.Context, URL string, linkId string) (*entities.Link, error)
}

type linkService struct {
	storage Storage
}

func NewService(storage Storage) *linkService {
	return &linkService{storage: storage}
}

func (s *linkService) Create(ctx context.Context, URL string, linkId string) (*entities.Link, error) {
	link, err := s.storage.Create(ctx, URL, linkId)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (s *linkService) GetById(ctx context.Context, id string) (*entities.Link, error) {
	link, err := s.storage.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	return link, nil
}
