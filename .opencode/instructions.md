# Taylors Project - OpenCode Instructions

## Project Overview
This project uses a text-based task tracking system that can be imported into various open-source project management tools.

## Task Format Specification

### File Location
- All tasks go in the `tasks/` directory
- File naming: `{task-id}.json` (e.g., `task-001.json`)

### Required Fields
Every task must include:
- `id`: Unique identifier (auto-generated)
- `title`: Brief task description
- `status`: Current state (todo, repeating, in_progress, review, done, blocked)
- `priority`: Importance level (low, medium, high, critical)
- `created`: Creation timestamp (ISO8601 format)
- `description`: Detailed task description

### People Tracking Fields
- `assignee`: Primary person responsible (single string)
- `involved_people`: Array of all people involved in the task (array of strings)

### Update Tracking
- `updates`: Array of update objects tracking task progress and ownership changes

### Date Format
Use ISO8601 format: `YYYY-MM-DDTHH:MM:SSZ`
Example: `2025-02-10T14:30:00Z`

### Status Values
- `todo`: Not started
- `repeating`: Recurring task
- `in_progress`: Currently being worked on
- `review`: Ready for review
- `done`: Completed
- `blocked`: Cannot proceed due to dependencies

## Code Guidelines
- Write clear json without additional comments
- Updates should follow ```update_template.json```
- always commit after making new tasks or updates

## Task Creation Process
1. Generate unique ID (next sequential number)
2. Fill all required fields
3. Add `assignee` and `involved_people` if known
7. Save as JSON file in `tasks/` directory

## Update Process
When updating a task:
1. Create a new update object using the template in `.opencode/update_template.json`
2. Fill required fields: `date`, `description`, `creator`
3. If responsibility changes, add `new_owner` field
4. If status changes, add `status_change` with `from` and `to` values
5. If priority changes, add `priority_change` with `from` and `to` values
6. Add the update object to the task's `updates` array

### Example Update
```json
{
  "date": "2025-02-10T14:30:00Z",
  "description": "Completed initial research phase",
  "creator": "alice@example.com",
  "new_owner": "bob@example.com",
  "status_change": "done"
}
```

## Import Compatibility
This format is designed to be easily converted to:
- CSV for spreadsheet import
- JSON for API-based tools
- YAML for configuration-based systems
- Markdown for documentation
