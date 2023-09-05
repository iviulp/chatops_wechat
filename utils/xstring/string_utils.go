package xstring

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"hash/fnv"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

// TransBytesToMarkdownStr 将数据转为markdown格式
func TransBytesToMarkdownStr(raw string) string {
	output := fmt.Sprintf("```\n%s\n```", raw)
	output = strings.Replace(output, "\\", "\\\\", -1)
	output = strings.Replace(output, "\"", "\\\"", -1)
	return output
}

// GenLogoAscii 生成ascii
func GenLogoAscii(text string, color string) {
	myFigure := figure.NewColorFigure(text, "", color, true)
	myFigure.Print()
}

const keyLength = 10

// GenerateRandomStr
// @Description: 生成随机字符串
// @return string
func GenerateRandomStr() string {
	rand.Seed(time.Now().UnixNano())
	var builder strings.Builder
	for i := 0; i < keyLength; i++ {
		builder.WriteRune(rune(rand.Intn(26) + 97)) // generate random lowercase letter
	}

	return builder.String()
}

// chunkSize represents the size of each chunk
const chunkSize = 1024 * 1024

// HashFile returns the hash value of a file
func HashFile(filePath string) int64 {
	data, _ := os.ReadFile(filePath)
	return HashData(data)
}

// HashData returns the hash value of a byte slice
func HashData(data []byte) int64 {
	// Compute hash of the data
	h := fnv.New64a()
	_, err := h.Write(data)
	if err != nil {
		return 0
	}
	hash := h.Sum64()

	return int64(hash)
}

// HashFileConcurrently returns the hash value of a file using concurrent processing
func HashFileConcurrently(filePath string) int64 {
	data, _ := os.ReadFile(filePath)
	return HashDataConcurrently(data)
}

// HashDataConcurrently returns the hash value of a byte slice using concurrent processing
func HashDataConcurrently(data []byte) int64 {
	// Calculate the number of chunks
	numChunks := len(data) / chunkSize
	if len(data)%chunkSize != 0 {
		numChunks++
	}

	// Create a hash function
	h := fnv.New64a()

	// Process each chunk concurrently
	var wg sync.WaitGroup
	wg.Add(numChunks)
	for i := 0; i < numChunks; i++ {
		go func(i int) {
			defer wg.Done()
			start := i * chunkSize
			end := start + chunkSize
			if end > len(data) {
				end = len(data)
			}
			chunk := data[start:end]
			_, err := h.Write(chunk)
			if err != nil {
				return
			}
		}(i)
	}
	wg.Wait()

	// Return the hash value
	hash := h.Sum64()
	return int64(hash)
}
