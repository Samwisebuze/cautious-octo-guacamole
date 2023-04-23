package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

/*
	TODO:
		- I'd like to add tests but I'm short on time. I've explored some edge cases like invalid tickers without a rate, malformed inputs
		- Present v1 and explore requirements with customer. do you need to support more tickers? different splits? This could be generalized to support that.
		- Revisit tech-debt with the New constructor methods. I think CryptoTicker is dependent on ExchangeRates, but doesn't need to be
*/

const (
	COINBASE_URL = "https://api.coinbase.com/v2/exchange-rates?currency=USD"
	OUTPUT_FMT   = "%s => %f %s\n"
)

type USD int64

func toUSD(value float64) USD {
	return USD(math.Round((value * 100) + 0.05))
}

// Note: stolen from SO for a more type safe solution
// https://stackoverflow.com/questions/20596428/how-to-represent-currency-in-go
// Float64 converts a USD to float64
func (m USD) Float64() float64 {
	x := float64(m)
	x = x / 100
	return x
}

// Multiply safely multiplies a USD value by a float64, rounding
// to the nearest cent.
func (m USD) Multiply(f float64) USD {
	x := (float64(m) * f) + 0.5
	return USD(x)
}

// String returns a formatted USD value
func (m USD) String() string {
	x := float64(m)
	x = x / 100
	return fmt.Sprintf("$%.2f", x)
}

type CryptoTicker struct {
	symbol  string
	usdRate float64
}

func NewCryptoTicker(symbol string, exchangeRates ExchangeRates) *CryptoTicker {
	rate, err := exchangeRates.getRateFor(symbol)
	if err != nil {
		log.Fatal(err)
	}

	return &CryptoTicker{symbol, rate}
}

func (cryptoTicker *CryptoTicker) ExchangeToUSD(value USD) float64 {
	return value.Float64() * cryptoTicker.usdRate
}

type USDToCryptoSplitUseCase struct {
	usd       USD
	Primary   *CryptoTicker
	Secondary *CryptoTicker
}

func NewUSDToCryptoSplitUseCase(amount string, pTicker *CryptoTicker, sTicker *CryptoTicker) *USDToCryptoSplitUseCase {
	cleanAmount := strings.Trim(amount, "$")
	fAmount, err := strconv.ParseFloat(cleanAmount, 64)
	if err != nil {
		panic(fmt.Sprintf("%s is not a valid USD amount.", cleanAmount))
	}

	usd := toUSD(fAmount)

	return &USDToCryptoSplitUseCase{usd, pTicker, sTicker}
}

func (useCase *USDToCryptoSplitUseCase) SplitValue(splitPercentage float64) (amount USD, remainder USD) {
	return useCase.usd.Multiply(splitPercentage), useCase.usd.Multiply(1.0 - splitPercentage)
}

func ParseCmdLineArgs(args []string) (usdAmount string, primaryTicker string, secondaryTicker string) {
	usdAmount = args[0]
	primaryTicker = args[1]
	secondaryTicker = args[2]
	return
}

// .prog amount Ticker1 Ticker2
func main() {
	rates, err := getExchangeRates()
	if err != nil {
		log.Fatal(err)
	}

	vargs := os.Args[1:]
	usdAmount, primaySymbol, secondarySymbol := ParseCmdLineArgs(vargs)
	useCase := NewUSDToCryptoSplitUseCase(usdAmount, NewCryptoTicker(primaySymbol, rates), NewCryptoTicker(secondarySymbol, rates))

	primaryAmountSplit, secondaryAmountSplit := useCase.SplitValue(0.7)

	primaryTickerAmount := useCase.Primary.ExchangeToUSD(primaryAmountSplit)
	secondaryTickerAmount := useCase.Secondary.ExchangeToUSD(secondaryAmountSplit)

	fmt.Printf(OUTPUT_FMT, primaryAmountSplit.String(), primaryTickerAmount, useCase.Primary.symbol)
	fmt.Printf(OUTPUT_FMT, secondaryAmountSplit.String(), secondaryTickerAmount, useCase.Secondary.symbol)

	os.Exit(0)
}

type CoinbaseApiResponse struct {
	Data ExchangeRates
}

type ExchangeRates struct {
	Currency string
	Rates    map[string]string
}

func (exchange *ExchangeRates) getRateFor(ticker string) (float64, error) {
	if rate, exists := exchange.Rates[ticker]; exists {
		return strconv.ParseFloat(rate, 64)
	}

	err := fmt.Errorf("we cannont find an Exchange Rate for %s, please try another crypto ticker", ticker)
	return 0, err
}

func getExchangeRates() (exchangeRates ExchangeRates, err error) {
	res, err := http.Get(COINBASE_URL)
	if err != nil {
		return *new(ExchangeRates), err
	}

	var apiRes CoinbaseApiResponse
	defer res.Body.Close()
	bRes, err := io.ReadAll(res.Body)
	if err != nil {
		return *new(ExchangeRates), err
	}

	err = json.Unmarshal(bRes, &apiRes)
	return apiRes.Data, err
}
