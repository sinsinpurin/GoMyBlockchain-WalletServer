package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/sinsinpurin/gomyblockchain"
	"github.com/sinsinpurin/gomyblockchain-walletserver/httpreq"
	"github.com/sinsinpurin/gomyblockchain-walletserver/responsetypes"
)

// Port ポート番号の設定
var Port string

var blockChainCache gomyblockchain.BlockChainServer

var responseWallet responsetypes.WalletType

// BlockChainServerAPIURL ブロックチェーンサーバーのURL
const BlockChainServerAPIURL = "http://0.0.0.0:8085"

func main() {
	flag.StringVar(&Port, "p", ":8081", "Set Port Number")
	flag.Parse()
	e := echo.New()

	e.Use(middleware.CORS())

	// wallet を生成し，JSONで返します
	e.POST("/wallet", func(c echo.Context) error {
		wallet := gomyblockchain.GenerateWallet()
		resWallet := &responsetypes.WalletType{
			PrivateKey: hex.EncodeToString(wallet.PrivateKey),
			PublicKey:  hex.EncodeToString(wallet.PublicKey),
			Address:    wallet.Address,
		}
		return c.JSON(http.StatusOK, resWallet)
	})

	e.POST("/transaction", func(c echo.Context) error {
		param := new(responsetypes.TransactionParams)
		if err := c.Bind(param); err != nil {
			return err
		}
		bprivateKey, _ := hex.DecodeString(param.SenderPrivateKey)
		bpublicKey, _ := hex.DecodeString(param.SenderPublicKey)
		// signature 作成
		wallet := &gomyblockchain.Wallet{
			PrivateKey: bprivateKey,
			PublicKey:  bpublicKey,
			Address:    param.SenderAddress,
		}
		iValue, _ := strconv.Atoi(param.Value)
		transaction := gomyblockchain.CreateTransaction(param.SenderAddress, param.RecipientAddress, uint64(iValue))
		bsignature := gomyblockchain.GenerateSignature(*wallet, transaction)
		// response レスポンスの作成
		transactionWithSig := &responsetypes.TransactionType{
			RecipientAddress: param.RecipientAddress,
			SenderAddress:    param.SenderAddress,
			Value:            uint64(iValue),
			SenderPublicKey:  param.SenderPublicKey,
			Signature:        hex.EncodeToString(bsignature),
		}
		resp, err := httpreq.ReqPOST(BlockChainServerAPIURL+"/transactions", transactionWithSig)
		if err != nil {
			return c.JSON(http.StatusBadGateway, nil)
		}
		return c.JSON(http.StatusOK, resp.Body)
	})

	e.GET("/wallet/amount", func(c echo.Context) error {
		var query responsetypes.GetRequestWalletAmountType
		query.Address = c.QueryParam("Address")
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, BlockChainServerAPIURL+"/amount", nil)
		if err != nil {
			panic(err)
		}
		params := req.URL.Query()
		params.Add("Address", query.Address)
		req.URL.RawQuery = params.Encode()
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		var responseT responsetypes.GetResponseWalletAmountType
		respBodyB, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(respBodyB, &responseT); err != nil {
			log.Fatal(err)
		}
		fmt.Println(responseT.Amount)
		return c.JSON(http.StatusOK, responseT.Amount)
	})

	e.Logger.Fatal(e.Start(Port))
}
