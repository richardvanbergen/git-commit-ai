package adapters

type QueryableLLM interface {
	Summerize(changes string) string
}
