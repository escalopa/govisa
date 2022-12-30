package main

import (
	"context"
	"log"

	"github.com/escalopa/govisa/pkg/logger"

	"github.com/escalopa/govisa/pkg/config"
	"github.com/escalopa/govisa/pkg/security"
	"github.com/escalopa/govisa/telegram/internal/adapters/redis"
	"github.com/escalopa/govisa/telegram/internal/adapters/server"
	"github.com/escalopa/govisa/telegram/internal/application"
	"github.com/escalopa/govisa/telegram/internal/handlers"

	bt "github.com/SakoDroid/telego"
	cfg "github.com/SakoDroid/telego/configs"
)

func main() {
	c := config.NewConfig()

	bot, err := bt.NewBot(cfg.Default(c.Get("BOT_TOKEN")))
	if err != nil {
		log.Fatal(err)
	}
	err = bot.Run()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create encryptor
	encrypt, err := security.NewEncrypter(c.Get("ENCRYPT_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	// Create Redis client
	cache, err := redis.NewClient(c.Get("REDIS_URL"))
	if err != nil {
		log.Fatal(err)
	}

	// Create UserCache
	uc, err := redis.NewUserCache(cache)
	if err != nil {
		log.Fatal(err)
	}

	// Create Logger instance, logs to stdout and Log file
	l, err := logger.New(c.Get("TG_LOG_FILE"))
	if err != nil {
		log.Fatal(err)
	}

	// Connect to Server endpoint
	srv, err := server.NewServer(c.Get("SERVER_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}

	// Create Application UseCase
	app, err := application.New(uc, srv, encrypt)
	if err != nil {
		log.Fatal(err)
	}

	run(bot, app, l, ctx)
}

func run(bot *bt.Bot, app *application.UseCase, l *log.Logger, ctx context.Context) {

	//The general update channel.
	updateChannel := bot.GetUpdateChannel()
	h := handlers.NewBotHandler(bot, app, l, ctx)
	h.Register()

	//Monitors any other update.
	for {
		update := <-*updateChannel
		if update.Message == nil {
			continue
		}
		if update.Message.Chat.Type == "private" {
			h.Help(update)
		} else {
			h.Public(update)
		}
	}
}
