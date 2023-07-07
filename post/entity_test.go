package post

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPost(t *testing.T) {
	title := "Test Title"
	body := "Test Body"

	post := NewPost(title, body)

	assert.NotNil(t, post)
	assert.Equal(t, title, post.Title)
	assert.Equal(t, body, post.Body)
	_, err := uuid.Parse(post.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, 0, post.Votes)
}
