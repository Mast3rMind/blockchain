package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/anacrolix/utp"
	chord "github.com/euforia/go-chord"
	"github.com/ipkg/blockchain"
	"github.com/ipkg/go-mux"
)

// CFG is the global config.
var CFG = &config{
	DialTimeout: 6 * time.Second,
	RpcTimeout:  3 * time.Second,
	MaxConnIdle: 300 * time.Second,
}

type config struct {
	Chord *chord.Config

	BindAddr  string
	AdvAddr   string
	JoinAddrs string
	// http bind address for admin interface
	AdminAddr string

	DialTimeout time.Duration
	RpcTimeout  time.Duration
	MaxConnIdle time.Duration
	// how often to poll for txs in milliseconds
	TxPollMsec int
}

// Peers as provided in JoinAddrs
func (c *config) Peers() []string {
	out := []string{}
	for _, p := range strings.Split(c.JoinAddrs, ",") {
		h := strings.TrimSpace(p)
		if len(h) > 0 {
			out = append(out, h)
		}
	}

	return out
}

func (c *config) Bootstrap() bool {
	return c.JoinAddrs == ""
}

func initDeps(addr string) (*mux.Mux, *chord.Config, *chord.UTPTransport, error) {
	conf := chord.DefaultConfig(addr)
	conf.StabilizeMin = time.Duration(1 * time.Second)
	conf.StabilizeMax = time.Duration(5 * time.Second)

	ln, err := utp.NewSocket("udp", addr)
	if err != nil {
		return nil, nil, nil, err
	}

	mx := mux.NewMux(ln, ln.Addr())
	go mx.Serve()
	sock := mx.Listen(72)

	trans, err := chord.InitUTPTransport(sock, CFG.DialTimeout, CFG.RpcTimeout, CFG.MaxConnIdle)
	if err != nil {
		return nil, nil, nil, err
	}
	return mx, conf, trans, nil
}

func buildTx(kp *blockchain.ECDSAKeypair, prevHash, payload []byte) (*blockchain.Tx, error) {
	//tx := blockchain.NewTransaction(kp.Public, nil, payload)
	tx := blockchain.NewTx(prevHash, payload)
	//tx.Header.Nonce = tx.GenerateNonce(stx.TRANSACTION_POW)
	err := tx.Sign(kp)
	return tx, err
}

func readStdin() chan string {

	cb := make(chan string)
	sc := bufio.NewScanner(os.Stdin)

	go func() {
		if sc.Scan() {
			cb <- sc.Text()
		}
	}()

	return cb
}

func initChordRing(cfg *chord.Config, trans *chord.UTPTransport) (*chord.Ring, error) {
	if CFG.Bootstrap() {
		log.Println("[chord] Creating ring")
		return chord.Create(cfg, trans)
	}

	log.Println("[chord] Joining ring")
	peers := CFG.Peers()
	for _, peer := range peers {
		ring, err := chord.Join(cfg, trans, peer)
		if err == nil {
			log.Printf("[chord] Joined: %s", peer)
			return ring, err
		}
		log.Printf("[chord] Failed to join: %s", peer)
	}
	//return chord.Join(cfg, trans, CFG.JoinAddr)
	return nil, fmt.Errorf("all peers exhausted")
}

func init() {
	flag.StringVar(&CFG.JoinAddrs, "j", "", "Comma delimited list of existing peers")
	flag.StringVar(&CFG.AdminAddr, "a", "", "Admin HTTP bind address")
	flag.StringVar(&CFG.BindAddr, "b", "127.0.0.1:45454", "Bind address")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

}

func main() {

	// chord init
	mx, ccfg, ctrans, err := initDeps(CFG.BindAddr)
	if err != nil {
		log.Fatal(err)
	}

	ring, err := initChordRing(ccfg, ctrans)
	if err != nil {
		log.Fatal(err)
	}

	// signator
	kp, err := blockchain.GenerateECDSAKeypair()
	if err != nil {
		log.Fatal(err)
	}

	// store
	store := blockchain.NewInMemBlockStore()

	// blockchain transport
	btrans := blockchain.NewChordTransport(mx.Listen(73), ccfg, ring)

	chain, err := blockchain.NewBlockchain(kp, store, btrans, &stateMachine{}, CFG.Peers()...)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[blockchain] Starting")
	go chain.Start()

	// admin server
	go func() {
		if len(CFG.AdminAddr) == 0 {
			return
		}

		s := &server{store}
		log.Printf("Starting Admin server on %s", CFG.AdminAddr)
		if err := http.ListenAndServe(CFG.AdminAddr, s); err != nil {
			log.Println("ERR", err)
		}
	}()

	fmt.Println("Type something and hit enter")
	for {
		str := <-readStdin()

		if str == "chain.length" {
			fmt.Println(store.BlockCount())
			continue
		} else if str == "block.last" {
			fmt.Printf("%+v\n", store.LastBlock())
			continue
		}

		lth := store.LastTx().Hash()

		tx, err := buildTx(kp, lth, []byte(str))
		if err != nil {
			log.Println("ERR", err)
			continue
		}

		// TEMPORARY
		if CFG.Bootstrap() {
			chain.QueueTransactions(tx)
			continue
		}

		if err = btrans.BroadcastTransaction(tx); err != nil {
			log.Println("ERR", err)
			continue
		}

		//log.Printf("Submitted tx=%x", tx.Hash())
		//lth = store.LastTx().Hash()
	}
}
