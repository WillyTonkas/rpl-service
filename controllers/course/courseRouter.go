package course

import (
	"rpl-service/models"
)

const BaseURL = "/courses"

var ExistsEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            BaseURL + "/course/{id}",
	HandlerFunction: Exists,
	IsProtected:     true,
}

var CreateCourseEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            BaseURL + "/course",
	HandlerFunction: Create,
	IsProtected:     true,
}

var EnrollToCourseEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            BaseURL + "/enroll",
	HandlerFunction: EnrollToCourse,
	IsProtected:     true,
}

var StudentExistsEndPoint = models.Endpoint{
	Method:          models.GET,
	Path:            BaseURL + "/course/{id}/is-enrolled",
	HandlerFunction: StudentExists,
	IsProtected:     true,
}

var DeleteStudentEndpoint = models.Endpoint{
	Method:          models.DELETE,
	Path:            BaseURL + "/course/{id}/student",
	HandlerFunction: DeleteStudent,
	IsProtected:     true,
}

var Endpoints = []models.Endpoint{
	ExistsEndpoint,
	CreateCourseEndpoint,
	EnrollToCourseEndpoint,
	StudentExistsEndPoint,
	DeleteStudentEndpoint,
}
