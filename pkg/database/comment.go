package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/salty-max/go-comments/pkg/comment"
	uuid "github.com/satori/go.uuid"
)

type CommentRow struct {
	ID     string
	Slug   sql.NullString
	Body   sql.NullString
	Author sql.NullString
}

func convertRowToComment(c CommentRow) comment.Comment {
	return comment.Comment{
		ID:     c.ID,
		Slug:   c.Slug.String,
		Body:   c.Body.String,
		Author: c.Author.String,
	}
}

func (d *Database) GetComments(ctx context.Context) ([]comment.Comment, error) {
	var cmtRow CommentRow
	var comments []comment.Comment

	_, err := d.Client.ExecContext(ctx, `SELECT pg_sleep(16)`)
	if err != nil {
		return []comment.Comment{}, err
	}

	rows, err := d.Client.QueryContext(
		ctx,
		`
		SELECT id, slug, body, author FROM comments	
		`,
	)
	if err != nil {
		return []comment.Comment{}, fmt.Errorf("error fetching comments: %w", err)
	}

	for rows.Next() {
		if err := rows.Scan(&cmtRow.ID, &cmtRow.Slug, &cmtRow.Body, &cmtRow.Author); err != nil {
			return []comment.Comment{}, fmt.Errorf("error fetching comments: %w", err)
		}
		comments = append(comments, convertRowToComment(cmtRow))
	}

	if err := rows.Close(); err != nil {
		return []comment.Comment{}, fmt.Errorf("failed to close rows: %w", err)
	}

	return comments, nil
}

func (d *Database) GetComment(ctx context.Context, uuid string) (comment.Comment, error) {
	var cmtRow CommentRow

	row := d.Client.QueryRowContext(
		ctx,
		`
		SELECT id, slug, body, author
		FROM comments
		WHERE id = $1
		`,
		uuid,
	)

	err := row.Scan(&cmtRow.ID, &cmtRow.Slug, &cmtRow.Body, &cmtRow.Author)
	if err != nil {
		return comment.Comment{}, fmt.Errorf("error fetching comment by uuid: %w", err)
	}

	return convertRowToComment(cmtRow), nil
}

func (d *Database) CreateComment(ctx context.Context, cmt comment.Comment) (comment.Comment, error) {
	cmt.ID = uuid.NewV4().String()
	cmtRow := CommentRow{
		ID:     cmt.ID,
		Slug:   sql.NullString{String: cmt.Slug, Valid: true},
		Body:   sql.NullString{String: cmt.Body, Valid: true},
		Author: sql.NullString{String: cmt.Author, Valid: true},
	}

	rows, err := d.Client.NamedQueryContext(
		ctx,
		`
		INSERT INTO comments
		(id, slug, body, author)
		VALUES (:id, :slug, :body, :author)
		`,
		cmtRow,
	)
	if err != nil {
		return comment.Comment{}, fmt.Errorf("error inserting comment: %w", err)
	}

	if err := rows.Close(); err != nil {
		return comment.Comment{}, fmt.Errorf("failed to close rows: %w", err)
	}

	return convertRowToComment(cmtRow), nil
}

func (d *Database) DeleteComment(ctx context.Context, id string) error {
	_, err := d.Client.ExecContext(
		ctx,
		`
		DELETE FROM comments
		WHERE id = $1
		`,
		id,
	)
	if err != nil {
		return fmt.Errorf("error deleting comment: %w", err)
	}

	return nil
}

func (d *Database) UpdateComment(
	ctx context.Context,
	id string,
	cmt comment.Comment,
) (comment.Comment, error) {
	cmtRow := CommentRow{
		ID:     id,
		Slug:   sql.NullString{String: cmt.Slug, Valid: true},
		Body:   sql.NullString{String: cmt.Body, Valid: true},
		Author: sql.NullString{String: cmt.Author, Valid: true},
	}

	rows, err := d.Client.NamedQueryContext(
		ctx,
		`
		UPDATE comments SET
		slug = :slug,
		body = :body,
		author = :author
		WHERE id = :id	
		`,
		cmtRow,
	)
	if err != nil {
		return comment.Comment{}, fmt.Errorf("error updating comment: %w", err)
	}

	if err := rows.Close(); err != nil {
		return comment.Comment{}, fmt.Errorf("failed to close rows: %w", err)
	}

	return convertRowToComment(cmtRow), nil
}
