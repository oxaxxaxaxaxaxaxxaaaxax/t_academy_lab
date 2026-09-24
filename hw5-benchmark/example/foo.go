package example

type Person struct {
	Entity
	Name string
	Age  int
}

func (p Person) GetName() string {
	return p.Name
}

func (p Person) SetName(name string) {
	p.Name = name
}
