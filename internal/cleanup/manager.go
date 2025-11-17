package cleanup

type Manager struct {
	Tasks map[string]Task
}

type Task struct {
	Order   int
	Name    string
	Handler func() error
}

func NewManager() *Manager {
	return &Manager{
		Tasks: map[string]Task{},
	}
}

func (m *Manager) Add(name string, handler func() error) {
	t := &Task{
		Order:   len(m.Tasks),
		Name:    name,
		Handler: handler,
	}
	m.Tasks[name] = *t
}

func (m *Manager) Remove(name string) {
	delete(m.Tasks, name)
}

func (m *Manager) Clean() error {
	return nil
}
