package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getPlaylistsCmd)
}

var getPlaylistsCmd = &cobra.Command{
	Use:   "getplaylists",
	Short: "Get all playlists for a user",
	Long:  `uuh`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		getPlaylists(args[0])
	},
}

func getPlaylists(userId string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	manager := &SpotifyTokenManager{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
	}

	type SimplifiedPlaylistObject struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Owner       struct {
			ID   string          `json:"id"`
			Rest json.RawMessage `json:"-"`
		} `json:"owner"`
		Tracks struct {
			Href  string `json:"href"`
			Total int    `json:"total"`
		} `json:"tracks"`
		Rest json.RawMessage `json:"-"`
	}

	type PartialResponse struct {
		Href   string                     `json:"href"`
		Limit  int                        `json:"limit"`
		Next   string                     `json:"next"`
		Offset int                        `json:"offset"`
		Total  int                        `json:"total"`
		Items  []SimplifiedPlaylistObject `json:"items"`
		Rest   json.RawMessage            `json:"-"`
	}
	var url = fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", userId)

	resp, err := manager.DoSpotifyAPIRequest(url)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var partialResp PartialResponse
	if err := json.Unmarshal(body, &partialResp); err != nil {
		panic(err)
	}

	fmt.Println("Status:", resp.Status)
	playlists := partialResp.Items
	nextURL := partialResp.Next

	file, _ := os.Create("playlists.csv")

	for nextURL != "" {
		resp, err := manager.DoSpotifyAPIRequest(nextURL)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var nextResp PartialResponse
		if err := json.Unmarshal(body, &nextResp); err != nil {
			panic(err)
		}

		playlists = append(playlists, nextResp.Items...)
		nextURL = nextResp.Next
	}

	partialResp.Items = playlists

	fmt.Println("Playlists:", playlists)

	defer file.Close()
	gocsv.MarshalFile(&playlists, file)
}
