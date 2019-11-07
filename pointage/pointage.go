package pointage

import (
	"fmt"
	"github.com/extrame/xls"
	"github.com/goodsign/monday"
	"github.com/tealeg/xlsx"
	"time"
)

// School activities
var ActivityCodes = [3]string{"matin", "repas", "soir"}

// Full name of each school activity
var ActivityTitles = map[string]string{
	"matin": "Périscolaire du matin",
	"repas": "Restauration",
	"soir":  "Périscolaire du soir",
}

// Source XLS headers for students and activities
const HeaderLine = 7
const FirstNameCode = "Nom"
const LastNameCode = "Prénom"

var ActivityHeaders = map[string]string{
	"matin": "GARDERIE MATIN 4J",
	"repas": "FORFAIT RESTAU 4J",
	"soir":  "GARDERIE SOIR 4J",
}

// Export parameters
const ExportTitle = "%s de la classe de %s pour le mois de %s"

// School days
func IsSchoolOpen(day string) bool {
	switch day {
	case
		"Monday",
		"Tuesday",
		"Thursday",
		"Friday":
		return true
	}
	return false
}

type Student struct {
	FirstName, LastName   string
	IsFlatRateForActivity map[string]bool
}

type ClassOfStudents struct {
	Name     string
	Students []Student
}

func GetColIndex(row *xls.Row, text string) int {
	for i := row.FirstCol(); i < row.LastCol(); i++ {
		if row.Col(i) == text {
			return i
		}
	}
	panic(fmt.Sprintf("Ligne %d incorrecte, impossible de trouver : %s", HeaderLine, text))
}

// Parse school data (XLS file)
func Load(studentsFileName string) []ClassOfStudents {
	xlFile, err := xls.Open(studentsFileName, "utf-8")
	if err != nil {
		panic(err)
	}

	classesOfStudents := make([]ClassOfStudents, xlFile.NumSheets())
	for classId := 0; classId < xlFile.NumSheets(); classId++ {
		classData := xlFile.GetSheet(classId)
		classesOfStudents[classId].Name = classData.Name

		// Read and validate headers
		headerData := classData.Row(HeaderLine - 1)
		var firstNameIndex = GetColIndex(headerData, FirstNameCode)
		var lastNameIndex = GetColIndex(headerData, LastNameCode)
		var activityIndexes = make(map[string]int, len(ActivityHeaders))
		for activityName, activityCode := range ActivityHeaders {
			activityIndexes[activityName] = GetColIndex(headerData, activityCode)
		}

		// Read all students data
		for row := HeaderLine; row < int(classData.MaxRow); row++ {
			rowData := classData.Row(row)
			if rowData.Col(0) != "" {
				student := Student{
					FirstName: rowData.Col(firstNameIndex),
					LastName:  rowData.Col(lastNameIndex),
				}
				student.IsFlatRateForActivity = make(map[string]bool, len(ActivityHeaders))
				for _, activityName := range ActivityCodes {
					student.IsFlatRateForActivity[activityName] = (rowData.Col(activityIndexes[activityName]) != "")
				}
				classesOfStudents[classId].Students = append(classesOfStudents[classId].Students, student)
			}
		}

		// Add a blank line (in case of new student)
		student := Student{FirstName: " ", LastName: " "}
		student.IsFlatRateForActivity = make(map[string]bool, len(ActivityHeaders))
		for _, activityName := range ActivityCodes {
			student.IsFlatRateForActivity[activityName] = false
		}
		classesOfStudents[classId].Students = append(classesOfStudents[classId].Students, student)
	}
	return classesOfStudents
}

// Style Method
func ThinBorder() xlsx.Border {
	b := xlsx.NewBorder("thin", "thin", "thin", "thin")
	b.TopColor, b.BottomColor = "FF000000", "FF000000"
	b.BottomColor, b.LeftColor = "FF000000", "FF000000"
	return *b
}

// Style Method
func ThickBorder() xlsx.Border {
	b := xlsx.NewBorder("thick", "thin", "thin", "thin")
	b.TopColor, b.BottomColor = "FF000000", "FF000000"
	b.BottomColor, b.LeftColor = "FF000000", "FF000000"
	return *b
}

// Style Method
func DefaultCellStyleOdd() *xlsx.Style {
	style := xlsx.NewStyle()
	style.ApplyBorder = true
	style.Border = ThinBorder()
	style.Font = *xlsx.NewFont(12, "Arial")
	return style
}

// Style Method
func DefaultCellStyleEven() *xlsx.Style {
	style := DefaultCellStyleOdd()
	style.ApplyFill = true
	style.Fill = *xlsx.NewFill(xlsx.Solid_Cell_Fill, "FFEEEEEE", "FFDDDDDD")
	return style
}

// Style Method
func AddTickBorder(input *xlsx.Style) *xlsx.Style {
	style := input
	style.ApplyBorder = true
	style.Border = ThickBorder()
	return style
}

// Style Method
func TitleCellStyle() *xlsx.Style {
	style := xlsx.NewStyle()
	f := xlsx.NewFont(16, "Arial")
	f.Bold = true
	style.Font = *f
	return style
}

// Style Method
func HeaderCellStyle() *xlsx.Style {
	style := xlsx.NewStyle()
	style.ApplyAlignment = true
	style.Alignment = xlsx.Alignment{
		Horizontal: "center",
		Vertical:   "center",
	}
	style.ApplyFill = true
	style.Fill = *xlsx.NewFill(xlsx.Solid_Cell_Fill, "FF666666", "FFFFFFFF")
	f := xlsx.NewFont(10, "Arial")
	f.Color = "FFFFFFFF"
	f.Bold = true
	style.Font = *f
	return style
}

// Style Method
func HighlightedCellStyle() *xlsx.Style {
	style := DefaultCellStyleOdd()
	style.ApplyFill = true
	style.Fill = *xlsx.NewFill(xlsx.Solid_Cell_Fill, "FFAAAAAA", "FF000000")
	return style
}

// Export school tally sheets
func Export(classesOfStudents []ClassOfStudents, date time.Time, exportFileName string) {
	exportFile := xlsx.NewFile()
	month := monday.Format(date, "January 2006", monday.LocaleFrFR)

	for _, activity := range ActivityCodes {
		for _, classOfStudents := range classesOfStudents {
			// Create a new sheet
			sheet, _ := exportFile.AddSheet(fmt.Sprintf("%s - %s", classOfStudents.Name, activity))

			// Add some title
			row1 := sheet.AddRow()
			row1.AddCell().SetString(fmt.Sprintf(ExportTitle, ActivityTitles[activity], classOfStudents.Name, month))
			row1.SetHeight(24)

			for _, cell := range row1.Cells {
				cell.SetStyle(TitleCellStyle())
			}

			// Add two header lines
			row2 := sheet.AddRow()
			row2.AddCell()
			row2.AddCell()
			row2.AddCell()

			row3 := sheet.AddRow()
			row3.AddCell()
			row3.AddCell().SetString("Nom")
			row3.AddCell().SetString("Prénom")

			iterDate := date
			nbOpenDays := 0
			for iterDate.Month() == date.Month() {
				if IsSchoolOpen(iterDate.Format("Monday")) {
					row2.AddCell().SetString(monday.Format(iterDate, "Mon", monday.LocaleFrFR))
					row3.AddCell().SetString(iterDate.Format("02"))
					nbOpenDays++
				}
				iterDate = iterDate.AddDate(0, 0, 1)
			}

			row2.AddCell().SetString("Total")
			row3.AddCell().SetString(monday.Format(date, "Jan", monday.LocaleFrFR))

			for _, cell := range row2.Cells {
				cell.SetStyle(HeaderCellStyle())
			}
			for _, cell := range row3.Cells {
				cell.SetStyle(HeaderCellStyle())
			}

			// Set column width
			sheet.SetColWidth(1, 1, 4.0)
			sheet.SetColWidth(2, 3, 30)
			sheet.SetColWidth(4, nbOpenDays+3, 5)
			sheet.SetColWidth(nbOpenDays+4, nbOpenDays+5, 8.0)

			// Add students
			for num, student := range classOfStudents.Students {
				row := sheet.AddRow()
				row.SetHeight(18)
				row.AddCell().SetInt(num + 1)
				row.AddCell().SetString(student.FirstName)
				row.AddCell().SetString(student.LastName)

				for i := 0; i < nbOpenDays+1; i++ {
					row.AddCell()
				}

				for _, cell := range row.Cells {
					if student.IsFlatRateForActivity[activity] {
						cell.SetStyle(HighlightedCellStyle())
					} else if num%2 == 0 {
						cell.SetStyle(DefaultCellStyleEven())
					} else {
						cell.SetStyle(DefaultCellStyleOdd())
					}
				}

				if student.IsFlatRateForActivity[activity] {
					row.Cells[len(row.Cells)-1].SetString("Forfait")
				}
			}

			// Add a border on the last cells
			for _, row := range sheet.Rows[3:] {
				row.Cells[len(row.Cells)-1].SetStyle(AddTickBorder(row.Cells[len(row.Cells)-1].GetStyle()))
			}

			// Add a footer notice
			footer1 := sheet.AddRow()
			footer1.AddCell().SetString("Pointer chaque jour tous les enfants, quelque soit la couleur de la ligne.")
			footer2 := sheet.AddRow()
			footer2.AddCell().SetString("Le dernier jour du mois, comptabiliser le total pour les lignes hors forfait et déposer la fiche dans la banette 'Trésorier OGEC'")
		}
	}

	exportFile.Save(exportFileName)
}
