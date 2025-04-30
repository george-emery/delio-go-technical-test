# Stock Price Checker Technical Test

## The Task

Create a CLI application in Go that checks the current price of AAPL and MSFT shares and compares them to the previous day's closing price using the Finnhub API.

## Requirements

1. Create a CLI application that:

   - Fetches current prices for AAPL and MSFT from Finnhub API
   - Compares them with previous day's closing prices
   - Calculates P&L (Profit and Loss) for 10 shares of each stock
   - Displays results in the terminal

2. Technical Requirements:
   - Use Go 1.22 or later
   - Implement proper error handling
   - Write unit tests
   - Use the Finnhub API (https://finnhub.io/docs/api/quote)

## Getting Started

1. Fork this template repository
2. Set up your development environment:

   ```bash
   go mod init github.com/yourusername/go-technical-test
   ```

3. Get a Finnhub API key from https://finnhub.io

## Implementation Tips

- Consider using a CLI framework like Cobra or urfave/cli
- Use goroutines for concurrent API calls
- Implement proper error handling and logging
- Write tests for your implementation
- Consider adding benchmarks

## Bonus Points

- Implement concurrent API calls
- Add benchmarks
- Add configuration options (e.g., number of shares, stock symbols)
- Add proper logging
- Implement rate limiting for API calls
