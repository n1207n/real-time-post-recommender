package post

import (
	"github.com/brianvoe/gofakeit/v6"
	"math"
	"math/rand"
)

const (
	MIN_PARAGRAPH = 5
	MAX_PARAGRAPH = 10
	MIN_SENTENCE  = 5
	MAX_SENTENCE  = 10
	MIN_WORD      = 10
	MAX_WORD      = 20
)

func GeneratePostBody() string {
	paragraphSize := int(math.Ceil(rand.Float64() * (MAX_PARAGRAPH - MIN_PARAGRAPH)))
	sentenceSize := int(math.Ceil(rand.Float64() * (MAX_SENTENCE - MIN_SENTENCE)))
	wordSize := int(math.Ceil(rand.Float64() * (MAX_WORD - MIN_WORD)))

	return gofakeit.Paragraph(paragraphSize, sentenceSize, wordSize, "\n")
}

func GeneratePostTitle() string {
	wordSize := int(math.Ceil(rand.Float64() * (MAX_WORD - MIN_WORD)))
	return gofakeit.Sentence(wordSize)
}
