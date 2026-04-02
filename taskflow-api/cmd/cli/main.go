package main

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "taskflow-cli",
	Short: "TaskFlow API CLI",
	Long:  "A simple CLI for interacting with the TaskFlow API",
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("erreur lors de l'exécution de la commande: %v", err)
	}
}
