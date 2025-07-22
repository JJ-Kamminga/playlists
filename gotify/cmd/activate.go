package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(activateCmd)
}

var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Get a fresh bearer token",
	Long:  `Get a fresh bearer token and verify it`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("App version none")
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		manager := &SpotifyTokenManager{
			ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		}

		// make a simple request to verify bearer token works
		resp, err := manager.DoSpotifyAPIRequest("https://api.spotify.com/v1/browse/categories")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("Status:", resp.Status)
	},
}
