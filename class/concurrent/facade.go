package main

import "fmt"

type Facade struct {
	UserSvc    UserSvc
	ArticleSvc ArticleSvc
	CommentSvc CommentSvc
}

// 用户登录
func (f *Facade) login(name, password string) int {
	user := f.UserSvc.GetUser(name)
	if password == user.password {
		fmt.Println("登录成功！！！")
	}
	return user.id
}

func (f *Facade) CreateArticle(userId int, title, content string) *Article {
	articleId := 12345
	article := f.ArticleSvc.Create(articleId, title, content, userId)
	return article
}

func (f *Facade) CreateComment(articleId int, userId int, comment string) *Comment {
	commentId := 12345
	cm := f.CommentSvc.Create(commentId, comment, articleId, userId)
	return cm
}

// 用户服务
type UserSvc struct {
}

type User struct {
	id       int
	name     string
	password string
}

func (user *UserSvc) GetUser(name string) *User {
	if name == "zhangsan" {
		return &User{
			id:       12345,
			name:     "zhangsan",
			password: "zhangsan",
		}
	} else {
		return &User{}
	}
}

// 文章服务
type ArticleSvc struct {
}

type Article struct {
	articleId int
	title     string
	content   string
	authorId  int
}

func (articleSvc *ArticleSvc) Create(articleId int, title string, content string, userId int) *Article {
	return &Article{
		articleId: articleId,
		title:     title,
		content:   content,
		authorId:  userId,
	}
}

// 评论服务
type CommentSvc struct {
}

type Comment struct {
	commentId int
	comment   string
	articleId int
	userId    int
}

func (commentSvc *CommentSvc) Create(commentId int, comment string, articleId int, userId int) *Comment {
	return &Comment{
		commentId: commentId,
		comment:   comment,
		articleId: articleId,
		userId:    userId,
	}
}

func main() {
	f := &Facade{}
	userId := f.login("zhangsan", "zhangsan")
	fmt.Println("登录成功,当前用户Id", userId)

	title := "go设计模式外观模式"
	content := "外观模式是结构模式的一种。。。。"
	article := f.CreateArticle(userId, title, content)
	fmt.Println("文章发表成功,文章id", article.articleId)

	comment := f.CreateComment(article.articleId, userId, "介绍的很详细")
	fmt.Println("评论提交成功,评论id", comment.commentId)
}
