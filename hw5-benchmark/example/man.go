package example

type Manager struct {
	Employee
	Department string
}

func (m Manager) GetDepartment() string {
	return m.Department
}

func (m Manager) ChangeDepartment(dep string) {
	m.Department = dep
}
