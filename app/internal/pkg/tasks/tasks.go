package tasks

type tasksDB interface {
}

type Facade struct {
	tasksDB tasksDB
}

func New(crawlerDB tasksDB) Facade {
	return Facade{
		tasksDB: crawlerDB,
	}
}
