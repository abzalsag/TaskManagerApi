package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type ExternalClient struct {
	client *http.Client
}

func NewExternalClient() *ExternalClient {
	return &ExternalClient{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (c *ExternalClient) GetPost(
	ctx context.Context,
	id int,
) (Post, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://jsonplaceholder.typicode.com/posts/"+strconv.Itoa(id),
		nil,
	)
	if err != nil {
		return Post{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return Post{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Post{}, fmt.Errorf(
			"external API returned status: %d",
			resp.StatusCode,
		)
	}

	var post Post

	err = json.NewDecoder(resp.Body).Decode(&post)
	if err != nil {
		return Post{}, err
	}

	return post, nil
}
