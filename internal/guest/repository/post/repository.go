package post

import (
	"context"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type Repository struct {
	db database.Client
}

func New(db database.Client) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) List(ctx context.Context, in entity.ListPostsInput) (entity.ListPostsOutput, error) {
	// TODO: get posts sql
	// sql := `
	// `

	// TODO: count posts
	// TODO: get posts

	return entity.ListPostsOutput{
		Posts: []entity.Post{
			{
				ID:          "uuid-1",
				Description: "hello-world-2",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          "uuid-2",
				Description: "hello-world-2",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
		Total: 2,
	}, nil
}

func (r *Repository) Get(ctx context.Context, postID string) (entity.Post, error) {
	// sql := `
	// 	SELECT
	//            id,
	//            user_id,
	//            description
	// 	FROM posts
	//        WHERE id = $1
	//    `
	// q := database.Query{
	// 	Name: "post_repository.Get",
	// 	Sql:  sql,
	// }
	//
	// row, err := r.db.DB().QueryContext(ctx, q, postID)
	// if err != nil {
	// 	return entity.Post{}, executeSQLError(err)
	// }
	// defer row.Close()
	//
	// var post Post
	// if err := pgxscan.ScanOne(&post, row); err != nil {
	// 	if errors.Is(err, pgx.ErrNoRows) {
	// 		return entity.Post{}, apperror.PostNotFoundID(postID)
	// 	}
	//
	// 	return entity.Post{}, scanRowError(err)
	// }

	// return post.ToEntity(), nil

	// example error
	return entity.Post{}, apperror.PostNotFoundID(postID)

	// return entity.Post{
	// 	ID:          postID,
	// 	Description: "the best place i've ever visited",
	// 	CreatedAt:   time.Now(),
	// 	UpdatedAt:   time.Now(),
	// }, nil
}
