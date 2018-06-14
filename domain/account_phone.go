package domain

type AccountPhone struct {
	_     int `db:"robot,account_phone"`
	Id    int `id:"id"`
	VccId string
	B     string `index:"idx_a_b,1"`
}
