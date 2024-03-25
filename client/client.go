package client

import (
	"context"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	blockchainClient  pactus.BlockchainClient
	networkClient     pactus.NetworkClient
	transactionClient pactus.TransactionClient
	conn              *grpc.ClientConn
}

func NewClient(endpoint string) (*Client, error) {
	conn, err := grpc.Dial(endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		blockchainClient:  pactus.NewBlockchainClient(conn),
		networkClient:     pactus.NewNetworkClient(conn),
		transactionClient: pactus.NewTransactionClient(conn),
		conn:              conn,
	}, nil
}

func (c *Client) GetBlockchainInfo(ctx context.Context) (*pactus.GetBlockchainInfoResponse, error) {
	blockchainInfo, err := c.blockchainClient.GetBlockchainInfo(ctx, &pactus.GetBlockchainInfoRequest{})
	if err != nil {
		return nil, err
	}
	return blockchainInfo, nil
}

func (c *Client) GetBlockchainHeight(ctx context.Context) (uint32, error) {
	blockchainInfo, err := c.blockchainClient.GetBlockchainInfo(ctx, &pactus.GetBlockchainInfoRequest{})
	if err != nil {
		return 0, err
	}
	return blockchainInfo.LastBlockHeight, nil
}

func (c *Client) GetNetworkInfo(ctx context.Context) (*pactus.GetNetworkInfoResponse, error) {
	networkInfo, err := c.networkClient.GetNetworkInfo(ctx, &pactus.GetNetworkInfoRequest{})
	if err != nil {
		return nil, err
	}

	return networkInfo, nil
}

func (c *Client) LastBlockTime(ctx context.Context) (uint32, uint32, error) {
	info, err := c.blockchainClient.GetBlockchainInfo(ctx, &pactus.GetBlockchainInfoRequest{})
	if err != nil {
		return 0, 0, err
	}

	lastBlockTime, err := c.blockchainClient.GetBlock(ctx, &pactus.GetBlockRequest{
		Height:    info.LastBlockHeight,
		Verbosity: pactus.BlockVerbosity_BLOCK_INFO,
	})

	return lastBlockTime.BlockTime, info.LastBlockHeight, err
}

func (c *Client) GetBalance(ctx context.Context, address string) (int64, error) {
	account, err := c.blockchainClient.GetAccount(ctx, &pactus.GetAccountRequest{
		Address: address,
	})
	if err != nil {
		return 0, err
	}

	return account.Account.Balance, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
