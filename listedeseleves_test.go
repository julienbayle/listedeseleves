package listedeseleves_test

import (
	"github.com/julienbayle/listedeseleves"
	"testing"
)

func TestMain(t *testing.T) {
	classesOfStudents := listedeseleves.Load("students-test.xls")
	if len(classesOfStudents) != 8 {
		t.Errorf("Classes count invalid : %d", len(classesOfStudents))
	}
}
