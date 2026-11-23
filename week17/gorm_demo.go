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

// Hook - hash password before saving
func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

const dsn = "host=localhost user=postgres password=password dbname=siwes_gorm sslmode=disable"

func main() {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// AutoMigrate - create/update tables from structs
	if err := db.AutoMigrate(&User{}, &Post{}, &Tag{}); err != nil {
		log.Fatal(err)
	}

	// Create user (BeforeCreate hook hashes password)
	user := User{Username: "michael_o", Email: "michealoluwatoba14@gmail.com", Password: "secret123"}
	if err := db.Create(&user).Error; err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created user: ID=%d Username=%s\n", user.ID, user.Username)

	// Create post with tags
	goTag := Tag{Name: "golang"}
	devTag := Tag{Name: "devops"}
	post := Post{Title: "Learning Go at SIWES", Body: "Week 17...", UserID: user.ID, Tags: []Tag{goTag, devTag}}
	if err := db.Create(&post).Error; err != nil {
		log.Fatal(err)
	}

	// Eager load - preload user with their posts and tags
	var loaded User
	db.Preload("Posts.Tags").First(&loaded, user.ID)
	fmt.Printf("\nUser: %s\n", loaded.Username)
	for _, p := range loaded.Posts {
		fmt.Printf("  Post: %s | Tags:", p.Title)
		for _, t := range p.Tags {
			fmt.Printf(" #%s", t.Name)
		}
		fmt.Println()
	}

	// Update - partial update
	db.Model(&user).Updates(User{Email: "newemail@example.com"})

	// Query with conditions (before the soft delete, so the post is found)
	var posts []Post
	db.Where("title LIKE ?", "%Go%").Find(&posts)
	fmt.Printf("\nPosts matching 'Go': %d\n", len(posts))

	// Soft delete - sets DeletedAt, row stays in DB
	db.Delete(&post)

	// Normal queries no longer see it...
	var visible []Post
	db.Where("title LIKE ?", "%Go%").Find(&visible)
	fmt.Printf("After soft delete (normal query): %d\n", len(visible))

	// ...but Unscoped still can
	var all []Post
	db.Unscoped().Where("title LIKE ?", "%Go%").Find(&all)
	fmt.Printf("After soft delete (Unscoped):     %d\n", len(all))

	// Raw SQL when needed
	var count int64
	db.Raw("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&count)
	fmt.Println("Active users:", count)
}