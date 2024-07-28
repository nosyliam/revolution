package actions

import "github.com/nosyliam/revolution/pkg/common"

type logicAction struct {
	name  string
	logic func() error
}

func (a *logicAction) Execute(macro *common.Macro) error {
	return a.logic()
}

func Logic(name string) common.Action {
	return &logicAction{name: name}
}

type (
	conditionType int
	PredicateFunc func() bool
)

const (
	ifConditionType conditionType = iota
	elseifConditionType
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

func (a *conditionalAction) Execute(macro *common.Macro) error {
	var lastCond conditionType
	for _, cond := range a.conditions {
		if cond.Predicate() {
			for _, exec := range cond.Exec {
				if err := exec(macro); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func Condition(conds ...interface{}) common.Action {
	var activeCond *condition
	var conditions []*condition
	for _, cond := range conds {
		if fn, ok := cond.(Predicate); ok {
			if activeCond != nil {
				conditions = append(conditions, activeCond)
			}
			activeCond = &condition{}
			activeCond.Type = fn.Type()
			if activeCond.Type != elseConditionType {
				activeCond.Predicate = fn.Predicate()
			} else {
				activeCond.Predicate = func() bool { return true }
			}
		} else if fn, ok := cond.(common.Action); ok {
			activeCond.Exec = append(activeCond.Exec, func(macro *common.Macro) error { return fn.Execute(macro) })
		} else if fn, ok := cond.(func() error); ok {
			activeCond.Exec = append(activeCond.Exec, func(macro *common.Macro) error { return fn() })
		} else if fn, ok := cond.(func(macro *common.Macro) error); ok {
			activeCond.Exec = append(activeCond.Exec, fn)
		}
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

type elseifPredicate struct {
	predicate PredicateFunc
}

func (p *elseifPredicate) Type() conditionType      { return elseifConditionType }
func (p *elseifPredicate) Predicate() PredicateFunc { return p.predicate }

func Elseif(predicate PredicateFunc) Predicate {
	return &elseifPredicate{predicate}
}

type elsePredicate struct{}

func (p *elsePredicate) Type() conditionType      { return elseConditionType }
func (p *elsePredicate) Predicate() PredicateFunc { return nil }

func Else() Predicate {
	return &elsePredicate{}
}

func NotNil(obj interface{}) PredicateFunc {
	return func() bool { return obj != nil }
}

func Nil(obj interface{}) PredicateFunc {
	return func() bool { return obj == nil }
}

type retryAction struct{}

func (r *retryAction) Execute(*common.Macro) error { return common.RetrySignal }
func Retry() common.Action                         { return &retryAction{} }

type restartAction struct{}

func (r *restartAction) Execute(*common.Macro) error { return common.RestartSignal }
func Restart() common.Action                         { return &restartAction{} }

type terminateAction struct{}

func (r *terminateAction) Execute(*common.Macro) error { return common.TerminateSignal }
func Terminate() common.Action                         { return &terminateAction{} }
