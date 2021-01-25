package main

import (
	"fmt"
	"os"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/gocolly/colly/v2"
	"github.com/puckzxz/dismand"
)

func labCommand(ctx *dismand.Context, args []string) {
	c := colly.NewCollector()
	embed := &disgord.Embed{}

	c.OnHTML("#notesImg", func(e *colly.HTMLElement) {
		embed.Image = &disgord.EmbedImage{
			URL: e.Attr("src"),
		}
	})
	c.OnHTML(".comment-content", func(e *colly.HTMLElement) {
		embed.Description = fmt.Sprintf("```%s```", e.Text)
	})
	c.Visit("https://www.poelab.com/wfbra/")

	embed.Title = fmt.Sprintf("PoE Uber Lab -  %s", time.Now().Format("02-Jan-06"))
	embed.Color = 4030808

	_, err := ctx.Client.SendMsg(ctx.Message.ChannelID, embed)

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
