package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	guildID = "1202350149964140584" // Your Guild-ID here

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Antwortet mit Pong!",
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

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds

	dg.AddHandler(interactionCreate)

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		for _, v := range commands {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, v)
			if err != nil {
				fmt.Printf("Command konnte nicht registriert werden: %v", err)
			}
		}
		fmt.Println("Slash-Commands registriert.")
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
	/*case "wtf":
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hey du! Ich bin ein digitaler Klon " +
				"von xn4k, der existiert, bis er sein Bewusstsein " +
				"in die Cloud hochladen kann. Sollte ich jemals Amok" +
				" laufen, richte bitte alle Beschwerden direkt an xn4k" +
				" – sie werden garantiert spätestens bis zum Ende der Welt bearbeitet. 🌎🤖",
		},
	})*/

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
	}
}
