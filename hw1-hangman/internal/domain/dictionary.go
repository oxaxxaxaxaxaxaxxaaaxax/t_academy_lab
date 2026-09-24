package domain

//go:generate mockgen -source=dictionary.go -destination=../mocks/mock_rand.go -package=mocks

import (
	"math/rand"
	"time"
)

type WordRand struct{}

func NewWordRand() *WordRand {
	return &WordRand{}
}

var rng = rand.New(rand.NewSource(time.Now().Unix()))

func GetRandInt(n int) int {
	return rng.Intn(n)
}

func (r *WordRand) GetRandomWord(category, level int) ([]rune, error) {
	words := dict[Category(category)][Difficulty(level)]
	return []rune(words[GetRandInt(len(words))]), nil
}

func (r *WordRand) GetRandomDifficulty() int {
	return GetRandInt(DifficultyCount())
}

func (r *WordRand) GetRandomCategory() int {
	return GetRandInt(CategoryCount())
}

var helpDescription = map[string]string{
	"cat": "^._.^", "dog": "not a cat", "fox": "cunning orange",
	"wolf": "gray hunter", "bear": "honey enjoyer",
	"lion": "savannah king", "frog": "green jumper",
	"fish": "O{",

	"parrot": "talking bird", "giraffe": "long neck",
	"jaguar": "spotted big cat",
	"monkey": "almost human", "turtle": "tortilla",

	"chameleon": "color-changing lizard", "orangutan": "red ape",
	"echidna": "almost hedgehog", "axolotl": "smiling salamander",

	"sun": "our star", "moon": "earth's satellite", "star": "make a wish", "comet": "icy tail",
	"earth": "our planet", "mars": "red planet",

	"galaxy": "star island", "shuttle": "spaceplane", "saturn": "with a ring",
	"uranus": "ice giant", "space": "the universe",

	"interstellar": "between stars",
	"astrobiology": "life in space", "exoplanet": "planet elsewhere",
	"gravitational": "due to gravity",
	"ionosphere":    "charged sky layer",

	"fantasy": "magic worlds", "comedy": "humorous plot",
	"drama": ";( sad", "horror": "scares",

	"romance":   "love story",
	"adventure": "quest", "memoir": "about life ",
	"detective": "Agatha Christie",

	"psychodrama": "psycho + drama", "tragicomedy": "almost comedy but :(",
	"surrealism": "dreamlike logic",
	"steampunk":  "gears and steam", "cyberpunk": "high tech",
}

var dict = map[Category]map[Difficulty][]string{
	Animals: {
		Easy: {"cat", "dog", "fox",
			"wolf", "bear",
			"lion", "frog",
			"fish"},
		Medium: {"parrot", "giraffe",
			"jaguar",
			"monkey", "turtle"},
		Hard: {"chameleon", "orangutan",
			"echidna", "axolotl"},
	},

	Space: {
		Easy: {"sun", "moon", "star", "comet",
			"earth", "mars"},
		Medium: {"galaxy", "shuttle", "saturn",
			"uranus", "cosmos"},
		Hard: {"interstellar",
			"astrobiology", "exoplanet",
			"gravitational",
			"ionosphere"},
	},

	LiteratureGenre: {
		Easy: {"fantasy", "comedy",
			"drama", "horror"},
		Medium: {"romance",
			"adventure", "memoir",
			"detective"},

		Hard: {"psychodrama", "tragicomedy",
			"surrealism",
			"steampunk", "cyberpunk"},
	},
}

func GetCategoryName(category int) string {
	return categoryNames[category]
}

func GetLevelDifficulty(level int) string {
	return difficultyLevel[level]
}

func GetValByRune(r rune) (int, bool) {
	if r < '0' || r > '9' {
		return 0, false
	}
	return int(r - '0'), true
}

func CategoryCount() int {
	return len(difficultyLevel)
}

func DifficultyCount() int {
	return len(categoryNames)
}

func GetHelpDescription(word string) string {
	return helpDescription[word]
}
