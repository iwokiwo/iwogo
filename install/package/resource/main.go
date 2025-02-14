package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run resource/main.go warga/warga.json")
	}

	jsonFile := os.Args[1]

	// Run generate_dto.go
	cmd1 := exec.Command("go", "run", "resource/dto/generate_dto.go", jsonFile)
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	if err := cmd1.Run(); err != nil {
		log.Fatal("Error running generate_dto.go:", err)
	}

	// Run generate_repo.go
	cmd2 := exec.Command("go", "run", "resource/module/generate_repo.go", jsonFile)
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	if err := cmd2.Run(); err != nil {
		log.Fatal("Error running generate_repo.go:", err)
	}

	log.Println("Code generation completed successfully.")
}
