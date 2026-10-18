package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	_ "github.com/lib/pq"
	"github.com/muhammad-romli/ArachNet/internal/config"
	"github.com/muhammad-romli/ArachNet/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	CommandsMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	function, ok := c.CommandsMap[cmd.Name]
	if !ok {
		return fmt.Errorf("Your command doesn't exist in register CommandsMap")
	}
	err := function(s, cmd)
	if err != nil {
		return fmt.Errorf("error happen when running command: %w", err)
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.CommandsMap[name] = f
}

func getFeedByUrl(s *state, url string) (database.Feed, error) {
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return database.Feed{}, fmt.Errorf("fail getting feed info from database by url: %w", err)
	}
	return feed, nil
}

func getFeedByName(s *state, name string) (database.Feed, error) {
	feed, err := s.db.GetFeedByName(context.Background(), name)
	if err != nil {
		return database.Feed{}, fmt.Errorf("fail getting feed info from database by name: %w", err)
	}
	return feed, nil
}

func getFeedByAll(s *state, query string) (database.Feed, error) {
	feed, errURL := getFeedByUrl(s, query)
	if errURL != nil {
		var errName error
		feed, errName = getFeedByName(s, query)
		if errName != nil {
			return database.Feed{}, fmt.Errorf("fail getting feed by either name and url:\n error by Url: %w\n error by Name: %w", errURL, errName)
		}
	}
	return feed, nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 { // main.go already take care of the Args of the command so param start from 0
		return fmt.Errorf("Login function need 1 parameter at least, needed: user_name")
	}
	name := cmd.Args[0]
	user, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("failed getting user data, name might not exist: %w", err)
	}
	s.cfg.SetUser(name)
	fmt.Println("User has been set:")
	fmt.Println("User ID:", user.ID)
	fmt.Println("User Name:", user.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Register function need 1 parameter at least, needed: user_name")
	}
	id := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now()
	name := cmd.Args[0]
	userParam := database.CreateUserParams{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Name:      name,
	}

	user, err := s.db.CreateUser(context.Background(), userParam)
	if err != nil {
		return fmt.Errorf("failed inserting user to database %w", err)
	}
	s.cfg.SetUser(name)
	fmt.Printf("User is created, congratulations!\n")
	fmt.Println("User ID:", user.ID)
	fmt.Println("User Name:", user.Name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.ResetTable(context.Background()); err != nil {
		return fmt.Errorf("failed resetting the table: %w", err)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetListUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed getting list of users: %w", err)
	}
	for _, u := range users {
		fmt.Printf("* %s", u.Name)
		if u.Name == s.cfg.CurrentUserName {
			fmt.Printf(" (current)")
		}
		fmt.Printf("\n")
	}
	return nil
}

func handlerAggregate(s *state, cmd command) error {
	rssStruct, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error while fetching data: %w", err)
	}
	processXML(rssStruct)
	for _, item := range rssStruct.Channel.Item {
		fmt.Println(item.Title)
		fmt.Println(item.Description)
		fmt.Printf("\n----------------------------------------------------\n")
	}
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 { // main.go already take care of the Args of the command so param start from 0
		return fmt.Errorf("addfeed command need 2 parameters at least, name and url")
	}
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("fail creating feed: %w", err)
	}
	fmt.Println("Feed successfully created")
	fmt.Println("Feed ID:", feed.ID)
	fmt.Println("Feed Name:", feed.Name)
	fmt.Println("Feed Url:", feed.Url)
	tempCmd := command{
		Args: []string{feed.Url},
	}
	fmt.Println("==============================================")
	if err = handlerFollow(s, tempCmd, user); err != nil {
		return fmt.Errorf("fail following created feed, feed is created: %w", err)
	}
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("fail when getting list feeds from database: %w", err)
	}
	fmt.Print("\n")
	for _, f := range feeds {
		fmt.Println("Feed name:", f.Name)
		fmt.Println("URL:", f.Url)
		fmt.Println("Username:", f.UserName)
		fmt.Println("\n----------------------------------------------")
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Follow function need 1 parameter at least, needed: feed_url or feed_name")
	}
	query := cmd.Args[0]

	feed, err := getFeedByAll(s, query)
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("fail inserting follow information to database: %w", err)
	}
	fmt.Println("Congratulations you succesfully follow the feed")
	fmt.Println("Your ID:", feedFollow.UserID)
	fmt.Println("Followed feed ID:", feedFollow.FeedID)
	fmt.Println("Followed feed URL:", cmd.Args[0])
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feedsFollowed, err := s.db.ListFollowing(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed getting list of followed feeds: %w", err)
	}
	fmt.Println("==================================================")
	fmt.Println("Your ID:", user.ID)
	fmt.Printf("Feed that you have followed:\n")
	for _, f := range feedsFollowed {
		fmt.Println("-------------------------------------------------")
		fmt.Println("Feed ID:", f.FeedID)
		fmt.Println("Feed Name:", f.FeedName)
	}
	fmt.Println("==================================================")
	return nil
}

func handlerUnfollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("unfollow function need 1 parameter at least, needed: feed_url or feed_name")
	}
	query := cmd.Args[0]
	feed, err := getFeedByAll(s, query)
	if err != nil {
		return err
	}
	params := database.UnfollowingFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	if err := s.db.UnfollowingFeed(context.Background(), params); err != nil {
		return fmt.Errorf("fail unfollowing feed: %w", err)
	}
	return nil
}
