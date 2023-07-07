package post

import (
	"github.com/brianvoe/gofakeit/v6"
	"math"
	"math/rand"
)

type FakerInterface interface {
	Paragraph(paragraphs, sentences, words int, separator string) string
	Sentence(wordCount int) string
}

type Faker struct{}

func (f *Faker) Paragraph(paragraphs, sentences, words int, separator string) string {
	return gofakeit.Paragraph(paragraphs, sentences, words, separator)
}

func (f *Faker) Sentence(wordCount int) string {
	return gofakeit.Sentence(wordCount)
}

func GeneratePostBody(minParagraph int, maxParagraph int, minSentence int, maxSentence int, minWord int, maxWord int, faker FakerInterface) string {
	paragraphSize := int(math.Ceil(rand.Float64() * float64(maxParagraph-minParagraph)))
	sentenceSize := int(math.Ceil(rand.Float64() * float64(maxSentence-minSentence)))
	wordSize := int(math.Ceil(rand.Float64() * float64(maxWord-minWord)))

	return faker.Paragraph(paragraphSize, sentenceSize, wordSize, "\n")
}

func GeneratePostTitle(minWord int, maxWord int, faker FakerInterface) string {
	wordSize := int(math.Ceil(rand.Float64() * float64(maxWord-minWord)))
	return faker.Sentence(wordSize)
}
