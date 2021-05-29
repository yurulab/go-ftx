package realtime_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-numb/go-ftx/realtime"
	"golang.org/x/sync/errgroup"
)

func TestConnect(t *testing.T) {
	eg, ctx := errgroup.WithContext(context.Background())
	ctxConnect, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := make(chan realtime.Response)
	eg.Go(func() error {
		return realtime.Connect(ctxConnect, ch, []string{"ticker"}, []string{"BTC-PERP", "ETH-PERP"}, nil)
	})

	go func() {
		for {
			select {
			case v := <-ch:
				switch v.Type {
				case realtime.TICKER:
					fmt.Printf("%s	%+v\n", v.Symbol, v.Ticker)

				case realtime.TRADES:
					fmt.Printf("%s	%+v\n", v.Symbol, v.Trades)
					for i := range v.Trades {
						if v.Trades[i].Liquidation {
							fmt.Printf("-----------------------------%+v\n", v.Trades[i])
						}
					}

				case realtime.ORDERBOOK:
					fmt.Printf("%s	%+v\n", v.Symbol, v.Orderbook)

				case realtime.UNDEFINED:
					fmt.Printf("%s	%s\n", v.Symbol, v.Results.Error())
				}
			}
		}
	}()

	if err := eg.Wait(); err != nil {
		close(ch)
		cancel()
		fmt.Println("Please restart process.")
	}
}

func TestConnectForPrivate(t *testing.T) {
	eg, ctx := errgroup.WithContext(context.Background())
	ctxConnect, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := make(chan realtime.Response)
	eg.Go(func() error {
		return realtime.ConnectForPrivate(ctxConnect, ch, "", "", []string{"orders", "fills"}, nil)
	})

	go func() {
		for {
			select {
			case v := <-ch:
				switch v.Type {
				case realtime.ORDERS:
					fmt.Printf("%d	%+v\n", v.Type, v.Orders)
				case realtime.FILLS:
					fmt.Printf("%d	%+v\n", v.Type, v.Fills)

				case realtime.UNDEFINED:
					fmt.Printf("UNDEFINED %s	%s\n", v.Symbol, v.Results.Error())
				}
			case <-ctxConnect.Done():
				return
			}
		}
	}()

	if err := eg.Wait(); err != nil {
		close(ch)
		cancel()
		fmt.Println("Please restart process.")
	}
}
