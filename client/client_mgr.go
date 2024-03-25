package client

import (
	"context"
	"errors"
	"sync"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type Mgr struct {
	valMapLock sync.RWMutex
	valMap     map[string]*pactus.PeerInfo

	ctx     context.Context
	clients []Client
}

func NewClientMgr(ctx context.Context) *Mgr {
	return &Mgr{
		clients:    make([]Client, 0),
		valMap:     make(map[string]*pactus.PeerInfo),
		valMapLock: sync.RWMutex{},
		ctx:        ctx,
	}
}

// AddClient should call before Start.
func (cm *Mgr) AddClient(c Client) {
	cm.clients = append(cm.clients, c)
}

// NOTE: local client is always the first client.
func (cm *Mgr) getLocalClient() *Client {
	return &cm.clients[0]
}

func (cm *Mgr) GetRandomClient() Client {
	for _, c := range cm.clients {
		return c
	}

	return Client{}
}

func (cm *Mgr) GetBlockchainInfo() (*pactus.GetBlockchainInfoResponse, error) {
	for _, c := range cm.clients {
		info, err := c.GetBlockchainInfo(cm.ctx)
		if err != nil {
			continue
		}
		return info, nil
	}

	return nil, errors.New("can't get blockchain info")
}

func (cm *Mgr) GetLastBlockTime() (uint32, uint32) {
	for _, c := range cm.clients {
		lastBlockTime, lastBlockHeight, err := c.LastBlockTime(cm.ctx)
		if err != nil {
			continue
		}

		return lastBlockTime, lastBlockHeight
	}

	return 0, 0
}

func (cm *Mgr) GetNetworkInfo() (*pactus.GetNetworkInfoResponse, error) {
	for _, c := range cm.clients {
		info, err := c.GetNetworkInfo(cm.ctx)
		if err != nil {
			continue
		}
		return info, nil
	}

	return nil, errors.New("unable to get network info")
}

func (cm *Mgr) GetBalance(addr string) (int64, error) {
	for _, c := range cm.clients {
		b, err := c.GetBalance(cm.ctx, addr)
		if err != nil {
			continue
		}

		return b, nil
	}

	return 0, errors.New("can't get balance")
}

func (cm *Mgr) GetCirculatingSupply() (int64, error) {
	height, err := cm.GetBlockchainInfo()
	if err != nil {
		return 0, err
	}
	minted := float64(height.LastBlockHeight) * 1e9
	staked := height.TotalPower
	warm := int64(630_000_000_000_000)

	var addr1Out int64 = 0
	var addr2Out int64 = 0
	var addr3Out int64 = 0
	var addr4Out int64 = 0
	var addr5Out int64 = 0 // warm wallet
	var addr6Out int64 = 0 // warm wallet

	balance1, err := cm.GetBalance("pc1z2r0fmu8sg2ffa0tgrr08gnefcxl2kq7wvquf8z")
	if err == nil {
		addr1Out = 8_400_000_000_000_000 - balance1
	}

	balance2, err := cm.GetBalance("pc1zprhnvcsy3pthekdcu28cw8muw4f432hkwgfasv")
	if err == nil {
		addr2Out = 6_300_000_000_000_000 - balance2
	}

	balance3, err := cm.GetBalance("pc1znn2qxsugfrt7j4608zvtnxf8dnz8skrxguyf45")
	if err == nil {
		addr3Out = 4_200_000_000_000_000 - balance3
	}

	balance4, err := cm.GetBalance("pc1zs64vdggjcshumjwzaskhfn0j9gfpkvche3kxd3")
	if err == nil {
		addr4Out = 2_100_000_000_000_000 - balance4
	}

	balance5, err := cm.GetBalance("pc1zuavu4sjcxcx9zsl8rlwwx0amnl94sp0el3u37g")
	if err == nil {
		addr5Out = 420_000_000_000_000 - balance5
	}

	balance6, err := cm.GetBalance("pc1zf0gyc4kxlfsvu64pheqzmk8r9eyzxqvxlk6s6t")
	if err == nil {
		addr6Out = 210_000_000_000_000 - balance6
	}

	circulating := (addr1Out + addr2Out + addr3Out + addr4Out + addr5Out + addr6Out + int64(minted)) - staked - warm
	return circulating, nil
}
