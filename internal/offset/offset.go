package offset

import (
	"log"
	"os"
	"strconv"
	"strings"
)

const offsetFile = "offset.txt"

func LoadOffset() int {
	data, err := os.ReadFile(offsetFile)
	if err != nil {
		if os.IsNotExist(err) {
			return 0
		}
		log.Printf("Error reading offset file: %v", err)
		return 0
	}
	offsetStr := strings.TrimSpace(string(data))
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		log.Printf("Error parsing offset: %v", err)
		return 0
	}
	return offset
}

func SaveOffset(offset int) {
	err := os.WriteFile(offsetFile, []byte(strconv.Itoa(offset)), 0644)
	if err != nil {
		log.Printf("Error writing offset file: %v", err)
	}
}
