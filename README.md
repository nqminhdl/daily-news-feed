# Daily News Feed

This is small application designed to keep me informed about the latest news and events for those interested domains.

## What this small application can do

- Define news by category (e.g technology, daily news, view from experity)
- Define retention for each category
  - [ ] Global Retention
  - [ ] Category Retention
- Support different backend (e.g file system, sqlite, mysql, postgres)
  - [x] File System
  - [x] Sqlite
- Receiving channels
  - Telgram:
    - [x] Support hide secret Telegram bot token and Channel ID
    - [x] Support multiple channels
  - Grafana Cloud:
    - [x] Support receiver to remote write metrics to Prometheus
    - [x] Visualize to dashboard

## Get started

- [How to configure](docs/how-to-configure.md)
- [How to create Telegram group and grant access to bot](docs/telegram.md)
- [How to push metrics to Grafana Cloud](docs/grafana-cloud.md)

## Build binary

```bash
make build app=<app_name>
```
