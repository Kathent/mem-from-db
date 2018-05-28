package manager

type indexReserve struct {
	tree *IndexTree
}


type indexReserveBuilder struct {
	res []interface{}
	i indexInfo
	c []columnInfo
}

func NewIndexReserveBuilder() *indexReserveBuilder {
	return &indexReserveBuilder{}
}

func (b *indexReserveBuilder) withRes(res []interface{}) *indexReserveBuilder {
	b.res = res
	return b
}

func (b *indexReserveBuilder) withIndexInfo(i indexInfo) *indexReserveBuilder {
	b.i = i
	return b
}

func (b *indexReserveBuilder) withColumnInfo(c []columnInfo) *indexReserveBuilder {
	b.c = c
	return b
}

func (b *indexReserveBuilder) build() *indexReserve {
	return &indexReserve{}
}