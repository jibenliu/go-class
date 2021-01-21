package main

import (
	"fmt"
	"sync"
)

type Subject interface {
	Register(o Observer)
	Deregister(o Observer) error
	NotifyObserver()
}

type Observer interface {
	Update(name, status string)
	GetID() int
}

const (
	TimeIsUp  = "time is up"
	IsEnd     = "is end"
	NotAtTime = "not at time"
)

type shirt struct {
	sync.Mutex
	customers []Observer
	status    string
	name      string
}

func NewShirt() *shirt {
	return &shirt{status: NotAtTime, name: "shirt"}
}

func (s *shirt) Register(o Observer) {
	s.Lock()
	defer s.Unlock()
	s.customers = append(s.customers, o)
	fmt.Printf("[%s] registered a new customer with ID[%d]\n", s.name, o.GetID())
}

func (s *shirt) Deregister(o Observer) error {
	s.Lock()
	s.Unlock()
	var index int
	var found bool
	id := o.GetID()
	for i, c := range s.customers {
		if c.GetID() == id {
			index = i
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("Customer %d not found\n", id)
	}
	s.customers = append(s.customers[:index], s.customers[index+1:]...)
	fmt.Printf("Removed the customer with ID[%d]\n", id)
	return nil
}

func (s *shirt) NotifyObservers() {
	s.Lock()
	defer s.Unlock()
	wg := sync.WaitGroup{}
	for _, c := range s.customers {
		wg.Add(1)
		go func(c Observer) {
			defer wg.Done()
			c.Update(s.name, s.status)
		}(c)
	}
	wg.Wait()
	fmt.Println("Finished notify customers")
}

type customer struct {
	ID             int
	wantItemStatus string
}

func NewCustomers(id int) *customer {
	return &customer{ID: id}
}

func (c *customer) Update(name, status string) {
	c.wantItemStatus = status
	fmt.Printf("Update: hi customer %d, the item[%s] you want is [%v] now\n", c.ID, name, c.wantItemStatus)
}

func (c *customer) GetID() int {
	return c.ID
}

func main()  {
	c1 := NewCustomers(1)
	c2 := NewCustomers(2)
	c3 := NewCustomers(3)
	c4 := NewCustomers(4)

	s := NewShirt()
	s.Register(c1)
	s.Register(c2)
	s.Register(c3)
	s.Register(c4)

	s.NotifyObservers()

	s.status = TimeIsUp
	s.NotifyObservers()
}