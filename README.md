# cmpe273-assignment1

To build a virtual trading stock system using Yahoo Finance API

1. Buy stocks - user will invest some budget on N number of stocks in varying percentages. The percentage should sum up to 100. The budget will be allocated to the stocks according to the percentage. When the user provides input, a trade ID is returned, confirming the user's transaction.
2. Check Portfolio - With returned trade ID in previous step, user can compare the loss/gain on each stock which he/she has invested.


Achieved Outcome:
1. Able to successfully generate trade ID on the server side, along with number of stocks and unvested amount for each stock. The total unvested amount also is returned.
2. With returned tradeID, the portfolio information was also returned.

To run:
go run yahFinServer.go
1. go run yahFinClient.go FB:50,IBM:50 5000
Returns tradeID xxxx
2. go run yahFinClient.go xxxx
Returns portfolio information
