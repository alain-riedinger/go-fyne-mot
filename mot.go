package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

// NbTirage is the number of letters in the tirage
const nbTirage = 10

// Mot is the class that holds the game
type Mot struct {
	voyelles  string
	consonnes string
}

// NewMot initializes the Mot structure
func NewMot() *Mot {
	m := new(Mot)

	// Initializes uniquely the random
	rand.Seed(time.Now().UnixNano())

	// Allocation of vocals frequency on 248 slots
	m.voyelles = ""
	m.voyelles += strings.Repeat("a", 48)
	m.voyelles += strings.Repeat("e", 92)
	m.voyelles += strings.Repeat("i", 43)
	m.voyelles += strings.Repeat("i", 33)
	m.voyelles += strings.Repeat("u", 29)
	m.voyelles += strings.Repeat("y", 3)
	// Allocation of consonants frequency on 248 slots
	m.consonnes = ""
	m.consonnes += strings.Repeat("b", 7)
	m.consonnes += strings.Repeat("c", 19)
	m.consonnes += strings.Repeat("d", 21)
	m.consonnes += strings.Repeat("f", 7)
	m.consonnes += strings.Repeat("h", 7)
	m.consonnes += strings.Repeat("j", 2)
	m.consonnes += strings.Repeat("k", 2)
	m.consonnes += strings.Repeat("l", 28)
	m.consonnes += strings.Repeat("m", 15)
	m.consonnes += strings.Repeat("n", 36)
	m.consonnes += strings.Repeat("p", 14)
	m.consonnes += strings.Repeat("q", 4)
	m.consonnes += strings.Repeat("r", 35)
	m.consonnes += strings.Repeat("s", 37)
	m.consonnes += strings.Repeat("t", 34)
	m.consonnes += strings.Repeat("v", 7)
	m.consonnes += strings.Repeat("w", 1)
	m.consonnes += strings.Repeat("x", 3)
	m.consonnes += strings.Repeat("z", 1)

	// Lets' shuffle a bit the seeds to avoid artefacts due to random
	m.voyelles = shuffle(m.voyelles)
	m.consonnes = shuffle(m.consonnes)

	return m
}

// GetPlaques returns a random tirage
// with nbVoyelles vocals
func (m *Mot) GetPlaques(nbVoyelles int) string {
	var chosenVoyelles []int
	var chosenConsonnes []int

	maxVoyelles := len(m.voyelles)
	maxConsonnes := len(m.consonnes)

	nbConsonnes := nbTirage - nbVoyelles

	tirage := ""
	// Choose randomly the letters
	for nbChosen := 0; nbChosen < nbVoyelles; {
		p := rand.Intn(maxVoyelles)
		if !contains(chosenVoyelles, p) {
			chosenVoyelles = append(chosenVoyelles, p)
			nbChosen++
		}
	}
	// Then the consonnes
	for nbChosen := 0; nbChosen < nbConsonnes; {
		p := rand.Intn(maxConsonnes)
		if !contains(chosenConsonnes, p) {
			chosenConsonnes = append(chosenConsonnes, p)
			nbChosen++
		}
	}

	// Compose the tirage, by mixing voyelles and consonnes
	for i := 0; i < len(chosenVoyelles); i++ {
		tirage += fmt.Sprintf("%c", m.voyelles[chosenVoyelles[i]])
	}
	for i := 0; i < len(chosenConsonnes); i++ {
		tirage += fmt.Sprintf("%c", m.consonnes[chosenConsonnes[i]])
	}
	// Lets' shuffle a bit, to make tirage a mix of both
	return shuffle(tirage)
}

// SolveTirage recursively finds the longest matching words in the dictionary
func (m *Mot) SolveTirage(dico map[[14]byte][]string, solution Solution) *Solution {
	// Stop searching for smaller words if something is already found
	if len(solution.Current) < solution.BestLen {
		// Remaining stub of letters is less than best yet found: nothing to hope
		return nil
	}

	// Retrieve combination from dictionary
	idx := calcIndex(solution.Current)
	mots, ok := dico[idx]
	if ok {
		// Combination found in the dictionary: return it !
		found := NewSolution()
		found.Current = solution.Current
		found.BestLen = len(solution.Current)
		if found.BestLen > solution.BestLen {
			found.Best = mots
		} else {
			if found.BestLen == solution.BestLen {
				// Solution with same length has been found: add them, if not yet in
				// found.Best = solution.Best
				found.Best = append(found.Best, solution.Best...)
				for _, w := range mots {
					if !contains(found.Best, w) {
						found.Best = append(found.Best, w)
					}
				}
			}
		}
		return found
	} else {
		// No matching combination in dictionary: search recursively with one letter less
		ln := len(solution.Current)
		bestSol := NewSolution()
		for i := 0; i < ln; i++ {
			iter := NewSolution()
			iter.Current = removeLetter(solution.Current, i)
			found := m.SolveTirage(dico, *iter)
			if found != nil {
				if found.BestLen > bestSol.BestLen {
					// A better solution has been found: replace it
					bestSol.BestLen = found.BestLen
					bestSol.Best = found.Best
				} else {
					if found.BestLen == bestSol.BestLen {
						// Solution with same length has been found: add them, if not yet in
						for _, w := range found.Best {
							if !contains(bestSol.Best, w) {
								bestSol.Best = append(bestSol.Best, w)
							}
						}
					}
				}
			}
		}
		if bestSol.BestLen > 0 {
			// A set of solution has been found: return it
			return bestSol
		}
	}

	return nil
}

func contains(s interface{}, elem interface{}) bool {
	arrV := reflect.ValueOf(s)
	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {
			// XXX - panics if slice element points to an unexported struct field
			// see https://golang.org/pkg/reflect/#Value.Interface
			if arrV.Index(i).Interface() == elem {
				return true
			}
		}
	}
	return false
}

func removeLetter(text string, rank int) string {
	result := text[0:rank] + text[rank+1:]
	return result
}

func shuffle(text string) string {
	rnText := []rune(text)
	rand.Shuffle(len(rnText), func(i, j int) {
		rnText[i], rnText[j] = rnText[j], rnText[i]
	})
	return string(rnText)
}
