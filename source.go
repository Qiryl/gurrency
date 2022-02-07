package gurrency

import (
	"fmt"
	"log"
	"sync"
)

// CurrencySource is an interface that must be implemented by a concrete source.
type CurrencySource interface {
	// GetRate returns the currency rate of the concrete source.
	// The return value is the address of the CurrencyRate variable.
	GetRate() (*CurrencyRate, error)
}

// CurrencyRate represents the information about the currency rate.
type CurrencyRate struct {
	// The name of the service that provides the information about the currency.
	ServiceName string

	// Base currency in relation to which rates are displayed.
	Base string

	// Currencies for which you need to get the ratio.
	Reference map[string]float32
}

// Generalized source that stores concrete source.
type Source struct {
	src CurrencySource
}

// NewSource is the constructor of the generalized source.
// The return value is the address of a Source variable that stores the value of the concrete source.
func NewSource(s CurrencySource) *Source {
	return &Source{src: s}
}

// GetRate prints the currency rates of the generalized source.
// The wg argument points to the sync.WaitGroup variable that is needed to synchronize goroutines.
func (s Source) GetRate(wg *sync.WaitGroup) {
	rates, err := s.src.GetRate()
	if err != nil {
		log.Fatalln(err)
	}
	out := fmt.Sprintf("Service Name: %s\n", rates.ServiceName)
	for symbol, rate := range rates.Reference {
		out += fmt.Sprintf("%s/%s: %.2f\n", rates.Base, symbol, rate)
	}
	fmt.Print(out)
	wg.Done()
}
