package main

import (
	"bufio"
	"fmt"
	_ "fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Cell struct {
	toKeyMaker, stepsMade, x, y int
	cellType                    string // D (default), S (Sentinel), A (Agent), P (perception), K (KeyMaker), B (backdoor)
	parentHistory               *Cell
	childrenHistory             []*Cell
	isVisited                   bool
}

func (cell *Cell) sum() int {
	return cell.toKeyMaker + cell.stepsMade
}

type Field struct {
	keyMakerCoordinates [2]int
	cells               map[[2]int]*Cell
	possibleCells       []*Cell
	currentPosition     [2]int
}

func (field *Field) toKeyMakerCount(x, y int) int {

	i := field.keyMakerCoordinates
	return int(math.Abs(float64(x-i[0]))) + int(math.Abs(float64(y-i[1])))
}

func (field *Field) constructor(xKeyMaker, yKeyMaker int) {
	field.cells = make(map[[2]int]*Cell)

	field.keyMakerCoordinates = [...]int{xKeyMaker, yKeyMaker}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			var cell = new(Cell)
			cell.toKeyMaker = field.toKeyMakerCount(i, j)
			cell.cellType = "D" // Set default cells in initializing
			cell.x, cell.y = i, j
			cell.isVisited = i == 0 && j == 0 // Вот это финт ушами!
			field.cells[[...]int{i, j}] = cell
		}
	}
}
func (field *Field) lookForCells() {
	currentX, currentY := field.currentPosition[0], field.currentPosition[1]
	variants := [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for i := 0; i < len(variants); i++ {
		toY := currentY + variants[i][1]
		toX := currentX + variants[i][0]
		toPair := [2]int{toX, toY}
		if field.limitValidity(toX, toY) {
			toCell := field.cells[toPair] // toCell - cell where can go
			toType := toCell.cellType     // toType - type of cell where can go
			if toType != "P" && toType != "S" && toType != "A" {
				if !toCell.isVisited {

					if field.cells[toPair].stepsMade > (field.cells[field.currentPosition].stepsMade+1) || field.cells[toPair].stepsMade == 0 {
						field.possibleCells = append(field.possibleCells, field.cells[toPair])
						field.cells[toPair].stepsMade = field.cells[field.currentPosition].stepsMade + 1
					}
				}
			}
		}
	}
}
func (*Field) limitValidity(x, y int) bool {
	return x >= 0 && x < 9 && y < 9 && y >= 0
}

func (field *Field) makeStep() {
	sort.Slice(field.possibleCells, func(i, j int) bool {
		if field.possibleCells[i].sum() == field.possibleCells[j].sum() {
			return field.possibleCells[i].toKeyMaker < field.possibleCells[j].toKeyMaker
		}
		return field.possibleCells[i].sum() < field.possibleCells[j].sum()
	})
	if field.areClose() {
		field.stepToGoal()
		return
	}
	//If cells are not close
	var stack []*Cell // 0,0 - is root
	// Firstly went to root
	for !(field.currentPosition[0] == 0 && field.currentPosition[1] == 0) {
		field.currentPosition = [2]int{field.cells[field.currentPosition].parentHistory.x, field.cells[field.currentPosition].parentHistory.y} // MAKE STEP
		fmt.Printf("m %d %d\n", field.cells[field.currentPosition].x, field.cells[field.currentPosition].y)
		field.input()
	}

	stack = append(stack, field.cells[field.currentPosition])

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if field.areClose() {
			field.stepToGoal()
			return
		}
		for _, child := range current.childrenHistory {
			stack = append(stack, child)
		}
		field.currentPosition = [2]int{current.x, current.y} // MAKE STEP
		fmt.Printf("m %d %d\n", current.x, current.y)
		field.input()
	}
}

func (field *Field) stepToGoal() {
	temp := -1
	for j := 0; j < len(field.cells[field.currentPosition].childrenHistory); j++ { // Check ion repetitions
		if field.cells[field.currentPosition].childrenHistory[j] == field.cells[[2]int{field.possibleCells[0].x, field.possibleCells[0].y}] {
			temp = 1
		}
	}

	if len(field.cells[field.currentPosition].childrenHistory) == 0 || temp == -1 {
		field.cells[field.currentPosition].childrenHistory = append(field.cells[field.currentPosition].childrenHistory, field.cells[[2]int{field.possibleCells[0].x, field.possibleCells[0].y}])
		field.cells[[2]int{field.possibleCells[0].x, field.possibleCells[0].y}].parentHistory = field.cells[field.currentPosition]
	}
	field.currentPosition = [2]int{field.possibleCells[0].x, field.possibleCells[0].y} // MAKE STEP
	fmt.Printf("m %d %d\n", field.possibleCells[0].x, field.possibleCells[0].y)
	field.possibleCells[0].isVisited = true
	field.possibleCells = append(field.possibleCells[:0], field.possibleCells[1:]...)
}

func (field *Field) less(i, j int) bool {
	return field.possibleCells[i].sum() < field.possibleCells[j].sum()
}

func (field *Field) areClose() bool {
	if field.currentPosition[0]-field.possibleCells[0].x == 1 || field.currentPosition[0]-field.possibleCells[0].x == -1 {
		if field.currentPosition[1] == field.possibleCells[0].y {
			return true
		}
	} else if field.currentPosition[1]-field.possibleCells[0].y == 1 || field.currentPosition[1]-field.possibleCells[0].y == -1 {
		if field.currentPosition[0] == field.possibleCells[0].x {
			return true
		}
	}
	return false
}

func (field *Field) input() {
	reader := bufio.NewReader(os.Stdin)
	var length int
	_, err := fmt.Scan(&length)
	if err != nil {
		return
	}
	for i := 0; i < length; i++ {
		inputStr, err2 := reader.ReadString('\n')
		if err2 != nil {
			return
		}
		temp := strings.Split(inputStr, " ")
		x, err3 := strconv.Atoi(temp[0])
		if err3 != nil {
			return
		}
		y, err4 := strconv.Atoi(temp[1])
		if err4 != nil {
			return
		}
		field.cells[[...]int{x, y}].cellType = string(temp[2][0])
	}
}

func main() {
	var field Field
	keyMaker := [2]int{6, 1}
	field.constructor(keyMaker[0], keyMaker[1]) // ВВОД С УЧЕТОМ ТОГО ЧТО НУЛЕВАЯ ТОЧКА ЭТО 0, 0 А НЕ 1, 1
	fmt.Println("m 0 0")
	for !(field.currentPosition[0] == keyMaker[0] && field.currentPosition[1] == keyMaker[1]) {
		field.input()
		field.lookForCells()
		field.makeStep()
	}
	fmt.Printf("e %d", field.cells[[2]int{keyMaker[0], keyMaker[1]}].stepsMade)
	d := field.cells[[2]int{0, 0}]
	d.sum()
}
