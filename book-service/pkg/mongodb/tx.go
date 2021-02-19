package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// Session wrap mongo session function needed to perform transaction operations
type Session interface {
	WithTransaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) (interface{}, error)) (interface{}, error)
	EndSession(context.Context)
}

// TXRepository wrap mongo function needed to perform transaction operations
type TXRepository interface {
	StartSession() (Session, error)
}

type txRepository struct {
	db *mongo.Database
}

// NewTXRepository returns new NewTXRepository
func NewTXRepository(db *mongo.Database) TXRepository {
	return &txRepository{
		db: db,
	}
}

func (r *txRepository) StartSession() (Session, error) {
	sess, err := r.db.Client().StartSession()
	return &session{
		session: sess,
	}, err
}

type session struct {
	session mongo.Session
}

func (s *session) EndSession(ctx context.Context) {
	s.session.EndSession(ctx)
}

func (s *session) WithTransaction(
	ctx context.Context,
	fn func(sessCtx mongo.SessionContext) (interface{}, error),
) (interface{}, error) {
	return s.session.WithTransaction(ctx, fn)
}
