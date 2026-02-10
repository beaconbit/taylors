# Task Tracking System

## Overview
This directory contains a text-based task tracking system configured for use with OpenCode and compatible with various open-source project management tools.

## Structure
```
Taylors/
├── .opencode/              # OpenCode configuration
│   ├── config.json        # Task format specification
│   ├── instructions.md    # Project guidelines
│   └── task_template.json # Template for new tasks
├── tasks/                 # Task files
│   └── task-001.json     # Sample task
└── README-tasks.md       # This file
```

## Quick Start

### Creating a New Task
1. Use the template in `.opencode/task_template.json`
2. Fill in required fields:
   - `id`: Next sequential number (task-002, task-003, etc.)
   - `title`: Brief description
   - `description`: Detailed information
   - `status`: todo, in_progress, review, done, or blocked
   - `priority`: low, medium, high, or critical
   - `created`: Current timestamp in ISO8601 format
3. Save as `tasks/task-{id}.json`

### Example Task Creation
```json
{
  "id": "task-002",
  "title": "Your task title",
  "description": "Detailed description...",
  "status": "todo",
  "priority": "medium",
  "created": "2025-02-10T14:30:00Z",
  "updated": "2025-02-10T14:30:00Z",
  "assignee": "person@example.com",
  "involved_people": ["person@example.com", "other@example.com"],
  "updates": [],
  "tags": ["category1", "category2"],
  "dependencies": ["task-001"]
}
```

### Example Update
```json
{
  "date": "2025-02-10T15:30:00Z",
  "description": "Completed initial phase, ready for review",
  "creator": "person@example.com",
  "new_owner": "reviewer@example.com",
  "status_change": {
    "from": "in_progress",
    "to": "review"
  }
}
```

## OpenCode Integration
OpenCode will automatically:
- Read configuration from `.opencode/config.json`
- Follow instructions in `.opencode/instructions.md`
- Use the task template for consistency
- Understand the task format and required fields

## Tool Compatibility
This format can be converted to:
- **CSV** for spreadsheet import (Excel, Google Sheets)
- **JSON** for API-based tools (Planka, Taiga, OpenProject)
- **YAML** for configuration-based systems
- **Markdown** for documentation

## Sample Task
See `tasks/task-001.json` for a complete example including:
- People tracking (`assignee` and `involved_people`)
- Update history with responsibility transfers
- Status change tracking
- Complete workflow example

## Next Steps
1. Create additional tasks using the template
2. Use OpenCode to manage and update tasks
3. Convert to other formats as needed for specific tools
