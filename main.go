package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/kehiy/pactatus/client"
	"github.com/pactus-project/pactus/util"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

var rpcNodes = []string{"181.214.208.165:50051", "bootstrap1.pactus.org:50051", "bootstrap2.pactus.org:50051", "bootstrap3.pactus.org:50051", "bootstrap4.pactus.org:50051", "151.115.110.114:50051", "188.121.116.247:50051"}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cmgr := client.NewClientMgr(ctx)

	for _, rn := range rpcNodes {
		c, e := client.NewClient(rn)
		if e != nil {
			continue
		}
		cmgr.AddClient(c)
	}

	b, err := bot.New("6632503496:AAEr6zh6btUazt74T5xe9UI_gi2A31MBD10", nil)
	if err != nil {
		panic(err)
	}

	go PostUpdates(ctx, b, cmgr)

	b.Start(ctx)
}

func PostUpdates(ctx context.Context, b *bot.Bot, cmgr *client.Mgr) {
	time.Sleep(10 * time.Second)

	for {
		status, lbt, lbh, td := networkHealth(cmgr)
		bi, err := cmgr.GetBlockchainInfo()
		if err != nil {
			panic(err)
		}

		cs, err := cmgr.GetCirculatingSupply()
		if err != nil {
			panic(err)
		}

		msg := makeMessage(bi, cs, td, status, lbt, lbh)
		b.SendMessage(ctx, makeMessageParams(msg))

		time.Sleep(5 * time.Minute)
	}
}

func makeMessage(b *pactus.GetBlockchainInfoResponse, c, timeDiff int64, status, lastBlkTime string, lastBlkH uint32) string {
	var s strings.Builder

	s.WriteString("Pactus Network Status Update 🔴\n")
	s.WriteString("Blockchain Info\n")
	s.WriteString(fmt.Sprintf("%s is Last Block Height⛓️\n", formatNumber(int64(lastBlkH))))
	s.WriteString(fmt.Sprintf("%v Active Accounts👤\n", b.TotalAccounts))
	s.WriteString(fmt.Sprintf("%v Total Validators🕵️\n", b.TotalValidators))
	s.WriteString(fmt.Sprintf("%v Total PAC Staked (network power)🦾\n", formatNumber(int64(util.ChangeToCoin(b.TotalPower)))))
	s.WriteString(fmt.Sprintf("%v Committee Power🦾\n", formatNumber(int64(util.ChangeToCoin(b.CommitteePower)))))
	s.WriteString(fmt.Sprintf("%v PAC is in Circulating🔄\n\n", formatNumber(int64(util.ChangeToCoin(c)))))

	s.WriteString("Network Status🧑🏻‍⚕️\n")
	s.WriteString(fmt.Sprintf("Network is %s\n", status))
	s.WriteString(fmt.Sprintf("%s is Last Block Time\n", lastBlkTime))
	s.WriteString(fmt.Sprintf("%v Time Difference\n", timeDiff))

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

func makeMessageParams(t string) *bot.SendMessageParams {
	return &bot.SendMessageParams{
		ChatID: "@pactatus",
		Text:   t,
	}
}

func formatNumber(num int64) string {
	numStr := strconv.FormatInt(num, 10)

	var formattedNum string
	for i, c := range numStr {
		if (i > 0) && (len(numStr)-i)%3 == 0 {
			formattedNum += ","
		}
		formattedNum += string(c)
	}

	return formattedNum
}
