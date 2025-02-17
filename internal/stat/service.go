package stat

import (
	"log"
	"short-link/pkg/event"
)

type ServiceDeps struct {
	EventBus       *event.EventBus
	StatRepository *StatRepository
}

type Service struct {
	EventBus       *event.EventBus
	StatRepository *StatRepository
}

func NewStatService(deps *ServiceDeps) *Service {
	return &Service{
		EventBus:       deps.EventBus,
		StatRepository: deps.StatRepository,
	}
}

func (s *Service) AddClick() {
	for msg := range s.EventBus.Subscribe() {
		switch msg.Type {
		case event.LinkVisitedEvent:
			id, ok := msg.Data.(uint)
			if !ok {
				log.Println("Error casting link visited event: ", msg.Data)
				continue
			}
			s.StatRepository.AddClick(id)
		}
	}
}
