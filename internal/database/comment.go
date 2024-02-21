package database

import (
	"database/sql"
	"forum/internal/models"
)

type CommentRepoImpl struct {
	db *sql.DB
}

func CreateNewCommentDB(db *sql.DB) *CommentRepoImpl {
	return &CommentRepoImpl{db}
}

func (cmnt *CommentRepoImpl) CreateCommentRepo(comment *models.Comment) (int64, error) {
	result, err := cmnt.db.Exec(`
	INSERT INTO comments (user_id, post_id, content, created_time, likes_counter, dislikes_counter) VALUES (?, ?, ?, ?, ?, ?);`,
		comment.UserID, comment.PostID, comment.Content, comment.CreatedTime, comment.LikesCounter, comment.DislikeCounter)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func (cmnt *CommentRepoImpl) GetAlCommentsForPost(postID int) ([]*models.Comment, error) {
	comments := []*models.Comment{}

	rows, err := cmnt.db.Query("SELECT * FROM comments ORDER BY created_time DESC")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var comment models.Comment
		err = rows.Scan(&comment.CommentID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedTime, &comment.LikesCounter, &comment.DislikeCounter)
		if err != nil {
			return nil, err
		}

		if postID == comment.PostID {
			comments = append(comments, &comment)
		}

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (cmnt *CommentRepoImpl) UpdateLikesCounter(commentID, valueToAdd int) error {
	_, err := cmnt.db.Exec("UPDATE comments SET likes_counter = likes_counter + ? WHERE id = ?", valueToAdd, commentID)
	if err != nil {
		return err
	}
	return nil
}

func (cmnt *CommentRepoImpl) UpdateDislikesCounter(commentID, valueToAdd int) error {
	_, err := cmnt.db.Exec("UPDATE comments SET dislikes_counter = dislikes_counter + ? WHERE id = ?", valueToAdd, commentID)
	if err != nil {
		return err
	}
	return nil
}

func (cmnt *CommentRepoImpl) GetCommentReaction(commentID, userID int) int {
	var reaction int
	if err := cmnt.db.QueryRow(
		`SELECT reaction FROM comment_votes WHERE comment_id = ? AND user_id = ?`,
		commentID, userID).Scan(&reaction); err != nil {
		return 0
	}
	return reaction
}

func (cmnt *CommentRepoImpl) AddReactionToCommentVotes(commentID, userID, reaction int) error {
	_, err := cmnt.db.Exec(`
		INSERT INTO comment_votes (comment_id, user_id,reaction) VALUES (?, ?, ?);`,
		commentID, userID, reaction)
	if err != nil {
		return err
	}
	return nil
}

func (cmnt *CommentRepoImpl) DeleteReactionFromCommentVotes(commentID, userID int) error {
	_, err := cmnt.db.Exec("DELETE FROM comment_votes WHERE comment_id = ? AND user_id = ?", commentID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (cmnt *CommentRepoImpl) UpdateReactionInCommentVotes(commentID, userID, newReaction int) error {
	_, err := cmnt.db.Exec("UPDATE comment_votes SET reaction = ? WHERE comment_id = ? AND user_id = ?", newReaction, commentID, userID)
	if err != nil {
		return err
	}
	return nil
}
