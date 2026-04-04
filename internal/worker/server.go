package worker

import (
	"log"

	"github.com/hibiken/asynq"
)

type Worker struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

func NewWorker(redisAddr, redisPassword string, concurrency int) *Worker {
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr, Password: redisPassword},
		asynq.Config{
			Concurrency: concurrency,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	return &Worker{
		server: server,
		mux:    asynq.NewServeMux(),
	}
}

func (w *Worker) RegisterHandler(taskType string, handler asynq.HandlerFunc) {
	w.mux.HandleFunc(taskType, handler)
}

func (w *Worker) Start() error {
	log.Printf("Starting asynq worker server...")
	return w.server.Start(w.mux)
}

func (w *Worker) Stop() {
	log.Printf("Stopping asynq worker server...")
	w.server.Shutdown()
}
