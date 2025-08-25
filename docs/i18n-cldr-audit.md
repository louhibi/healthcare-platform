# CLDR i18n Audit and Remediation Plan

Date: 2025-08-08
Scope: frontend/patient-portal (Vue 3 + vue-i18n v9)

## Summary
- vue-i18n v9 is configured (Intl-based → CLDR-backed) with locale messages, numberFormats, and datetimeFormats.
- Several UI areas bypass vue-i18n and CLDR by using hardcoded strings and `toLocale*` APIs.
- No pluralization rules present. Arabic (ar-MA) requires CLDR plural categories.
- Backend translation merge logic is incorrect and may not update runtime messages.

## What’s good
- createI18n configured with legacy: false, per-locale `numberFormats` and `datetimeFormats`.
- Locale, lang, and dir switching handled (`setLocale`, RTL support).

## Gaps (CLDR correctness)
1) Pluralization
- No ICU/CLDR plural messages in locale files.

2) Numbers and dates
- Components use native `toLocaleDateString`/`toLocaleTimeString`/`toLocaleString` instead of vue-i18n `d()`/`n()`.

3) Hardcoded strings
- Many labels/messages are inline English and not using `t()`.

4) Weekday/month names
- Hardcoded arrays for weekdays (e.g., `['Sun', 'Mon', ...]`). Should use Intl/CLDR or translations.

5) Validation messages
- Hardcoded English in `useFormConfig.validateField` instead of `t('validation.*')` keys.

6) Backend merge bug
- `useI18n.js` writes into `availableLocales` from `useVueI18n()` (an array of codes), not the message store. Should use `i18n.global.mergeLocaleMessage`/`setLocaleMessage`.

## Remediation plan (actionable)

### A. Fix backend translation merge
- File: `src/composables/useI18n.js`
- Replace the custom merge with vue-i18n APIs:

```js
import i18n from '@/i18n'

function mergeBackendTranslations(locale, translations) {
  const nested = flatToNested(translations) // keep your flat→nested logic
  const existing = i18n.global.getLocaleMessage(locale) || {}
  const merged = deepMerge(existing, nested)
  i18n.global.setLocaleMessage(locale, merged)
}
```

### B. Use CLDR-backed formatters everywhere
- Prefer `d(date, formatKey)` and `n(number, formatKey)` from `useI18n()`.
- Replace:
  - `new Date(x).toLocaleString()` → `d(new Date(x), 'long')`
  - `toLocaleDateString('en-US', ...)` → `d(date, 'short'|'long')`
  - Currency/amounts → `n(amount, 'currency')`

Examples:
- Files to update:
  - `src/components/AppointmentStatusModal.vue` (formatDateTime)
  - `src/components/AppointmentCalendar.vue` (formatTime)
  - `src/components/AvailabilityCalendar.vue` (currentMonthYear)
  - `src/views/Patients.vue` (DOB rendering)

### C. Add plural-safe messages
- Locale keys to add (example):

```json
// en-US.json
{
  "common": {
    "items": "{count, plural, one {# item} other {# items}}"
  },
  "availability": {
    "workingDays": "{count, plural, one {# working day} other {# working days}}"
  }
}
```

```json
// fr-FR.json
{
  "common": {
    "items": "{count, plural, one {# élément} other {# éléments}}"
  },
  "availability": {
    "workingDays": "{count, plural, one {# jour ouvré} other {# jours ouvrés}}"
  }
}
```

```json
// ar-MA.json (CLDR categories: zero, one, two, few, many, other)
{
  "common": {
    "items": "{count, plural, zero {# عنصر} one {# عنصر} two {# عنصران} few {# عناصر} many {# عنصرًا} other {# عنصر}}"
  },
  "availability": {
    "workingDays": "{count, plural, zero {# يوم عمل} one {# يوم عمل} two {# يومان عمل} few {# أيام عمل} many {# يوم عمل} other {# يوم عمل}}"
  }
}
```

Usage:
```vue
{{ t('availability.workingDays', { count: workingDaysCount }) }}
```

### D. Localize hardcoded UI text
- Move strings to locale files and use `t()`.
- Key spots:
  - `DynamicPatientForm.vue`: header title/description; select placeholder (use `t('common.select')`); buttons (reset/submit); counters (“fields enabled”, “required”).
  - `AppointmentStatusModal.vue`: status labels/messages/placeholders.
  - `AvailabilityCalendar.vue`: legend labels, “This Month:”, weekday headers, “Total working days”.

### E. Weekday/month names via CLDR
- Generate weekday headers with Intl using current locale:

```js
const weekdays = [...Array(7).keys()].map(i =>
  new Intl.DateTimeFormat(currentLocale, { weekday: 'short' })
    .format(new Date(1970, 0, 4 + i)) // 1970-01-04 is a Sunday
)
```

Or add translated arrays in locale files.

### F. Validation messages via i18n
- File: `src/composables/useFormConfig.js`
- Replace hardcoded messages with `t('validation.*')` and interpolate field label when needed, e.g.:

```js
if (field.is_required && isEmpty(value)) {
  errors.push(t('validation.required'))
}
```

### G. DynamicPatientForm specific improvements
- Replace placeholders like `Enter ${display_name}` with either:
  - `field.placeholder_text` (if provided from backend/localization), or
  - `t('forms.patient.fields.<name>.placeholder')` fallback, or
  - a generic `t('common.enterField', { field: field.display_name })` you add to locales.
- For the empty option on selects, use `t('common.select')` and the field label: e.g., `t('common.selectField', { field: field.display_name })`.

### H. Optional CLDR niceties
- Add utilities for `Intl.RelativeTimeFormat` ("in 3 days", "2 hours ago") and `Intl.ListFormat` (A, B, and C) if needed.

## Concrete TODO checklist
- [ ] Fix backend translation merge (`useI18n.js` → `setLocaleMessage`/`mergeLocaleMessage`).
- [ ] Replace `toLocale*` usages with `d()`/`n()` in the listed components.
- [ ] Introduce pluralization keys in `en-US.json`, `fr-FR.json`, `ar-MA.json` and use them.
- [ ] Localize hardcoded strings in `DynamicPatientForm.vue`, `AppointmentStatusModal.vue`, `AvailabilityCalendar.vue`, `AppointmentCalendar.vue`, `Patients.vue`.
- [ ] Switch validation messages to i18n.
- [ ] Weekday headers via Intl or locale messages.
- [ ] Add select/placeholder common keys (`common.selectField`, `common.enterField`).

## Minimal code examples

Expose d/n/t in components:
```js
import { useI18n } from 'vue-i18n'
const { t, d, n, locale } = useI18n()
```

Use in templates:
```vue
{{ d(new Date(appointment.updated_at), 'long') }}
{{ n(total, 'currency') }}
{{ t('availability.workingDays', { count: workingDaysCount }) }}
```

DynamicPatientForm select placeholder:
```vue
<option value="">
  {{ t('common.selectField', { field: field.display_name }) }}
</option>
```

Add common keys:
```json
// en-US.json
{
  "common": {
    "selectField": "Select {field}",
    "enterField": "Enter {field}"
  }
}
```

## Notes
- vue-i18n v9 supports ICU MessageFormat; ensure the message compiler is enabled (default when bundling with Vite).
- For Arabic plural correctness, prefer ICU messages with all categories.
