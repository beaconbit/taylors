#!/bin/bash
# Setup script for Discord Message Logger

echo "=== Discord Message Logger Setup ==="

# Install required packages
echo "1. Installing Python packages..."
pip3 install discord.py

# Make script executable
echo "2. Making scripts executable..."
chmod +x discord_message_logger.py
chmod +x setup_discord_logger.sh

# Create database directory
echo "3. Creating database..."
python3 -c "
import sqlite3
conn = sqlite3.connect('discord_messages.db')
cursor = conn.cursor()
cursor.execute('''
    CREATE TABLE IF NOT EXISTS messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        discord_message_id TEXT UNIQUE,
        channel_id TEXT,
        author_id TEXT,
        author_name TEXT,
        content TEXT,
        timestamp DATETIME,
        attachments_count INTEGER,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )
''')
cursor.execute('CREATE INDEX IF NOT EXISTS idx_timestamp ON messages(timestamp)')
conn.commit()
conn.close()
print('Database created: discord_messages.db')
"

echo ""
echo "=== NEXT STEPS ==="
echo "1. Update discord_message_logger.py with your:"
echo "   - DISCORD_TOKEN (from Discord Developer Portal)"
echo "   - CHANNEL_ID (right-click channel → Copy ID with Developer Mode enabled)"
echo ""
echo "2. Test the script:"
echo "   python3 discord_message_logger.py"
echo ""
echo "3. Set up hourly cron job:"
echo "   crontab -e"
echo "   Add: 0 * * * * /usr/bin/python3 /full/path/to/discord_message_logger.py"
echo ""
echo "4. Query the database:"
echo "   sqlite3 discord_messages.db"
echo "   .tables"
echo "   SELECT * FROM messages ORDER BY timestamp DESC LIMIT 5;"