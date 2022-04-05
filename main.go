package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

type Post struct {
	ID          int64
	Title       string
	Body        string
	isPublished bool
}

func main() {
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "recordings",
		AllowNativePasswords: true,
	}
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	postID, err := addPost(Post{
		Title:       "Section Post",
		Body:        "Section Post body with simple text",
		isPublished: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added post: %v\n", postID)

	posts, err := postsByTitle("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Posts found: %v\n", posts)

	post, err := postByID(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Post found: %v\n", post)
	// update post
	updatePost(1, Post{
		Title:       "Docker and kubernetes",
		Body:        "Frank Doe",
		isPublished: true,
	})
	//delete post
	deletePost(7)
}

func postsByTitle(title string) ([]Post, error) {
	var posts []Post

	rows, err := db.Query("SELECT * FROM posts WHERE title = ?", title)
	if err != nil {
		return nil, fmt.Errorf("postsByTitle %q: %v", title, err)
	}
	defer rows.Close()
	for rows.Next() {
		var alb Post
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Body, &alb.isPublished); err != nil {
			return nil, fmt.Errorf("postsByBody %q: %v", title, err)
		}
		posts = append(posts, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postsByBody %q: %v", title, err)
	}
	return posts, nil
}

func postByID(id int64) (Post, error) {
	var alb Post
	row := db.QueryRow("SELECT * FROM posts WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Body, &alb.isPublished); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("postsById %d: no such post", id)
		}
		return alb, fmt.Errorf("postsById %d: %v", id, err)
	}
	return alb, nil
}

func addPost(alb Post) (int64, error) {
	result, err := db.Exec("INSERT INTO posts (title, body, is_published) VALUES (?, ?, ?)", alb.Title, alb.Body, alb.isPublished)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}

func updatePost(id int64, alb Post) (int64, error) {
	result, err := db.Exec("UPDATE posts SET title=?, body=?, is_published=? WHERE id=?", alb.Title, alb.Body, alb.isPublished, id)
	if err != nil {
		return 0, fmt.Errorf("updateAlbum: %v", err)
	}
	id, err = result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("updateAlbum: %v", err)
	}
	return id, nil
}

func deletePost(id int64) error {
	_, err := db.Exec("DELETE FROM posts WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("deleteAlbum: %v", err)
	}
	return nil
}
