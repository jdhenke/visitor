package visitor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

var _ Visitor[int] = mpgVisitor{}

type mpgVisitor struct{}

func (mpgVisitor) VisitBus(b Bus) (int, error) {
	return b.MPG, nil
}

func (mpgVisitor) VisitCar(c Car) (int, error) {
	return c.MPG, nil
}

func (mpgVisitor) VisitBike(Bike) (int, error) {
	return math.MaxInt, nil
}

func TestVisitor(t *testing.T) {
	mpgVisitor := mpgVisitor{}
	for i, tc := range []struct {
		t       *Transport
		wantMPG int
		wantID  string
	}{
		{
			t: TransportFromBus(Bus{
				Number: 123,
				MPG:    50,
			}),
			wantMPG: 50,
			wantID:  "123",
		},
		{
			t: TransportFromCar(Car{
				LicensePlate: "CO-AYE-YOO",
				MPG:          30,
			}),
			wantMPG: 30,
			wantID:  "CO-AYE-YOO",
		},
		{
			t: TransportFromBike(Bike{
				Barcode: "ABC123",
			}),
			wantMPG: math.MaxInt,
			wantID:  "ABC123",
		},
	} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			{
				// Using the non-generic accept, this does force at least _handling_ each case at compile time, but
				// there's no compile time check that each case properly sets mpg.
				var mpg int
				require.NoError(t, Accept(tc.t, func(b Bus) error {
					mpg = b.MPG
					return nil
				}, func(c Car) error {
					mpg = c.MPG
					return nil
				}, func(b Bike) error {
					mpg = math.MaxInt //
					return nil
				}))
				assert.Equal(t, tc.wantMPG, mpg)
			}
			{
				// Using the generic accept, this forces handling **and** returning something of the valid type for each
				// case at compile time. This is convenient for one of use cases, but would be annoying to have to
				// specify all these functions for a repeated one.
				id, err := AcceptGeneric(tc.t, func(b Bus) (string, error) {
					return fmt.Sprint(b.Number), nil
				}, func(c Car) (string, error) {
					return c.LicensePlate, nil
				}, func(b Bike) (string, error) {
					return b.Barcode, nil
				})
				require.NoError(t, err)
				assert.Equal(t, tc.wantID, id)
			}
			{
				// Using this requires defining a dedicated type which can be a little annoying for one off cases,
				// although the previous example can be annoying for repeated uses of the same visitor.
				mpg, err := AcceptVisitorType[int](tc.t, mpgVisitor)
				require.NoError(t, err)
				assert.Equal(t, tc.wantMPG, mpg)
			}

			// Perhaps a library could offer both approaches? You wouldn't need the one-off function version if you
			// could anonymously implement interfaces in Go like you can in Java. Perhaps you could code-gen, for each
			// interface you wanted, a type-checked constructor that forced all methods to be present as functions
			// that assigned them to an inner struct then used them in their implementation of that interface. It could
			// look like this.
			{
				id, err := AcceptVisitorType[string](tc.t, NewVisitor[string](
					func(b Bus) (string, error) {
						return fmt.Sprint(b.Number), nil
					}, func(c Car) (string, error) {
						return c.LicensePlate, nil
					}, func(b Bike) (string, error) {
						return b.Barcode, nil
					},
				))
				require.NoError(t, err)
				assert.Equal(t, tc.wantID, id)
			}
		})
	}
}
