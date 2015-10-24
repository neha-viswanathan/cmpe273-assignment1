package main

//import package
import (
  "fmt"
  "strings"
  "strconv"
  "net"
  "net/rpc/jsonrpc"
  "os"
  )

  //struct definitions
type StockRequest struct {
  StockMap map[string]float64
}

type StockResponse struct {
  TradeID string
  Symbol []string
  PricePerStock []float64
  NumOfStock []int
  UnvestedAmount []float64
}

type TradeIDRequest struct {
  TradeID string
}

type TradeIDResponse struct {
  Symbol []string
  CurrentValue []float64
  ChangeInValue []float64
  NumOfStock []int
  UnvestedAmount []float64
}

//main function
func main() {

  //variable declaration
  stock := make(map[string]float64)
  var stocks string
  var unvestedAmount float64

  //create a connection to server
  connect, err := net.Dial("tcp", "127.0.0.1:1111")
  if err != nil {
    panic("Failed to connect : " + err.Error())
  }
  defer connect.Close()

  //create new client to handle requests
  client := jsonrpc.NewClient(connect)

  //check for command line inputs
  if len(os.Args) > 2 {
    arg1 := strings.Split(os.Args[1], ",")
    arg2 := os.Args[2]
    totalBudget, err := strconv.Atoi(arg2)
    if (err!= nil) {
      panic("Failed to convert : " + err.Error())
    }

    //iterate over all the stock Symbols and the respective percentage allocated to each Symbol
    for count := 0; count < len(arg1); count++ {
      stockSym := strings.Split(arg1[count],":")
      percent := strings.TrimSuffix(stockSym[1],"%")
      symPercent, err := strconv.Atoi(percent)
      if err != nil {
        panic("Failed to convert percentage: " + err.Error())
        //return
      }
      //fmt.Println(stockSym, percent, symPercent)
      //allocate amount to each Symbol based on budget and percentage per Symbol
      symbolAmount := float64(totalBudget * symPercent / 100 )
        fmt.Println("amt = ",symbolAmount)
      //store Symbol and its amount in a map
      stock[stockSym[0]] = symbolAmount
    }

    //initialize stockReq
    var stockReq *StockRequest
    var stockResp StockResponse
    stockReq = &StockRequest{stock}

    //client makes a remote call to server
  err = client.Call("YahooFinance.QueryYahooFinance", stockReq, &stockResp)
    if err != nil {
      fmt.Errorf("Error in RPC : ", err)
    }

    for iter:=range stockResp.Symbol{
      //fmt.Println(iter)
      stocks += fmt.Sprintf("%s:%d:$%f", stockResp.Symbol[iter], stockResp.NumOfStock[iter], stockResp.PricePerStock[iter])
    //  fmt.Println(stocks)
      unvestedAmount += stockResp.UnvestedAmount[iter]
    }

    fmt.Println("Trade ID :: ", stockResp.TradeID)
    fmt.Println("Stocks :: ", stocks)
    fmt.Println("Total Unvested Amount :: ", unvestedAmount)
} else {
  //fmt.Println("Else")
  readTradeID := os.Args[1]
  var tradeIDReq *TradeIDRequest
  tradeIDReq = &TradeIDRequest{readTradeID}

  var tradeIDResp TradeIDResponse
  err = client.Call("YahooFinance.GetPortfolioDetails", tradeIDReq, &tradeIDResp)
  if err != nil {
      fmt.Errorf("Error while retrieving portfolio details :: ", err)
  }

  var portfolio, finalCurrentPrice string
  var totalUnvestedAmount float64

  for iter := range tradeIDResp.Symbol {
    if tradeIDResp.ChangeInValue[iter] > 0 {
      portfolio += fmt.Sprintf("%s:%d:+%f, ", tradeIDResp.Symbol[iter], tradeIDResp.NumOfStock[iter], tradeIDResp.CurrentValue[iter])
    } else {
      portfolio += fmt.Sprintf("%s:%d:-%f, ", tradeIDResp.Symbol[iter], tradeIDResp.NumOfStock[iter], tradeIDResp.CurrentValue[iter])
    }

    finalCurrentPrice += fmt.Sprintf("%s:$%f, ", tradeIDResp.Symbol[iter],  tradeIDResp.CurrentValue[iter])
    totalUnvestedAmount += tradeIDResp.UnvestedAmount[iter]
  }

  fmt.Println("Portfolio :: ", portfolio)
  fmt.Println("Current Stock Price :: ", finalCurrentPrice)
  fmt.Println("Total Unvested Amount :: $", totalUnvestedAmount)
}
}
