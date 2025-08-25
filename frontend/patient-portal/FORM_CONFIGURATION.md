# Form Configuration System

This document describes the dynamic form configuration system implemented in the Vue.js frontend for the Healthcare Platform.

## Overview

The form configuration system allows administrators to dynamically configure which fields are visible and required in various forms (patient registration, appointment booking, etc.) without code changes. The system consists of:

1. **API Client** - Handles communication with backend form configuration endpoints
2. **Pinia Store** - Manages form configuration state and caching
3. **Composable** - Provides reusable form configuration logic
4. **Admin Interface** - UI for configuring form fields
5. **Dynamic Form Components** - Forms that adapt based on configuration

## Architecture

### Backend API Endpoints (Already Implemented)

- `GET /api/forms/types` - Get all form types
- `GET /api/forms/{formType}/metadata` - Get complete form configuration
- `GET /api/forms/{formType}/fields` - Get field configurations
- `PUT /api/forms/{formType}/fields/{fieldId}` - Update single field configuration
- `PUT /api/forms/{formType}/fields` - Update multiple field configurations
- `POST /api/forms/{formType}/reset` - Reset form to defaults

### Frontend Components

#### 1. API Client (`/src/api/forms.js`)

Handles all communication with backend form configuration endpoints.

```javascript
import { formsApi } from '@/api/forms'

// Get form fields
const fields = await formsApi.getFormFields('patient')

// Update field configuration
await formsApi.updateFormField('patient', fieldId, { is_enabled: false })
```

#### 2. Pinia Store (`/src/stores/formConfig.js`)

Manages form configuration state with caching and real-time updates.

```javascript
import { useFormConfigStore } from '@/stores/formConfig'

const formConfigStore = useFormConfigStore()

// Load form configuration
await formConfigStore.loadFormFields('patient')

// Get enabled fields
const enabledFields = formConfigStore.getEnabledFields('patient')
```

#### 3. Composable (`/src/composables/useFormConfig.js`)

Provides reusable form configuration logic for components.

```javascript
import { useFormConfig } from '@/composables/useFormConfig'

const {
  fields,
  enabledFields,
  requiredFields,
  isFieldEnabled,
  isFieldRequired,
  validateForm
} = useFormConfig('patient')
```

#### 4. Admin Interface (`/src/components/admin/FormConfigManager.vue`)

Complete admin interface for configuring form fields:

- View all form types
- Toggle field visibility
- Set field requirements
- Reorder fields
- Reset to defaults
- Live preview

#### 5. Dynamic Form Component (`/src/components/DynamicPatientForm.vue`)

Example implementation showing how forms adapt to configuration:

- Renders only enabled fields
- Enforces required field validation
- Groups fields by category
- Provides debug information

## Usage

### For Administrators

1. **Access Form Configuration**
   - Navigate to Admin Dashboard (`/admin`)
   - Click "Configure Forms" or select "Form Configuration" tab
   - Select the form type you want to configure

2. **Configure Fields**
   - Toggle fields on/off using the switches
   - Change required/optional status for enabled fields
   - Reorder fields using the up/down arrows
   - Use categories to organize fields

3. **Preview Changes**
   - Click "Show Preview" to see how the form will look
   - Test the form functionality using the Form Demo (`/form-demo`)

4. **Save and Reset**
   - Changes are saved automatically when modified
   - Use "Reset to Defaults" to restore original configuration

### For Developers

#### Creating a Dynamic Form

```vue
<template>
  <form @submit.prevent="handleSubmit">
    <div v-for="category in fieldsByCategory" :key="category">
      <h3>{{ category }}</h3>
      <div v-for="field in enabledFields" :key="field.field_id">
        <label>
          {{ field.display_name }}
          <span v-if="field.is_required">*</span>
        </label>
        <input
          v-model="formData[field.name]"
          :type="field.field_type"
          :required="field.is_required"
          :disabled="!field.is_enabled"
        />
      </div>
    </div>
  </form>
</template>

<script setup>
import { useFormConfig } from '@/composables/useFormConfig'

const {
  enabledFields,
  fieldsByCategory,
  formData,
  validateForm,
  initializeFormData
} = useFormConfig('patient')

await initializeFormData()

const handleSubmit = () => {
  if (validateForm()) {
    // Submit form data
  }
}
</script>
```

#### Integrating with Existing Forms

1. Replace hardcoded field lists with dynamic field loading
2. Use composable for field validation and rendering
3. Apply conditional rendering based on field configuration
4. Implement proper error handling and loading states

## Field Configuration Schema

Each field has the following configuration:

```json
{
  "field_id": 1,
  "name": "first_name",
  "display_name": "First Name",
  "field_type": "text",
  "is_enabled": true,
  "is_required": true,
  "is_core": true,
  "category": "Personal Information",
  "description": "Patient first name",
  "options": [],
  "sort_order": 1
}
```

### Field Properties

- **field_id**: Unique identifier for the field
- **name**: Internal field name (used in API/database)
- **display_name**: Human-readable label shown to users
- **field_type**: Type of input (text, email, select, etc.)
- **is_enabled**: Whether the field is visible in forms
- **is_required**: Whether the field is mandatory
- **is_core**: Whether the field cannot be disabled (core business logic)
- **category**: Grouping category for organization
- **description**: Help text or additional information
- **options**: Array of options for select/multiselect fields
- **sort_order**: Order in which field appears

### Supported Field Types

- `text` - Single line text input
- `email` - Email address input with validation
- `tel` - Phone number input
- `url` - URL input with validation
- `number` - Numeric input
- `date` - Date picker
- `textarea` - Multi-line text input
- `select` - Dropdown selection
- `multiselect` - Multiple selection dropdown
- `boolean` - Checkbox input

## Features

### Core Features

- ✅ Dynamic field visibility control
- ✅ Required/optional field configuration
- ✅ Field reordering with drag-and-drop interface
- ✅ Category-based field organization
- ✅ Core field protection (cannot be disabled)
- ✅ Form preview functionality
- ✅ Reset to defaults
- ✅ Real-time configuration updates

### Validation Features

- ✅ Required field validation
- ✅ Type-specific validation (email, phone, etc.)
- ✅ Custom validation messages
- ✅ Real-time validation feedback
- ✅ Form-level validation status

### Admin Features

- ✅ Clean, professional admin interface
- ✅ Visual indicators for field status
- ✅ Bulk operations support
- ✅ Configuration history (via backend)
- ✅ Role-based access control

## Testing

### Demo Application

Access the form demo at `/form-demo` to test the dynamic form functionality:

1. Configure fields in the admin panel
2. Navigate to the form demo
3. Observe how the form adapts to your configuration
4. Test form validation and submission

### Manual Testing Checklist

- [ ] Form loads with correct field configuration
- [ ] Only enabled fields are visible
- [ ] Required fields show validation errors
- [ ] Field order matches configuration
- [ ] Core fields cannot be disabled
- [ ] Form preview matches actual form
- [ ] Configuration changes persist
- [ ] Reset to defaults works correctly

## Integration Points

### Existing Components

The form configuration system integrates with:

- **Patient Registration** (`/src/views/Patients.vue`)
- **Appointment Booking** (various appointment components)
- **Admin Dashboard** (`/src/views/AdminDashboard.vue`)

### Future Enhancements

Planned improvements include:

- Field validation rules configuration
- Conditional field display (show field X if field Y has value Z)
- Multi-step form configuration
- Field templating and reuse
- Import/export configuration
- A/B testing support

## Troubleshooting

### Common Issues

1. **Fields not loading**
   - Check network requests for API errors
   - Verify authentication and permissions
   - Check browser console for JavaScript errors

2. **Configuration not saving**
   - Ensure admin role access
   - Check for validation errors in field updates
   - Verify network connectivity

3. **Form not updating**
   - Clear browser cache
   - Check if form is using dynamic configuration
   - Verify reactive data binding

### Debug Information

Enable debug mode in the dynamic form component to see:
- Current form data values
- Active validation errors
- Field configuration status
- API request/response details

## Security Considerations

- Form configuration is restricted to admin users only
- Core business fields cannot be disabled
- Field validation is enforced on both client and server
- Configuration changes are logged for audit trails
- Sensitive field data is protected during configuration

## Performance

The system is optimized for performance with:

- Configuration caching in Pinia store
- Lazy loading of form configurations
- Efficient re-rendering with Vue 3 reactivity
- Minimal API requests through smart caching
- Optimized bundle size with tree-shaking