package comparator

type Comparator interface {
	Compare(c Comparator) int
}

type intComparator int

func (ic intComparator) Compare(c Comparator) int {
	if val, ok := c.(intComparator); ok {
		return int(ic) - int(val)
	}

	return 1
}

func NewComparator(t string, val interface{}) Comparator {
	if t == "int" {
		realVal, _ := val.(int)
		return intComparator(realVal)
	}

	return nil
}

type KeyValueComparator struct {
	Keys []Comparator
}

func (kvc *KeyValueComparator) Compare(c Comparator) int {
	if val, ok := c.(*KeyValueComparator); ok {
		for idx, v := range kvc.Keys {
			compareVal := v.Compare(val.Keys[idx])
			if compareVal != 0 {
				return compareVal
			}
		}

		return 0
	}

	return 1
}
