package cmd

import (
	"bytes"
	"errors"
	"github.com/Finnhub-Stock-API/finnhub-go/v2"
	"net/http"
	"sync"
	"testing"
)

func TestPrintData(t *testing.T) {
	input := SymbolResponse{
		Symbol:       "TEST",
		CurrentPrice: 69.69,
		ClosePrice:   420.420,
	}
	shares := int64(10)
	longOutput := false
	var output bytes.Buffer

	err := printData(&output, input, shares, longOutput)

	if output.String() != "TEST: -3507.300049\n" || err != nil { // Cursed floating point numbers.
		// I'd work out a way to make this make more sense if this was anything fancier
		t.Errorf(`output to stdOut: %s, expected: TEST: -3507.300049`, output.String())
	}
}

func TestPrintDataLongOutput(t *testing.T) {
	input := SymbolResponse{
		Symbol:       "TEST",
		CurrentPrice: 69.69,
		ClosePrice:   420.420,
	}
	shares := int64(10)
	longOutput := true
	var output bytes.Buffer

	err := printData(&output, input, shares, longOutput)

	if output.String() != "TEST: It's worth $69.690002 now, was worth $420.420013 last close, with 10 shares you made $-3507.300049 since last close\n" || err != nil { // Cursed floating point numbers.
		// I'd work out a way to make this make more sense if this was anything fancier
		t.Errorf(`output to stdOut: %s, expected: TEST: It's worth $69.690002 now, was worth $420.420013 last close, with 10 shares you made $-3507.300049 since last close`, output.String())
	}
}

func TestFetchData(t *testing.T) {
	symbol := "TEST"
	output := make(chan SymbolResponse, 1)
	wg := sync.WaitGroup{}
	finnhubQuoter := mockFinnhubGetQuote{}

	wg.Add(1)
	go fetchData(symbol, output, finnhubQuoter, &wg)
	wg.Wait()
	close(output)

	if len(output) != 1 {
		t.Errorf(`output len: %d, expected: 1`, len(output))
	}

	expectedOutput := SymbolResponse{
		Symbol:       "TEST",
		CurrentPrice: float32(123.123),
		ClosePrice:   float32(456.456),
	}
	for response := range output {
		if response != expectedOutput {
			t.Errorf(`output: %v, expected: %v`, output, expectedOutput)
		}
	}
}

func TestFetchDataError(t *testing.T) {
	symbol := "TEST"
	output := make(chan SymbolResponse, 1)
	wg := sync.WaitGroup{}
	finnhubQuoter := mockFinnhubGetErrorQuote{}

	wg.Add(1)
	go fetchData(symbol, output, finnhubQuoter, &wg)
	wg.Wait()
	close(output)

	if len(output) != 1 {
		t.Errorf(`output len: %d, expected: 1`, len(output))
	}

	expectedOutput := SymbolResponse{
		Symbol: "TEST",
		Error:  errors.New("401 Unauthorised"),
	}
	for response := range output {
		if response.Error.Error() != expectedOutput.Error.Error() || response.Symbol != expectedOutput.Symbol {
			t.Errorf(`output: %v, expected: %v`, response, expectedOutput)
		}
	}
}

type mockFinnhubGetQuote struct{}

func (f mockFinnhubGetQuote) GetQuote(symbol string) (finnhub.Quote, *http.Response, error) {
	C := float32(123.123)
	Pc := float32(456.456)
	return finnhub.Quote{
		C:  &C,
		Pc: &Pc,
	}, &http.Response{}, nil
}

type mockFinnhubGetErrorQuote struct{}

func (f mockFinnhubGetErrorQuote) GetQuote(symbol string) (finnhub.Quote, *http.Response, error) {
	return finnhub.Quote{}, &http.Response{}, errors.New("401 Unauthorised")
}
