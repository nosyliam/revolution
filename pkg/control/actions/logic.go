package actions

type logicAction struct {
	name string
}

func (a *subroutineAction) Execute(deps *Dependencies) error {
	return deps.Exec(a.name)
}

func Subroutine(name string) Action {
	return &subroutineAction{name: name}
}
