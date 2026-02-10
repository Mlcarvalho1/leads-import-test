# Lead Import – Specification

This document describes **what** the lead import feature must do. Implementation details (language, framework, validation order) are up to you.

---

## 1. Endpoint

**`POST /import`** – `multipart/form-data`

### Input

| Part   | Type        | Description |
|--------|-------------|-------------|
| `file` | file        | CSV or Excel file with leads to import |
| `data` | JSON string | Import metadata (see below) |

**`data` fields:**

| Field        | Type            | Required | Rules |
|-------------|-----------------|----------|-------|
| `name`       | string          | yes      | max 255 chars, must be unique per account/company (among non-deleted imports) |
| `account_id` | integer         | yes      | must reference a valid messaging account the user has access to |
| `source_id`  | integer         | yes      | must reference a valid, non-deleted lead source |
| `tag_ids`    | array of ints   | no       | max 5 elements, each must be a valid tag |

### Authorization

The user must have the `IMPORT_LEADS` permission on the `LEADS` module.

---

## 2. File Rules

The file (CSV or Excel) must have **exactly 5 columns** with these headers in order:

| Index | Header  | Required | Rules |
|-------|---------|----------|-------|
| 0     | `name`  | yes      | string, max 255 |
| 1     | `phone` | yes      | must be a valid phone number (parseable into dial code, country code, DDD, MSI) |
| 2     | `cpf`   | no       | must be a valid CPF if present |
| 3     | `email` | no       | valid email, max 255 |
| 4     | `tags`  | no       | comma-separated tag names (e.g. `"tag1, tag2"`), max 5 tags per row, max 255 chars total |

**Limits:** (this was for node, in golang we can improove it later)
- At least 1 data row (besides the header).
- Max **5,000** data rows per file.

---

## 3. Rate Limiting

Per company + account, within a **1-hour window**:

- Max **5,000 total leads** across all imports.
- Max **5 imports**.

---

## 4. What to Do (high-level flow)

### 4.1 Synchronous (on request)

1. Validate the file structure, row data, source, account, import name uniqueness, and rate limit.
2. Create one `lead_imports` record with status `PROCESSING`.
3. Respond with the import ID. **Do not block** – actual lead creation happens asynchronously.

### 4.2 Asynchronous (worker/job)

1. **Filter duplicates** – For each lead's phone (normalized), check if a lead, patient, or chat already exists for the same account/company/country. If so, skip it and count as `total_existing`.

2. **Resolve tags** – Collect all distinct tag names from file rows + the `tag_ids` from the request. Find or create tags by name within the company.

3. **Create leads** – For each non-duplicate lead:
   - Optionally validate the phone with the WhatsApp provider.
   - Create a chat in MongoDB.
   - Insert the lead in PostgreSQL (with `channel = IMPORT`, linked `source_id`, `import_id`, `chat_id`).
   - Update the chat with the `lead_id`.
   - Create `chat_tags` entries linking the chat to its tags.
   - Process in chunks (e.g. 100) with limited concurrency.

4. **Finish** – Update the `lead_imports` record with:
   - `status`: `FINISHED` or `FAILED`
   - `total_created`, `total_existing`, `total_errors`
   - Clean up unused tags, clear caches, emit `lead:import-finished` event.

---

## 5. Database Models

All PostgreSQL tables use schema `amigocare` unless noted.

### 5.1 `lead_imports`

| Column         | Type    | Notes |
|----------------|---------|-------|
| id             | int PK  | auto increment |
| name           | string  | required |
| status         | string  | `PROCESSING`, `FINISHED`, or `FAILED` |
| total_created  | int     | default 0 |
| total_existing | int     | default 0 |
| total_errors   | int     | default 0 |
| is_deleted     | bool    | default false |
| creator_id     | int FK  | -> users |
| company_id     | int FK  | -> companies |
| source_id      | int FK  | -> lead_sources |
| account_id     | int FK  | -> messaging_accounts |
| created_at     | timestamp | |
| updated_at     | timestamp | |

### 5.2 `amigocare_leads`

| Column                              | Type    | Notes |
|--------------------------------------|---------|-------|
| id                                   | int PK  | auto increment |
| name                                 | string  | nullable |
| email                                | string  | nullable |
| cpf                                  | string  | max 11, nullable |
| contact_cellphone                    | string  | max 25, required |
| contact_cellphone_dial_code          | string  | max 25, default "55" |
| contact_cellphone_country_code       | string  | max 25, default "BR" |
| source_id                            | int FK  | -> lead_sources |
| channel_id                           | int FK  | -> lead_channels (use the IMPORT channel) |
| chat_id                              | string  | nullable, MongoDB ObjectID |
| import_id                            | int FK  | -> lead_imports |
| company_id                           | int FK  | -> companies |
| amigocare_messaging_account_id       | int FK  | -> messaging_accounts |
| creator_id                           | int FK  | -> users |
| patient_id                           | int FK  | nullable |
| attendance_id                        | int FK  | nullable |
| converted_by                         | int FK  | nullable |
| converted_at                         | timestamp | nullable |
| converted_from                       | string  | nullable |
| is_deleted                           | bool    | default false |
| created_at                           | timestamp | |
| updated_at                           | timestamp | |

### 5.3 `lead_sources`

| Column     | Type    | Notes |
|-----------|---------|-------|
| id        | int PK  | auto increment |
| name      | string  | max 255 |
| is_deleted| bool    | default false |

### 5.4 `lead_channels`

| Column     | Type    | Notes |
|-----------|---------|-------|
| id        | int PK  | auto increment |
| name      | string  | max 100 |
| is_deleted| bool    | default false |

### 5.5 `tags`

| Column       | Type    | Notes |
|-------------|---------|-------|
| id          | int PK  | auto increment |
| name        | string  | max 255 |
| is_deleted  | bool    | default false |
| company_id  | int FK  | -> companies |
| creator_id  | int FK  | -> users |
| destroyer_id| int FK  | nullable |
| created_at  | timestamp | |
| updated_at  | timestamp | |

### 5.6 `chat_tags`

| Column       | Type    | Notes |
|-------------|---------|-------|
| id          | int PK  | auto increment |
| chat_id     | string  | required (MongoDB ObjectID) |
| tag_id      | int FK  | -> tags |
| lead_id     | int FK  | -> amigocare_leads |
| company_id  | int FK  | -> companies |
| creator_id  | int FK  | -> users |
| destroyer_id| int FK  | nullable |
| is_deleted  | bool    | default false |
| created_at  | timestamp | |
| updated_at  | timestamp | |

### 5.7 `patients` (read-only, for duplicate check)

Only needed to query by phone + company to detect existing contacts:

| Column                          | Type    |
|----------------------------------|---------|
| id                               | int PK  |
| company_id                       | int FK  |
| contact_cellphone                | string  |
| contact_cellphone_dial_code      | string  |
| contact_cellphone_country_code   | string  |
| deleted_at                       | timestamp (nullable) |

### 5.8 Referenced tables (must exist, not managed by this feature)

- **companies**
- **users**
- **messaging_accounts**
- **attendances** (optional)

### 5.9 MongoDB

- **chats** – Created per lead during import. The `chat_id` stored in `amigocare_leads` and `chat_tags` references the MongoDB document `_id`. Also used for duplicate checking by phone + account.

---

## 6. Summary

| Step | What happens | Where |
|------|-------------|-------|
| Receive request | Validate file, metadata, permissions, rate limit | API (sync) |
| Create import record | Insert `lead_imports` with `PROCESSING` | PostgreSQL |
| Respond to client | Return import ID | API (sync) |
| Filter duplicates | Check leads, patients, chats by phone | PostgreSQL + MongoDB (async) |
| Resolve tags | Find or create tags by name per company | PostgreSQL (async) |
| Create leads | Insert chats, leads, chat_tags | MongoDB + PostgreSQL (async) |
| Finish | Update import status and counts, notify | PostgreSQL (async) |
