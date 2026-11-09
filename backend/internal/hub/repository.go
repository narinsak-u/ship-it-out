// Package hub provides the hub service, HTTP handler, and focused persistence contract.
package hub

import "context"

type Repository interface {
	FindAll(context.Context) ([]Hub, error)
	FindByID(context.Context, string) (*Hub, error)
	Create(context.Context, *Hub) error
	Save(context.Context, *Hub) error
	Delete(context.Context, string) error
}
