package main

import (
	"fmt"
	"log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Models
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Posts    []Post
}

type Post struct {
	gorm.Model
	Title  string `gorm:"not null"`
	Body   string
	UserID uint
	Tags   []Tag `gorm:"many2many:post_tags;"`
}

type Tag struct {
	gorm.Model
	Name string `gorm:"uniqueIndex"`
}

// Hook — hash password before saving
func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil { return err }
	u.Password = string(hashed)
	return nil
}

const dsn = "host=localhost user=postgres password=password dbname=siwes_db sslmode=disable"

func main() {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil { log.Fatal(err) }

	// AutoMigrate — create/update tables from structs
	db.AutoMigrate(&User{}, &Post{}, &Tag{})

	// Create user (BeforeCreate hook hashes password)
	user := User{Username: "david_o", Email: "david@example.com", Password: "secret123"}
	db.Create(&user)
	fmt.Printf("Created user: ID=%d Username=%s\n", user.ID, user.Username)

	// Create post with tags
	goTag  := Tag{Name: "golang"}
	devTag := Tag{Name: "devops"}
	post   := Post{Title: "Learning Go at SIWES", Body: "Week 17...", UserID: user.ID, Tags: []Tag{goTag, devTag}}
	db.Create(&post)

	// Eager load — preload user with their posts and tags
	var loaded User
	db.Preload("Posts.Tags").First(&loaded, user.ID)
	fmt.Printf("\nUser: %s\n", loaded.Username)
	for _, p := range loaded.Posts {
		fmt.Printf("  Post: %s | Tags:", p.Title)
		for _, t := range p.Tags { fmt.Printf(" #%s", t.Name) }
		fmt.Println()
	}

	// Update — partial update
	db.Model(&user).Updates(User{Email: "newemail@example.com"})

	// Soft delete — sets DeletedAt, row stays in DB
	db.Delete(&post)

	// Query with conditions
	var posts []Post
	db.Where("title LIKE ?", "%Go%").Find(&posts)
	fmt.Printf("\nPosts matching 'Go': %d\n", len(posts))

	// Raw SQL when needed
	var count int64
	db.Raw("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&count)
	fmt.Println("Active users:", count)
}