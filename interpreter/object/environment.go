package object

func NewEnvironment() *Environment {
	store := make(map[string]Object)
	return &Environment{store: store}
}

type Environment struct {
	store map[string]Object
}

func (env *Environment) Get(name string) (Object, bool) {
	obj, ok := env.store[name]
	return obj, ok
}

func (env *Environment) Set(name string, val Object) Object {
	env.store[name] = val
	return val
}
