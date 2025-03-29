package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"sort"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	guildID = "1202350149964140584" // Your Guild-ID here

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "msgcount",
			Description: "Zeigt die Top 5 aktivsten User",
		},
		{
			Name:        "ping",
			Description: "Antwortet mit Pong!",
		},
		{
			Name:        "gibmir",
			Description: "Gibt dir böse!",
		},
		{
			Name:        "wtf",
			Description: "Antwortet mit!",
		},
		{
			Name:        "ask",
			Description: "Frage DeepSeek etwas!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "frage",
					Description: "Your question to DeepSeek",
					Required:    true,
				},
			},
		},
	}
)

var msgCount = make(map[string]int)

// DeepSeek API Response-Struktur
type DeepSeekAPIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func frageDeepSeek(prompt string) (string, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		fmt.Println("Bitte setze die OPENROUTER_API_KEY Umgebungsvariable.")
		fmt.Println("Aktueller Wert der Umgebungsvariable:", apiKey)
	}

	url := "https://openrouter.ai/api/v1/chat/completions"

	payload, _ := json.Marshal(map[string]interface{}{
		"model": "deepseek/deepseek-chat-v3-0324:free",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	})

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Debugging-Ausgabe (hilfreich!)
	fmt.Println("DeepSeek API Antwort:", string(body))

	var response DeepSeekAPIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("DeepSeek lieferte keine Antwort")
	}

	return response.Choices[0].Message.Content, nil
}

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		fmt.Println("Bitte setze die DISCORD_TOKEN Umgebungsvariable.")
		return
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		fmt.Println("Bitte setze die OPENROUTER_API_KEY Umgebungsvariable.")
		return
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Fehler beim Erstellen der Session:", err)
		return
	}
	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds

	dg.AddHandler(interactionCreate)

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		validName := regexp.MustCompile(`^[\w-]{1,32}$`)

		for _, cmd := range commands {
			if !validName.MatchString(cmd.Name) {
				fmt.Printf("❌ Ungültiger Command-Name: '%s' – wird übersprungen\n", cmd.Name)
				continue
			}

			_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)
			if err != nil {
				fmt.Printf("Command konnte nicht registriert werden: %v\n", err)
			} else {
				fmt.Printf("✅ Slash-Command '%s' registriert.\n", cmd.Name)
			}
		}

		fmt.Println("Alle gültigen Slash-Commands verarbeitet.")
	})

	err = dg.Open()
	if err != nil {
		fmt.Println("Fehler beim Öffnen der Verbindung:", err)
		return
	}

	fmt.Println("Bot läuft. Drücke STRG+C zum Beenden.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	dg.Close()
}

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "ping":
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Pong! 🏓",
			},
		})
	case "wtf":
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Hey du! Ich bin ein digitaler Klon " +
					"von xn4k, der existiert, bis er sein Bewusstsein " +
					"in die Cloud hochladen kann. Sollte ich jemals Amok" +
					" laufen, richte bitte alle Beschwerden direkt an xn4k" +
					" – sie werden garantiert spätestens bis zum Ende der Welt bearbeitet. 🌎🤖",
			},
		})

	case "gibmir":
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "legs spread, wrists bound, eyes blindfolded," +
					" soft whimpers, silk sliding, collar tight, thighs" +
					" trembling, breath held, nails digging in, tongue" +
					" tracing, floor wet, ceiling fan spinning, mirror foggy" +
					", candle dripping, knees bruised, throat sore, red marks," +
					" mouth full, overstimulated, can’t speak, can't think, don't" +
					" stop, don't you dare, hair pulled, bitten lips, stuck" +
					" in position, tears on cheeks, shaking legs, bed frame" +
					" knocking, neighbors hearing, don’t care, spit dripping," +
					" fingers everywhere, soaked sheets, locked doors, screaming" +
					" silence, nightstand shaking, legs twitching, no mercy, more," +
					" deeper, slower, harder, grip tighter, lights flickering, air" +
					" heavy, skin flushed, music echoing, time frozen, the edge," +
					" again, again, can't move, being held down, whispering filth," +
					" moaning please, tongue inside, legs on shoulders, hips grinding," +
					" face buried, can’t breathe, want more, aching, dripping, pulsing," +
					" obsessed, ruined, used, worshipped, destroyed, begging, shaking," +
					" craving, owned, devoured, filled, full, gone",
			},
		})

	case "ask":
		frage := i.ApplicationCommandData().Options[0].StringValue()

		// Discord direkt mitteilen, dass die Antwort etwas dauert.
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})

		// Asynchron Abfrage an DeepSeek
		go func() {
			antwort, err := frageDeepSeek(frage)
			if err != nil {
				antwort = fmt.Sprintf("Fehler: %v", err)
			}
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &antwort,
			})
		}()
	case "msgcount":
		type userStat struct {
			ID    string
			Count int
		}

		// Map in Slice umwandeln
		var stats []userStat
		for id, count := range msgCount {
			stats = append(stats, userStat{ID: id, Count: count})
		}

		// Nach Anzahl sortieren
		sort.Slice(stats, func(i, j int) bool {
			return stats[i].Count > stats[j].Count
		})

		// Nur Top 5
		max := 5
		if len(stats) < 5 {
			max = len(stats)
		}

		output := "🏆 **Top Nachrichtenschreiber:**\n"
		for i := 0; i < max; i++ {
			user, err := s.User(stats[i].ID)
			if err != nil {
				continue
			}
			output += fmt.Sprintf("%d. %s – %d Nachrichten\n", i+1, user.Username, stats[i].Count)
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: output,
			},
		})

	}

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Bot-Nachrichten ignorieren
	if m.Author.Bot {
		return
	}

	msgCount[m.Author.ID]++
}
