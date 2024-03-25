package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PACZone/pactus-status/client"
	"github.com/PACZone/pactus-status/config"
	"github.com/PACZone/pactus-status/utils"
	"github.com/go-telegram/bot"
	"github.com/pactus-project/pactus/util"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type StatusChecker struct {
	ctx   context.Context
	cfg   config.Config
	cmgr  *client.Mgr
	tgbot *bot.Bot
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Panic(err)
	}

	ctx := context.Background()

	log.Println("starting")

	cmgr := client.NewClientMgr(ctx)

	for _, rn := range cfg.RPCNodes {
		c, e := client.NewClient(rn)
		if e != nil {
			log.Printf("error: %v adding client %s\n", e, rn)
			continue
		}
		cmgr.AddClient(*c)
		log.Printf("client added %s\n", rn)
	}

	b, err := bot.New(cfg.BotToken, bot.WithAllowedUpdates(bot.AllowedUpdates{}))
	if err != nil {
		panic(err)
	}

	sc := StatusChecker{
		ctx:   ctx,
		cfg:   cfg,
		cmgr:  cmgr,
		tgbot: b,
	}

	go sc.postUpdates()

	b.Start(ctx)
}

func (sc *StatusChecker) postUpdates() {
	for {
		log.Println("posting new update!")
		status, lbt, lbh, td := networkHealth(sc.cmgr)
		bi, err := sc.cmgr.GetBlockchainInfo()
		if err != nil {
			panic(err)
		}

		log.Println("got network health and Blockcahin info successfully")

		cs, err := sc.cmgr.GetCirculatingSupply()
		if err != nil {
			panic(err)
		}

		log.Println("got circ supply successfully")

		price := utils.GetPACPrice(sc.cfg.PriceAPI)
		log.Println("got price successfully")

		msg := makeMessage(bi, cs, td, status, lbt, price, lbh)
		_, err = sc.tgbot.EditMessageText(sc.ctx, utils.MakeMessageParams(msg, 37))
		if err != nil {
			log.Printf("can't post updates: %v\n", err)
		}
		log.Println("updated posted successfully")

		time.Sleep(7 * time.Second)
	}
}

func makeMessage(b *pactus.GetBlockchainInfoResponse, c, timeDiff int64, status, lastBlkTime string, price float64, lastBlkH uint32) string {
	var s strings.Builder

	mcap := float64(util.ChangeToCoin(c)) * price
	fdv := float64(util.ChangeToCoin(c+b.TotalPower)) * price
	tvl := float64(util.ChangeToCoin(b.TotalPower)) * price

	s.WriteString("🟢 Pactus Network Status Update\n\n")
	s.WriteString(fmt.Sprintf("⛓️ %s Last Block Height\n\n", utils.FormatNumber(int64(lastBlkH))))
	s.WriteString(fmt.Sprintf("👤 %v Accounts\n\n", utils.FormatNumber(int64(b.TotalAccounts))))
	s.WriteString(fmt.Sprintf("🕵️ %v Validators\n\n", utils.FormatNumber(int64(b.TotalValidators))))
	s.WriteString(fmt.Sprintf("🦾 %v PAC Staked\n\n", utils.FormatNumber(int64(util.ChangeToCoin(b.TotalPower)))))
	s.WriteString(fmt.Sprintf("🦾 %v PAC Committee Power\n\n", utils.FormatNumber(int64(util.ChangeToCoin(b.CommitteePower)))))
	s.WriteString(fmt.Sprintf("🔄 %v PAC Circulating Supply\n\n", utils.FormatNumber(int64(util.ChangeToCoin(c)))))
	s.WriteString(fmt.Sprintf("🪙 %v PAC Total Supply\n\n", utils.FormatNumber(int64(util.ChangeToCoin(c+b.TotalPower)))))

	s.WriteString(fmt.Sprintf("📊 %v$ Market Cap\n\n", utils.FormatNumber(int64(mcap))))
	s.WriteString(fmt.Sprintf("💹 %v$ Fully Diluted Value (FDV)\n\n", utils.FormatNumber(int64(fdv))))
	s.WriteString(fmt.Sprintf("🔒 %v$ Total Value Locked (TVL)\n\n", utils.FormatNumber(int64(tvl))))

	s.WriteString(fmt.Sprintf("📈 Exbitron Price %v$ \n\n", price))

	s.WriteString(fmt.Sprintf("Network is %s\n%s is The LastBlock time and there is %v seconds passed from last block", status, lastBlkTime, timeDiff))

	return s.String()
}

func networkHealth(cmgr *client.Mgr) (string, string, uint32, int64) {
	lastBlockTime, lastBlockHeight := cmgr.GetLastBlockTime()
	lastBlockTimeFormatted := time.Unix(int64(lastBlockTime), 0).Format("02/01/2006, 15:04:05")
	currentTime := time.Now()

	timeDiff := (currentTime.Unix() - int64(lastBlockTime))

	healthStatus := true
	if timeDiff > 15 {
		healthStatus = false
	}

	var status string
	if healthStatus {
		status = "Healthy✅"
	} else {
		status = "UnHealthy❌"
	}

	return status, lastBlockTimeFormatted, lastBlockHeight, timeDiff
}
