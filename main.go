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

	tb "gopkg.in/telebot.v3"
)

func main() {
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
			createPagesPoll(bot, chat, *currentPage, *group.PagesTopic)
		}
		if group.AthkarTopic != nil {
			createAthkarPoll(bot, chat, *group.AthkarTopic)
		}
	}
}

func createPagesPoll(bot *tb.Bot, chat *tb.Chat, num, topic int) {
	if num == -1 {
		startDate := time.Date(2025, 10, 7, 0, 0, 0, 0, time.UTC)
		today := time.Now().UTC()
		diff := today.Sub(startDate)
		num = int(diff.Hours()/24)*2 + 1
	}

	poll := tb.Poll{
		Type:      tb.PollRegular,
		Question:  fmt.Sprintf("صفحى %d و %d ", num, num+1),
		Anonymous: false,
	}
	poll.AddOptions("تم", "لسا بس اليوم أكيد إن شاء الله")
	// SendOptions has ThreadID for forum topics
	if _, err := poll.Send(bot, chat, &tb.SendOptions{ThreadID: topic}); err != nil {
		log.Fatal(err)
	}
}

func createAthkarPoll(bot *tb.Bot, chat *tb.Chat, topic int) {
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
	if _, err := poll.Send(bot, chat, &tb.SendOptions{ThreadID: topic}); err != nil {
		log.Fatal(err)
	}

}
