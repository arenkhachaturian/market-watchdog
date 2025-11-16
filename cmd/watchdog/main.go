package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/joho/godotenv"

    inMemoryAlerts "github.com/arenkhachaturian/market-watchdog/internal/store/memory"
    bot "github.com/arenkhachaturian/market-watchdog/internal/bot"
    "github.com/arenkhachaturian/market-watchdog/internal/outbox"
    "github.com/arenkhachaturian/market-watchdog/internal/core"
    "github.com/arenkhachaturian/market-watchdog/internal/fetcher"
    "github.com/arenkhachaturian/market-watchdog/internal/notifier"
	"github.com/arenkhachaturian/market-watchdog/internal/store"
	"github.com/arenkhachaturian/market-watchdog/internal/config"
)


func main() {
    if err := godotenv.Load(".env"); err != nil {
        log.Printf("[warn] failed to load .env: %v", err)
    }

    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    token := os.Getenv("TELEGRAM_TOKEN")
    if token == "" {
        log.Fatal("[fatal] TELEGRAM_TOKEN is not set")
    }

    log.Println("[startup] market-watchdog booting")

	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("[fatal] failed to load config.json: %v", err)
	}

    // ----------------------------------
    // ALERT REPO
    // ----------------------------------
    alertRepo := inMemoryAlerts.NewAlerts(cfg.DefaultCooldownMin)
    log.Println("[startup] AlertRepo initialized")

    // ----------------------------------
    // TELEGRAM BOT
    // ----------------------------------
    tgBot, err := bot.NewTelegramBot(token, alertRepo)
    if err != nil {
        log.Fatalf("[fatal] failed to create Telegram bot: %v", err)
    }
    log.Println("[startup] Telegram bot created")

    // ----------------------------------
    // SENDER (Telegram)
    // ----------------------------------
    sender, err := notifier.NewTelegramSender(token)
    if err != nil {
        log.Fatalf("[fatal] failed to create Telegram sender: %v", err)
    }
    log.Println("[startup] Telegram sender created")

    // ----------------------------------
    // OUTBOX
    // ----------------------------------
    ob := outbox.New(cfg.OutboxMaxRetry)
    log.Println("[startup] Outbox initialized")

    // ----------------------------------
    // FETCHER + EVALUATOR
    // ----------------------------------
    fetch := fetcher.NewCoinGeckoFetcher()
    eval := core.NewEvaluator(fetch)
    log.Println("[startup] Fetcher + Evaluator initialized")

    // ----------------------------------
    // BOT GOROUTINE
    // ----------------------------------
    go func() {
        log.Println("[bot] starting Telegram listener...")
        if err := tgBot.Run(ctx); err != nil {
            log.Printf("[bot] stopped: %v", err)
        }
    }()

    // ----------------------------------
    // WATCHDOG GOROUTINE
    // ----------------------------------
    go func() {
        log.Println("[watchdog] starting loop...")
		ticker := time.NewTicker(time.Duration(cfg.PollIntervalSeconds) * time.Second)
        defer ticker.Stop()

        for {
            select {
            case <-ctx.Done():
                log.Println("[watchdog] shutdown signal received")
                return
            case now := <-ticker.C:
                runWatcherCycle(ctx, now, alertRepo, eval, ob, sender)
            }
        }
    }()

    log.Println("[startup] system running. Press Ctrl+C to stop.")
    <-ctx.Done()
    log.Println("[shutdown] exiting")
}

func runWatcherCycle(
    ctx context.Context,
    now time.Time,
    repo store.AlertRepo,
    eval *core.Evaluator,
    ob *outbox.Outbox,
    sender notifier.Sender,
) {
    rules, err := repo.ListActive(ctx)
    if err != nil {
        log.Printf("[repo] list error: %v", err)
        return
    }

    if len(rules) == 0 {
        log.Println("[watchdog] no rules found")
        return
    }

    triggered, err := eval.EvaluateRules(rules, now)
    if err != nil {
        log.Printf("[evaluator] error: %v", err)
        return
    }

    if len(triggered) == 0 {
        log.Println("[watchdog] no triggers (conditions/cooldown)")
        return
    }

    for _, t := range triggered {
        log.Printf("[watchdog] TRIGGERED: ID=%d coin=%s threshold=%.4f comp=%v",
            t.ID, t.Coin, t.Threshold, t.Comparator)
    }

    core.EnqueueMatches(ob, triggered)
	core.DeliverOnce(ctx, ob, sender)

    log.Println("[watchdog] cycle complete")
}
