package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/machinebox/graphql"
)

const Endpoint = "https://api.github.com/graphql"

type Status struct {
	Message string `json:"message"`
}
type User struct {
	Login  string  `json:"login"`
	Status *Status `json:"status"`
}
type Response struct {
	User User `json:"user"`
}

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("Please create a Github PAT and export it as `GITHUB_TOKEN`.")
	}
	if len(os.Args) < 1 {
		log.Fatal("Usage: `finger [github username]`")
	}
	user := os.Args[1]
	client := graphql.NewClient(Endpoint)
	req := graphql.NewRequest(`
		query ($login: String!) {
			user (login: $login) {
				login
				status {
					message
				}
			}
		}
`)

	// set any variables
	req.Var("login", user)
	req.Header.Add("Authorization", "bearer "+token)
	var resp Response
	c, cf := context.WithTimeout(context.Background(), 4*time.Second)
	defer cf()

	if err := client.Run(c, req, &resp); err != nil {
		log.Fatal(err)
	}
	if resp.User.Status != nil {
		fmt.Printf("@%s: %s\n", user, resp.User.Status.Message)
		return
	}
	fmt.Printf("No status for %s\n", user)

}

/*
query {
  user(login: "bketelsen") {
		login
    status {
      message
    }
  }
}
*/
