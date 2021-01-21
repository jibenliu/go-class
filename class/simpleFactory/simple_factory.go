package simpleFactory

type Product interface {
	SetName(name string)
	GetName() string
}

type Product1 struct {
	name string
}

func (p1 *Product1) SetName(name string) {
	p1.name = name
}

func (p1 *Product1) GetName() string {
	return "产品1的name为" + p1.name
}

type Product2 struct {
	name string
}

func (p2 *Product2) SetName(name string) {
	p2.name = name
}

func (p2 *Product2) GetName() string {
	return "产品2的name为" + p2.name
}

type productType int

const (
	p1 productType = iota
	p2
)

type productFactory struct {

}

func (pf *productFactory) Create(productType productType) Product {
	switch productType {
	case p1:
		return &Product1{}
	case p2:
		return &Product2{}
	default:
		return nil
	}
}
