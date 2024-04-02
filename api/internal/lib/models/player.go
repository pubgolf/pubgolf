package models

import (
	"fmt"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

// Player contains queryable information about a player.
type Player struct {
	ID     PlayerID
	Name   string
	Events []EventRegistration
}

// PlayerParams contains writable information about a player.
type PlayerParams struct {
	Name string
}

// EventRegistration contains info about a player's relationship to an event.
type EventRegistration struct {
	EventKey        string
	ScoringCategory ScoringCategory
}

// Proto converts a Player model to a protobuf representation suitable for return from an RPC handler.
func (p Player) Proto() (*apiv1.Player, error) {
	events := make([]*apiv1.EventRegistration, len(p.Events))

	for i, e := range p.Events {
		cat, err := e.ScoringCategory.ProtoEnum()
		if err != nil {
			return nil, fmt.Errorf("convert scoring category: %w", err)
		}

		events[i] = &apiv1.EventRegistration{
			EventKey:        e.EventKey,
			ScoringCategory: cat,
		}
	}

	return &apiv1.Player{
		Id: p.ID.String(),
		Data: &apiv1.PlayerData{
			Name: p.Name,
		},
		Events: events,
	}, nil
}
