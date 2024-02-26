package database

import (
	"database/sql"
	"forum/internal/models"
)

type PostRepoImpl struct {
	db *sql.DB
}

func CreateNewPostDB(db *sql.DB) *PostRepoImpl {
	return &PostRepoImpl{db}
}

func (postObj *PostRepoImpl) CreatePostRepo(post *models.Post) (int64, error) {
	result, err := postObj.db.Exec(`
		INSERT INTO posts (user_id, title, content, created_time, likes_counter, dislikes_counter, image_path) VALUES (?, ?, ?, ?, ?, ?, ?);`,
		post.UserID, post.Title, post.Content, post.CreatedTime, post.LikesCounter, post.DislikeCounter, post.ImagePath)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func (postObj *PostRepoImpl) GetAllPosts() ([]*models.Post, error) {
	posts := []*models.Post{}
	rows, err := postObj.db.Query("SELECT * FROM posts ORDER BY created_time DESC")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.CreatedTime, &post.LikesCounter, &post.DislikeCounter, &post.ImagePath)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (postObj *PostRepoImpl) GetCategoriesByPostID(postID int) ([]string, error) {
	categories := []string{}
	rows, err := postObj.db.Query("SELECT category_name FROM post_category WHERE post_id = ?", postID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var category string
		if err = rows.Scan(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (postObj *PostRepoImpl) CreatePostCategory(categories []string, postID int) (int64, error) {
	var err error
	var result sql.Result

	for _, category := range categories {
		result, err = postObj.db.Exec(`
		INSERT INTO post_category (category_name, post_id) VALUES (?, ?);`,
			category, postID)
		if err != nil {
			return -1, err
		}

	}
	return result.LastInsertId()
}

func (postObj *PostRepoImpl) UpdateLikesCounter(postID, valueToAdd int) error {
	_, err := postObj.db.Exec("UPDATE posts SET likes_counter = likes_counter + ? WHERE id = ?", valueToAdd, postID)
	if err != nil {
		return err
	}
	return nil
}

func (postObj *PostRepoImpl) UpdateDislikesCounter(postID, valueToAdd int) error {
	_, err := postObj.db.Exec("UPDATE posts SET dislikes_counter = dislikes_counter + ? WHERE id = ?", valueToAdd, postID)
	if err != nil {
		return err
	}
	return nil
}

func (postObj *PostRepoImpl) GetReaction(postID, userID int) (int, error) {
	var reaction int
	if err := postObj.db.QueryRow(
		`SELECT reaction FROM post_votes WHERE post_id = ? AND user_id = ?`,
		postID, userID).Scan(&reaction); err != nil {
		return 0, err
	}
	return reaction, nil
}

func (postObj *PostRepoImpl) AddReactionToPostVotes(postID, userID, reaction int) error {
	_, err := postObj.db.Exec(`
		INSERT INTO post_votes (post_id, user_id,reaction) VALUES (?, ?, ?);`,
		postID, userID, reaction)
	if err != nil {
		return err
	}
	return nil
}

func (postObj *PostRepoImpl) DeleteFromPostVotes(postID, userID int) error {
	_, err := postObj.db.Exec("DELETE FROM post_votes WHERE post_id = ? AND user_id = ?", postID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (postObj *PostRepoImpl) UpdateReactionInPostVotes(postID, userID, newReaction int) error {
	_, err := postObj.db.Exec("UPDATE post_votes SET reaction = ? WHERE post_id = ? AND user_id = ?", newReaction, postID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (postObj *PostRepoImpl) GetPostByID(postID int) (*models.Post, error) {
	post := &models.Post{}

	if err := postObj.db.QueryRow(
		`SELECT id, user_id, title, content, created_time, likes_counter, dislikes_counter FROM posts WHERE id = ?`,
		postID).Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.CreatedTime, &post.LikesCounter, &post.DislikeCounter); err != nil {
		return nil, err
	}
	return post, nil
}

func (postObj *PostRepoImpl) GetPostsByCategory(category string) ([]*models.Post, error) {
	posts := []*models.Post{}

	rows, err := postObj.db.Query(`
	SELECT * FROM posts WHERE id IN (SELECT post_id FROM post_category WHERE category_name = ?) ORDER BY created_time DESC
	`, category)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.CreatedTime, &post.LikesCounter, &post.DislikeCounter)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (postObj *PostRepoImpl) GetPostsByUserId(userID int) ([]*models.Post, error) {
	posts := []*models.Post{}
	rows, err := postObj.db.Query(`
	SELECT * FROM posts WHERE user_id = ? ORDER BY created_time DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.CreatedTime, &post.LikesCounter, &post.DislikeCounter)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (postObj *PostRepoImpl) GetPostsByLikes(userID int) ([]*models.Post, error) {
	posts := []*models.Post{}
	rows, err := postObj.db.Query(`
	SELECT * FROM posts WHERE id IN (SELECT post_id FROM post_votes WHERE user_id = ?) ORDER BY created_time DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.CreatedTime, &post.LikesCounter, &post.DislikeCounter)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (postObj *PostRepoImpl) DeletePostByID(postID int) error {
	_, err := postObj.db.Exec("DELETE FROM posts WHERE id = ? ", postID)
	if err != nil {
		return err
	}
	return nil
}

func (postObj *PostRepoImpl) DeletePostCategoryByPostID(postID int) error {
	_, err := postObj.db.Exec("DELETE FROM post_category WHERE post_id = ? ", postID)
	if err != nil {
		return err
	}
	return nil
}

func (postObj *PostRepoImpl) DeleteAllPostVotesByPostID(postID int) error {
	_, err := postObj.db.Exec("DELETE FROM post_votes WHERE post_id = ? ", postID)
	if err != nil {
		return err
	}
	return nil
}
