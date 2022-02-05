package discord

import (
	"bytes"
	"image"
	"image/color"
	_ "image/jpeg"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/anthonynsimon/bild/transform"
	"github.com/bwmarrin/discordgo"
	"github.com/ftqo/kirby/database"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

const (
	PfpSize = 256

	width   = 848
	height  = 477
	margin  = 15
	res     = 1
	titlSz  = 100
	stitlSz = 80
)

type welcomeMessageInfo struct {
	mention   string
	nickname  string
	username  string
	guildName string
	avatarURL string
}

func generateWelcomeMessage(gw database.GuildWelcome, wi welcomeMessageInfo) discordgo.MessageSend {
	var msg discordgo.MessageSend

	r := strings.NewReplacer("%mention%", wi.mention, "%nickname", wi.nickname, "%username%", wi.username, "%guild%", wi.guildName)
	gw.Text = r.Replace(gw.Text)
	gw.ImageText = r.Replace(gw.ImageText)

	msg.Content = gw.Text

	switch gw.Type {
	case "embed":
		log.Println("embedded welcome messages not implemented, sending plain")
	case "image":
		cv := canvas.New(width, height)
		ctx := canvas.NewContext(cv)
		resp, err := http.Get(wi.avatarURL)
		if err != nil {
			log.Printf("failed to get avatar URL: %v", err)
		}
		defer resp.Body.Close()

		pfpBuf := &bytes.Buffer{}
		_, err = io.Copy(pfpBuf, resp.Body)
		if err != nil {
			log.Printf("failed to copy pfp to buffer: %v", err)
		}
		rawPfp, _, err := image.Decode(pfpBuf)
		if err != nil {
			log.Printf("failed to decode profile picture: %v", err)
		}
		var pfp image.Image
		if rawPfp.Bounds().Max.X != PfpSize {
			pfp = image.Image(transform.Resize(rawPfp, PfpSize, PfpSize, transform.Linear))
		} else {
			pfp = rawPfp
		}

		ctx.DrawImage(0, 0, h.Images[gw.Image], res)

		// BACKGROUND LOADED

		ctx.SetFillColor(color.RGBA{50, 45, 50, 130})
		ctx.DrawPath(margin, margin, canvas.Rectangle(width-(2*margin), height-(2*margin)))

		// BACKGROUND OVERLAY LOADED

		ctx.DrawImage(width/2-PfpSize/2, height/2-PfpSize/2, pfp, res)

		// PFP LOADED

		coolvetica := canvas.NewFontFamily("coolvetica")
		err = coolvetica.LoadFont(h.Fonts["coolvetica"], 0, canvas.FontRegular)
		if err != nil {
			log.Printf("failed to load font: %v", err)
		}
		coolFace := coolvetica.Face(titlSz, canvas.White, canvas.FontRegular, canvas.FontNormal)
		ctx.DrawText(width/2, height/2, canvas.NewTextLine(coolFace, gw.ImageText, canvas.Center))

		// TITLE LOADED

		buf := &bytes.Buffer{}
		cw := renderers.JPEG()
		cw(buf, cv)
		f := &discordgo.File{
			Name:        "welcome_" + wi.nickname + ".jpg",
			ContentType: "image/jpeg",
			Reader:      buf,
		}
		msg.Files = append(msg.Files, f)
	}
	return msg
}
