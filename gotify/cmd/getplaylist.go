package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getPlaylistCmd)
}

var getPlaylistCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of my app's",
	Long:  `All software has versions. This is my app's`,
	Run: func(cmd *cobra.Command, args []string) {
		getPlaylist()
	},
}

func getPlaylist() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	manager := &SpotifyTokenManager{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
	}

	resp, err := manager.DoSpotifyAPIRequest("https://api.spotify.com/v1/playlists")

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
}
