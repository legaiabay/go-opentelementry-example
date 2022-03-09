package main

import (
	"fmt"
	otel "test-otel/pkg/opentelementry"
	"time"
)

func main() {
	otel.Tracer = otel.Init("info payroll")
	defer otel.Tracer.End()

	foo()
	bar()
	pye()
}

func foo() {
	span := otel.Tracer.Trace("get payroll income", "")
	defer span.End()

	time.Sleep(1 * time.Second)

	fmt.Println("trace foo")
}

func bar() {
	span := otel.Tracer.Trace("get payroll deduction", "")
	defer span.End()

	time.Sleep(2 * time.Second)

	fmt.Println("trace bar")
}

func pye() {
	span := otel.Tracer.Trace("calculate payroll tax", "")
	defer span.End()

	time.Sleep(1 * time.Second)

	fmt.Println("trace pye")
}
