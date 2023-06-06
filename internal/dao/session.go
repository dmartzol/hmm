package dao

import (
	"context"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"

	"github.com/dmartzol/hmm/internal/dao/postgres"
	"github.com/dmartzol/hmm/internal/hmm"
)

type SessionService struct {
	DB *postgres.DB
}

func NewSessionService(db *sqlx.DB) *SessionService {
	ss := SessionService{
		DB: &postgres.DB{DB: db},
	}
	return &ss
}

func (ss SessionService) Create(ctx context.Context, email, password string) (*hmm.Session, error) {
	ctx, span := otel.Tracer("dao").Start(ctx, "SessionService.CreateSession")
	defer span.End()

	session, err := ss.DB.CreateSession(ctx, email, password)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (ss SessionService) SessionFromToken(token string) (*hmm.Session, error) {
	session, err := ss.DB.SessionFromToken(token)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (ss SessionService) ExpireSession(token string) (*hmm.Session, error) {
	session, err := ss.DB.ExpireSession(token)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (ss SessionService) UpdateSession(token string) (*hmm.Session, error) {
	session, err := ss.DB.UpdateSession(token)
	if err != nil {
		return nil, err
	}
	return session, nil
}
