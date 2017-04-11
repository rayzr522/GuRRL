package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"net/http"

	"io/ioutil"

	"encoding/json"

	"github.com/bwmarrin/discordgo"
)

// Command Represents a single command
type Command struct {
	Name        string
	Description string
}

// Discord variables
var (
	token string
	bot   *discordgo.Session
	user  *discordgo.User
	err   error
)

// Command variables
var (
	prefix   string
	commands = []Command{
		{Name: "info", Description: "Shows info about GuRRL"},
		{Name: "help", Description: "Gives you help with using GuRRL"},
		{Name: "cat", Description: "Fetches a random cat picture"},
	}
)

func init() {
	flag.StringVar(&token, "t", "", "The Discord bot token")
	flag.StringVar(&prefix, "p", "^", "The bot's command prefix")
	flag.Parse()

	if token == "" {
		fmt.Println("Please provide a bot token with -t!")
		os.Exit(1)
	}
}

func main() {
	bot, err = discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
	}

	bot.AddHandler(messageCreate)

	err = bot.Open()
	if err != nil {
		fmt.Println("Failed to connect to Discord:", err)
		return
	}

	user, err = bot.User("@me")
	if err != nil {
		fmt.Println("Failed to get bot user:", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	<-make(chan struct{})
	return
}

func isValidCommand(str string) bool {
	for _, element := range commands {
		if str == element.Name {
			return true
		}
	}
	return false
}

func createEmbed(title string, desc string) discordgo.MessageEmbed {
	return discordgo.MessageEmbed{
		Title:       title,
		Description: desc,
		Color:       16711935,
	}
}

func readURL(url string) (body string, err error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || m.Author.ID == user.ID {
		return
	}

	msg := m.ContentWithMentionsReplaced()

	if len(msg) < len(prefix) {
		return
	}

	if msg[:len(prefix)] != prefix {
		return
	}

	split := strings.Split(msg[len(prefix):], " ")

	command := strings.ToLower(split[0])

	if !isValidCommand(command) {
		return
	}

	// args := split[1:]

	if command == "info" {
		embed := createEmbed("", `
***Why the name?***
__GuRRL__ stands for **G**olang **u**sing **R**eal **R**adioactive **L**asers. This bot is a specially crafted well-oiled bot-ing machine, and it's written in [go](http://github.com/golang/go)!

***What can you do?***
Do `+prefix+`help for a list of commands

[:computer: _Read the source, Luke!_](https://github.com/Rayzr522/GuRRL)
`)

		bot.ChannelMessageSendEmbed(m.ChannelID, &embed)
	} else if command == "help" {
		output := ""
		for _, command := range commands {
			output += "`" + prefix + command.Name + "` \u2794 " + command.Description + "\n\n"
		}

		bot.ChannelMessageSend(m.ChannelID, output)
	} else if command == "cat" {
		body, err := readURL("http://random.cat/meow")
		if err != nil {
			fmt.Println("Failed to load cat image:", err)
			bot.ChannelMessageSend(m.ChannelID, "Failed to load cat image!")
			return
		}

		data := new(CatData)
		err = json.NewDecoder(strings.NewReader(body)).Decode(data)
		if err != nil {
			fmt.Println("JSON decoding failed:", err)
			bot.ChannelMessageSend(m.ChannelID, "Failed to decode JSON!")
			return
		}

		embed := createEmbed("", "")
		embed.Image = &discordgo.MessageEmbedImage{URL: data.URL}

		bot.ChannelMessageSendEmbed(m.ChannelID, &embed)
	}
}

// CatData The structure of data returned from http://random.cat/meow
type CatData struct {
	URL string `json:"file"`
}
