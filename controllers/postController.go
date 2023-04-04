package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	database "hello/db"
	"hello/models"
	"net/http"
	"net/smtp"
	"time"
)

var postCollection *mongo.Collection = database.OpenCollection(database.Client, "post")

func CreateUserPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var a string
		r := make([]string, 0, 0)
		g := make([]models.Answers, 0, 0)
		if c.PostForm("Frontend") == "on" {
			r = append(r, "Frontend")
		}
		if c.PostForm("Backend") == "on" {
			r = append(r, "Backend")
		}
		if c.PostForm("design") == "on" {
			r = append(r, "Design")
		}
		if c.PostForm("SQL") == "on" {
			r = append(r, "SQL")
		}
		if c.PostForm("Nosql") == "on" {
			r = append(r, "NoSQL")
		}
		if c.PostForm("game") == "on" {
			r = append(r, "GameDev")
		}
		a, _ = c.Cookie("username")
		local, _ := time.LoadLocation("Asia/Almaty")
		timea := time.Now().In(local)
		post1 := models.Post{
			Username:  a,
			Title:     c.PostForm("title"),
			Text:      c.PostForm("text"),
			Answers:   g,
			Tags:      r,
			CreatedAt: timea,
		}
		resultInsertionNumber, insertErr := postCollection.InsertOne(ctx, post1)
		if insertErr != nil {
			msg := fmt.Sprintf("Post item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		fmt.Println(resultInsertionNumber)
		y, _ := FindAllSubs(context.Background())
		var k int
		for _, _ = range y {
			k++
		}
		q := make([]string, k)
		for z, w := range y {
			q[z] = *w.Email
		}
		auth := smtp.PlainAuth(
			"",
			"ddev05702@gmail.com",
			"pjudcojtdhlfpshs",
			"smtp.gmail.com",
		)
		msg := "Hello User,we have new post in our site."
		err := smtp.SendMail("smtp.gmail.com:587",
			auth,
			"akmagambetovaanel0@gmail.com",
			q,
			[]byte(msg),
		)
		if err != nil {
			fmt.Println(err)
		}
		c.Redirect(303, "/forum")
	}

}

func Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "createpost.html", nil)
	}
}
func SeeSinglePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		a, _ := c.Params.Get("title")
		b, _ := c.Cookie("username")
		fmt.Println(a)
		x, _ := FindOne(ctx, a)
		if x.Username != b {
			FindAndUpdateView(context.Background(), a)
		}
		c.HTML(http.StatusOK, "singlePost.html", x)
		defer cancel()

	}
}
func SeeAllPosts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		a, _ := FindAllPost(ctx)
		c.HTML(http.StatusOK, "ForumPosts.html", a)
		defer cancel()
	}
}
func LeaveAnswer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var a string
		a, _ = c.Cookie("username")
		b, _ := c.Params.Get("title")
		FindAndUpdateCount(context.Background(), b)
		answer := models.Answers{
			Text:      c.PostForm("answers"),
			Username:  a,
			CreatedAt: time.Now(),
		}
		filter := bson.D{{"title", b}}
		update := bson.D{{"$push", bson.D{{"answers", answer}}}}
		resultInsertionNumber, insertErr := postCollection.UpdateOne(ctx, filter, update)
		if insertErr != nil {
			msg := fmt.Sprintf("Post item was not updated")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		fmt.Println(resultInsertionNumber)
		c.Redirect(303, "/forum")
	}
}
