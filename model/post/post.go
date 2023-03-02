package post

import (
	"AirExcludes/db"
	"errors"
	"strings"
)

type Post struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
type Posts []Post

func GetPosts() Posts {
	var posts Posts
	con := db.ConnectDB()
	defer con.Close()
	var query = `SELECT Id, Name FROM krasecology.dbo.Svc_AirPosts where visible = 1 order by SortOrder`
	rows, err := con.Query(query)
	if err != nil {
		return posts
	}
	defer rows.Close()
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.Id, &post.Name); err != nil {
			return posts
		}
		posts = append(posts, post)
	}
	return posts
}

func (p Posts) FindPost(shortPostName string) (Post, error) {
	for _, post := range p {
		if strings.Contains(post.Name, shortPostName) {
			return post, nil
		}
	}
	return Post{}, errors.New("Пост " + shortPostName + " не найден")
}
