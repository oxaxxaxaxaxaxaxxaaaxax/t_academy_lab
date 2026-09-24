package benchmark

import (
	"reflect"
	"testing"
)

type Student struct{ n string }

func (s Student) Name() string { return s.n }

// Вариант для интерфейсного вызова
type Namer interface {
	Name() string
}

func BenchmarkStudentName_Direct(b *testing.B) {
	s := Student{n: "alice"}
	var sink string

	b.ResetTimer()

	for b.Loop() {
		sink = s.Name()
	}

	_ = sink
}

func BenchmarkStudentName_Interface(b *testing.B) {
	s := Student{n: "alice"}
	var i Namer = s
	var sink string

	b.ResetTimer()

	for b.Loop() {
		sink = i.Name()
	}

	_ = sink
}

func BenchmarkStudentName_MethodValue(b *testing.B) {
	s := Student{n: "alice"}
	f := s.Name
	var sink string

	b.ResetTimer()

	for b.Loop() {
		sink = f()
	}

	_ = sink
}

func BenchmarkStudentName_MethodExpression(b *testing.B) {
	s := Student{n: "alice"}
	f := Student.Name
	var sink string

	b.ResetTimer()

	for b.Loop() {
		sink = f(s)
	}

	_ = sink
}

func BenchmarkStudentName_Reflect(b *testing.B) {
	s := Student{n: "alice"}
	var i Namer = s
	v := reflect.ValueOf(i)
	method := v.MethodByName("Name")
	var sink string

	b.ResetTimer()

	for b.Loop() {
		sink = method.Call(nil)[0].String()
	}

	_ = sink
}
