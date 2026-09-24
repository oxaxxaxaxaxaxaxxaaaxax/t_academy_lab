package adapter

var pictures = []string{
	`  +---+
  |   |
  |   
  |   
  |   
  |   
==+=====`,
	`  +---+
  |   |
  |   O
  |   
  |   
  |   
==+=====`,
	`  +---+
  |   |
  |   O
  |   |
  |   
  |   
==+=====`,
	`  +---+
  |   |
  |   O
  |  /|
  |   
  |   
==+=====`,
	`  +---+
  |   |
  |   O
  |  /|\
  |   
  |   
==+=====`,
	`  +---+
  |   |
  |   O
  |  /|\
  |  / 
  |   
==+=====`,
	`  +---+
  |   |
  |   O
  |  /|\
  |  / \
  |   
==+=====`,
}

func GetHangmanPicture(idx int) string {

	return pictures[idx]
}
