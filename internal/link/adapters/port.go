package adapters

type MappingPort[S any, T any] interface {
	Map(in S) (T, error)
	Reverse(in T) (S, error)
}
