package main

import (
    "context"
    "fmt"
    "log"
    // "strings"
    "math/big"
    // "math"

    // "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/crypto"
    // "github.com/ethereum/go-ethereum"
    // "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"

    "subscriber/generated/wyvern_exchange"
    "subscriber/generated/weth9"
)

// func printLog(vLog types.Log){
    // fmt.Println(vLog.BlockHash.Hex())
    // fmt.Println(vLog.BlockNumber)
    // fmt.Println(vLog.TxHash.Hex())
    // // parse log OrdersMatched (bytes32 buyHash, bytes32 sellHash, index_topic_1 address maker, index_topic_2 address taker, uint256 price, index_topic_3 bytes32 metadata)

    // contractAbi, err := abi.JSON(strings.NewReader(string(wyvern_exchange.WyvernexchangeABI)))
    // event, err := contractAbi.Unpack("OrdersMatched", vLog.Data)
    // if err != nil{
        // log.Fatal(err)
    // }
    // // bind.NewBoundContract(address, contractAbi, )
    // wyvern_exchange.FilterOrdersMatched(opts, taker, []string{maker}, []string{taker}, metadta)

    // fmt.Println(event[0].BuyHash[:])
    // fmt.Println(event[0].SellHash[:])

    // var topics [4]string
    // for i:= range vLog.Topics{
        // topics[i] = vLog.Topics[i].Hex()
    // }

    // fmt.Println(topics[0])
// }

func main(){
    rpcUri := "ws://localhost:8545"
    client, err := ethclient.Dial(rpcUri)
    if err != nil{
        log.Fatal(err)
    }

    // headers := make(chan *types.Header)
    // sub, err := client.SubscribeNewHead(context.Background(), headers);
    // if err!= nil{
        // log.Fatal(err)
    // }

    // for {
        // select{
        // case err:=<-sub.Err():
            // log.Fatal(err)
        // case header:=<-headers:
            // fmt.Println(header.Hash().Hex())
            // block, err := client.BlockByHash(context.Background(), header.Hash())
            // if err!=nil{
                // log.Fatal(err)
            // }

            // // print block
            // // fmt.Println(block.Hash().Hex())
            // // fmt.Println(block.Number().Uint64())
            // // fmt.Println(block.Time())
            // // fmt.Println(block.Nonce())
            // // fmt.Println(len(block.Transactions()))
            // for _, tx := range block.Transactions(){
                // fmt.Println(tx.Hash().Hex())
                // receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
                // if err!=nil{
                    // log.Fatal(err)
                // }
                // fmt.Println(receipt.Status)
                // fmt.Println(receipt.Logs)
                // // fmt.Println(receipt)
            // }
        // }
    // }
    contractAddress := common.HexToAddress("0x7Be8076f4EA4A4AD08075C2508e481d6C946D12b")
    tokenAddress := common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
    // query := ethereum.FilterQuery{
        // Addresses: []common.Address{contractAddress},
    // }

    // logs := make(chan types.Log)
    // sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
    // if err != nil{
        // log.Fatal(err)
    // }
    // for {
        // select {
        // case err :=<-sub.Err():
            // log.Fatal(err)
        // case vLog:= <-logs:
            // fmt.Println(vLog)
        // }
    // }
    wyvernExchange, err := wyvern_exchange.NewWyvernexchange(contractAddress, client)
    if err!=nil{
        log.Fatal(err)
    }

    WETH9, err := weth9.NewWeth9(tokenAddress, client)
    sink := make(chan *weth9.Weth9Deposit)
    sub, err := WETH9.Weth9Filterer.WatchDeposit(nil, sink, []common.Address{})

    // call
    name, err := wyvernExchange.Name(nil)
    if err!=nil{
        log.Fatal(err)
    }

    fmt.Println(name)

    // send
    key:= "057868c4074f55bb4346bda4b855f74e1a3bb2db96a0fdffb51c688ed1df1c0b"
    privateKey, err := crypto.HexToECDSA(key)
    if err!=nil{
        log.Fatal(err)
    }
    auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1))

    if err!=nil{
        log.Fatal(err)
    }

    session := &wyvern_exchange.WyvernexchangeSession{Contract: wyvernExchange, CallOpts: bind.CallOpts{
        Pending: true,
    }, TransactOpts: bind.TransactOpts{From: auth.From, Signer: auth.Signer}}
    // WETH9Session := &weth9.WETH9Session{Contract: WETH9, CallOpts:bind.CallOpts{pending:true,}, TransactOpts: bind.TransactOpts{From: auth.From, Signer: auth.Signer}}
    fmt.Println(session.Name())

    auth.Value = big.NewInt(100000000000000)
    fmt.Println(WETH9.Weth9Caller.Name(nil))
    // tx, err:=WETH9.Weth9Transactor.Deposit(auth)

    // if err!=nil{
        // log.Fatal(err)
    // }
    // fmt.Println(tx.Hash())

    // receipt, err:= client.TransactionReceipt(context.Background(), tx.Hash())
    // if err!=nil{
        // log.Fatal(err)
    // }
    // fmt.Println(receipt.Logs)
    // fmt.Println(receipt.Status)



    // eth balance
    ethBalance, err:= client.BalanceAt(context.Background(), auth.From, nil)
    fmt.Println("eth balance: ", ethBalance)

    // weth balance
    wethBalance, err:= WETH9.Weth9Caller.BalanceOf(nil, auth.From)
    fmt.Println("weth balance: ", wethBalance)

    defer sub.Unsubscribe()
    for {
        select{
        case err:=<-sub.Err():
            log.Fatal(err)
        case event:=<-sink:
            if err!=nil{
                log.Fatal(err)
            }
            fmt.Println(event.Dst)
            fmt.Println(event.Wad)
            fmt.Println(event.Raw.TxHash)
            fmt.Println(event.Raw.BlockHash)
            fmt.Println(event.Raw.BlockNumber)
        }
    }
}
