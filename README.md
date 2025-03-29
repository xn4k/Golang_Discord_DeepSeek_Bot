# 🚀 **Discord x DeepSeek AI Bot**

Ein smarter Discord-Bot, der die Kraft von **DeepSeek (via OpenRouter)** nutzt, um auf deine Fragen zu antworten, coole Features bietet und natürlich mit Humor punktet!

---

## 💡 **Was kann dieser Bot?**

- **Slash-Befehle:**
  - `/ping`: Prüft, ob der Bot aktiv ist. Antwort: Pong! 🏓
  - `/ask <Frage>`: Fragt DeepSeek AI und erhält eine intelligente Antwort.
  - `/wtf`: Gibt dir eine witzige Beschreibung des Bots selbst.

---

## 🛠️ **Technologien & Services**

- **Programmiersprache:** Go (Golang)
- **Discord-Bibliothek:** discordgo
- **AI-Modell:** DeepSeek (via OpenRouter API)

---

## ⚙️ **Einrichtung & Start**

### 📋 **Voraussetzungen**

- Go installiert (https://golang.org/dl)
- Discord-Bot erstellt (https://discord.com/developers/applications)
- API-Key bei OpenRouter.ai erstellt (https://openrouter.ai)

### 📥 **Installation**

Clone zuerst dieses Repository:

```bash
git clone <dein_repo_url>
cd dein_repo_ordner
```

Module installieren:

```bash
go mod tidy
```

### 🔑 **Umgebungsvariablen setzen**

Setze diese Variablen:

**Windows PowerShell:**
```powershell
$env:DISCORD_TOKEN = "dein_discord_token"
$env:OPENROUTER_API_KEY = "dein_openrouter_api_key"
```

**Linux / MacOS:**
```bash
export DISCORD_TOKEN="dein_discord_token"
export OPENROUTER_API_KEY="dein_openrouter_api_key"
```

### 🚀 **Starten des Bots**

```bash
go run main.go
```

---

## 📌 **Slash-Commands**

- `/ping`
  - Checkt den Status des Bots.

- `/ask <frage>`
  - Stelle DeepSeek AI jede beliebige Frage.

- `/wtf`
  - Erfahre, wer oder was dieser Bot eigentlich ist.

---

## ⚠️ **Wichtig!**

- Teile niemals deinen Discord- oder OpenRouter API-Key öffentlich!
- Stelle sicher, dass dein OpenRouter-Konto Guthaben besitzt (kostenlose Nutzung mit gewissen Modellen möglich).

---

## 🤖 **Über diesen Bot**

> Hey du fragst dich, WTF ich bin? Ich bin eine digitale Version von xn4k, gecodet und erschaffen, um seine gottesfürchtige Programmierer-Präsenz zu imitieren, bis die Singularität eintrifft und er seinen Verstand in eine KI klonen kann. Sollte ich jemals Amok laufen, leite alle Beschwerden direkt an xn4k weiter – sie werden garantiert spätestens bis zum Ende der Welt bearbeitet. 🌎🤖✨
