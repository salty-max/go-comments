package comment

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
)

var (
	ErrFetchingComment = errors.New("failed to fetch comment by id")
	ErrCreatingComment = errors.New("failed to create comment")
	ErrUpdatingComment = errors.New("failed to update comment")
	ErrDeletingComment = errors.New("failed to delete comment")
	ErrNotImplemented  = errors.New("not implemented")
)

// Comment - a representation of the comment
// structure for the service
type Comment struct {
	ID     string
	Slug   string
	Body   string
	Author string
}

// Store - defines all methods that the service
// needs in order to operate
type Store interface {
	GetComment(context.Context, string) (Comment, error)
	CreateComment(context.Context, Comment) (Comment, error)
	UpdateComment(context.Context, string, Comment) (Comment, error)
	DeleteComment(context.Context, string) error
}

// Service - the struct on which all
// logic will be built on top of
type Service struct {
	Store Store
}

// NewService - returns a pointer to a new service
func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

func (s *Service) GetComment(ctx context.Context, id string) (Comment, error) {
	cmt, err := s.Store.GetComment(ctx, id)
	if err != nil {
		log.Error(err)
		return Comment{}, ErrFetchingComment
	}

	return cmt, nil
}

func (s *Service) CreateComment(ctx context.Context, cmt Comment) (Comment, error) {
	insertedCmt, err := s.Store.CreateComment(ctx, cmt)
	if err != nil {
		log.Error(err)
		return Comment{}, ErrCreatingComment
	}

	return insertedCmt, nil
}

func (s *Service) UpdateComment(ctx context.Context, id string, cmt Comment) (Comment, error) {
	updatedCmt, err := s.Store.UpdateComment(ctx, id, cmt)
	if err != nil {
		log.Error(err)
		return Comment{}, ErrUpdatingComment
	}

	return updatedCmt, nil
}

func (s *Service) DeleteComment(ctx context.Context, id string) error {
	return s.Store.DeleteComment(ctx, id)
}
