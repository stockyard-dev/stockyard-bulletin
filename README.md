# Stockyard Bulletin

**Email list and newsletter tool — manage subscribers, send campaigns via your own SMTP**

Part of the [Stockyard](https://stockyard.dev) family of self-hosted developer tools.

## Quick Start

```bash
docker run -p 9260:9260 -v bulletin_data:/data ghcr.io/stockyard-dev/stockyard-bulletin
```

Or with docker-compose:

```bash
docker-compose up -d
```

Open `http://localhost:9260` in your browser.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `9260` | HTTP port |
| `DATA_DIR` | `./data` | SQLite database directory |
| `BULLETIN_LICENSE_KEY` | *(empty)* | Pro license key |

## Free vs Pro

| | Free | Pro |
|-|------|-----|
| Limits | 1 list, 500 subscribers | Unlimited lists and subscribers |
| Price | Free | $4.99/mo |

Get a Pro license at [stockyard.dev/tools/](https://stockyard.dev/tools/).

## Category

Creator & Small Business

## License

Apache 2.0
