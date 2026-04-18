# Nova Create-a-thon 2026

## Requirements

- Go 1.25+
- MySQL

## Setup

1. Create the database:
```bash
mysql -u root -e "CREATE DATABASE IF NOT EXISTS nova;"
```

2. Copy and configure environment variables:
```bash
cp .env.example .env
```

3. Start the backend (runs migrations automatically):
```bash
make run
```

4. (Optional) Seed with sample data:
```bash
make seed
```

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DB_USERNAME` | `root` | MySQL user |
| `DB_PASSWORD` | — | MySQL password |
| `DB_HOST` | `127.0.0.1` | MySQL host |
| `DB_PORT` | `3306` | MySQL port |
| `DB_NAME` | `nova` | Database name |
| `APP_ENV` | — | Set to `production` to block seeding |

## Resetting the Database

```bash
mysql -u root -e "DROP DATABASE nova; CREATE DATABASE nova;"
make run   # re-runs migrations
make seed  # re-seeds
```

## Seeding on a Deployed Backend

If only the `backend/` directory is deployed (e.g. DigitalOcean), place a `.env` in the backend root and run:

```bash
cd backend && go run ./cmd/seed/main.go
```

## License

[MIT](LICENSE)