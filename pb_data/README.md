# pb_data/

PocketBase runtime data directory. Created and managed automatically by PocketBase.

## Contents

| Path | Description |
|------|-------------|
| `data.db` | SQLite database (all collections, records, auth data) |
| `storage/` | Uploaded files organized by collection and record ID |
| `backups/` | Database backup archives |
| `logs/` | Request and error logs |

## Git

This entire directory is gitignored. **Do not commit `pb_data/` contents.**

Only this `README.md` is tracked so the directory intent is documented in the repo.

For production, back up `data.db` and `storage/` independently (e.g., via PocketBase's built-in backup API or external tooling).
