package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/frankrap/bybit-api/rest"
	"github.com/rluisr/tvbit-bot/pkg/domain"
	"github.com/rluisr/tvbit-bot/pkg/external/bybit"
	"github.com/thrasher-corp/gocryptotrader/common/crypto"
)

type (
	BybitRepository struct {
		BaseURL      string
		APIKey       string
		APISecretKey string
		Client       *rest.ByBit
	}
)

func (r *BybitRepository) Set(req domain.TV) {
	r.Client, r.BaseURL = bybit.Init(req)
	r.APIKey = req.APIKey
	r.APISecretKey = req.APISecretKey
}

func (r *BybitRepository) CreateOrder(req domain.TV) (*domain.TVOrder, error) {
	orderHistory := &domain.TVOrder{
		Type:   req.Order.Type,
		Symbol: req.Order.Symbol,
		Side:   req.Order.Side,
		QTY:    req.Order.QTY,
		TP:     req.Order.TP,
		SL:     req.Order.SL,
	}

	if !strings.Contains(req.Order.Symbol, "PERP") {
		_, _, order, err := r.Client.LinearCreateOrder(req.Order.Side, req.Order.Type, req.Order.Price, req.Order.QTY, "ImmediateOrCancel", req.Order.TP.(float64), req.Order.SL.(float64), false, false, "", req.Order.Symbol)
		if err != nil {
			return nil, err
		}

		entryPrice, err := order.Price.Float64()
		if err != nil {
			return nil, err
		}
		orderHistory.OrderID = order.OrderId
		orderHistory.Price = entryPrice
		return orderHistory, nil
	}

	params := map[string]interface{}{}
	params["symbol"] = req.Order.Symbol
	params["orderType"] = req.Order.Type
	params["orderFilter"] = "Order"
	params["orderQty"] = req.Order.QTY
	params["side"] = req.Order.Side
	params["timeInForce"] = "ImmediateOrCancel"
	params["takeProfit"] = strconv.FormatFloat(req.Order.TP.(float64), 'f', -1, 64)
	params["stopLoss"] = strconv.FormatFloat(req.Order.SL.(float64), 'f', -1, 64)

	var order domain.BybitPerpResponseOrder
	_, resp, err := r.Client.SignedRequest(http.MethodPost, "perpetual/usdc/openapi/private/v1/place-order", params, &order)
	if err != nil {
		return nil, fmt.Errorf("SignedRequest err: %w, body: %s", err, string(resp))
	}

	orderPrice, err := strconv.ParseFloat(order.Result.OrderPrice, 64)
	if err != nil {
		// bybit returns fucking error response sometimes but order is success.
		return nil, err
	}
	orderHistory.OrderID = order.Result.OrderID
	orderHistory.Price = orderPrice

	return orderHistory, nil
}

// GetCurrentPrice returns mark price
func (r *BybitRepository) GetCurrentPrice(symbol string) (float64, error) {
	var isPerp bool

	var tickersURL string
	if strings.Contains(symbol, "PERP") {
		isPerp = true
		tickersURL = fmt.Sprintf("perpetual/usdc/openapi/public/v1/tick?symbol=%s", symbol)
	} else {
		tickersURL = fmt.Sprintf("derivatives/v3/public/tickers?category=linear&symbol=%s", symbol)
	}

	var markPriceStr string

	if !isPerp {
		var ticker domain.BybitDerivTicker
		_, resp, err := r.Client.PublicRequest(http.MethodGet, tickersURL, nil, &ticker)
		if err != nil {
			return 0, fmt.Errorf("PublicRequest err: %w, body: %s", err, string(resp))
		}
		markPriceStr = ticker.Result.List[0].MarkPrice
	} else {
		var ticker domain.BybitPerpTicker
		_, resp, err := r.Client.PublicRequest(http.MethodGet, tickersURL, nil, &ticker)
		if err != nil {
			return 0, fmt.Errorf("PublicRequest err: %w, body: %s", err, string(resp))
		}
		markPriceStr = ticker.Result.MarkPrice
	}

	return strconv.ParseFloat(markPriceStr, 64)
}

func (r *BybitRepository) CalculateTPSL(req domain.TV, value interface{}, isType string) (float64, error) {
	var err error

	str, isOK := value.(string)
	if !isOK {
		return 0, fmt.Errorf("failed to assertion to string: %v", value)
	}

	if str == "" || str == "0" {
		return 0, nil
	}

	currentPrice, err := r.GetCurrentPrice(req.Order.Symbol)
	if err != nil {
		return 0, err
	}

	if strings.Contains(str, "%") {
		var price float64

		if req.Order.Price == 0 {
			price = currentPrice
		} else {
			price = req.Order.Price
		}

		tpslStr := strings.Replace(value.(string), "%", "", 1)
		tpsl, err := strconv.ParseFloat(tpslStr, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse to float64 %s", tpslStr)
		}
		tpsl = tpsl * 0.1

		switch req.Order.Side {
		case "Buy":
			if isType == "TP" {
				return (price * tpsl * req.Order.QTY) + price, nil
			}
			return math.Abs((price * tpsl * req.Order.QTY) - price), nil
		case "Sell":
			if isType == "SL" {
				return (price * tpsl * req.Order.QTY) + price, nil
			}
			return math.Abs((price * tpsl * req.Order.QTY) - price), nil
		default:
			return 0, fmt.Errorf("unknown request: %v / isType: %s", req, isType)
		}
	}

	if strings.Contains(str, "+") {
		inputPriceString := strings.Replace(str, "+", "", 1)
		inputPrice, err := strconv.Atoi(inputPriceString)
		if err != nil {
			return 0, fmt.Errorf("convert string to int err %w", err)
		}
		if req.Order.Side == "Sell" {
			return currentPrice - float64(inputPrice), nil
		}
		return currentPrice + float64(inputPrice), nil
	}
	if strings.Contains(str, "-") {
		inputPriceString := strings.Replace(str, "-", "", 1)
		inputPrice, err := strconv.Atoi(inputPriceString)
		if err != nil {
			return 0, fmt.Errorf("convert string to int err %w", err)
		}
		if req.Order.Side == "Sell" {
			return currentPrice + float64(inputPrice), nil
		}
		return currentPrice - float64(inputPrice), nil
	}

	f64, isOK := value.(float64)
	if !isOK {
		return 0, fmt.Errorf("failed to assertion to float64: %v", value)
	}

	return f64, nil
}

func (r *BybitRepository) GetWalletInfoUSDC() (*domain.BybitWallet, error) {
	var wallet domain.BybitWallet

	resp, err := r.signedRequestWithHeaderWithoutBody(http.MethodPost, "option/usdc/openapi/private/v1/query-wallet-balance", &wallet)
	if err != nil {
		return nil, fmt.Errorf("signedRequestWithHeaderWithoutBody err: %w, body: %s", err, string(resp))
	}

	return &wallet, nil
}

func (r *BybitRepository) signedRequestWithHeaderWithoutBody(method, path string, result interface{}) (string, error) {
	nowTimeInMilli := strconv.FormatInt(time.Now().UnixMilli(), 10)
	var hmacSigned []byte
	payload := strings.NewReader(`{}`)
	signInput := nowTimeInMilli + r.APIKey + "5000" + "{}"
	hmacSigned, err := crypto.GetHMAC(crypto.HashSHA256, []byte(signInput), []byte(r.APISecretKey))
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s%s", r.BaseURL, path)

	client := &http.Client{}
	req, _ := http.NewRequest(method, url, payload)
	req.Header.Add("X-BAPI-API-KEY", r.APIKey)
	req.Header.Add("X-BAPI-SIGN", crypto.HexEncodeToString(hmacSigned))
	req.Header.Add("X-BAPI-SIGN-TYPE", "2")
	req.Header.Add("X-BAPI-TIMESTAMP", nowTimeInMilli)
	req.Header.Add("X-BAPI-RECV-WINDOW", "5000")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(b, result)

	return string(b), err
}
