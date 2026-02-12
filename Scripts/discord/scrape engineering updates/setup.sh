#!/bin/bash
# Setup script for Go Discord Message Scraper

echo "=== Go Discord Message Scraper Setup ==="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "ERROR: Go is not installed. Please install Go 1.21 or later."
    echo "Visit: https://golang.org/dl/"
    exit 1
fi

echo "1. Installing Go modules..."
go mod download

echo "2. Building the scraper..."
go build -o discord-scraper main.go

echo "3. Creating .env file from example..."
if [ ! -f .env ]; then
    cp .env.example .env
    echo "Created .env file. Please edit it with your Discord credentials."
else
    echo ".env file already exists."
fi

echo ""
echo "=== NEXT STEPS ==="
echo "1. Edit .env file with your credentials:"
echo "   nano .env"
echo "   - Set DISCORD_TOKEN (from Discord Developer Portal)"
echo "   - Set DISCORD_CHANNEL_ID (right-click channel → Copy ID)"
echo ""
echo "2. Test the scraper:"
echo "   ./discord-scraper"
echo ""
echo "3. Set up hourly cron job:"
echo "   crontab -e"
echo "   Add: 0 * * * * cd /full/path/to/directory && ./discord-scraper"
echo ""
echo "4. Query the database:"
echo "   sqlite3 discord_messages.db"
echo "   .tables"
echo "   SELECT * FROM messages ORDER BY timestamp DESC LIMIT 5;"