package domain

type Difficulty int

const (
	Easy Difficulty = iota
	Medium
	Hard
)

var difficultyLevel = [...]string{
	"Easy   <------ press 0",
	"Medium <------ press 1",
	"Hard   <------ press 2",
}
