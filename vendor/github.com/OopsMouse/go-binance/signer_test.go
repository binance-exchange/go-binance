package binance

import (
	"testing"
)

func TestSuccess(t *testing.T) {
	signer := &HmacSigner{
		Key: []byte("NhqPtmdSJYdKjVHjA7PZj4Mge3R5YNiP1e3UZjInClVN65XAbvqqM6A7H5fATj0j"),
	}
	queryString := []byte("symbol=LTCBTC&side=BUY&type=LIMIT&timeInForce=GTC&quantity=1&price=0.1&recvWindow=5000&timestamp=1499827319559")

	s := signer.Sign(queryString)
	if s != "c8db56825ae71d6d79447849e617115f4a920fa2acdcab2b053c4b2838bd6b71" {
		t.Errorf("signer returned invalid signature: %s", s)
	}
}
