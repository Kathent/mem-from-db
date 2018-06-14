package domain

type Sample struct {
	_  int    `db:"robot,sample"`
	Id int    `id:"id"`
	A  int    `index:"idx_a_b,0"`
	B  string `index:"idx_a_b,1"`
}
