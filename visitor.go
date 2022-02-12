package visitor

import "fmt"

type Bus struct {
	Number int
	MPG    int
}

type Car struct {
	LicensePlate string
	MPG          int
}
type Bike struct {
	Barcode string
}

type Transport struct {
	bus  *Bus
	car  *Car
	bike *Bike
}

func TransportFromBus(b Bus) *Transport {
	return &Transport{
		bus: &b,
	}
}

func TransportFromCar(c Car) *Transport {
	return &Transport{
		car: &c,
	}
}

func TransportFromBike(b Bike) *Transport {
	return &Transport{
		bike: &b,
	}
}

func Accept(
	t *Transport,
	visitBus func(b Bus) error,
	visitCar func(c Car) error,
	visitBike func(b Bike) error,
) error {
	if t.bus != nil {
		return visitBus(*t.bus)
	}
	if t.car != nil {
		return visitCar(*t.car)
	}
	if t.bike != nil {
		return visitBike(*t.bike)
	}
	return fmt.Errorf("no type of tranpsortation set")
}

func AcceptGeneric[T any](
	t *Transport,
	visitBus func(b Bus) (T, error),
	visitCar func(c Car) (T, error),
	visitBike func(b Bike) (T, error),
) (T, error) {
	if t.bus != nil {
		return visitBus(*t.bus)
	}
	if t.car != nil {
		return visitCar(*t.car)
	}
	if t.bike != nil {
		return visitBike(*t.bike)
	}
	var zero T
	return zero, fmt.Errorf("no type of tranpsortation set")
}

type Visitor[T any] interface {
	VisitBus(b Bus) (T, error)
	VisitCar(c Car) (T, error)
	VisitBike(b Bike) (T, error)
}

//func (t *Transport) AcceptVisitorType[T any](v Visitor[T]) (T, error) {
// ^INVALID
//
// ./visitor.go:15:27: syntax error: method must have no type parameters
//
// Interestingly, just making this a function, not a method, is valid.
//
// https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods

func AcceptVisitorType[T any](t *Transport, v Visitor[T]) (T, error) {
	if t.bus != nil {
		return v.VisitBus(*t.bus)
	}
	if t.car != nil {
		return v.VisitCar(*t.car)
	}
	if t.bike != nil {
		return v.VisitBike(*t.bike)
	}
	var zero T
	return zero, fmt.Errorf("not type of tranpsortation set")
}

// Ideally the following would all be auto-generated.

func NewVisitor[T any](
	visitBus func(b Bus) (T, error),
	visitCar func(c Car) (T, error),
	visitBike func(b Bike) (T, error),
) Visitor[T] {
	return &protoVisitor[T]{
		visitBus:  visitBus,
		visitCar:  visitCar,
		visitBike: visitBike,
	}
}

type protoVisitor[T any] struct {
	visitBus  func(b Bus) (T, error)
	visitCar  func(c Car) (T, error)
	visitBike func(b Bike) (T, error)
}

func (p *protoVisitor[T]) VisitBus(b Bus) (T, error) {
	return p.visitBus(b)
}

func (p *protoVisitor[T]) VisitCar(c Car) (T, error) {
	return p.visitCar(c)
}

func (p *protoVisitor[T]) VisitBike(b Bike) (T, error) {
	return p.visitBike(b)
}
