package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Message struct {
	MessageID     string
	SenderID      string
	Content       string
	Timestamp     string
	EditTimestamp string
	ChannelID     string
}

type JsonMessage struct {
	ID              string        `json:"id"`
	Type            string        `json:"type"`
	Timestamp       string        `json:"timestamp"`
	TimestampEdited string        `json:"timestampEdited"`
	IsPinned        bool          `json:"isPinned"`
	Content         string        `json:"content"`
	Author          Author        `json:"author"`
	Attachments     []interface{} `json:"attachments"`
	Embeds          []interface{} `json:"embeds"`
	Stickers        []interface{} `json:"stickers"`
	Reactions       []interface{} `json:"reactions"`
	Mentions        []interface{} `json:"mentions"`
}

type Author struct {
	ID            string      `json:"id"`
	Name          interface{} `json:"name"`
	Discriminator interface{} `json:"discriminator"`
	Nickname      interface{} `json:"nickname"`
	Color         interface{} `json:"color"`
	IsBot         bool        `json:"isBot"`
	AvatarURL     interface{} `json:"avatarUrl"`
}

type Guild struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IconUrl string `json:"iconUrl"`
}

type Channel struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	CategoryID string      `json:"categoryId"`
	Category   string      `json:"category"`
	Name       string      `json:"name"`
	Topic      interface{} `json:"topic"`
}

type DateRange struct {
	After  interface{} `json:"after"`
	Before interface{} `json:"before"`
}

type JsonOutput struct {
	Guild        Guild         `json:"guild"`
	Channel      Channel       `json:"channel"`
	DateRange    DateRange     `json:"dateRange"`
	Messages     []JsonMessage `json:"messages"`
	MessageCount int           `json:"messageCount"`
}

var (
	versionFlag = flag.Bool("v", false, "prints the version")
	dbPath      = flag.String("in", "input.db", "path to the input SQLite database")
	jsonPath    = flag.String("out", "output.json", "path to the output JSON file")
	channelID   = flag.String("channel", "", "channel ID")
)

var channelName string

const version = "0.2"

func main() {
	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		return
	}

	if *channelID == "" {
		log.Fatal("Channel ID is required")
	}

	fmt.Println("Opening database...")
	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Getting messages from database...")

	rows, err := db.Query(`
	SELECT m.message_id, m.sender_id, m.channel_id, m.text, m.timestamp, COALESCE(e.edit_timestamp, ''), c.name
	FROM messages m
	LEFT JOIN edit_timestamps e ON m.message_id = e.message_id
	LEFT JOIN channels c ON m.channel_id = c.id
	WHERE m.channel_id = ?
`, *channelID)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	messages := make([]Message, 0)
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.MessageID, &msg.SenderID, &msg.ChannelID, &msg.Content, &msg.Timestamp, &msg.EditTimestamp, &channelName); err != nil {
			log.Fatal(err)
		}

		messages = append(messages, msg)
	}

	jsonMessages := make([]JsonMessage, 0)
	for _, msg := range messages {
		jsonMsg := JsonMessage{
			ID:              msg.MessageID,
			Type:            "Default",
			Timestamp:       msg.Timestamp,
			TimestampEdited: msg.EditTimestamp,
			IsPinned:        false,
			Content:         msg.Content,
			Author: Author{
				ID:            msg.SenderID,
				Name:          nil,
				Discriminator: nil,
				Nickname:      nil,
				Color:         nil,
				IsBot:         false,
				AvatarURL:     nil,
			},
			Attachments: []interface{}{},
			Embeds:      []interface{}{},
			Stickers:    []interface{}{},
			Reactions:   []interface{}{},
			Mentions:    []interface{}{},
		}
		jsonMessages = append(jsonMessages, jsonMsg)
	}

	jsonOutput := JsonOutput{
		Guild: Guild{
			ID:      "0",
			Name:    "Direct Messages",
			IconUrl: "null",
		},
		Channel: Channel{
			ID:         *channelID,
			Type:       "DirectTextChat",
			CategoryID: "0",
			Category:   "Private",
			Name:       channelName,
			Topic:      nil,
		},
		DateRange: DateRange{
			After:  nil,
			Before: nil,
		},
		Messages:     jsonMessages,
		MessageCount: len(jsonMessages),
	}

	fmt.Println("Formatting JSON...")
	jsonOutputBytes, err := json.MarshalIndent(jsonOutput, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Writing to file...")
	file, err := os.Create(*jsonPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(jsonOutputBytes)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done.")
}
