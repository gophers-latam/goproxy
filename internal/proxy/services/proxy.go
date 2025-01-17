package proxy

import (
	domain "github.com/kenriortega/ngonx/internal/proxy/domain"
)

// ProxyService interface service for proxy repository funcionalities
type ProxyService interface {
	SaveSecretKEY(string, string, string) error
	GetKEY(string, string) (string, error)
}

// DefaultProxyService struct for management proxy repository
type DefaultProxyService struct {
	repo domain.ProxyRepository
}

// NewProxyService return new DefaultProxyService
func NewProxyService(repository domain.ProxyRepository) DefaultProxyService {
	return DefaultProxyService{repo: repository}
}

// SaveSecretKEY save secret key
func (s DefaultProxyService) SaveSecretKEY(engine, key, apikey string) (string, error) {

	err := s.repo.SaveKEY(engine, key, apikey)
	if err != nil {
		return "failed", err
	}
	return "ok", nil
}

// GetKEY get key
func (s DefaultProxyService) GetKEY(engine, key string) (string, error) {
	result, err := s.repo.GetKEY(engine, key)
	if err != nil {
		return "failed", err
	}
	return result, nil
}
