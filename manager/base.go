package manager

import "github.com/orcaman/concurrent-map"

type Manager struct {
	m    cmap.ConcurrentMap
	conf TableConfig
}

type TableConfig struct {
	Name string
}
