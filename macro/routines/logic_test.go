package routines

import (
	. "github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/control"
	. "github.com/nosyliam/revolution/pkg/control/actions"
	"github.com/stretchr/testify/assert"
	"testing"
)

func exec(macro *Macro, actions []Action) *Macro {
	stop := make(chan struct{})
	pause := make(chan (<-chan struct{}))
	err := make(chan string)
	status := make(chan string)
	control.ExecuteRoutine(macro, actions, stop, pause, status, err)
	return macro
}

func Increment(macro *Macro) error {
	macro.Results.RetryCount++
	return nil
}

func Decrement(macro *Macro) error {
	macro.Results.RetryCount--
	return nil
}

func Test_Conditionals(t *testing.T) {
	macro := &Macro{Results: &ActionResults{}}
	exec(macro, []Action{
		Condition(
			If(True),
			Increment,
		),
		Terminate(),
	})
	assert.Equal(t, 1, macro.Results.RetryCount)
	exec(macro, []Action{
		Condition(
			If(False),
			Increment,
			Else(),
			Decrement,
		),
		Terminate(),
	})
	assert.Equal(t, 0, macro.Results.RetryCount)
	exec(macro, []Action{
		Logic(Increment),
		Condition(
			If(LessThan(RetryCount, 10)),
			StepBack(),
		),
		Condition(
			If(Equal(RetryCount, 10)),
			Subroutine(
				Logic(Increment),
				Condition(
					If(LessThan(RetryCount, 10)),
					StepBack(),
				),
				Terminate(),
			),
		),
		Terminate(),
	})
	assert.Equal(t, 10, macro.Results.RetryCount)
}

func Test_Assertions(t *testing.T) {
	macro := &Macro{Results: &ActionResults{}}
	exec(macro, []Action{
		Condition(
			If(Equal(RetryCount, 0)),
			Increment,
		),
		Condition(
			If(NotEqual(RetryCount, 1)),
			Increment,
		),
		Terminate(),
	})
	assert.Equal(t, 1, macro.Results.RetryCount)
	exec(macro, []Action{
		Condition(
			If(GreaterThan(RetryCount, 0)),
			Increment,
		),
		Condition(
			If(LessThan(RetryCount, 3)),
			Increment,
		),
		Condition(
			If(GreaterThanEq(RetryCount, 3)),
			Increment,
		),
		Condition(
			If(LessThanEq(RetryCount, 4)),
			Increment,
		),
		Terminate(),
	})
	assert.Equal(t, 5, macro.Results.RetryCount)
	exec(macro, []Action{
		Condition(
			If(And(GreaterThan(RetryCount, 3), LessThanEq(RetryCount, 5))),
			Increment,
		),
		Condition(
			If(Or(And(Equal(RetryCount, 6), LessThanEq(RetryCount, 6)), LessThanEq(RetryCount, 1))),
			Increment,
		),
		Terminate(),
	})
	assert.Equal(t, 7, macro.Results.RetryCount)
}
