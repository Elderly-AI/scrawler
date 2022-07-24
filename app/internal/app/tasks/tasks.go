package tasks

type Facade interface {
}

type Implementation struct {
	facade Facade
}

func New(facade Facade) Implementation {
	return Implementation{
		facade: facade,
	}
}
