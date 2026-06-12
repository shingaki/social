package db

import (
	"SOCIAL/internal/store"
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
)

var usernames = []string{
	"Sam", "Patty", "Alex", "Jamie", "Chris", "Taylor", "Jordan", "Casey",
	"Riley", "Morgan", "Drew", "Cody", "Avery", "Quinn", "Jesse", "Dakota",
	"Reese", "Blake", "Corey", "Skyler", "Rowan", "Parker", "Emerson", "Harley",
	"Bailey", "Cameron", "Elliot", "Finley", "Hayden", "Kendall", "Logan",
	"Micah", "Noel", "Oakley", "Peyton", "Robin", "Sage", "Shawn", "Terry",
	"Tracy", "Tyler", "Wesley", "Zion", "Andy", "Bobby", "Danny", "Freddy", "Johnny", "Lenny", "Tommy",
}

var titles = []string{
	"Quick Wins for Productivity",
	"Mastering Daily Habits",
	"Simple Ways to Stay Focused",
	"Boost Your Workflow Today",
	"Small Changes, Big Results",
	"How to Stay Organized",
	"The Power of Consistency",
	"Work Smarter, Not Harder",
	"Tips for Better Time Management",
	"Getting More Done Faster",
	"Habits That Actually Stick",
	"Level Up Your Routine",
	"Focus Like a Pro",
	"Declutter Your Mind",
	"Stay Motivated Every Day",
	"Plan Less, Do More",
	"Beat Procrastination Now",
	"Minimalism for Productivity",
	"Win Your Morning Routine",
	"Stay Sharp and Productive",
}

var contents = []string{
	"Small habits create big change over time. Start with one simple improvement today and build momentum.",
	"Productivity isn’t about doing more—it’s about doing what matters most with clarity and focus.",
	"Consistency beats intensity. Showing up every day will take you further than bursts of effort.",
	"Take control of your day by planning just three key tasks. Simplicity drives execution.",
	"Distractions are everywhere, but your attention is valuable. Protect it like your most important asset.",
	"Success often comes from refining the basics. Master the small things and results will follow.",
	"Your environment shapes your behavior. Set up your space to support your goals.",
	"Stop waiting for motivation. Action creates motivation, not the other way around.",
	"Clear goals turn effort into progress. Define what success looks like before you begin.",
	"Break large tasks into smaller steps. Progress becomes easier when the path is visible.",
	"Focus on progress, not perfection. Done is always better than perfect.",
	"Time is limited, so spend it intentionally. Prioritize what truly moves you forward.",
	"Momentum builds when you take action daily. Even small wins add up quickly.",
	"Clarity reduces stress. Know what you’re working toward and why it matters.",
	"Your mindset determines your direction. Choose growth over comfort every time.",
	"Eliminate what doesn’t serve you. Simplicity creates space for meaningful work.",
	"Energy management matters as much as time management. Work when you’re at your best.",
	"Reflection helps you improve. Take time to review what worked and what didn’t.",
	"Discipline creates freedom. The more consistent you are, the more control you gain.",
	"Start before you feel ready. Progress begins the moment you take the first step.",
}

var tags = []string{
	"productivity", "tech", "lifestyle", "coding", "golang", "react", "backend", "frontend", "cloud",
	"aws", "docker", "devops", "api", "design", "startup", "career", "learning", "tutorial", "security", "database",
}

var comments = []string{
	"Great post! I'm glad you enjoyed it.",
	"I learned a lot from this post. Thanks for sharing!",
	"This is a great resource for learning about productivity.",
	"I'm glad you enjoyed reading this post.",
	"I'm impressed with your knowledge of productivity.",
	"I've been following this post for a while now.",
	"This post is a great way to stay motivated.",
	"I'm glad you found this post useful.",
	"I'm excited to learn more about productivity.",
}

// Prep seed data
func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	// Create users
	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)
	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("error creating user: ", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("error creating post: ", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("error creating comment: ", err)
			return
		}
	}
	log.Println("seed data created")
}

func generateUsers(num int) (users []*store.User) {
	users = make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Role: store.Role{
				Name: "user",
			},
		}
	}
	return users
}

// Create posts
// Need posts to be related to a user
func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}
	return posts
}

// Create comments
// Need comments to be related to a post and user
func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)

	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}
	return cms
}
