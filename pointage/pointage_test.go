package pointage_test

import (
	"github.com/julienbayle/listedeseleves/pointage"
	"testing"
)

func TestLoad(t *testing.T) {
	classesOfStudents := pointage.Load("students-test.xls")
	if len(classesOfStudents) != 8 {
		t.Errorf("Classes count invalid : %d", len(classesOfStudents))
	}
}
