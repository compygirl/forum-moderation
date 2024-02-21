package migration

import (
	"context"
	"database/sql"
)

func CreateDb(dbName, dbPath string, ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "data/forum.db")
	if err != nil {
		return nil, err
	}
	// db.SetMaxIdleConns(100)
	if err = CreateAllTables(ctx, db); err != nil {
		return nil, err
	}
	return db, nil
}

func CreateAllTables(ctx context.Context, db *sql.DB) error {
	trans, err := db.Begin()
	if err != nil {
		trans.Rollback()
		return err
	}

	// create user table
	if _, err := db.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            firstName TEXT,
			secondName TEXT,
			usernames TEXT UNIQUE,
			email TEXT UNIQUE,
			password TEXT,
			role TEXT
        );
    `); err != nil {
		trans.Rollback()
		return err
	}

	// create session table
	if _, err := db.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS sessions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER UNIQUE,
			token TEXT UNIQUE,
			exp_time DATE,
			FOREIGN KEY (user_id) REFERENCES user (id)
        );
    `); err != nil {
		trans.Rollback()
		return err
	}

	// create post table
	if _, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			title TEXT, 
			content TEXT,
			created_time DATE,
			likes_counter INTEGER, 
			dislikes_counter INTEGER,
			image_path TEXT,
			FOREIGN KEY (user_id) REFERENCES users (id)
		)
	`); err != nil {
		trans.Rollback()
		return err
	}

	// create post_category
	if _, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS post_category(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER,
			category_name TEXT,
			FOREIGN KEY (post_id) REFERENCES posts (id)
		)
	`); err != nil {
		trans.Rollback()
		return err
	}

	// create post_votes
	if _, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS post_votes(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER,
			user_id INTEGER,
			reaction INTEGER,
			FOREIGN KEY (post_id) REFERENCES posts (id)
			FOREIGN KEY (user_id) REFERENCES users (id)
		)
	`); err != nil {
		trans.Rollback()
		return err
	}

	// create comments
	// create post table
	if _, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER,
			user_id INTEGER, 
			content TEXT,
			created_time DATE,
			likes_counter INTEGER, 
			dislikes_counter INTEGER,
			FOREIGN KEY (post_id) REFERENCES posts (id)
			FOREIGN KEY (user_id) REFERENCES users (id)
		)
	`); err != nil {
		trans.Rollback()
		return err
	}

	// create comment_votes
	if _, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS comment_votes(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			comment_id INTEGER,
			user_id INTEGER,
			reaction INTEGER,
			FOREIGN KEY (comment_id) REFERENCES comments (id)
			FOREIGN KEY (user_id) REFERENCES users (id)
		)
	`); err != nil {
		trans.Rollback()
		return err
	}

	return err
}
