package celery

import "context"

type Client struct {
	broker     Broker
	backend    Backend
	workerPool *WorkerPool
}

func NewCeleryClient(
	broker Broker,
	backend Backend,
	numWorkers int,
) (*Client, error) {
	return &Client{
		broker:     broker,
		backend:    backend,
		workerPool: NewWorkerPool(broker, backend, numWorkers),
	}, nil
}

func (c *Client) Register(name string, task Task) {
	c.workerPool.RegisterTask(name, task)
}

func (c *Client) StartWorkers(ctx context.Context) {
	c.workerPool.Start(ctx)
}

func (c *Client) StopWorkers() {
	c.workerPool.Stop()
}
