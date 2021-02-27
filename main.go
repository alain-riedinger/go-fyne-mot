// [Fyne toolkit documentation for developers | Develop using Fyne](https://developer.fyne.io/index.html)

package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Create the application and window
	myApp := app.New()
	myWindow := myApp.NewWindow("Le Mot le plus Long")

	// Load the dictionary
	dico := loadStrictDico("fr-mlpl-flat-strict.txt")

	// Global string for textual solution
	var solToDisplay *Solution
	const solEmpty = "\n\n\n\n\n"

	// Create the items of the window
	const nbPlaques = 10
	var btnPlaques [nbPlaques](*widget.Button)
	for i := 0; i < len(btnPlaques); i++ {
		btnPlaques[i] = widget.NewButton("_", tapped)
	}
	// Items are displayed horizontally in a grid
	gridPlaques := container.New(layout.NewGridLayout(10),
		btnPlaques[0], btnPlaques[1], btnPlaques[2], btnPlaques[3], btnPlaques[4], btnPlaques[5], btnPlaques[6], btnPlaques[7], btnPlaques[8], btnPlaques[9])

	lblNbVoyelles := widget.NewLabel("# Voy.:")
	entryNbVoyelles := widget.NewEntry()
	progress := widget.NewProgressBar()
	// Button is on the left, with its default size
	// Progress bar is on the left, stretched to use all the remaining space
	gridNbVoyelles := container.New(layout.NewGridLayout(2), lblNbVoyelles, entryNbVoyelles)
	gridTirage := container.New(layout.NewBorderLayout(nil, nil, gridNbVoyelles, nil), gridNbVoyelles, progress)

	// Set default number of voyelles
	entryNbVoyelles.SetText("4")

	// Separator makes the layout nice
	separator1 := widget.NewSeparator()

	// Trigger for solution display
	// Text is forced on several lines
	txtSolution := widget.NewTextGridFromString(solEmpty)
	scrollSolution := container.NewScroll(txtSolution)
	scrollSolution.SetMinSize(txtSolution.MinSize())
	btnSolution := widget.NewButton("Solution?", func() {
		// Output final result
		txtSol := ""
		if solToDisplay != nil {
			txtSol = fmt.Sprintf("Best words found: %d letters\n", solToDisplay.BestLen)
			for _, w := range solToDisplay.Best {
				txtSol += strings.ToUpper(w) + "\n"
			}
			lBest := len(solToDisplay.Best)
			lDef := len(solEmpty)
			if lDef-(lBest+1) > 0 {
				txtSol += strings.Repeat("\n", lDef-(lBest+1))
			}
		} else {
			txtSol = fmt.Sprintf("Best words found: %d letters\n", solToDisplay.BestLen)
		}
		txtSolution.SetText(txtSol)
	})

	// Separator makes the layout nice
	separator2 := widget.NewSeparator()

	// Channel for stopping the time of the progress bar
	// stop := make(chan bool)

	// Buttons with actions
	newGame := widget.NewButton("Play!", func() {
		mot := NewMot()

		txtSolution.SetText(solEmpty)

		nbVoyelles, _ := strconv.Atoi(entryNbVoyelles.Text)
		plaques := mot.GetPlaques(nbVoyelles)
		for i := 0; i < len(btnPlaques); i++ {
			btnPlaques[i].SetText(strings.ToUpper(fmt.Sprintf("%c", plaques[i])))
		}

		// go timer(stop, progress)
		countup(progress)

		// Solution is searched during chrono time (it's shorter so no cheating)
		sol := make(chan *Solution)
		go findSolution(mot, dico, plaques, sol)
		solToDisplay = <-sol
		close(sol)
	})
	quit := widget.NewButton("Quit", func() {
		myApp.Quit()
	})
	gridActions := container.New(layout.NewGridLayout(2), newGame, quit)

	// Compose the window with the items
	// Items are horizontally stacked, first parameter is on top, and so on downwards
	myWindow.SetContent(container.New(layout.NewVBoxLayout(),
		gridPlaques,
		gridTirage,
		separator1,
		btnSolution,
		scrollSolution,
		separator2,
		gridActions))

	myWindow.SetFixedSize(true)

	// Trigger the progress bar
	progress.Min = 0
	progress.Max = 30
	progress.TextFormatter = func() string {
		// No percent displayed to avoid distraction
		return ""
	}

	// Master loop that runs the widow
	myWindow.ShowAndRun()
}

func timer(stop chan bool, pb *widget.ProgressBar) {
	for t := 0.0; t <= pb.Max; t++ {
		select {
		case <-stop:
			return
		default:
			time.Sleep(1 * time.Second)
			pb.SetValue(t)
		}
	}
}

// tapped is a dummy function, for the plaques and tirage to be buttons
func tapped() {
	// Nothing to do
}

func findSolution(mot *Mot, dico map[[14]byte][]string, plaques string, sol chan *Solution) {
	// Initialize the recursive search root structure
	solution := NewSolution()
	solution.Current = plaques

	// Start the recursive resolution
	found := mot.SolveTirage(dico, *solution)

	// Send solution to the channel, while execution in parallel
	sol <- found
}

func countup(pb *widget.ProgressBar) {
	up := 0.0
	timer := time.Tick(1 * time.Second)
	for up < pb.Max {
		<-timer
		pb.SetValue(up)
		up += 1.0
	}
}
