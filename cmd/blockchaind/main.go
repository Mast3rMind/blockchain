package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ipkg/blockchain"
)

var (
	cfg = blockchain.DefaultConfig()
)

func init() {
	address := flag.String("addr", fmt.Sprintf("127.0.0.1:%d", blockchain.BLOCKCHAIN_DEFAULT_PORT), "Public facing ip address")
	dataDir := flag.String("data-dir", "blockchain-data", "Data directory")
	seedList := flag.String("seeds", "", "List of seed nodes")

	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg.SeedNodes = []string{}
	seeds := strings.Split(*seedList, ",")
	for _, s := range seeds {
		if s = strings.TrimSpace(s); s != "" {
			cfg.SeedNodes = append(cfg.SeedNodes, s)
		}
	}

	if *dataDir != cfg.DataDir {
		cfg.DataDir = *dataDir
	}
	if *address != cfg.BindAddr {
		cfg.BindAddr = *address
	}
}

func main() {
	core := blockchain.NewCore(cfg)
	core.Start()
	//
	// Read a line from stdin and submit it as a transaction
	//
	for {
		str := <-ReadStdin()
		tx := core.SubmitTransaction([]byte(str))
		log.Println("Submitted", tx.String())
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
