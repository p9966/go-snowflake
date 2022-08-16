package gosnowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	workerIDBits     = uint64(5)
	dataCenterIDBits = uint64(5)
	sequenceBits     = uint64(12)

	maxWorkerID     = int64(-1) ^ (int64(-1) << workerIDBits)
	maxDataCenterID = int64(-1) ^ (int64(-1) << dataCenterIDBits)
	maxSequence     = int64(-1) ^ (int64(-1) << sequenceBits)

	timeLeft = uint8(22)
	dataLeft = uint8(17)
	workLeft = uint8(12)
)

var startTime int64 = 1645539742000

type Worker struct {
	mu           sync.Mutex
	LastStamp    int64
	WorkerID     int64
	DataCenterID int64
	Sequence     int64
}

func SetStartTime(s int64) {
	startTime = s
}

func NewWorker(workerID, dataCenterID int64) *Worker {
	return &Worker{
		WorkerID:     workerID,
		LastStamp:    0,
		Sequence:     0,
		DataCenterID: dataCenterID,
	}
}

func (w *Worker) NextID() (uint64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.nextID()
}

func (w *Worker) getMilliSeconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func (w *Worker) nextID() (uint64, error) {
	timeStamp := w.getMilliSeconds()
	if timeStamp < w.LastStamp {
		return 0, errors.New("time is moving backwards,waiting until")
	}

	if w.LastStamp == timeStamp {

		w.Sequence = (w.Sequence + 1) & maxSequence

		if w.Sequence == 0 {
			for timeStamp <= w.LastStamp {
				timeStamp = w.getMilliSeconds()
			}
		}
	} else {
		w.Sequence = 0
	}

	w.LastStamp = timeStamp
	id := ((timeStamp - startTime) << timeLeft) |
		(w.DataCenterID << dataLeft) |
		(w.WorkerID << workLeft) |
		w.Sequence

	return uint64(id), nil
}
