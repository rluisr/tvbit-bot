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

	"github.com/jinzhu/copier"

	"github.com/frankrap/bybit-api/rest"
	"github.com/rluisr/tvbit-bot/pkg/domain"
	"github.com/rluisr/tvbit-bot/pkg/external/bybit"
	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/common/crypto"
)

type (
	BybitRepository struct {
		BaseURL      string
		APIKey       string
		APISecretKey string
		HTTPClient   *http.Client
		Client       *rest.ByBit
	}
)

func (r *BybitRepository) Set(req domain.TV) {
	r.Client, r.BaseURL = bybit.Init(req, r.HTTPClient)
	r.APIKey = req.APIKey
	r.APISecretKey = req.APISecretKey
}

// GetPositions return an array of positions.
// This func uses deep copy if the symbol is PERP/USDC so return positions are not exactly but only used for GetActiveOrderCount
func (r *BybitRepository) GetPositions(symbol string) (*[]rest.LinearPosition, error) {
	var positions []rest.LinearPosition

	if strings.Contains(symbol, "PERP") {
		usdcPositions, err := r.getUSDCPosition()
		if err != nil {
			return nil, err
		}

		err = copier.Copy(&positions, &usdcPositions.Result.DataList)
		if err != nil {
			return nil, fmt.Errorf("GetPositions failed deep copy %w", err)
		}
	} else {
		_, _, linearPositions, err := r.Client.LinearGetPosition(symbol)
		if err != nil {
			return nil, err
		}
		positions = linearPositions
	}

	return &positions, nil
}

func (r *BybitRepository) CreateOrder(req domain.TV) (string, error) {
	var orderID string

	if strings.Contains(req.Order.Symbol, "PERP") {
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
			return "", fmt.Errorf("SignedRequest err: %w, body: %s", err, string(resp))
		}
		orderID = order.Result.OrderID
	} else {
		_, _, order, err := r.Client.LinearCreateOrder(req.Order.Side, req.Order.Type, req.Order.Price, req.Order.QTY, "ImmediateOrCancel", req.Order.TP.(float64), req.Order.SL.(float64), false, false, "", req.Order.Symbol)
		if err != nil {
			return "", err
		}
		orderID = order.OrderId
	}

	return orderID, nil
}

// FetchOrder is set entry price
func (r *BybitRepository) FetchOrder(req *domain.TV, orderID string) error {
	entryPrice, err := r.getEntryPrice(req)
	if err != nil {
		return err
	}

	req.Order.EntryPrice = *entryPrice
	req.Order.OrderID = orderID

	return nil
}

func (r *BybitRepository) GetActiveOrderCount(req *domain.TV, positions *[]rest.LinearPosition) int {
	var activeOrder int

	if strings.Contains(req.Order.Symbol, "PERP") {
		activeOrder = len(*positions)
	} else {
		for _, position := range *positions {
			if position.Side == req.Order.Side {
				if position.Size > 0 {
					activeOrder = 1
				}
			}
		}
	}

	return activeOrder
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

	currentPrice, err := r.getMarkPrice(req.Order.Symbol)
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

// getMarkPrice returns mark price
func (r *BybitRepository) getMarkPrice(symbol string) (float64, error) {
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

func (r *BybitRepository) getEntryPrice(req *domain.TV) (*decimal.Decimal, error) {
	var (
		entryPrice decimal.Decimal
		err        error
	)

	if strings.Contains(req.Order.Symbol, "PERP") {
		var positions domain.BybitUSDCPositions

		resp, err := r.signedRequestWithHeader(http.MethodPost, "option/usdc/openapi/private/v1/query-position", []byte("{\"category\":\"PERPETUAL\"}"), &positions)
		if err != nil {
			return nil, fmt.Errorf("signedRequest err: %w, body: %s", err, resp)
		}

		for _, position := range positions.Result.DataList {
			entryPrice, err = decimal.NewFromString(position.EntryPrice)
			if err != nil {
				return nil, err
			}
		}
	} else {
		_, _, positions, err := r.Client.LinearGetPosition(req.Order.Symbol)
		if err != nil {
			return nil, err
		}

		for _, position := range positions {
			if position.Side == req.Order.Side {
				entryPrice = decimal.NewFromFloat(position.EntryPrice)
			}
		}
	}

	return &entryPrice, err
}

func (r *BybitRepository) getUSDCPosition() (*domain.BybitUSDCPositions, error) {
	var positions domain.BybitUSDCPositions

	resp, err := r.signedRequestWithHeader(http.MethodPost, "option/usdc/openapi/private/v1/query-position", []byte("{\"category\":\"PERPETUAL\"}"), &positions)
	if err != nil {
		return nil, fmt.Errorf("signedRequest err: %w, body: %s", err, resp)
	}

	return &positions, nil
}

func (r *BybitRepository) GetWalletInfoUSDC() (*domain.BybitWallet, error) {
	var wallet domain.BybitWallet

	resp, err := r.signedRequestWithHeader(http.MethodPost, "option/usdc/openapi/private/v1/query-wallet-balance", []byte("{}"), &wallet)
	if err != nil {
		return nil, fmt.Errorf("signedRequest err: %w, body: %s", err, resp)
	}

	return &wallet, nil
}

func (r *BybitRepository) GetWalletInfoDeriv() (*rest.Balance, error) {
	_, _, wallet, err := r.Client.GetWalletBalance("USDT")
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *BybitRepository) signedRequestWithHeader(method, path string, body []byte, result interface{}) (string, error) {
	var payload io.Reader

	nowTimeInMilli := strconv.FormatInt(time.Now().UnixMilli(), 10)

	if body == nil {
		// TODO have to using query params. it's not working when query params are passed.
		payload = strings.NewReader("")
	} else {
		payload = strings.NewReader(string(body))
	}

	signInput := nowTimeInMilli + r.APIKey + "5000" + string(body)
	hmacSigned, err := crypto.GetHMAC(crypto.HashSHA256, []byte(signInput), []byte(r.APISecretKey))
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s%s", r.BaseURL, path)

	req, _ := http.NewRequest(method, url, payload)
	req.Header.Add("X-BAPI-API-KEY", r.APIKey)
	req.Header.Add("X-BAPI-SIGN", crypto.HexEncodeToString(hmacSigned))
	req.Header.Add("X-BAPI-SIGN-TYPE", "2")
	req.Header.Add("X-BAPI-TIMESTAMP", nowTimeInMilli)
	req.Header.Add("X-BAPI-RECV-WINDOW", "5000")

	resp, err := r.HTTPClient.Do(req)
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
