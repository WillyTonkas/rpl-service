package mappers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"rpl-service/models"
)

func SolveExerciseRequestBody(exerciseCode, testScript string) io.Reader {
	return bytes.NewReader([]byte(`{"script":"` + exerciseCode + "\", \"tests\":" + testScript + `"`))
}

func SolveExerciseResponseToResult(response *http.Response, exerciseName, testName string) models.ExerciseResult {
	var exerciseResponse struct {
		Stdout     string `json:"stdout"`
		Stderr     string `json:"stderr"`
		ReturnCode int    `json:"returncode"`
	}
	err := json.NewDecoder(response.Body).Decode(&exerciseResponse)
	if err != nil {
		return models.ExerciseResult{}
	}

	return models.ExerciseResult{
		ExerciseName: exerciseName,
		TestName:     testName,
		TestPassed:   exerciseResponse.ReturnCode == 0,
		Stdout:       exerciseResponse.Stdout,
		Stderr:       exerciseResponse.Stderr,
	}
}
