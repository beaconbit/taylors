# Taylors Project - OpenCode Instructions

## Project Overview
This project uses a text-based task tracking system that can be imported into various open-source project management tools.

## Task Format Specification

### File Location and Naming
- All tasks go directly in the topic directory (e.g., `factory projects/`, `info management/`, `recovery/`)
- **IMPORTANT: Task files MUST be named after the task title, NOT the task ID**
- File naming convention: Convert task title to lowercase, replace spaces with underscores, and add `.json` extension
- Example: Task title "Clean bag extend storage line" becomes `clean_bag_extend_storage_line.json`
- Never use task ID in filename (e.g., do NOT use `task-001.json`)

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
4. **CRITICAL: Save file with title-based filename, NOT task ID**
   - Convert title to lowercase
   - Replace spaces with underscores
   - Add `.json` extension
   - Example: `rail_system_map_grid.json`

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

## Automated Reports

### Generate Synopsis
When instructed to "generate synopsis":
1. Read all task files from topic directories
2. Create `synopsis.md` at project root with:
   - Current date
   - High priority tasks only
3. For each high priority task, include:
   - Task title and ID and path (directory path)
   - description
   - Last 1 update (most recent)
4. Use template from `example_synopsis.md` as guide

### Generate Recurring Tasks Report
When instructed to "generate recurring":
1. Find all tasks with status "repeating"
2. Filter for tasks with `next_occurrence` date within current week
3. Create `recurring.md` at project root with:
   - Week range (Monday to Sunday)
   - Each recurring task with:
     - Title and ID
     - Frequency and next occurrence
     - Directory location
     - Placeholders for outcome, resulting tasks, and next occurrence
4. Use template from `example_recurring.md` as guide

### Process Recurring Updates
When `recurring.md` has been populated with outcomes:
1. Read `recurring.md` file
2. For each completed recurring task:
   - Add update to original task file with outcome
   - Update `next_occurrence` field based on specified time
   - Create new tasks from "resulting_tasks" section
   - Assign new tasks to appropriate directories
