# Internationalization

Guide to the multi-language support and translation management system.

[[Home|← Back to Home]]

---

## Overview

Safecast New Map supports **29 languages** with a comprehensive internationalization (i18n) system:

- PostgreSQL-backed translations
- Admin UI for live translation editing
- Incremental seeding from embedded translations
- Language selection via URL parameter or browser preference
- Filtered translations payload for performance

---

## Supported Languages

| Code | Language | Native Name |
|------|----------|-------------|
| ar | Arabic | العربية |
| bg | Bulgarian | Български |
| cs | Czech | Čeština |
| da | Danish | Dansk |
| de | German | Deutsch |
| el | Greek | Ελληνικά |
| en | English | English |
| es | Spanish | Español |
| fa | Persian | فارسی |
| fi | Finnish | Suomi |
| fr | French | Français |
| he | Hebrew | עברית |
| hi | Hindi | हिन्दी |
| hu | Hungarian | Magyar |
| id | Indonesian | Bahasa Indonesia |
| it | Italian | Italiano |
| ja | Japanese | 日本語 |
| ko | Korean | 한국어 |
| ms | Malay | Bahasa Melayu |
| nl | Dutch | Nederlands |
| no | Norwegian | Norsk |
| pl | Polish | Polski |
| pt | Portuguese | Português |
| ru | Russian | Русский |
| sv | Swedish | Svenska |
| th | Thai | ไทย |
| tr | Turkish | Türkçe |
| uk | Ukrainian | Українська |
| vi | Vietnamese | Tiếng Việt |
| zh | Chinese | 中文 |

---

## Language Selection

### Priority Order

1. **URL parameter** (`?lang=ja`)
2. **Browser Accept-Language header**
3. **English fallback**

### URL Parameter

```
/?lang=en    - English
/?lang=ja    - Japanese
/?lang=de    - German
/?lang=fr    - French
/?lang=ru    - Russian
```

### Browser Language Detection

The server reads the `Accept-Language` header:
```
Accept-Language: ja,en-US;q=0.9,en;q=0.8
```

Priority: Japanese → US English → English

---

## Translation Architecture

### Database Storage

Translations stored in PostgreSQL table:

```sql
CREATE TABLE translations (
  id BIGSERIAL PRIMARY KEY,
  language_code TEXT NOT NULL,
  key TEXT NOT NULL,
  value TEXT NOT NULL,
  UNIQUE(language_code, key)
);

CREATE INDEX idx_translations_lang ON translations(language_code);
CREATE INDEX idx_translations_key ON translations(key);
```

### Embedded Translations

Default translations embedded in binary:
- `translations.json` - Master translation file
- Seeded on first startup
- Incremental updates on restart

### Memory Cache

Translations loaded into memory at startup:
- Fast lookup
- No database queries during rendering
- Reload via admin panel

### Performance Optimization

**Filtered payload:**
- Only active language + English fallback sent to browser
- ~30KB vs ~850KB for all 30 languages
- Keeps AI widget within Claude's token limits

---

## Translation Management

### Admin Panel

Access: `/admin/translations`

**Features:**
- Filter by language (29 languages)
- Search across keys and values
- Inline editing with save
- Add new translation keys
- Delete translation keys
- **Reload into Memory** button

### Translation Workflow

1. **Select language** from dropdown
2. **Search** for specific key or value
3. **Edit** translation inline
4. **Save** changes to database
5. **Reload into Memory** to apply changes

### Add New Translation Key

**Via Admin Panel:**
1. Click "Add New Key"
2. Enter key name (e.g., `map.legend.title`)
3. Enter translations for each language
4. Save

**Via API:**
```bash
POST /api/admin/translations?password=admin-password
Content-Type: application/json

{
  "key": "new.translation.key",
  "translations": {
    "en": "English text",
    "ja": "日本語テキスト",
    "fr": "Texte français",
    "de": "Deutscher Text"
  }
}
```

### Delete Translation Key

**Via Admin Panel:**
1. Find key in list
2. Click delete button
3. Confirm deletion

**Via API:**
```bash
DELETE /api/admin/translations/{key}?password=admin-password
```

### Reload Translations

After editing translations:
1. Click **Reload into Memory** button
2. Changes applied immediately
3. No server restart required

---

## Translated Components

All UI components are fully translated:

### Map Interface
- Map legend
- Layer controls
- Zoom controls
- Search bar
- Coordinate input dialog

### AI Assistant
- Chat interface
- Tool descriptions
- Error messages
- Help text

### Authentication
- Login modal
- Register modal
- Forgot password
- Email verification

### Spectrum Viewer
- Chart labels
- Isotope names
- Download options
- Energy scale

### Profile Page
- User information
- Upload history
- API key section
- Account actions

### Admin Panel
- User management
- Upload management
- MCP analytics
- Real-time sensors
- Translation management

---

## Translation File Format

### translations.json

```json
{
  "map.legend.title": {
    "en": "Radiation Level",
    "ja": "放射線量",
    "de": "Strahlungswert",
    "fr": "Niveau de radiation"
  },
  "map.legend.safe": {
    "en": "Safe",
    "ja": "安全",
    "de": "Sicher",
    "fr": "Sûr"
  },
  "map.legend.warning": {
    "en": "Warning",
    "ja": "警告",
    "de": "Warnung",
    "fr": "Avertissement"
  }
}
```

### Key Naming Convention

Use dot notation for namespacing:
- `map.legend.title` - Map legend title
- `auth.login.button` - Login button text
- `spectrum.download.csv` - Spectrum download CSV label
- `admin.users.create` - Admin create user action

---

## Incremental Seeding

### How It Works

On startup:
1. Load embedded `translations.json`
2. For each key in JSON:
   - Check if key exists in database
   - If not, insert with `ON CONFLICT DO NOTHING`
   - If exists, preserve existing translation

### Preserve Edits

Database translations take priority:
- Admin edits are preserved
- Embedded JSON only adds new keys
- No overwrites of existing translations

### Add New Keys

**Update `translations.json`:**
```json
{
  "new.key": {
    "en": "New English text",
    "ja": "新しい日本語テキスト"
  }
}
```

**Restart server:**
```bash
systemctl restart safecast
```

**Or reload via admin panel**

---

## Migration

### Add Translations Table

**PostgreSQL:**
```bash
psql -d safecast -f migrations/add_translations_table.sql
```

### Seed UI Translations

```bash
psql -d safecast -f migrations/add_ui_translations.sql
```

Adds extended UI translations:
- AI widget
- Auth modals
- Search bar
- Spectrum viewer
- Profile page

### Seed All Translations

```bash
psql -d safecast -f migrations/seed_translations.sql
```

---

## Translation Best Practices

### Brand Names

**Rule:** "Safecast" must remain untranslated in all languages.

```json
{
  "app.name": {
    "en": "Safecast",
    "ja": "Safecast",
    "de": "Safecast"
  }
}
```

### Technical Terms

Keep technical terms consistent:
- Radiation units (µSv/h, µR)
- Device names (bGeigie, RadiaCode)
- API terminology

### Context

Provide context for translators:
```json
{
  "map.coloring.safecast": {
    "en": "Scientific",
    "ja": "科学的",
    "_comment": "Coloring scheme name - scientific gradient"
  }
}
```

### Length

Consider UI space constraints:
- German translations are often longer
- Japanese/Chinese are often shorter
- Test with longest translations

### Plurals

Handle plurals appropriately:
```json
{
  "track.count": {
    "en": "{count} track(s)",
    "ja": "{count}件のトラック"
  }
}
```

---

## Contributing Translations

### Add New Language

1. **Request language addition** via GitHub issue
2. **Provide translations** for all keys
3. **Test rendering** with new language
4. **Submit pull request**

### Update Existing Translation

1. **Edit via admin panel** (if you have access)
2. **Or submit pull request** with `translations.json` changes
3. **Include native speaker review**

### Translation Quality

**Review process:**
1. Native speaker review
2. Context verification
3. UI testing
4. Merge after approval

---

## Troubleshooting

### Translation Not Appearing

**Check database:**
```sql
SELECT * FROM translations
WHERE language_code = 'ja'
  AND key = 'map.legend.title';
```

**Reload translations:**
1. Go to `/admin/translations`
2. Click "Reload into Memory"

**Check logs:**
```bash
grep "translation" /var/log/safecast.log
```

### Missing Translations

**Seed missing keys:**
```bash
psql -d safecast -f migrations/seed_translations.sql
```

**Or add manually:**
```sql
INSERT INTO translations (language_code, key, value)
VALUES ('ja', 'new.key', '新しい値')
ON CONFLICT (language_code, key) DO NOTHING;
```

### Language Not Available

**Check supported languages:**
```sql
SELECT DISTINCT language_code FROM translations ORDER BY language_code;
```

**Add language:**
1. Add to `translations.json`
2. Run migration
3. Restart server

### Performance Issues

**Check payload size:**
```bash
# Check response size
curl -H "Accept-Language: ja" http://localhost:8765/ | wc -c
```

**Should be ~30KB, not ~850KB**

**Verify filtering:**
- Only active language + English should be sent
- Check server logs for errors

---

## API Reference

### Get Translations

```bash
GET /api/translations?lang=ja
```

Returns filtered translations for specified language.

### Admin: List Translations

```bash
GET /api/admin/translations?password=admin-password&lang=ja
```

### Admin: Update Translation

```bash
PUT /api/admin/translations?password=admin-password
Content-Type: application/json

{
  "language_code": "ja",
  "key": "map.legend.title",
  "value": "新しい放射線量"
}
```

### Admin: Add Translation Key

```bash
POST /api/admin/translations?password=admin-password
Content-Type: application/json

{
  "key": "new.key",
  "translations": {
    "en": "English",
    "ja": "日本語"
  }
}
```

### Admin: Delete Translation Key

```bash
DELETE /api/admin/translations/{key}?password=admin-password
```

### Admin: Reload Translations

```bash
POST /api/admin/translations/reload?password=admin-password
```

---

## See Also

- [Map Features](Map-Features) - Map UI translations
- [Admin Panel](Admin-Panel) - Translation management interface
- [API Documentation](API-Documentation) - Translation API endpoints
- [Development](Development) - Adding new translations
- [Configuration Reference](Configuration-Reference) - Language settings
