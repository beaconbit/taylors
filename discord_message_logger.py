#!/usr/bin/env python3
"""
Discord Message Logger
Pulls messages from the last hour and stores in SQLite database
Run hourly via cron: 0 * * * * /usr/bin/python3 /path/to/discord_message_logger.py
"""

import discord
import sqlite3
import asyncio
import datetime
import os
from datetime import timedelta, timezone

# Configuration
DISCORD_TOKEN = "YOUR_BOT_TOKEN_HERE"  # Replace with your bot token
CHANNEL_ID = 123456789012345678  # Replace with your channel ID
DATABASE_PATH = "discord_messages.db"

class DiscordMessageLogger:
    def __init__(self):
        self.intents = discord.Intents.default()
        self.intents.message_content = True  # Required to read message content
        self.intents.messages = True
        self.client = discord.Client(intents=self.intents)
        self.setup_database()
        
        # Connect event handlers
        self.client.event(self.on_ready)
        
    def setup_database(self):
        """Create SQLite database and messages table"""
        self.conn = sqlite3.connect(DATABASE_PATH)
        self.cursor = self.conn.cursor()
        
        # Create messages table if it doesn't exist
        self.cursor.execute('''
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
        
        # Create index for faster lookups
        self.cursor.execute('''
            CREATE INDEX IF NOT EXISTS idx_timestamp 
            ON messages(timestamp)
        ''')
        
        self.conn.commit()
        print(f"Database setup complete at {DATABASE_PATH}")
    
    async def on_ready(self):
        """Called when bot connects to Discord"""
        print(f'Logged in as {self.client.user} (ID: {self.client.user.id})')
        print('------')
        
        # Get the channel
        channel = self.client.get_channel(CHANNEL_ID)
        if not channel:
            print(f"Error: Could not find channel with ID {CHANNEL_ID}")
            await self.client.close()
            return
            
        # Calculate time threshold (1 hour ago)
        one_hour_ago = datetime.datetime.now(timezone.utc) - timedelta(hours=1)
        print(f"Fetching messages since: {one_hour_ago}")
        
        try:
            # Fetch messages from the last hour
            messages_fetched = 0
            async for message in channel.history(limit=None, after=one_hour_ago):
                await self.save_message(message)
                messages_fetched += 1
                
            print(f"Successfully fetched and saved {messages_fetched} messages")
            
        except Exception as e:
            print(f"Error fetching messages: {e}")
        finally:
            await self.client.close()
    
    async def save_message(self, message):
        """Save a single message to the database"""
        try:
            # Check if message already exists
            self.cursor.execute(
                "SELECT 1 FROM messages WHERE discord_message_id = ?",
                (str(message.id),)
            )
            
            if self.cursor.fetchone():
                print(f"Message {message.id} already exists, skipping")
                return
            
            # Insert new message
            self.cursor.execute('''
                INSERT INTO messages 
                (discord_message_id, channel_id, author_id, author_name, 
                 content, timestamp, attachments_count)
                VALUES (?, ?, ?, ?, ?, ?, ?)
            ''', (
                str(message.id),
                str(message.channel.id),
                str(message.author.id),
                str(message.author),
                message.content,
                message.created_at.isoformat(),
                len(message.attachments)
            ))
            
            self.conn.commit()
            print(f"Saved message from {message.author}: {message.content[:50]}...")
            
        except Exception as e:
            print(f"Error saving message {message.id}: {e}")
            self.conn.rollback()
    
    def run(self):
        """Start the Discord client"""
        self.client.run(DISCORD_TOKEN)

def main():
    # Check if token is set
    if DISCORD_TOKEN == "YOUR_BOT_TOKEN_HERE":
        print("ERROR: Please replace DISCORD_TOKEN with your actual bot token")
        print("1. Go to https://discord.com/developers/applications")
        print("2. Create a bot and copy the token")
        print("3. Update DISCORD_TOKEN in the script")
        return
    
    # Check if channel ID is set
    if CHANNEL_ID == 123456789012345678:
        print("ERROR: Please replace CHANNEL_ID with your actual channel ID")
        print("1. Enable Developer Mode in Discord Settings")
        print("2. Right-click your channel → Copy ID")
        print("3. Update CHANNEL_ID in the script")
        return
    
    # Run the logger
    logger = DiscordMessageLogger()
    logger.run()

if __name__ == "__main__":
    main()