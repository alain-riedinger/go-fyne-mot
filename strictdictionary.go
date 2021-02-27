package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
)

// MAX_LEN is the maximum length of a valid word
const MAX_LEN = 10

func parseUnmunchedDico(dicoPath string, lineStart int, lineEnd int) {
	// Unmunched dictionary file to be opened for reading
	fis, err := os.Open(dicoPath)
	if err != nil {
		log.Printf("parseUnmunchedDico - os.Open - Error: %s", err)
	}
	defer fis.Close()

	// Strict list of words, respecting all the rules
	dicoExt := filepath.Ext(dicoPath)
	dicoName := dicoPath[0 : len(dicoPath)-len(dicoExt)]
	outPath := dicoName + "-strict" + dicoExt
	out, err := os.Create(outPath)
	if err != nil {
		log.Printf("parseUnmunchedDico - os.Create - Error: %s", err)
	}
	defer out.Close()

	// Internal map for unicity of parsed words
	var dico = make(map[string]string)

	// Loop through all the lines
	sc := bufio.NewScanner(fis)
	for line := 1; sc.Scan(); line++ {
		if lineStart > 0 && line < lineStart {
			// The first lines are skipped
			continue
		} else {
			if lineEnd > 0 && lineEnd < line {
				// Processing is over: exit the loop without parsing last lines
				break
			}

			parsed := parseLine(sc.Text())
			if parsed != "" {
				_, ok := dico[parsed]
				if !ok {
					dico[parsed] = "present"

					out.WriteString(parsed + "\n")
				}
			}
		}
	}
}

func parseLine(line string) string {
	parsedLine := ""

	// Go UTF8 and Unicode interaction needs a string and rune equivalent:
	// à   â   ä   é   è   ê   ë   î   ï   ô   ö   ù   û   ü   ç
	// 224 226 228 233 232 234 235 238 239 244 246 249 251 252 231
	rline := []rune(line)

	length := 0
	for i := 0; i < len(rline); i++ {
		c := rline[i]
		if c == '/' {
			// Process derived words: accepted, but not added and finishes parsing
			// Must be checked outside of the "switch" to "break" the for loop
			break
		}

		switch c {
		// Process accented characters: accepted and added unaccented
		case 'a', 'à', 'â', 'ä':
			parsedLine += "a"
			length++
		case 'e', 'é', 'è', 'ê', 'ë':
			parsedLine += "e"
			length++
		case 'i', 'î', 'ï':
			parsedLine += "i"
			length++
		case 'o', 'ô', 'ö':
			parsedLine += "o"
			length++
		case 'u', 'ù', 'û', 'ü':
			parsedLine += "u"
			length++
		// Process modified characters: accepted and added
		case 'c', 'ç':
			parsedLine += "c"
			length++
		// Process plain characters: accepted and added
		case 'b', 'd', 'f', 'g', 'h', 'j', 'k', 'l', 'm', 'n', 'p', 'q', 'r', 's', 't', 'v', 'w', 'x', 'y', 'z':
			parsedLine += string(c)
			length++
		// Process compound word separator: accepted, but not added
		case '-':
			continue
		// Process any other character: refused and finishes parsing
		default:
			return ""
		}
	}
	if len(parsedLine) > MAX_LEN {
		// Word is too long: no need to store it
		return ""
	}
	return parsedLine
}

func loadStrictDico(dicoPath string) map[[14]byte][]string {
	dico := make(map[[14]byte][]string)

	// Strict processed dictionary file to be opened for reading
	fis, err := os.Open(dicoPath)
	if err != nil {
		log.Printf("loadStrictDico - os.Open - Error: %s", err)
	}
	defer fis.Close()

	// Loop through all the words
	sc := bufio.NewScanner(fis)
	for sc.Scan() {
		word := sc.Text()
		idx := calcIndex(word)
		_, ok := dico[idx]
		if !ok {
			// Add a new key / list to the dictionary if not yet existing
			var s []string
			dico[idx] = s
		}
		dico[idx] = append(dico[idx], word)
	}

	return dico
}

// Computes the index of a given word
// <Nb><Mask of 13 bytes>
//   <Nb>, nb of characters of the word
//   <Mask of 13 bytes>, occurences of each of 26 letters
//                       - ordered after
//                       - half bytes, recomposed
func calcIndex(word string) [14]byte {
	// Count the occurences of letters
	var counts [26]byte
	for i := 0; i < len(word); i++ {
		letter := word[i] - 'a'
		counts[letter]++
	}

	// Array of bytes: counts are grouped 2 by 2
	var idx [14]byte
	idx[0] = byte(len(word))
	for r := 0; r < 13; r++ {
		idx[1+r] = byte((int(counts[2*r]) << 4) + (int(counts[2*r+1])))
	}
	return idx
}
