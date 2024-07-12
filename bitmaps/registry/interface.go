package registry

type Registry interface {
	RegisterPng(name string, data []byte)
}
