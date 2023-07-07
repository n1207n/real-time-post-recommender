package post

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockFaker struct {
	mock.Mock
}

func (m *MockFaker) Paragraph(paragraphs, sentences, words int, separator string) string {
	args := m.Called(paragraphs, sentences, words, separator)
	return args.String(0)
}

func (m *MockFaker) Sentence(wordCount int) string {
	args := m.Called(wordCount)
	return args.String(0)
}

func TestGeneratePostBody(t *testing.T) {
	mockFaker := new(MockFaker)

	mockFaker.On("Paragraph", mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("string")).
		Return("Mocked Body")

	body := GeneratePostBody(5, 10, 5, 10, 10, 20, mockFaker)

	assert.Equal(t, "Mocked Body", body)

	mockFaker.AssertExpectations(t)
}

func TestGeneratePostTitle(t *testing.T) {
	mockFaker := new(MockFaker)

	mockFaker.On("Sentence", mock.AnythingOfType("int")).
		Return("Mocked Title")

	title := GeneratePostTitle(10, 20, mockFaker)

	assert.Equal(t, "Mocked Title", title)

	mockFaker.AssertExpectations(t)
}
