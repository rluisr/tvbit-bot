/*
 *
 * tvbit-bot
 * Copyright (C) 2022  rluisr(Takuya Hasegawa)
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 * /
 */

package external

import (
	"context"

	"github.com/rluisr/tvbit-bot/pkg/domain"

	"github.com/shopspring/decimal"

	"github.com/adhocore/gronx/pkg/tasker"
)

func cron() {
	task := tasker.New(tasker.Option{
		Verbose: true,
	})

	/*
		Set PnL to an order history.
		Update a record if createdAt of saved order history on DB is lower than last closed order(GetClosedOrderLast)'s createdAt.
		Why this implementation?
		An order_id of open and close is difference, so we can't fetch a closed PnL with open order_id.
	*/
	task.Task("* * * * *", func(ctx context.Context) (int, error) {
		settings, err := tvController.Interactor.TVRepository.GetSettings()
		if err != nil {
			return 0, err
		}

		for _, setting := range settings {
			switch setting.CEX {
			case "bybit":
				tvController.Bybit(domain.TV{
					IsTestNet:    setting.IsTestnet,
					APIKey:       setting.APIKey,
					APISecretKey: setting.APISecretKey,
				})

				savedPLNullOrders, err := tvController.Interactor.TVRepository.GetPLNullOrders(setting.ID)
				if err != nil {
					return 0, err
				}

				for _, savedOrder := range *savedPLNullOrders {
					closedOrder, err := tvController.Interactor.BybitRepository.GetClosedOrderLast(savedOrder.Symbol)
					if err != nil {
						return 0, err
					}

					if savedOrder.CreatedAt.Unix() < closedOrder.Result.Data[0].CreatedAt {
						savedOrder.PL = decimal.NewFromFloat(closedOrder.Result.Data[0].ClosedPnl)
						err = tvController.Interactor.TVRepository.UpdateOrder(&savedOrder)
						if err != nil {
							return 0, err
						}
						break
					}
				}
			}
		}
		return 0, nil
	})

	/*
		Save wallet balance
	*/
	task.Task("0 * * * *", func(ctx context.Context) (int, error) {
		settings, err := tvController.Interactor.TVRepository.GetSettings()
		if err != nil {
			return 0, err
		}

		var walletHistories []domain.WalletHistory

		for _, setting := range settings {
			switch setting.CEX {
			case "bybit":
				tvController.Bybit(domain.TV{
					IsTestNet:    setting.IsTestnet,
					APIKey:       setting.APIKey,
					APISecretKey: setting.APISecretKey,
				})

				// USDC
				bybitUSDCWallet, err := tvController.Interactor.BybitRepository.GetWalletInfoUSDC()
				if err != nil {
					return 1, err
				}

				balance, err := decimal.NewFromString(bybitUSDCWallet.Result.WalletBalance)
				if err != nil {
					return 1, err
				}
				totalRPL, err := decimal.NewFromString(bybitUSDCWallet.Result.TotalRPL)
				if err != nil {
					return 1, err
				}

				walletHistories = append(walletHistories, domain.WalletHistory{
					SettingID: setting.ID,
					Type:      "usdc",
					Balance:   balance,
					TotalRPL:  totalRPL,
				})

				// Deriv USDT
				bybitDerivWallet, err := tvController.Interactor.BybitRepository.GetWalletInfoDeriv()
				if err != nil {
					return 1, err
				}
				walletHistories = append(walletHistories, domain.WalletHistory{
					SettingID: setting.ID,
					Type:      "usdt",
					Balance:   decimal.NewFromFloat(bybitDerivWallet.Equity),
					TotalRPL:  decimal.NewFromFloat(bybitDerivWallet.CumRealisedPnl),
				})
			}
		}

		err = tvController.Interactor.TVRepository.SaveWalletHistories(walletHistories)
		if err != nil {
			return 1, err
		}
		return 0, nil
	})

	task.Run()
}
