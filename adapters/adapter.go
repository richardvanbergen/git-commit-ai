package adapters

type QueryableLLM interface {
	Summerize(changes string, query string) string
}
