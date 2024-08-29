package main

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Calculator represents a simple calculator
type Calculator struct {
	display   string
	operation string
	operand1  float64
	operand2  float64
}

// NewCalculator creates a new Calculator instance
func NewCalculator() *Calculator {
	return &Calculator{
		display: "",
	}
}

// Add adds input to the calculator display and stores operands and operation
func (c *Calculator) Add(input string) {
	if input == "+" || input == "-" || input == "*" || input == "/" || input == "√" || input == "log" {
		if input == "√" || input == "log" {
			c.operand1, _ = strconv.ParseFloat(c.display, 64)
			c.operation = input
			c.display = ""
			c.Calculate()
		} else {
			c.operand1, _ = strconv.ParseFloat(c.display, 64)
			c.operation = input
			c.display = ""
		}
	} else {
		c.display += input
	}
}

// Calculate performs the operation based on stored operand and operator
func (c *Calculator) Calculate() error {
	if c.operation != "√" && c.operation != "log" {
		c.operand2, _ = strconv.ParseFloat(c.display, 64)
	}

	switch c.operation {
	case "+":
		c.display = fmt.Sprintf("%v", c.operand1+c.operand2)
	case "-":
		c.display = fmt.Sprintf("%v", c.operand1-c.operand2)
	case "*":
		c.display = fmt.Sprintf("%v", c.operand1*c.operand2)
	case "/":
		if c.operand2 == 0 {
			c.display = "Error"
			return fmt.Errorf("Не делится на ноль")
		}
		c.display = fmt.Sprintf("%v", c.operand1/c.operand2)
	case "√":
		if c.operand1 < 0 {
			c.display = "Error"
			return fmt.Errorf("square root of negative number")
		}
		c.display = fmt.Sprintf("%v", math.Sqrt(c.operand1))
	case "log":
		if c.operand1 <= 0 {
			c.display = "Error"
			return fmt.Errorf("logarithm of non-positive number")
		}
		c.display = fmt.Sprintf("%v", math.Log10(c.operand1))
	default:
		return fmt.Errorf("unknown operation")
	}
	c.operation = ""
	return nil
}

// Clear resets the calculator
func (c *Calculator) Clear() {
	c.display = ""
	c.operation = ""
	c.operand1 = 0
	c.operand2 = 0
}

func main() {
	a := app.New()
	w := a.NewWindow("Калькулятор")

	calc := NewCalculator()

	display := widget.NewLabel(calc.display)
	buttons := [][]string{
		{"7", "8", "9", "/"},
		{"4", "5", "6", "*"},
		{"1", "2", "3", "-"},
		{"0", ".", "=", "+"},
		{"√", "log", "C", "Разность дат"},
		{"Расчёт силы"},
	}

	grid := container.NewGridWithRows(7,
		display,
	)

	for _, row := range buttons {
		rowContainer := container.NewHBox()
		for _, buttonLabel := range row {
			button := widget.NewButton(buttonLabel, func(label string) func() {
				return func() {
					switch label {
					case "=":
						err := calc.Calculate()
						if err != nil {
							display.SetText("Error")
						} else {
							display.SetText(calc.display)
						}
					case "C":
						calc.Clear()
						display.SetText(calc.display)
					case "Разность дат":
						openDateCalculator(a)
					case "Расчёт силы":
						openPhysicsCalculator(a)
					default:
						calc.Add(label)
						display.SetText(calc.display)
					}
				}
			}(buttonLabel))
			rowContainer.Add(button)
		}
		grid.Add(rowContainer)
	}

	w.SetContent(grid)
	w.ShowAndRun()
}

// Open a new window to calculate the difference between two dates
func openDateCalculator(a fyne.App) {
	w := a.NewWindow("Разность дат")

	label1 := widget.NewLabel("Введите первую дату (YYYY-MM-DD):")
	date1 := widget.NewEntry()
	label2 := widget.NewLabel("Введите вторую дату (YYYY-MM-DD):")
	date2 := widget.NewEntry()
	result := widget.NewLabel("Разность: ")

	calculateBtn := widget.NewButton("Разность", func() {
		d1, err1 := time.Parse("2006-01-02", date1.Text)
		d2, err2 := time.Parse("2006-01-02", date2.Text)
		if err1 != nil || err2 != nil {
			dialog.ShowError(fmt.Errorf("Неверный формат"), w)
			return
		}
		diff := d2.Sub(d1)
		result.SetText(fmt.Sprintf("Разница: %v дней", diff.Hours()/24))
	})

	w.SetContent(container.NewVBox(label1, date1, label2, date2, calculateBtn, result))
	w.Resize(fyne.NewSize(300, 200))
	w.Show()
}

// Open a new window to calculate simple physics formulas
func openPhysicsCalculator(a fyne.App) {
	w := a.NewWindow("Калькулятор по физике")

	label1 := widget.NewLabel("Масса (kg):")
	mass := widget.NewEntry()
	label2 := widget.NewLabel("Ускорение (m/s^2):")
	acceleration := widget.NewEntry()
	result := widget.NewLabel("Сила: ")

	calculateBtn := widget.NewButton("Калькулятор Силы (F = m * a)", func() {
		m, err1 := strconv.ParseFloat(mass.Text, 64)
		a, err2 := strconv.ParseFloat(acceleration.Text, 64)
		if err1 != nil || err2 != nil {
			dialog.ShowError(fmt.Errorf("Неверный ввод"), w)
			return
		}
		force := m * a
		result.SetText(fmt.Sprintf("Сила: %v N", force))
	})

	w.SetContent(container.NewVBox(label1, mass, label2, acceleration, calculateBtn, result))
	w.Resize(fyne.NewSize(300, 200))
	w.Show()
}
