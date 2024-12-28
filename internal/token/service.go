package token

import (
	"github.com/nmarsollier/commongo/cache"
	"github.com/nmarsollier/commongo/errs"
	"github.com/nmarsollier/commongo/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenService interface {
	Create(userID string) (*Token, error)
	Validate(tokenString string) (*Token, error)
	Invalidate(tokenString string) error
	Find(tokenID string) (*Token, error)
}

func NewTokenService(
	log log.LogRusEntry,
	cache cache.Cache[Token],
	repository TokenRepository,
) TokenService {
	return &tokenService{
		log:        log,
		cache:      cache,
		repository: repository,
	}
}

type tokenService struct {
	log        log.LogRusEntry
	cache      cache.Cache[Token]
	repository TokenRepository
}

func (s *tokenService) Create(userID string) (*Token, error) {
	_id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		s.log.Error(err)
		return nil, errs.Unauthorized
	}

	token, err := s.repository.Insert(_id)
	if err != nil {
		return nil, err
	}

	tokenString, err := Encode(token)
	if err != nil {
		return nil, err
	}

	s.cache.Add(tokenString, token)

	return token, nil
}

func (s *tokenService) Validate(tokenString string) (*Token, error) {
	if token, err := s.cache.Get(tokenString); err == nil {
		return token, err
	}

	tokenID, _, err := ExtractPayload(tokenString)
	if err != nil {
		return nil, err
	}

	token, err := s.repository.FindByID(tokenID)
	if err != nil || !token.Enabled {
		return nil, errs.Unauthorized
	}

	return token, nil
}

func (s *tokenService) Invalidate(tokenString string) error {
	tokenID, _, err := ExtractPayload(tokenString)
	if err != nil {
		return errs.Unauthorized
	}

	_id, err := primitive.ObjectIDFromHex(tokenID)
	if err != nil {
		s.log.Error(err)
		return errs.Unauthorized
	}

	if err = s.repository.Delete(_id); err != nil {
		return err
	}

	s.cache.Remove(tokenString)

	return nil
}

// Find busca un token en la db
func (s *tokenService) Find(tokenID string) (*Token, error) {
	token, err := s.repository.FindByID(tokenID)
	if err != nil {
		return nil, err
	}

	tokenString, err := Encode(token)
	if err != nil {
		return nil, err
	}

	s.cache.Add(tokenString, token)

	return token, nil
}
