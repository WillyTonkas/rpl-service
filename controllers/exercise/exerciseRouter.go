package exercise

import (
	"rpl-service/models"
)

const BaseURL = "/exercise" // May change

var SolveExerciseEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            BaseURL + "/{exerciseId}",
	HandlerFunction: SolveExercise,
	IsProtected:     true,
}

var CreateExerciseEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            BaseURL,
	HandlerFunction: CreateExercise,
	IsProtected:     true,
}

var GetExerciseEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            BaseURL + "/{exerciseId}",
	HandlerFunction: FindExercise,
	IsProtected:     true,
}

var Endpoints = []models.Endpoint{
	SolveExerciseEndpoint,
	CreateExerciseEndpoint,
	GetExerciseEndpoint,
}
