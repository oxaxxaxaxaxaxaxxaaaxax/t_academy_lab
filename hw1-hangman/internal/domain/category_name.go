package domain

type Category int

const (
	Animals Category = iota
	Space
	LiteratureGenre
)

var categoryNames = [...]string{
	"Animals         <------ press 0",
	"Space           <------ press 1",
	"LiteratureGenre <------ press 2",
}
