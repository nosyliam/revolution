package actions

import (
	"github.com/nosyliam/revolution/pkg/common"
)

type logicAction struct {
	exec func(*common.Macro) error
}

func (a *logicAction) Execute(macro *common.Macro) error {
	return a.exec(macro)
}

func Logic(logic interface{}) common.Action {
	switch fn := logic.(type) {
	case func():
		return &logicAction{func(*common.Macro) error { fn(); return nil }}
	case func() error:
		return &logicAction{func(*common.Macro) error { return fn() }}
	case func(macro *common.Macro) error:
		return &logicAction{fn}
	default:
		panic("unknown logic type")
	}
}

type (
	conditionType int
	loopType      int
	PredicateFunc func(macro *common.Macro) bool
)

const (
	ifConditionType conditionType = iota
	elseConditionType
)

type condition struct {
	Predicate PredicateFunc
	Exec      []func(macro *common.Macro) error
	Type      conditionType
}

type conditionalAction struct {
	conditions []*condition
}

type Predicate interface {
	Predicate() PredicateFunc
	Type() conditionType
}

type Loop interface {
	Predicate() PredicateFunc
	Type() loopType
}

func (a *conditionalAction) Execute(macro *common.Macro) error {
	for _, cond := range a.conditions {
		if cond.Predicate(macro) {
			for _, exec := range cond.Exec {
				if err := exec(macro); err != nil {
					return err
				}
			}
			return nil
		}
	}
	return nil
}

func Condition(conds ...interface{}) common.Action {
	var activeCond *condition
	var conditions []*condition
	for _, cond := range conds {
		switch fn := cond.(type) {
		case Predicate:
			if activeCond != nil {
				conditions = append(conditions, activeCond)
			}
			activeCond = &condition{}
			activeCond.Type = fn.Type()
			if activeCond.Type != elseConditionType {
				activeCond.Predicate = fn.Predicate()
			} else {
				activeCond.Predicate = func(*common.Macro) bool { return true }
			}
		case common.Action:
			activeCond.Exec = append(activeCond.Exec, func(macro *common.Macro) error { return fn.Execute(macro) })
		case func() error:
			activeCond.Exec = append(activeCond.Exec, func(macro *common.Macro) error { return fn() })
		case func(macro *common.Macro) error:
			activeCond.Exec = append(activeCond.Exec, fn)
		default:
			panic("unknown predicate type")
		}
	}
	if activeCond != nil {
		conditions = append(conditions, activeCond)
	}
	return &conditionalAction{conditions}
}

type ifPredicate struct {
	predicate PredicateFunc
}

func (p *ifPredicate) Type() conditionType      { return ifConditionType }
func (p *ifPredicate) Predicate() PredicateFunc { return p.predicate }

func If(predicate PredicateFunc) Predicate {
	return &ifPredicate{predicate}
}

type elsePredicate struct{}

func (p *elsePredicate) Type() conditionType      { return elseConditionType }
func (p *elsePredicate) Predicate() PredicateFunc { return nil }

func Else() Predicate {
	return &elsePredicate{}
}

func And(fns ...PredicateFunc) PredicateFunc {
	return func(macro *common.Macro) bool {
		for _, fn := range fns {
			if !fn(macro) {
				return false
			}
		}
		return true
	}
}

func Or(fns ...PredicateFunc) PredicateFunc {
	return func(macro *common.Macro) bool {
		for _, fn := range fns {
			if fn(macro) {
				return true
			}
		}
		return false
	}
}

type stepBackAction struct{}

func (r *stepBackAction) Execute(*common.Macro) error { return common.StepBackSignal }
func StepBack() common.Action                         { return &stepBackAction{} }

type restartAction struct{}

func (r *restartAction) Execute(*common.Macro) error { return common.RestartSignal }
func Restart() common.Action                         { return &restartAction{} }

type terminateAction struct{}

func (r *terminateAction) Execute(*common.Macro) error { return common.TerminateSignal }
func Terminate() common.Action                         { return &terminateAction{} }
