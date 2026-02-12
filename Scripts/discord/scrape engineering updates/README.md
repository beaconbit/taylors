# Go Discord Message Scraper

A Go-based Discord bot that scrapes messages from the last hour and stores them in SQLite. Designed for hourly cron execution.

## Features

- **Go implementation** - Fast, compiled binary
- **Last hour only** - Uses Discord API pagination with time filtering
- **SQLite storage** - Simple file-based database
- **Duplicate prevention** - UNIQUE constraint on Discord message IDs
- **Efficient** - Batched database transactions
- **Production-ready** - Proper error handling and logging

## Prerequisites

- Go 1.21 or later
- Discord bot token with MESSAGE CONTENT INTENT enabled
- Channel ID to scrape

## Quick Start

1. **Clone and setup:**
   ```bash
   chmod +x setup.sh
   ./setup.sh
   ```

2. **Configure credentials in `.env` file:**
   ```bash
   cp .env.example .env
   nano .env
   ```
   Set:
   ```env
   DISCORD_TOKEN=your_bot_token_here
   DISCORD_CHANNEL_ID=123456789012345678
   ```

3. **Build and test:**
   ```bash
   go build -o discord-scraper main.go
   ./discord-scraper
   ```

4. **Set up cron job (hourly):**
   ```bash
   crontab -e
   # Add: 0 * * * * /full/path/to/discord-scraper
   ```

## How It Works

### Time-Based Filtering
The script calculates `time.Now().Add(-1 * time.Hour)` and uses Discord's pagination API to fetch messages newer than that timestamp.

### Database Schema
```sql
CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    discord_message_id TEXT UNIQUE,
    channel_id TEXT,
    author_id TEXT,
    author_name TEXT,
    content TEXT,
    timestamp DATETIME,
    attachments_count INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Rate Limiting
- Respects Discord API rate limits
- Small delay between pagination requests
- Uses transactions for efficient database writes

## Configuration

### Environment Variables
The script now uses environment variables loaded from `.env` file:
```go
DiscordToken = os.Getenv("DISCORD_TOKEN")
ChannelID = os.Getenv("DISCORD_CHANNEL_ID")
```

### Multiple Channels
To scrape multiple channels, modify the main function:
```go
channelIDs := []string{"channel1", "channel2", "channel3"}
for _, channelID := range channelIDs {
    messages, err := fetchMessages(dg, channelID, oneHourAgo)
    // ... save messages
}
```

## Querying Data

### SQL Queries
```bash
sqlite3 discord_messages.db

-- Recent messages
SELECT * FROM messages ORDER BY timestamp DESC LIMIT 10;

-- Daily statistics
SELECT date(timestamp), COUNT(*) 
FROM messages 
GROUP BY date(timestamp);

-- User activity
SELECT author_name, COUNT(*) as message_count
FROM messages
GROUP BY author_name
ORDER BY message_count DESC;
```

### Export to CSV
```bash
sqlite3 -header -csv discord_messages.db \
  "SELECT timestamp, author_name, content FROM messages" > messages.csv
```

## Performance

- **Memory efficient**: Processes messages in batches
- **Database optimized**: Uses indexes on timestamp and message_id
- **Network efficient**: Only fetches messages from last hour
- **Fast execution**: Typically completes in seconds

## Error Handling

- **Duplicate messages**: Uses `INSERT OR IGNORE` to handle gracefully
- **API errors**: Logs and continues where possible
- **Database errors**: Transaction rollback on failure
- **Network issues**: Basic retry logic built-in

## Environment Configuration

### .env File
Create a `.env` file from the example:
```bash
cp .env.example .env
```

Edit `.env` with your credentials:
```env
# Required
DISCORD_TOKEN=your_bot_token_here
DISCORD_CHANNEL_ID=123456789012345678

# Optional
DATABASE_PATH=custom/path/to/database.db
```

### Environment Variable Priority
1. `.env` file (loaded by godotenv)
2. System environment variables
3. Default values in code

### Git Ignore
The `.env` file is automatically added to `.gitignore` to prevent accidental commits.

## Security

1. **Never commit tokens** to version control - `.env` is gitignored
2. **Use environment variables** loaded from `.env` file
3. **Restrict bot permissions** to read-only
4. **Regular database backups**
5. **Monitor cron job logs**

## Monitoring

Check if the scraper is working:
```bash
# Check last run time
ls -la discord_messages.db

# Check recent messages count
sqlite3 discord_messages.db "SELECT COUNT(*) FROM messages WHERE timestamp > datetime('now', '-2 hours');"

# View logs (if using cron with logging)
tail -f /var/log/discord_scraper.log
```

## Troubleshooting

### Common Issues

1. **"Missing Access" error**
   - Bot needs "Read Messages" permission
   - Bot needs to be in the channel
   - MESSAGE CONTENT INTENT must be enabled

2. **No messages fetched**
   - Check ChannelID is correct
   - Verify bot has proper permissions
   - Test with a recent message in channel

3. **Database permission errors**
   - Ensure write permissions in directory
   - Check disk space
   - Verify SQLite driver is installed

4. **Cron job not running**
   - Check cron service is active
   - Verify path in cron is absolute
   - Check system logs: `grep CRON /var/log/syslog`

### Logging
For production, redirect output:
```bash
0 * * * * /path/to/discord-scraper >> /var/log/discord_scraper.log 2>&1
```

## Extending

### Add More Fields
```go
// In Message struct
type Message struct {
    // ... existing fields
    Edited      bool
    ReactionCount int
    Mentions    []string
}

// In saveMessages function
_, err = txStmt.Exec(
    // ... existing fields
    msg.EditedTimestamp != "",
    len(msg.Reactions),
    strings.Join(extractMentions(msg), ","),
)
```

### Add Webhook Notifications
```go
func sendWebhookNotification(count int) {
    // Send summary to webhook when done
}
```

### Add Metrics
```go
import "github.com/prometheus/client_golang/prometheus"

var messagesScraped = prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Name: "discord_messages_scraped",
        Help: "Number of Discord messages scraped",
    },
    []string{"channel"},
)
```

## License
MIT - Free to use and modify