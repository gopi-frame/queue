package driver

type Driver interface {
	Open(options map[string]any) (Queue, error)
}
