package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/gocolly/colly/v2"
	"github.com/puckzxz/dismand"
)

type embedImage struct {
	URL string `json:"url"`
}

type embed struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Image       embedImage `json:"image"`
}

type webhookPayload struct {
	Color int64   `json:"color"`
	Embed []embed `json:"embeds"`
}

var webhookURL string = os.Getenv("WEBHOOK_URL")

func labCommand(ctx *dismand.Context, args []string) {
	c := colly.NewCollector()
	embed := embed{}

	c.OnHTML("#notesImg", func(e *colly.HTMLElement) {
		embed.Image.URL = e.Attr("src")
	})
	c.OnHTML(".comment-content", func(e *colly.HTMLElement) {
		embed.Description = fmt.Sprintf("```%s```", e.Text)
	})
	c.Visit("https://www.poelab.com/wfbra/")

	embed.Title = fmt.Sprintf("PoE Uber Lab -  %s", time.Now().Format("02-Jan-06"))
	wp := webhookPayload{}
	wp.Embed = append(wp.Embed, embed)
	wp.Color = 4030808

	data, err := json.Marshal(wp)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}

	_, err = http.Post(webhookURL, "application/json", bytes.NewBuffer(data))

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}
}

func main() {
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("TOKEN"),
		RejectEvents: []string{
			disgord.EvtPresenceUpdate,
			disgord.EvtGuildMemberAdd,
			disgord.EvtGuildMemberUpdate,
			disgord.EvtGuildMemberRemove,
		},
	})
	d := dismand.New(client, &dismand.Config{
		Prefix: "!",
	})

	defer client.Gateway().StayConnectedUntilInterrupted()

	d.On("lab", labCommand)

	client.Gateway().MessageCreate(d.MessageHandler)
}
