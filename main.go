package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sendReminder/utils"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	tb "gopkg.in/telebot.v3"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	currentPage := flag.Int("page", -1, "current page")
	token := os.Getenv("TELEGRAM_TOKEN")
	dbUrl := os.Getenv("DB_CONNECTION_STRING")
	ctx := context.Background()
	repo, err := utils.NewRepo(ctx, dbUrl)
	flag.Parse()

	bot, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	groups, err := repo.GetAll(ctx)
	for _, group := range groups {
		groupId, err := strconv.ParseInt(group.GroupID, 10, 64)
		if err != nil {
			continue
		}
		chat := &tb.Chat{ID: groupId} // supergroup ID
		if group.PagesTopic != nil {
			createPagesPoll(bot, chat, *currentPage, *group.PagesTopic, group.Type == "normal")
		}
		if group.AthkarTopic != nil {
			createAthkarPoll(bot, chat, *group.AthkarTopic, group.Type == "normal")
		}
	}
}

func createPagesPoll(bot *tb.Bot, chat *tb.Chat, num, topic int, normalType bool) {
	if num == -1 {
		startDate := time.Date(2025, 10, 7, 0, 0, 0, 0, time.UTC)
		today := time.Now().UTC()
		diff := today.Sub(startDate)
		num = int(diff.Hours()/24)*2 + 1
	}

	poll := tb.Poll{
		Type:      tb.PollRegular,
		Question:  fmt.Sprintf("صفحة %d و %d ", num, num+1),
		Anonymous: false,
	}
	poll.AddOptions("تم", "لسا بس اليوم أكيد إن شاء الله")
	// SendOptions has ThreadID for forum topics
	options := &tb.SendOptions{}
	if !normalType {
		options.ThreadID = topic
	}
	if _, err := poll.Send(bot, chat, options); err != nil {
		log.Fatal(err)
	}
	images := tb.Album{
		&tb.Photo{
			File: tb.FromDisk(fmt.Sprintf("pages/%03d.png", num)),
		},
		&tb.Photo{
			File: tb.FromDisk(fmt.Sprintf("pages/%03d.png", num+1)),
		},
	}
	bot.SendAlbum(chat, images, options)

	bot.Send(chat, "@The_Hallak, @mohammadnahhas, @ammar_alphin, @EyadBT, @Othman_Tomhe, @ABD_ULKARIM_JAMAL, @ali_h_h_13, @falmouine, @Besher_Zaiback, @Mhd0mart0meh, @Wael_Zaiback, @Jaafar_86, @The_Soviet_Cat, @M123459, @FutureHustler, @mahdi_alhamid, @JaberSH1, @kinannotfound.", options)
}

func createAthkarPoll(bot *tb.Bot, chat *tb.Chat, topic int, normalType bool) {
	data, err := os.ReadFile("athkar.txt")
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")
	{
		var filteredLines []string
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				filteredLines = append(filteredLines, line)
			}
		}

		lines = filteredLines
	}

	poll := tb.Poll{
		Type:            tb.PollRegular,
		Question:        "أذكار اليوم",
		Anonymous:       false,
		MultipleAnswers: true,
	}
	rand.Shuffle(len(lines), func(i, j int) {
		lines[i], lines[j] = lines[j], lines[i]
	})
	poll.AddOptions(lines[0:10]...)
	// SendOptions has ThreadID for forum topics
	options := &tb.SendOptions{}
	if !normalType {
		options.ThreadID = topic
	}
	if _, err := poll.Send(bot, chat, options); err != nil {
		log.Fatal(err)
	}
	bot.Send(chat, "@The_Hallak, @mohammadnahhas, @ammar_alphin, @EyadBT, @Othman_Tomhe, @ABD_ULKARIM_JAMAL, @ali_h_h_13, @falmouine, @Besher_Zaiback, @Mhd0mart0meh, @Wael_Zaiback, @Jaafar_86, @The_Soviet_Cat, @M123459, @FutureHustler, @mahdi_alhamid, @JaberSH1, @kinannotfound.", options)
}
