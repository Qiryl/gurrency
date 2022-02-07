package fixer

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Qiryl/gurrency"
	"github.com/pkg/errors"
)

var (
	ErrBadResponse        = errors.New("bad response from the service")
	ErrInvalidResponse    = errors.New("can't unmarshall incoming JSON")
	ErrServiceUnavailable = errors.New("can't perform get request")
)

// fixerSource contains endpoint with specified currencies.
// fixerSource implements CurrencySource interface.
type fixerSource struct {
	urlRate string // Endpoint of the source, for getting the latest rates.
}

// fixerResponse represents the information about the currency returned from the fixer.io.
type fixerResponse struct {
	Base      string             `json:"base"`
	Reference map[string]float32 `json:"rates"`
}

// NewFixerSource is the constructor of the fixer source.
// The return value is the address of a fixerSource variable.
func NewFixerSource(key, url, base, reference string) *fixerSource {
	ur := fmt.Sprintf("%s/api/latest?access_key=%s&base=%s&symbols=%s", url, key, base, reference)
	return &fixerSource{urlRate: ur}
}

// GetRate returns the rates from the fixer.io.
// The return value is the address of a CurrencyRate variable.
func (fs fixerSource) GetRate() (*gurrency.CurrencyRate, error) {
	var rate fixerResponse

	r, err := http.Get(fs.urlRate)
	if err != nil {
		return nil, errors.Wrap(ErrServiceUnavailable, err.Error())
	}
	if r.StatusCode != http.StatusOK {
		return nil, ErrBadResponse
	}
	if err = json.NewDecoder(r.Body).Decode(&rate); err != nil {
		return nil, errors.Wrap(ErrInvalidResponse, err.Error())
	}
	return &gurrency.CurrencyRate{ServiceName: "fixer.io", Base: rate.Base, Reference: rate.Reference}, nil
}
