package main

import (
	"fmt"
	otel "test-otel/pkg/opentelementry"
)

func main() {
	otel.Tracer = otel.Init("info payroll")

	otel.Tracer.Trace("main", "")
	defer otel.Tracer.End()

	foo()
}

func foo() {
	otel.Tracer.Trace("get payroll income", "")

	fmt.Println("trace foo")

	bar()
}

func bar() {
	otel.Tracer.Trace("get payroll deduction", "")

	fmt.Println("trace bar")

	pye()
}

func pye() {
	otel.Tracer.Trace("calculate payroll tax", "")

	fmt.Println("trace pye")
}
