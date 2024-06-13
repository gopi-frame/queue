package service

import (
	"github.com/gopi-frame/contract/config"
	"github.com/gopi-frame/contract/foundation"
	queuecontract "github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/queue"
)

type Provider struct {
	app foundation.Application
}

func (provider *Provider) Init(app foundation.Application) {
	provider.app = app
}

func (provider *Provider) Register() {
	provider.app.Bind("queue", func() any {
		return queue.NewSupervisor()
	})
}

func (provider *Provider) Boot() {
	supervisor := provider.app.Get("queue").(*queue.Supervisor)
	queueConfigRepo := provider.app.Get("config").(config.Repository)
	defaultQueue := queueConfigRepo.Get("default.queue").(string)
	failed := queueConfigRepo.Get("failed").(map[string]any)
	failedJobProvider := queue.OpenFailedJobProvider(failed["driver"].(string), failed)
	supervisor.SetFailedJobProvider(failedJobProvider)
	supervisor.SetDefaultQueue(defaultQueue)
	for name, config := range queueConfigRepo.Get("queues").(map[string]any) {
		config := config.(map[string]any)
		supervisor.Bind(name, func() queuecontract.Dispatcher {
			dispatcher := queue.OpenQueue(config["driver"].(string), config)
			dispatcher.FailedJobProvider(failedJobProvider)
			return dispatcher
		})
	}
}
