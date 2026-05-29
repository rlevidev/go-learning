package service

import (
	"context"
	"math/rand"
	"url-shortening/internal/db"
)

type UrlService struct {
	queries *db.Queries
}

func New(queries *db.Queries) *UrlService {
	return &UrlService{
		queries: queries,
	}
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"0123456789"

func generateShortCode() string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (s *UrlService) CreateShortURL(ctx context.Context, urlOriginal string) (*db.Url, error) {
	shortCode := generateShortCode()

	params := db.CreateURLParams{
		UrlOriginal: urlOriginal,
		ShortCode:   shortCode,
	}

	url, err := s.queries.CreateURL(ctx, params)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (s *UrlService) GetShortURL(ctx context.Context, urlShorted string) (*db.Url, error) {
	url, err := s.queries.GetURLByShortCode(ctx, urlShorted)
	if err != nil {
		return nil, err
	}

	if err := s.queries.IncrementAccessCount(ctx, urlShorted); err != nil {
		return nil, err
	}

	return &url, nil
}

func (s *UrlService) DeleteShortURL(ctx context.Context, urlShorted string) error {
	if err := s.queries.DeleteURL(ctx, urlShorted); err != nil {
		return err
	}
	return nil
}

func (s *UrlService) UpdateShortURL(ctx context.Context, urlShorted string, newUrlOriginal string) (*db.Url, error) {
	params := db.UpdateURLParams{
		ShortCode:   urlShorted,
		UrlOriginal: newUrlOriginal,
	}

	url, err := s.queries.UpdateURL(ctx, params)
	if err != nil {
		return nil, err
	}

	return &url, nil
}
