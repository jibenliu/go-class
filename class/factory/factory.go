package factory

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

type ProductFactory interface {
	Create() Product
}

type ProductFactory1 struct {
}

func (pf1 *ProductFactory1) Create() Product {
	return &Product1{}
}

type ProductFactory2 struct {
}

func (pf2 *ProductFactory2) Create() Product {
	return &Product2{}
}
