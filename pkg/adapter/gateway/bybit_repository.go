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
		OrderType:  bybit.OrderTypeMarket,
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

// CalculateTPSL returns TP and SL
// "tp" and "sl" are not allowed to be 0 or not specified.
func (r *BybitRepository) CalculateTPSL(req *domain.Order) error {
	// INFO: テストネットとメインでは価格差が大きく、テストネットで注文を行う際に TP/SL の範囲外になる可能性があり、注文が失敗することがある
	currentPrice, err := r.getPrice(req.Symbol)
	if err != nil {
		return err
	}

	var (
		tp string
		sl string
	)

	if strings.Contains(req.TP, "%") {
		tp, sl, err = r.calculateTPSLByPercentage(req, currentPrice)
	} else {
		tp, sl, err = r.calculateTPSLByFixedPrice(req, currentPrice)
	}
	if err != nil {
		return err
	}

	req.TP = tp
	req.SL = sl

	return nil
}

func (r *BybitRepository) calculateTPSLByPercentage(req *domain.Order, currentPrice float64) (tp string, sl string, err error) {
	var price float64
	if req.Price == "0" {
		price = currentPrice
	} else {
		price = utils.StringToFloat64(req.Price)
	}

	tpStr := strings.Replace(req.TP, "%", "", 1)
	tpF64, err := strconv.ParseFloat(tpStr, 64)
	if err != nil {
		return "0", "0", errors.New("failed to parse to float64 " + tpStr)
	}
	tpF64 *= 0.1

	slStr := strings.Replace(req.SL, "%", "", 1)
	slF64, err := strconv.ParseFloat(slStr, 64)
	if err != nil {
		return "0", "0", errors.New("failed to parse to float64 " + slStr)
	}
	slF64 *= 0.1

	qtyF64 := utils.StringToFloat64(req.QTY)

	switch req.Side {
	case "Buy":
		tp = utils.Float64ToString((price * tpF64 * qtyF64) + price)
		sl = utils.Float64ToString(math.Abs((price * slF64 * qtyF64) - price))
	case "Sell":
		tp = utils.Float64ToString(math.Abs((price * tpF64 * qtyF64) - price))
		sl = utils.Float64ToString((price * slF64 * qtyF64) + price)
	default:
		return "0", "0", errors.New("invalid side")
	}

	return tp, sl, nil
}

// calculateTPSLByFixedPrice returns TP and SL
// req.TP contains "+" or "-" and a number (e.g. "+100", "-100")
// it means the price difference from the current price
func (r *BybitRepository) calculateTPSLByFixedPrice(req *domain.Order, currentPrice float64) (tp string, sl string, err error) {
	var inputPrice float64

	// TP
	switch {
	case strings.Contains(req.TP, "+"):
		inputPriceStr := strings.Replace(req.TP, "+", "", 1)
		inputPrice, err = strconv.ParseFloat(inputPriceStr, 64)
		if err != nil {
			return "0", "0", err
		}
		tp = utils.Float64ToString(currentPrice + inputPrice)
	case strings.Contains(req.TP, "-"):
		inputPriceStr := strings.Replace(req.TP, "-", "", 1)
		inputPrice, err = strconv.ParseFloat(inputPriceStr, 64)
		if err != nil {
			return "0", "0", err
		}
		tp = utils.Float64ToString(currentPrice - inputPrice)
	default:
		tp = req.TP
	}

	// SL
	switch {
	case strings.Contains(req.SL, "+"):
		inputPriceStr := strings.Replace(req.SL, "+", "", 1)
		inputPrice, err = strconv.ParseFloat(inputPriceStr, 64)
		if err != nil {
			return "0", "0", err
		}
		sl = utils.Float64ToString(currentPrice + inputPrice)
	case strings.Contains(req.SL, "-"):
		inputPriceStr := strings.Replace(req.SL, "-", "", 1)
		inputPrice, err = strconv.ParseFloat(inputPriceStr, 64)
		if err != nil {
			return "0", "0", err
		}
		sl = utils.Float64ToString(currentPrice - inputPrice)
	default:
		sl = req.SL
	}

	return tp, sl, nil
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
