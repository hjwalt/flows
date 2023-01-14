package runtime

type ContextKey struct {
	Name string
}

func Context(name string) ContextKey {
	return ContextKey{Name: name}
}
