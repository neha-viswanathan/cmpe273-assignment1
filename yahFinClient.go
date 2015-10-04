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
  tradeID string
  symbol []string
  pricePerStock []float64
  numOfStock []int
  unvestedAmount []float64
}

type TradeIDRequest struct {
  tradeID string
}

type TradeIDResponse struct {
  symbol []string
  currentValue []float64
  changeInValue []float64
  numOfStock []int
  unvestedAmount []float64
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
    if err != nil {
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
      fmt.Println(stockSym, percent, symPercent)
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
  err = client.Call("YahooFinance.QueryYahooFinance", &stockReq, &stockResp)
    if err != nil {
      fmt.Errorf("Error in RPC : ", err)
    }

    for iter:=range stockResp.symbol{
      //fmt.Println(iter)
      stocks += fmt.Sprintf("%s:%d:$%f", stockResp.symbol[iter], stockResp.numOfStock[iter], stockResp.pricePerStock[iter])
    //  fmt.Println(stocks)
      unvestedAmount += stockResp.unvestedAmount[iter]
    }

    fmt.Println("Trade ID :: ", stockResp.tradeID)
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

  for iter := range tradeIDResp.symbol {
    if tradeIDResp.changeInValue[iter] > 0 {
      portfolio += fmt.Sprintf("%s:%d:+%f, ", tradeIDResp.symbol[iter], tradeIDResp.numOfStock[iter], tradeIDResp.currentValue[iter])
    } else {
      portfolio += fmt.Sprintf("%s:%d:-%f, ", tradeIDResp.symbol[iter], tradeIDResp.numOfStock[iter], tradeIDResp.currentValue[iter])
    }

    finalCurrentPrice += fmt.Sprintf("%s:$%f, ", tradeIDResp.symbol[iter],  tradeIDResp.currentValue[iter])
    totalUnvestedAmount += tradeIDResp.unvestedAmount[iter]
  }

  fmt.Println("Portfolio :: ", portfolio)
  fmt.Println("Current Stock Price :: ", finalCurrentPrice)
  fmt.Println("Total Unvested Amount :: $", totalUnvestedAmount)
}
}
