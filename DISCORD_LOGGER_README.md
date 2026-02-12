# Discord Message Logger

A Python script that pulls messages from a Discord channel from the last hour and stores them in a SQLite database. Designed to run hourly via cron.

## Files Created

1. `discord_message_logger.py` - Main script
2. `setup_discord_logger.sh` - Setup script
3. `query_messages.py` - Query/statistics script
4. `DISCORD_LOGGER_README.md` - This file

## Setup Instructions

### 1. Create Discord Bot
1. Go to https://discord.com/developers/applications
2. Click "New Application" → Name it (e.g., "Message Logger")
3. Go to "Bot" section → Click "Add Bot"
4. Copy the **Bot Token** (keep this secret!)
5. Under "Privileged Gateway Intents", enable:
   - **MESSAGE CONTENT INTENT** (required to read messages)

### 2. Get Channel ID
1. In Discord, go to User Settings → Advanced
2. Enable "Developer Mode"
3. Right-click your channel → "Copy ID"

### 3. Install & Configure
```bash
# Run setup script
chmod +x setup_discord_logger.sh
./setup_discord_logger.sh

# Edit main script with your credentials
nano discord_message_logger.py
```
Update these lines:
```python
DISCORD_TOKEN = "YOUR_BOT_TOKEN_HERE"  # Replace with your bot token
CHANNEL_ID = 123456789012345678  # Replace with your channel ID
```

### 4. Invite Bot to Server
1. In Discord Developer Portal → OAuth2 → URL Generator
2. Select scopes: `bot`
3. Select bot permissions: `Read Messages/View Channels`
4. Use generated URL to invite bot to your server
5. Make sure bot has access to the target channel

### 5. Test the Script
```bash
python3 discord_message_logger.py
```

### 6. Set Up Hourly Cron Job
```bash
crontab -e
```
Add this line (adjust path):
```
0 * * * * /usr/bin/python3 /home/ubuntu/Taylors/discord_message_logger.py
```

## How It Works

### Message Filtering (Last Hour Only)
The script uses Discord's `channel.history()` method with the `after` parameter:
```python
one_hour_ago = datetime.datetime.now(timezone.utc) - timedelta(hours=1)
async for message in channel.history(limit=None, after=one_hour_ago):
```

### Database Schema
```sql
CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    discord_message_id TEXT UNIQUE,  -- Discord's message ID
    channel_id TEXT,                 -- Channel ID
    author_id TEXT,                  -- User ID
    author_name TEXT,                -- Username
    content TEXT,                    -- Message text
    timestamp DATETIME,              -- When message was sent
    attachments_count INTEGER,       -- Number of attachments
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP  -- When logged
);
```

### Duplicate Prevention
- Uses `discord_message_id` as UNIQUE constraint
- Checks if message exists before inserting
- Prevents duplicate logging if script runs multiple times

## Querying Data

### Basic SQL Queries
```bash
sqlite3 discord_messages.db

-- Recent messages
SELECT * FROM messages ORDER BY timestamp DESC LIMIT 10;

-- Messages from specific user
SELECT * FROM messages WHERE author_name LIKE '%username%';

-- Daily summary
SELECT date(timestamp), COUNT(*) 
FROM messages 
GROUP BY date(timestamp);
```

### Using Query Script
```bash
python3 query_messages.py
```

## Troubleshooting

### Common Issues

1. **"Missing Access" error**
   - Bot needs to be invited to server with proper permissions
   - Bot needs access to the specific channel

2. **No messages fetched**
   - Check CHANNEL_ID is correct
   - Ensure bot has "Read Messages" permission
   - Verify MESSAGE CONTENT INTENT is enabled

3. **"Invalid token" error**
   - Regenerate bot token in Discord Developer Portal
   - Ensure token is copied correctly (no spaces)

4. **Script runs but no database updates**
   - Check database file permissions
   - Verify cron job is running (check syslog)
   - Test script manually: `python3 discord_message_logger.py`

### Logging
The script prints to console. For cron jobs, redirect output:
```
0 * * * * /usr/bin/python3 /path/to/discord_message_logger.py >> /var/log/discord_logger.log 2>&1
```

## Security Notes

1. **Keep bot token secret** - Never commit to version control
2. **Use environment variables** for production:
   ```python
   import os
   DISCORD_TOKEN = os.getenv('DISCORD_TOKEN')
   ```
3. **Restrict bot permissions** - Only grant necessary permissions
4. **Regular backups** - Backup SQLite database regularly

## Extending the Script

### Add More Fields
```python
# In save_message() function, add more fields:
'INSERT INTO messages (..., reaction_count, edited) VALUES (..., ?, ?)',
(len(message.reactions), message.edited_at is not None)
```

### Multiple Channels
```python
CHANNEL_IDS = [123456789, 987654321]  # List of channel IDs
for channel_id in CHANNEL_IDS:
    channel = self.client.get_channel(channel_id)
    # Fetch messages...
```

### Error Handling & Retries
Add retry logic for network issues:
```python
import tenacity

@tenacity.retry(stop=tenacity.stop_after_attempt(3), wait=tenacity.wait_exponential())
async def fetch_messages_with_retry(channel, after_time):
    return [msg async for msg in channel.history(limit=None, after=after_time)]
```