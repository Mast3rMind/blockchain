package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ipkg/blockchain/core"
)

var (
	cfg = &core.Config{}
)

func parseSeeds(slist string) []string {
	out := []string{}
	for _, s := range strings.Split(slist, ",") {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}

func init() {
	seedList := flag.String("seeds", "", "Seed nodes")
	flag.StringVar(&cfg.Addr, "ip", fmt.Sprintf("127.0.0.1:%s", core.BLOCKCHAIN_PORT), "Public facing ip address")
	flag.StringVar(&cfg.DataDir, "data-dir", "blockchain-data", "Data directory")
	flag.Parse()

	cfg.Seeds = parseSeeds(*seedList)
}

func main() {

	core.Start(cfg)

	for {
		str := <-ReadStdin()
		core.Core.Blockchain.TransactionsQueue <- core.CreateTransaction(str)
	}
}

func ReadStdin() chan string {

	cb := make(chan string)
	sc := bufio.NewScanner(os.Stdin)

	go func() {
		if sc.Scan() {
			cb <- sc.Text()
		}
	}()

	return cb
}
