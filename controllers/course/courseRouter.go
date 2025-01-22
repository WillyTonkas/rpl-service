package course

import (
	"rpl-service/models"
)

const BaseURL = "/courses"

var ExistsEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            BaseURL + "/course/{id}",
	HandlerFunction: Exists,
}

var CreateCourseEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            BaseURL + "/course",
	HandlerFunction: Create,
}

var EnrollToCourseEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            BaseURL + "/enroll",
	HandlerFunction: EnrollToCourse,
}

var StudentExistsEndPoint = models.Endpoint{
	Method:          models.GET,
	Path:            BaseURL + "/course/{id}/is-enrolled",
	HandlerFunction: StudentExists,
}

var DeleteStudentEndpoint = models.Endpoint{
	Method:          models.DELETE,
	Path:            BaseURL + "/course/{id}/student",
	HandlerFunction: DeleteStudent,
}
