//go:build integration

package database

import (
	"context"
	"testing"

	"github.com/salty-max/go-comments/pkg/comment"
	"github.com/stretchr/testify/assert"
)

func TestCommentDatabase(t *testing.T) {
	t.Run("gets comments", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		_, err = db.CreateComment(context.Background(), comment.Comment{
			Slug:   "slug",
			Body:   "body",
			Author: "author",
		})
		assert.NoError(t, err)

		cmts, err := db.GetComments(context.Background())
		assert.NoError(t, err)
		assert.NotEmpty(t, cmts)
	})

	t.Run("creates and gets a comment", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		cmt, err := db.CreateComment(context.Background(), comment.Comment{
			Slug:   "slug",
			Body:   "body",
			Author: "author",
		})
		assert.NoError(t, err)

		newCmt, err := db.GetComment(context.Background(), cmt.ID)
		assert.NoError(t, err)
		assert.Equal(t, "slug", newCmt.Slug)
	})

	t.Run("updates a comment", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		cmt, err := db.CreateComment(context.Background(), comment.Comment{
			Slug:   "slug",
			Body:   "body",
			Author: "author",
		})
		assert.NoError(t, err)

		updatedCmt, err := db.UpdateComment(context.Background(), cmt.ID, comment.Comment{
			Slug:   "test",
			Body:   "test de julien à nouveau",
			Author: "julien",
		})
		assert.NoError(t, err)
		assert.Equal(t, "julien", updatedCmt.Author)
	})

	t.Run("deletes a comment", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		cmt, err := db.CreateComment(context.Background(), comment.Comment{
			Slug:   "slug",
			Body:   "body",
			Author: "author",
		})
		assert.NoError(t, err)

		err = db.DeleteComment(context.Background(), cmt.ID)
		assert.NoError(t, err)

		_, err = db.GetComment(context.Background(), cmt.ID)
		assert.Error(t, err)
	})
}
