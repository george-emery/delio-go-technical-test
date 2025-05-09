package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

var rootCmd = &cobra.Command{
	Use:   "stock-check",
	Short: "Compare stocks from last close to now",
	Long:  `Compare stocks from last close to now using the Finnhub api`,
	RunE:  fetchAndPrintData,
}

type SymbolResponse struct {
	Symbol       string
	CurrentPrice float32
	ClosePrice   float32
	Error        error
}

// initialiseFinnhubClient Initialise the Finnhub Client using the provided api token
func initialiseFinnhubClient(apiToken string) (*finnhub.DefaultApiService, error) {
	log.Println("Setting up finnhub client") // Would set up a better logging system with configurable log levels
	// if I had more time
	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", apiToken)
	return finnhub.NewAPIClient(cfg).DefaultApi, nil
}

// Execute Run the rootCmd
func Execute() {
	rootCmd.Flags().Int64VarP(&Shares, "shares", "s", 10, "Number of shares to compare")
	rootCmd.Flags().BoolVarP(&LongOutput, "longOutput", "l", false, "Full calculation in response")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Error executing command", err)
	}
}

// fetchAndPrintData Pulls data from the finnhub api, then sends it to be printed to stdOut
func fetchAndPrintData(cmd *cobra.Command, args []string) (err error) {
	log.Println("Beginning fet and print command")
	if len(args) == 0 {
		log.Println("Using default args")
		args = []string{"AAPL", "MSFT"} // Decided to use args for the list of symbols, though I did consider using a
		// csv list as one of the flags
	}

	finnhubToken := os.Getenv("FINNHUB_TOKEN") // Probably should move this initialisation elsewhere if program gets bigger
	finnhubClient, err := initialiseFinnhubClient(finnhubToken)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	output := make(chan SymbolResponse, len(args))
	for _, symbol := range args {
		wg.Add(1)
		go fetchData(symbol, output, finnhubGetQuote{finnhubClient: finnhubClient}, &wg)
	}
	wg.Wait()
	close(output)
	for response := range output {
		err = printData(os.Stdout, response, Shares, LongOutput)
		if err != nil {
			return
		}
	}
	return
}

// fetchData Pulls data from the finnhub api, then puts the relevant data in a SymbolResponse and pushes that to the output channel
func fetchData(symbol string, output chan SymbolResponse, finnhubQuoteGetter finnhubGetQuoteInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Beginning fetch data for symbol %s", symbol)

	quote, _, err := finnhubQuoteGetter.GetQuote(symbol)
	if err != nil {
		output <- SymbolResponse{Symbol: symbol, Error: err}
		return
	}

	if *quote.C == 0 {
		output <- SymbolResponse{Symbol: symbol, Error: errors.New(fmt.Sprintf("%s: Doesn't exist\n", symbol))}
		return
		// I'm assuming that if a stock price is 0 then the stock doesn't exist, as the response
		// gives 200 no matter what and I couldn't find an api route to actually validate a specific stock symbol.
		// Depending on how often this code is rebooted you could pull a list of all symbols in the initialisation
		// and check against that, (also depends on how often new stocks are added?)
	}

	output <- SymbolResponse{
		Symbol:       symbol,
		CurrentPrice: *quote.C,
		ClosePrice:   *quote.Pc,
		Error:        nil,
	}
	return
}

// printData Prints the data in the response object to the writer object
func printData(writer io.Writer, response SymbolResponse, shares int64, longOutput bool) (err error) {
	log.Printf("Beginning print data for symbol %s", response.Symbol)
	if response.Error != nil {
		log.Println(fmt.Sprintf("Error while fetching finnhub data for symbol %s", response.Symbol), response.Error)
		_, err = fmt.Fprintf(writer, "%s:Error\n", response.Symbol)
		return
	}

	profitLoss := float32(shares) * (response.CurrentPrice - response.ClosePrice)

	if longOutput {
		_, err = fmt.Fprintf(
			writer,
			"%s: It's worth $%f now, was worth $%f last close, with %d shares you made $%f since last close\n",
			// Unsure if this should be 'made $-5' or 'lost $5' when in the negatives. I imagine it would depend on
			// who the intended user is
			response.Symbol,
			response.CurrentPrice,
			response.ClosePrice,
			shares,
			profitLoss,
		)
		return
	} else {
		_, err = fmt.Fprintf(writer, "%s: %f\n", response.Symbol, profitLoss)
		return
	}
}

// finnhubGetQuoteInterface interface for the finnhub api to allow mocking
type finnhubGetQuoteInterface interface {
	GetQuote(symbol string) (finnhub.Quote, *http.Response, error)
}

type finnhubGetQuote struct {
	finnhubClient *finnhub.DefaultApiService
}

func (f finnhubGetQuote) GetQuote(symbol string) (finnhub.Quote, *http.Response, error) {
	return f.finnhubClient.Quote(context.Background()).Symbol(symbol).Execute()
}
