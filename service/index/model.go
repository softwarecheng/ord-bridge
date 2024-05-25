package index

import "github.com/softwarecheng/ord-bridge/indexer"

type Model struct {
	indexer *indexer.Indexer
}

func NewModel(indexer *indexer.Indexer) *Model {
	return &Model{
		indexer: indexer,
	}
}
