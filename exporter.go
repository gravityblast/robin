package robin

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type exportJob struct {
	exporter Exporter
	item     Item
}

type exportQueue struct {
	sync.WaitGroup
	jobs chan *exportJob
}

func newExportQueue() *exportQueue {
	return &exportQueue{
		jobs: make(chan *exportJob),
	}
}

func (eq *exportQueue) push(j *exportJob) {
	eq.jobs <- j
}

func (eq *exportQueue) close() {
	close(eq.jobs)
}

func (eq *exportQueue) run(n int) {
	for i := 0; i < n; i++ {
		go func() {
			eq.Add(1)
			eq.runWorker()
		}()
	}
}

func (eq *exportQueue) runWorker() {
	for job := range eq.jobs {
		job.exporter.Export(job.item)
	}
	eq.Done()
}

type Exporter interface {
	Export(Item)
}

type stdoutExporter struct {
	encoder *json.Encoder
}

func NewStdoutExporter() Exporter {
	return &stdoutExporter{
		encoder: json.NewEncoder(os.Stdout),
	}
}

func (e *stdoutExporter) Export(item Item) {
	err := e.encoder.Encode(item)
	if err != nil {
		fmt.Printf("%s, item: %+v", err, item)
	}
}
