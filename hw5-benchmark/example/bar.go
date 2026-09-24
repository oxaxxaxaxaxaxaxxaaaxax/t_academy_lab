package example

type Employee struct {
	Person
	Position string
	Salary   int
}

func (e Employee) GetSalary() int {
	return e.Salary
}

func (e Employee) Promote(newPosition string, bonus int) {
	e.Position = newPosition
	e.Salary += bonus
}
