# Market Watchdog

A lightweight Go-based monitoring bot that tracks cryptocurrency market signals and issues alerts via Telegram in near real-time.

## Features
- Fetches live market data and evaluates custom rules (goroutines + WaitGroups).
- Uses templated notifications and an in-memory repository to manage state and cooldowns.
- Supports Telegram Bot API for push alerts (private chat or group).
- Designed for modular extension: add new fetchers, evaluators, or output channels with minimal code.

## Tech Stack
- Go modules (`go.mod`) for dependency management.
- Concurrency: goroutines + WaitGroups for non-blocking data fetch + evaluation.
- Templating: `text/template` or `html/template` for customizable alert messages.
- In-memory store for recent signals / cooldown logic.
- External integrations: Telegram Bot API, REST APIs for market data (e.g., CoinGecko).

## Getting Started
1. Clone the repository:
```bash
git clone https://github.com/arenkhachaturian/market-watchdog.git
cd market-watchdog
```

2. Configure the bot:
- Create a `.env` file with your `TELEGRAM_TOKEN`.

3. Build and run:
```bash
go build -o ./cmd/watchdog .
./watchdog
```

## 🛠️ Roadmap / TODO

- [ ] Add support for persistent storage (e.g., SQLite or PostgreSQL) instead of in-memory repository for signal history and cooldown tracking.  
- [ ] Implement advanced evaluator module with customizable rule engine (e.g., DSL or JSON-based conditions) for smarter alerts.  
- [ ] Expose REST API for external dashboards / integrations so users can query status and metrics.  
- [ ] Create configuration UI (web or CLI) to simplify setup of symbols, thresholds, alert channels, and cooldowns.  
- [ ] Add plugin architecture so users can easily add new market data fetchers, evaluators or output channels (e.g., Slack, Email).  
- [ ] Write end-to-end tests covering alert triggers, cooldown logic, and concurrency edge-cases.  
- [ ] Containerize whole application (Docker image + Helm chart) for easy deployment.  


## Usage
Run the compiled binary and watch for Telegram notifications according to your configured thresholds and intervals.

## Contribution & License
Contributions welcome. Please open issues or pull requests for enhancements (new fetcher, evaluator, etc.).  
Licensed under MIT.
