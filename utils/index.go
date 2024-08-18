package utils

import "fmt"

func PrintProgressBar(current, total int) {
	progress := float64(current) / float64(total) * 100
	barLength := 50
	progressBar := "["
	numHashes := int(progress / (100 / float64(barLength)))
	for i := 0; i < numHashes; i++ {
		progressBar += "#"
	}
	for i := numHashes; i < barLength; i++ {
		progressBar += " "
	}
	progressBar += "]"
	fmt.Printf("\r> Progress: %s %d%%", progressBar, int(progress))
}
