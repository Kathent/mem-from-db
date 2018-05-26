package manager

import "github.com/orcaman/concurrent-map"

func NewManager(conf TableConfig) *Manager {
	return &Manager{
		m:    cmap.New(),
		conf: conf,
	}
}

func (m *Manager) init() {

}
