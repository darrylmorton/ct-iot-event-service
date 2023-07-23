package test

import (
	"github.com/google/uuid"
	"testing"
)

func AssertEvent(expectedResult Event, actualResult Event, t *testing.T) {
	_, err := uuid.Parse(actualResult.Id)
	if err != nil {
		t.Errorf("invalid uuid: %v", actualResult.Id)
	}
	if expectedResult.DeviceName != actualResult.DeviceName {
		t.Errorf("expected: %+v, actual: %+v", expectedResult.DeviceName, actualResult.DeviceName)
	}
	if expectedResult.Type != actualResult.Type {
		t.Errorf("expected: %+v, actual: %+v", expectedResult.Type, actualResult.Type)
	}
	if expectedResult.Event != actualResult.Event {
		t.Errorf("expected: %+v, actual: %+v", expectedResult.Event, actualResult.Event)
	}
	if expectedResult.Read != actualResult.Read {
		t.Errorf("expected: %+v, actual: %+v", expectedResult.Read, actualResult.Read)
	}
}

func AssertEvents(expectedResult []Event, actualResult []Event, t *testing.T) {
	if len(expectedResult) != len(actualResult) {
		t.Errorf("expected: %v, actual: %v", len(expectedResult), len(actualResult))
	}

	for index, _ := range actualResult {
		AssertEvent(expectedResult[index], actualResult[index], t)
	}
}
