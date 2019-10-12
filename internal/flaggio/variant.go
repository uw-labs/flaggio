package flaggio

type Variant struct {
	ID             string
	Key            string
	Description    *string
	Value          interface{}
	DefaultWhenOn  bool
	DefaultWhenOff bool
}

type VariantList []*Variant

func (l VariantList) DefaultWhenOn() *Variant {
	for _, v := range l {
		if v.DefaultWhenOn {
			return v
		}
	}
	return nil
}

func (l VariantList) DefaultWhenOff() *Variant {
	for _, v := range l {
		if v.DefaultWhenOff {
			return v
		}
	}
	return nil
}
