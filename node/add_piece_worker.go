package node

import (
	"bytes"
	"context"
	"crypto/rand"
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/ipfs/go-merkledag"

	"github.com/filecoin-project/go-filecoin/proofs/sectorbuilder"
	go_sectorbuilder "github.com/filecoin-project/go-sectorbuilder"
)

var addPieceWorker AddPieceWorker

type AddPieceWorker struct {
	node *Node
	wg   sync.WaitGroup
}

func (w *AddPieceWorker) Start(ctx context.Context, node *Node) error {
	ok, err := strconv.ParseBool(os.Getenv("FIL_ADD_PIECE_WORKER"))
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	w.node = node

	maxBytes := go_sectorbuilder.GetMaxUserBytesPerStagedSector(node.sectorBuilder.(*sectorbuilder.RustSectorBuilder).SectorClass.SectorSize().Uint64())
	pieceData := make([]byte, maxBytes)
	_, err = io.ReadFull(rand.Reader, pieceData)
	if err != nil {
		return err
	}
	data := merkledag.NewRawNode(pieceData)

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
				allStaged, err := node.sectorBuilder.GetAllStagedSectors()
				if err != nil {
					log.Errorf("failed to get all staged sectors: %s", err)
					continue
				}
				if len(allStaged) > 1 {
					continue
				}

				sectorID, err := node.sectorBuilder.AddPiece(ctx, data.Cid(), uint64(len(pieceData)), bytes.NewReader(pieceData))
				if err != nil {
					log.Errorf("failed to add piece: %s", err)
					continue
				}
				log.Infof("added piece to sector %d", sectorID)
			}
		}
	}()

	return nil
}

func (w *AddPieceWorker) Stop() {
	w.wg.Wait()
}
