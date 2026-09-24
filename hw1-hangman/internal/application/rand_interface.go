package application

//go:generate mockgen -source=rand_interface.go -destination=../mocks/mock_rand.go -package=mocks

type Rand interface {
	GetRandomWord(category, level int) ([]rune, error)
	GetRandomCategory() int
	GetRandomDifficulty() int
}
