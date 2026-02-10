package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: png2ico <input.png> <output.ico>")
		return
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	pngData, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}

	icoFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Failed to create output file: %v\n", err)
		return
	}
	defer icoFile.Close()

	// ICO Header
	// Reserved (2 bytes)
	// Type (2 bytes, 1 = icon)
	// Count (2 bytes, 1 image)
	header := []byte{0, 0, 1, 0, 1, 0}
	if _, err := icoFile.Write(header); err != nil {
		fmt.Printf("Failed to write header: %v\n", err)
		return
	}

	// Image Directory Entry
	// Width (1 byte, 0 = 256, but we have 64)
	// Height (1 byte)
	// ColorCount (1 byte, 0 if >= 8bpp)
	// Reserved (1 byte, 0)
	// Planes (2 bytes, 1)
	// BitCount (2 bytes, 32)
	// BytesInRes (4 bytes)
	// ImageOffset (4 bytes)

	width := byte(64) // Assuming 64x64
	height := byte(64)

	entry := make([]byte, 16)
	entry[0] = width
	entry[1] = height
	entry[2] = 0
	entry[3] = 0
	binary.LittleEndian.PutUint16(entry[4:], 1)                    // Planes
	binary.LittleEndian.PutUint16(entry[6:], 32)                   // BitCount
	binary.LittleEndian.PutUint32(entry[8:], uint32(len(pngData))) // Size
	binary.LittleEndian.PutUint32(entry[12:], 22)                  // Offset (6 header + 16 entry)

	if _, err := icoFile.Write(entry); err != nil {
		fmt.Printf("Failed to write directory entry: %v\n", err)
		return
	}

	// Write PNG data
	if _, err := icoFile.Write(pngData); err != nil {
		fmt.Printf("Failed to write PNG data: %v\n", err)
		return
	}

	fmt.Println("Successfully created", outputPath)
}
