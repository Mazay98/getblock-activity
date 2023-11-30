package blockio

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"

	"getBlock/internal/config"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/metachris/eth-go-bindings/erc20"
	"github.com/ofen/getblock-go/eth"
	"go.uber.org/zap"
)

type Service struct {
	log    *zap.Logger
	client *ethclient.Client
	abi    *abi.ABI
	cfg    *config.Blockio
}
type aAp struct {
	address  string
	activity int
}

func New(ctx context.Context, log *zap.Logger, cfg *config.Blockio) (*Service, error) {
	conn, err := ethclient.DialContext(ctx, fmt.Sprintf("%s/%s", eth.Endpoint, cfg.Token))
	if err != nil {
		return nil, err
	}
	ABI, err := abi.JSON(strings.NewReader(erc20.Erc20ABI))
	if err != nil {
		return nil, err
	}

	return &Service{
		log:    log,
		client: conn,
		abi:    &ABI,
		cfg:    cfg,
	}, nil
}

// GetTopActivity top activity from blocks.
func (s *Service) GetTopActivity(ctx context.Context, positions int) {
	e := eth.New(s.cfg.Token)
	b, err := s.getLastBlock(ctx, e)

	if err != nil {
		s.log.Error("Failed get last block", zap.Error(err))
	}

	tr := make(chan string, 10000)
	defer close(tr)

	wg := sync.WaitGroup{}
	wg.Add(int(s.cfg.Blocks))

	addressTransactions := make(map[string]int)
	go func(tr <-chan string) {
		for t := range tr {
			addressTransactions[t]++
			wg.Done()
		}
	}(tr)

	for i := b.Int64(); i >= b.Int64()-s.cfg.Blocks; i-- {
		go s.parseBlockInfo(ctx, &wg, e, big.NewInt(i), tr)
	}
	wg.Wait()

	var addressActivityPairs []aAp
	for address, activity := range addressTransactions {
		addressActivityPairs = append(addressActivityPairs, aAp{address, activity})
	}
	sort.Slice(addressActivityPairs, func(i, j int) bool {
		return addressActivityPairs[i].activity > addressActivityPairs[j].activity
	})

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "№\tAddress\tActivity")

	for i := 0; i < positions; i++ {
		fmt.Fprintln(w, fmt.Sprintf("%d\t%s\t%d", i+1, addressActivityPairs[i].address, addressActivityPairs[i].activity))
	}

	w.Flush() //nolint:errcheck
}

func (s *Service) getLastBlock(ctx context.Context, ec *eth.Client) (*big.Int, error) {
	b, err := ec.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) getBlockInfo(ctx context.Context, ec *eth.Client, block *big.Int) (*eth.Block, error) {
	ib, err := ec.GetBlockByNumber(ctx, block, true)
	if err != nil {
		return nil, err
	}

	return ib, nil
}

func (s *Service) checkTransaction(tr *eth.Transaction) bool {
	if len(tr.Input) < 10 {
		return false
	}
	address := common.HexToAddress(tr.To)
	c, err := erc20.NewErc20(address, s.client)
	if err != nil {
		return false
	}

	_, err = c.Symbol(nil)
	if err != nil {
		return false
	}

	decodedSig, err := hex.DecodeString(tr.Input[2:10])
	if err != nil {
		return false
	}

	method, err := s.abi.MethodById(decodedSig)
	if err != nil {
		return false
	}

	if method.Name != "transfer" {
		return false
	}

	return true
}

func (s *Service) parseBlockInfo(ctx context.Context, wg *sync.WaitGroup, ec *eth.Client, b *big.Int, tr chan<- string) {
	defer wg.Done()

	bd, err := s.getBlockInfo(ctx, ec, b)
	s.log.Info("Check block", zap.String("id", b.String()))

	if err != nil {
		//todo push to retry if error , now possible < s.cfg.Blocks blocks
		s.log.Error("Failed get block info", zap.Error(err), zap.String("block", b.String()))
		return
	}
	for _, transaction := range bd.Transactions {
		ok := s.checkTransaction(&transaction)
		if ok {
			wg.Add(1)
			tr <- transaction.From
		}
	}
}
