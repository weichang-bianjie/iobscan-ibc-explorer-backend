package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/oneoff"
	"github.com/sirupsen/logrus"
)

type AddTransferDataTask struct {
}

var _ OneOffTask = new(AddTransferDataTask)

func (t *AddTransferDataTask) Name() string {
	return "add_transfer_data_task"
}

func (t *AddTransferDataTask) Switch() bool {
	return global.Config.Task.SwitchAddTransferDataTask
}

func (t *AddTransferDataTask) Run() int {
	return 1
}

func (t *AddTransferDataTask) RunWithParam(chainsStr string, height int64) int {
	return t.handle(chainsStr, height)
}

func (t *AddTransferDataTask) handle(chain string, height int64) int {
	chainMap, err := getAllChainMap()
	if err != nil {
		return -1
	}
	w := &syncTransferTxWorker{
		chainMap: chainMap,
	}

	logrus.Info("start handle chain:", chain)
	for {
		curH, size, err := t.DoChain(w, chain, height, defaultMaxHandlerTx)
		if err != nil {
			_ = txNewRepo.SaveTransfer([]*oneoff.TxNew{
				{Height: curH, ChainId: chain},
			})
			logrus.Error(err.Error())
			return -1
		}
		if size < defaultMaxHandlerTx {
			_ = txNewRepo.SaveTransfer([]*oneoff.TxNew{
				{Height: curH, ChainId: chain},
			})
			logrus.Info("finish handle chain:", chain)
			break
		}
	}

	return 1
}

func (t *AddTransferDataTask) DoChain(w *syncTransferTxWorker, chainId string, height, limit int64) (int64, int, error) {
	maxHeight := int64(-1)
	denomMap, err := w.getChainDenomMap(chainId)
	if err != nil {
		return maxHeight, 0, err
	}

	txList, err := w.getTxList(chainId, height, limit)
	if err != nil {
		return maxHeight, 0, err
	}
	total := len(txList)
	if err := t.handleChain(chainId, w, txList, denomMap); err != nil {
		return maxHeight, 0, err
	}
	if len(txList) > 0 {
		maxHeight = txList[len(txList)-1].Height
	}
	return maxHeight, total, nil
}

func (t *AddTransferDataTask) handleChain(chainId string, w *syncTransferTxWorker, txList []*entity.Tx, denomMap map[string]*entity.IBCDenom) error {
	if len(txList) == 0 {
		return nil
	}

	ibcTxList, ibcDenomList := w.handleSourceTx(chainId, txList, denomMap)
	if len(ibcDenomList) > 0 {
		if err := denomRepo.InsertBatch(ibcDenomList); err != nil {
			logrus.Errorf("task %s worker %s denomRepo.InsertBatch %s error, %v", w.taskName, w.workerName, chainId, err)
			return err
		}
	}
	if len(ibcTxList) > 0 {
		if err := ibcTxRepo.InsertBatch(ibcTxList); err != nil {
			logrus.Errorf("task %s worker %s ibcTxRepo.InsertBatch %s error, %v", w.taskName, w.workerName, chainId, err)
			return err
		}
	}
	return nil
}
