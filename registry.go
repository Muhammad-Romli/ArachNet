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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 { // main.go already take care of the Args of the command so param start from 0
		return fmt.Errorf("Login function need 1 parameter at least")
	}
	name := cmd.Args[0]
	user, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("failed getting user data, name might not exist: %w", err)
	}
	s.cfg.SetUser(name)
	fmt.Println("User has been set:")
	fmt.Println(user)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Register function need 1 parameter at least")
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

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) < 2 { // main.go already take care of the Args of the command so param start from 0
		return fmt.Errorf("login function need 2 parameters at least, name and url")
	}
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("fail getting user info: %w", err)
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
