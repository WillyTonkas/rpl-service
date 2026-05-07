package services

import "rpl-service/repositories"

// Intermediate layer. Calls repositories in order to interact with the database.
// Has most of the business logic.
type Service[T any] struct {
	Repository repositories.Repository[T]
}
