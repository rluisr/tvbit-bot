package gateway

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hirokisan/bybit/v2"
	"github.com/rluisr/tvbit-bot/pkg/domain"
	"github.com/rluisr/tvbit-bot/utils"
	"github.com/shopspring/decimal"
)

type (
	BybitRepository struct {
		BaseURL      string
		APIKey       string
		APISecretKey string
		HTTPClient   *http.Client
		Client       *bybit.Client
	}
)

func (r *BybitRepository) CreateOrder(req *domain.Order) error {
	orderParam := bybit.V5CreateOrderParam{
		Category:   bybit.CategoryV5Linear,
		Symbol:     bybit.SymbolV5(req.Symbol),
		Qty:        req.QTY,
		TakeProfit: &req.TP,
		StopLoss:   &req.SL,
		Side:       bybit.Side(req.Side),
	}

	if req.Type == "Limit" {
		orderParam.OrderType = bybit.OrderTypeLimit
		orderParam.Price = &req.Price
	}
	if req.Type == "Market" {
		orderParam.OrderType = bybit.OrderTypeMarket
	}

	if req.Side == "Buy" {
		buyHedge := bybit.PositionIdxHedgeBuy
		orderParam.PositionIdx = &buyHedge
	} else {
		sellHedge := bybit.PositionIdxHedgeSell
		orderParam.PositionIdx = &sellHedge
	}

	resp, err := r.Client.V5().Order().CreateOrder(orderParam)
	if err != nil {
		return fmt.Errorf("failed CreateOrder: %w", err)
	}
	req.OrderID = resp.Result.OrderID

	return nil
}

func (r *BybitRepository) FetchOrder(req *domain.Order) error {
	symbol := bybit.SymbolV5(req.Symbol)
	settle := bybit.CoinUSDT

	for i := 0; i < 10; i++ {
		order, err := r.Client.V5().Order().GetOpenOrders(bybit.V5GetOpenOrdersParam{
			Category:   bybit.CategoryV5Linear,
			Symbol:     &symbol,
			OrderID:    &req.OrderID,
			SettleCoin: &settle,
		})
		if err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				continue
			}
			return fmt.Errorf("failed GetOrder: %w", err)
		}

		if len(order.Result.List) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		entryPrice, err := decimal.NewFromString(order.Result.List[0].AvgPrice)
		if err != nil {
			return err
		}

		req.EntryPrice = entryPrice
		break
	}

	return nil
}

func (r *BybitRepository) CalculateTPSL(req domain.Order, value, isType string) (string, error) {
	if value == "" || value == "0" {
		return "0", nil
	}

	currentPrice, err := r.getPrice(req.Symbol)
	if err != nil {
		return "0", err
	}

	if strings.Contains(value, "%") {
		return r.calculateTPSLByPercentage(req, value, isType, currentPrice)
	}

	return r.calculateTPSLByFixedPrice(req, value, currentPrice)
}

func (r *BybitRepository) calculateTPSLByPercentage(req domain.Order, value, isType string, currentPrice float64) (string, error) {
	var price float64
	if req.Price == "0" {
		price = currentPrice
	} else {
		price = utils.StringToFloat64(req.Price)
	}

	tpslStr := strings.Replace(value, "%", "", 1)
	tpsl, err := strconv.ParseFloat(tpslStr, 64)
	if err != nil {
		return "0", errors.New("failed to parse to float64 " + tpslStr)
	}
	tpsl *= 0.1

	qtyF64 := utils.StringToFloat64(req.QTY)

	switch req.Side {
	case "Buy":
		if isType == "TP" {
			return utils.Float64ToString((price * tpsl * qtyF64) + price), nil
		}
		return utils.Float64ToString(math.Abs((price * tpsl * qtyF64) - price)), nil
	case "Sell":
		if isType == "SL" {
			return utils.Float64ToString((price * tpsl * qtyF64) + price), nil
		}
		return utils.Float64ToString(math.Abs((price * tpsl * qtyF64) - price)), nil
	default:
		return "0", errors.New("unknown request / isType: " + isType)
	}
}

func (r *BybitRepository) calculateTPSLByFixedPrice(req domain.Order, value string, currentPrice float64) (string, error) {
	var inputPrice float64
	var err error

	switch {
	case strings.Contains(value, "+"):
		inputPriceString := strings.Replace(value, "+", "", 1)
		inputPrice, err = strconv.ParseFloat(inputPriceString, 64)
		if err != nil {
			return "0", errors.New("convert string to float err " + err.Error())
		}
	case strings.Contains(value, "-"):
		inputPriceString := strings.Replace(value, "-", "", 1)
		inputPrice, err = strconv.ParseFloat(inputPriceString, 64)
		if err != nil {
			return "0", errors.New("convert string to float err " + err.Error())
		}
	default:
		return value, nil
	}

	if req.Side == "Sell" {
		return utils.Float64ToString(currentPrice + inputPrice), nil
	}
	return utils.Float64ToString(currentPrice - inputPrice), nil
}

// getPrice returns index price
func (r *BybitRepository) getPrice(symbol string) (float64, error) {
	v5symbol := bybit.SymbolV5(symbol)

	resp, err := r.Client.V5().Market().GetTickers(bybit.V5GetTickersParam{
		Category: bybit.CategoryV5Linear,
		Symbol:   &v5symbol,
	})
	if err != nil {
		return 0, fmt.Errorf("bybit GetTickers(): %w", err)
	}

	return utils.StringToFloat64(resp.Result.LinearInverse.List[0].IndexPrice), nil
}

func (r *BybitRepository) GetWalletBalance() (float64, error) {
	resp, err := r.Client.V5().Account().GetWalletBalance(bybit.AccountTypeUnified, []bybit.Coin{bybit.CoinUSDT})
	if err != nil {
		return 0, err
	}

	balance, err := strconv.ParseFloat(resp.Result.List[0].TotalAvailableBalance, 64)
	if err != nil {
		return 0, err
	}

	return balance, nil
}
