#!/usr/bin/env python3
"""
Query script for Discord messages database
"""

import sqlite3
import datetime
from tabulate import tabulate

def query_recent_messages(limit=10):
    """Query recent messages from database"""
    conn = sqlite3.connect('discord_messages.db')
    cursor = conn.cursor()
    
    cursor.execute('''
        SELECT 
            timestamp,
            author_name,
            substr(content, 1, 100) as preview,
            attachments_count
        FROM messages 
        ORDER BY timestamp DESC 
        LIMIT ?
    ''', (limit,))
    
    rows = cursor.fetchall()
    conn.close()
    
    if rows:
        print(f"\nLast {len(rows)} messages:\n")
        print(tabulate(rows, headers=['Timestamp', 'Author', 'Message Preview', 'Attachments'], tablefmt='grid'))
    else:
        print("No messages found in database")

def get_statistics():
    """Get database statistics"""
    conn = sqlite3.connect('discord_messages.db')
    cursor = conn.cursor()
    
    # Total messages
    cursor.execute('SELECT COUNT(*) FROM messages')
    total = cursor.fetchone()[0]
    
    # Messages today
    today = datetime.datetime.now().strftime('%Y-%m-%d')
    cursor.execute('SELECT COUNT(*) FROM messages WHERE date(timestamp) = ?', (today,))
    today_count = cursor.fetchone()[0]
    
    # Top authors
    cursor.execute('''
        SELECT author_name, COUNT(*) as message_count
        FROM messages
        GROUP BY author_name
        ORDER BY message_count DESC
        LIMIT 5
    ''')
    top_authors = cursor.fetchall()
    
    conn.close()
    
    print(f"\n=== Database Statistics ===")
    print(f"Total messages: {total}")
    print(f"Messages today: {today_count}")
    print(f"\nTop 5 authors:")
    for author, count in top_authors:
        print(f"  {author}: {count} messages")

if __name__ == "__main__":
    get_statistics()
    query_recent_messages(10)