package models

import (
	"math"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Diversificationable interface {
	GetValue() float32
	GetAssets() []Diversification
}

type Record struct {
	Id                    primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	UserId                primitive.ObjectID `json:"userId,omitempty" bson:"userId"`
	Date                  time.Time          `json:"date" bson:"date" binding:"required"`
	InvestedAmount        *float32           `json:"investedAmount,omitempty" bson:"investedAmount"`
	Stocks                Stocks             `json:"stocks" bson:"stocks" binding:"required,dive"`
	Cryptos               Cryptos            `json:"cryptos" bson:"cryptos" binding:"required"`
	StockDiversification  *[]Diversification `json:"stockDiversification,omitempty" `
	CryptoDiversification *[]Diversification `json:"cryptoDiversification,omitempty"`
	CryptosValue          *float32           `json:"cryptosValue,omitempty" bson:"cryptosValue"`
	StocksValue           *float32           `json:"stocksValue,omitempty" bson:"stocksValue"`
	Liquidity             *float32           `json:"liquidity" bson:"liquidity" binding:"required,min=0"`
	TotalInvested         float32            `json:"totalInvested,omitempty" bson:"totalInvested"`
}

type Diversification struct {
	Symbol  string  `json:"symbol"`
	Percent float32 `json:"percent"`
}

type Stock struct {
	Symbol   string   `json:"symbol" bson:"symbol" binding:"required"`
	Shares   float32  `json:"shares,omitempty" bson:"shares"`
	ValuedAt *float32 `json:"valuedAt" bson:"valuedAt" binding:"required,min=0"`
}

type Crypto struct {
	Symbol   string   `json:"symbol" bson:"symbol" binding:"required"`
	Coins    float32  `json:"coins,omitempty" bson:"coins"`
	ValuedAt *float32 `json:"valuedAt" bson:"valuedAt" binding:"required,min=0"`
}

type Stocks []Stock
type Cryptos []Crypto

type DeleteRecordBody struct {
	Id string `json:"id" binding:"required"`
}

type RecordBody struct {
	Id                    string            `json:"id,omitempty"`
	Date                  time.Time         `json:"date" binding:"required"`
	InvestedAmount        float32           `json:"investedAmount,omitempty"`
	Stocks                Stocks            `json:"stocks" binding:"required,dive"`
	Cryptos               Cryptos           `json:"cryptos" binding:"required"`
	StockDiversification  []Diversification `json:"stockDiversification" `
	CryptoDiversification []Diversification `json:"cryptoDiversification"`
	CryptosValue          float32           `json:"cryptosValue,omitempty"`
	StocksValue           float32           `json:"stocksValue,omitempty"`
	Liquidity             *float32          `json:"liquidity" binding:"required,min=0"`
	TotalInvested         float32           `json:"totalInvested,omitempty"`
}

func (stocks Stocks) GetValue() (sum float32) {
	for _, stock := range stocks {
		sum += *stock.ValuedAt
	}
	return
}
func (cryptos Cryptos) GetValue() (sum float32) {
	for _, crypto := range cryptos {
		sum += *crypto.ValuedAt
	}
	return
}
func (stock Stock) GetValue() (sum float32) {
	return *stock.ValuedAt
}
func (crypto Crypto) GetValue() (sum float32) {
	return *crypto.ValuedAt
}

func roundPercent(f float32) float32 {
	return float32(math.Round(float64(f * 100)))
}

func (record Record) GetStockDiversification(stocksValue float32) (diversification []Diversification) {
	diversification = make([]Diversification, 0, len(record.Stocks))
	for _, stock := range record.Stocks {
		div := Diversification{stock.Symbol, roundPercent(stock.GetValue() / stocksValue)}
		diversification = append(diversification, div)
	}
	return
}

func (record Record) GetCryptoDiversification(cryptosValue float32) (diversification []Diversification) {
	diversification = make([]Diversification, 0, len(record.Cryptos))
	for _, crypto := range record.Cryptos {
		div := Diversification{crypto.Symbol, roundPercent(crypto.GetValue() / cryptosValue)}
		diversification = append(diversification, div)
	}
	return
}

func (record *Record) GenerateStatistics() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		value := record.Stocks.GetValue()
		record.StocksValue = &value
		diversification := record.GetStockDiversification(*record.StocksValue)
		record.StockDiversification = &diversification
		wg.Done()
	}()
	go func() {
		value := record.Cryptos.GetValue()
		record.CryptosValue = &value
		diversification := record.GetCryptoDiversification(*record.CryptosValue)
		record.CryptoDiversification = &diversification
		wg.Done()
	}()
	wg.Wait()
	investedAmount := *record.CryptosValue + *record.StocksValue
	record.InvestedAmount = &investedAmount
	record.TotalInvested = investedAmount + *record.Liquidity
}

// func (recordBody RecordBody) Split() (record Record, info Info) {
// 	info.Id = primitive.NewObjectID()
// 	record.Date = recordBody.Date
// 	record.Id = primitive.NewObjectID()
// 	record.InfoId = info.Id
// 	info.InvestedAmount = recordBody.InvestedAmount
// 	info.Stocks = recordBody.Stocks
// 	info.Cryptos = recordBody.Cryptos
// 	info.GenerateStatistics()
// 	info.Liquidity = recordBody.Liquidity
// 	return
// }

// func ConcatRecord(record Record, info Info) RecordBody {
// 	recordBody := RecordBody{
// 		Id:                    record.Id.Hex(),
// 		Date:                  record.Date,
// 		InvestedAmount:        info.InvestedAmount,
// 		Stocks:                info.Stocks,
// 		Cryptos:               info.Cryptos,
// 		StockDiversification:  info.StockDiversification,
// 		CryptoDiversification: info.CryptoDiversification,
// 		CryptosValue:          info.CryptosValue,
// 		StocksValue:           info.StocksValue,
// 		Liquidity:             info.Liquidity,
// 	}
// 	return recordBody
// }
