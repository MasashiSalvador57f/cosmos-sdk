package bank

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc"
)

func unknownRequest(prefix string, ty interface{}) sdk.Result {
	errMsg := prefix + reflect.TypeOf(ty).Name()
	return sdk.ErrUnknownRequest(errMsg).Result()
}

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case ibc.MsgSend:
			return k.ibck.Send(func(p ibc.Payload) sdk.Result {
				switch p := msg.Payload.(type) {
				case PayloadCoins:
					return handlePayloadCoinsSend(ctx, k, p)
				default:
					return unknownRequest("Unrecognized ibc/bank payload type: ", p)
				}
			}, ctx, k.ibcStore(ctx), msg)
		case ibc.MsgReceive:
			return k.ibck.Receive(func(ctx sdk.Context, p ibc.Payload) (ibc.Payload, sdk.Result) {
				switch p := msg.Payload.(type) {
				case PayloadCoins:
					return handlePayloadCoinsReceive(ctx, k, p)
				case PayloadCoinsReceipt:
					return handlePayloadCoinsReceiptReceive(ctx, k, p)
				default:
					return nil, unknownRequest("Unrecognized ibc/bank payload type: ", p)
				}

			}, ctx, k.ibcStore(ctx), msg)
		// case ibc.MsgRelay
		default:
			return unknownRequest("Unrecognized ibc/bank Msg type: ", msg)
		}
	}
}

func handlePayloadCoinsSend(ctx sdk.Context, k Keeper, p PayloadCoins) sdk.Result {
	_, tags, err := k.bk.SubtractCoins(ctx, p.SrcAddr, p.Coins)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{Tags: tags}
}

func handlePayloadCoinsReceive(ctx sdk.Context, k Keeper, p PayloadCoins) (ibc.Payload, sdk.Result) {
	_, tags, err := k.bk.AddCoins(ctx, p.DestAddr, p.Coins)
	if err != nil {
		return PayloadCoinsReceipt{p, false}, err.Result()
	}
	return PayloadCoinsReceipt{p, true}, sdk.Result{Tags: tags}
}

func handlePayloadCoinsReceiptReceive(ctx sdk.Context, k Keeper, p PayloadCoinsReceipt) (_ ibc.Payload, _ sdk.Result) {
	if !p.Success {
		k.bk.AddCoins(ctx, p.SrcAddr, p.Coins)
	}
	return
}