package main

//import package
import (
  "fmt"
  "strconv"
  "net"
	"net/rpc"
  "net/http"
	"net/rpc/jsonrpc"
  "crypto/rand"
  "time"
  "encoding/json"
  "io/ioutil"
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

type Data struct {
  pricePerStock float64
  numOfStock int
  unvestedAmount float64
}

type StockData struct {
  List struct {
    Resources []struct {
      Resource struct {
        Fields struct {
          CompanyName string `json:"name"`
          StockPrice string `json:"price"`
          Symbol string `json:"symbol"`
          Timestamp string `json:"ts"`
          StockType string `json:"stockType"`
          UTCTime string `json:"utcTime"`
          Volume string `json:"volume"`
        } `json:"fields"`
      } `json:"resource"`
    } `json:"resources"`
  } `json:"list"`
}

const (
	timeout = time.Duration(time.Second * 10)
)

  //variable declaration
  var sDetail Data
  var tradeIDMap map[string]map[string] Data
  var tempMap1 map[string]Data
  var tempMap2 map[string]Data
  type YahooFinance int

  //trade ID generator
  func randStr(strSize int, randType string) string {
          var dictionary string

          if randType == "alphanum" {
                  dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
          }


          var bytes = make([]byte, strSize)
          rand.Read(bytes)
          for k, v := range bytes {
                  bytes[k] = dictionary[v%byte(len(dictionary))]
          }
          return string(bytes)
  }

  func (yF *YahooFinance) QueryYahooFinance(stockReq *StockRequest, stockResp *StockResponse) error {

    //generate trade ID
    tradeID := randStr(10, "alphanum")
  /*  if err != nil {
      fmt.Println("Error in generating Trade ID :: %v", err)
    }*/
    fmt.Println("TradeID :: ",tradeID)
    stockResp.tradeID = string(tradeID)
  //  fmt.Println("in here1")
    tradeIDMap = make(map[string]map[string] Data)
    tempMap1 = make(map[string]Data)
  //  fmt.Println("in here2")

    //Retrieve the price and number of stocks
    for key, value:= range stockReq.StockMap{
      priceStock := getPrice(key)
      stockCount := int(value/priceStock)
      unvestAmt := value - float64(stockCount)*priceStock

      stockResp.symbol = append(stockResp.symbol, key)
      stockResp.pricePerStock = append(stockResp.pricePerStock, priceStock)
      stockResp.numOfStock = append(stockResp.numOfStock, int(stockCount))
      stockResp.unvestedAmount = append(stockResp.unvestedAmount, unvestAmt)

      sDetail = Data{priceStock,int(stockCount),unvestAmt}

      tempMap1[key] = sDetail
    }
    tradeIDMap[tradeID]=tempMap1
    fmt.Println(tradeIDMap)
    return nil
  }

func getPrice(sym string) float64 {

  var stockD StockData
  //fmt.Println("in getPrice")
  //set timeout if no response
  client := http.Client{Timeout: timeout}
  yahooURL := fmt.Sprintf("http://finance.yahoo.com/webservice/v1/symbols/%s/quote?format=json", sym)
  //mt.Println(yahooURL)
	result, err := client.Get(yahooURL)
	if err != nil {
		fmt.Errorf("Unable to access Yahoo Finance API :: %v", err)
	}
	defer result.Body.Close()

	content, err := ioutil.ReadAll(result.Body)
	if err != nil {
		fmt.Errorf("JSON cannot be read :: %v", err)
	}

	err = json.Unmarshal(content, &stockD)
	if err != nil {
		fmt.Errorf("Error while parsing JSON data: %v", err)
	}
//  fmt.Println("Unmarshalled")
	price, err := strconv.ParseFloat(stockD.List.Resources[0].Resource.Fields.StockPrice, 64)
	if err != nil {
		fmt.Errorf("Price per Stock :: %v", err)
	}
  //fmt.Println(price)
	return price
}

func (yF *YahooFinance) GetPortfolioDetails(tradeIDReq *TradeIDRequest, tradeIDResp *TradeIDResponse) error {
  tempMap2 := tradeIDMap[tradeIDReq.tradeID]
  for k, v := range tempMap2 {
    newPrice := getPrice(k)
    fmt.Println("New Price :: ", newPrice)
    oldPrice := v.pricePerStock
    fmt.Println("Old Price :: ", oldPrice)
    diffInPrice := newPrice - oldPrice
    fmt.Println("Difference in Price :: ", diffInPrice)

    tradeIDResp.symbol = append(tradeIDResp.symbol, k)
    tradeIDResp.currentValue = append(tradeIDResp.currentValue, newPrice)
    tradeIDResp.changeInValue = append(tradeIDResp.changeInValue, diffInPrice)
    tradeIDResp.numOfStock = append(tradeIDResp.numOfStock, v.numOfStock)
    tradeIDResp.unvestedAmount = append(tradeIDResp.unvestedAmount, v.unvestedAmount)
  }

  return nil
}

  func main() {
    yFin := new(YahooFinance)
    server := rpc.NewServer()
    server.Register(yFin)

    conn, err := net.Listen("tcp", ":1111")
    if err != nil {
      fmt.Errorf("Error while listening on port 1111 :: ", err)
    }

    for {
      conn1, err := conn.Accept()
      if err != nil {
        fmt.Errorf("Error occurred :: ", err)
      }

      go server.ServeCodec(jsonrpc.NewServerCodec(conn1))
    }
  }
