# Form Configuration System - Testing Guide

## Overview
This guide provides comprehensive testing instructions for the dynamic form field configuration system implemented for the healthcare platform.

## System Components

### Backend Components ✅
- **Database Tables**: `form_types`, `field_definitions`, `entity_field_configurations`
- **Service Layer**: `FormConfigService` in user-service
- **API Endpoints**: Form configuration REST API
- **Migrations**: Database schema version 4

### Frontend Components ✅  
- **Admin Interface**: `FormConfigManager.vue`
- **Dynamic Forms**: Updated `Patients.vue` and `AppointmentEditModal.vue`
- **State Management**: `formConfig.js` Pinia store
- **Composables**: `useFormConfig.js`

## Testing Checklist

### 🔧 Backend API Testing

#### 1. Database Setup
```bash
# Start services
docker compose up -d

# Verify migrations applied
docker exec healthcare-user-db psql -U postgres -d user_service_db -c "SELECT version FROM schema_migrations ORDER BY version;"

# Check form configuration tables exist
docker exec healthcare-user-db psql -U postgres -d user_service_db -c "\dt" | grep -E "(form_types|field_definitions|entity_field_configurations)"
```

#### 2. API Endpoint Testing
```bash
# Test form types endpoint (Admin only)
curl -X GET "http://localhost:8080/api/forms/types" \
  -H "Authorization: Bearer <admin-jwt-token>" \
  -H "X-User-Role: admin" \
  -H "X-Healthcare-Entity-ID: 1"

# Test form metadata endpoint
curl -X GET "http://localhost:8080/api/forms/patient/metadata" \
  -H "Authorization: Bearer <admin-jwt-token>" \
  -H "X-User-Role: admin" \
  -H "X-Healthcare-Entity-ID: 1"

# Test field configuration update
curl -X PUT "http://localhost:8080/api/forms/patient/fields/1" \
  -H "Authorization: Bearer <admin-jwt-token>" \
  -H "X-User-Role: admin" \
  -H "X-Healthcare-Entity-ID: 1" \
  -H "Content-Type: application/json" \
  -d '{
    "is_enabled": false,
    "is_required": false,
    "custom_label": "Test Field",
    "sort_order": 5
  }'
```

#### 3. Core Field Protection Testing
```bash
# Attempt to disable core field (should fail)
curl -X PUT "http://localhost:8080/api/forms/patient/fields/1" \
  -H "Authorization: Bearer <admin-jwt-token>" \
  -H "X-User-Role: admin" \
  -H "X-Healthcare-Entity-ID: 1" \
  -H "Content-Type: application/json" \
  -d '{
    "is_enabled": false,
    "is_required": true
  }' 
# Expected: 400 Bad Request with "cannot disable core field" message
```

#### 4. Role-Based Access Testing
```bash
# Test with non-admin user (should fail)
curl -X GET "http://localhost:8080/api/forms/types" \
  -H "Authorization: Bearer <doctor-jwt-token>" \
  -H "X-User-Role: doctor" \
  -H "X-Healthcare-Entity-ID: 1"
# Expected: 403 Forbidden
```

### 🎨 Frontend Testing

#### 1. Admin Interface Testing
1. **Access Form Configuration**
   - Login as admin user
   - Navigate to `/admin`
   - Click "Configure Forms" button
   - Verify FormConfigManager loads

2. **Form Type Selection**
   - Select "Patient Registration" form
   - Verify patient fields load with categories
   - Select "Appointment Booking" form  
   - Verify appointment fields load

3. **Field Configuration**
   - Toggle field visibility on/off
   - Toggle required/optional status
   - Verify core fields cannot be disabled
   - Test field reordering (up/down arrows)
   - Test custom labels
   - Save changes and verify persistence

4. **Form Preview**
   - Use live preview panel
   - Verify changes reflect immediately
   - Test with different field configurations

5. **Reset Functionality**
   - Click "Reset to Defaults"
   - Confirm in modal dialog
   - Verify all settings reset to defaults

#### 2. Dynamic Form Testing
1. **Patient Forms Testing**
   - Navigate to `/patients`
   - Click "Add Patient"
   - Verify only enabled fields show
   - Verify required fields have asterisks
   - Test form validation respects configuration
   - Verify field order matches configuration

2. **Appointment Forms Testing**
   - Navigate to appointments section
   - Open appointment edit modal
   - Verify only enabled fields show
   - Test validation respects configuration

3. **Configuration Loading States**
   - Clear browser cache
   - Reload forms
   - Verify loading spinners show
   - Test error states (disconnect network)
   - Test retry functionality

### 🔒 Security Testing

#### 1. Authorization Testing
- Verify non-admin users cannot access `/api/forms/*` endpoints
- Test JWT token validation
- Verify healthcare entity isolation

#### 2. Input Validation Testing
- Test invalid form type names
- Test invalid field IDs
- Test malformed JSON payloads
- Verify SQL injection protection

### 🏢 Multi-Tenant Testing

#### 1. Entity Isolation Testing
```bash
# Test with different healthcare entities
curl -X GET "http://localhost:8080/api/forms/patient/metadata" \
  -H "Authorization: Bearer <admin-jwt-token>" \
  -H "X-User-Role: admin" \
  -H "X-Healthcare-Entity-ID: 1"

curl -X GET "http://localhost:8080/api/forms/patient/metadata" \
  -H "Authorization: Bearer <admin-jwt-token>" \
  -H "X-User-Role: admin" \
  -H "X-Healthcare-Entity-ID: 2"

# Verify configurations are isolated per entity
```

#### 2. Default Configuration Creation
- Create new healthcare entity
- Verify default field configurations created automatically
- Test that new entity gets all default field settings

## Performance Testing

### 1. Configuration Loading
- Test form configuration load times
- Verify caching is working
- Test concurrent admin users
- Monitor database query performance

### 2. Form Rendering
- Test large forms (25+ fields)
- Verify dynamic rendering performance
- Test form submission with many fields

## Error Scenario Testing

### 1. Database Connectivity
- Simulate database connection loss
- Verify proper error messages
- Test retry mechanisms

### 2. API Failures
- Simulate 500 errors from backend
- Test form behavior during API failures
- Verify error messages are user-friendly

### 3. Configuration Corruption
- Test with invalid JSON in configuration
- Test missing field definitions
- Verify graceful handling

## User Acceptance Testing

### 1. Admin Workflow
1. Admin logs in
2. Navigates to form configuration
3. Selects patient form
4. Hides "Blood Type" field
5. Makes "Insurance" field required
6. Saves changes
7. Tests patient form shows changes
8. Resets to defaults
9. Verifies reset worked

### 2. Form User Workflow
1. Staff user tries to add patient
2. Form loads with configured fields only
3. Required fields show proper validation
4. Hidden fields don't appear
5. Form submits successfully
6. Edit existing patient shows same configuration

## Expected Test Results

### ✅ Pass Criteria
- All API endpoints return expected responses
- Admin interface loads and functions correctly
- Forms render dynamically based on configuration
- Core fields cannot be disabled
- Non-admin users cannot access configuration
- Multi-tenant isolation works correctly
- Error states handled gracefully
- Performance is acceptable (<2s form load)

### ❌ Fail Criteria
- Forms show fallback/default fields
- Core fields can be disabled
- Non-admin users can access configuration
- Data leaks between healthcare entities
- Forms break when configuration unavailable
- Performance issues (>5s form load)

## Troubleshooting

### Common Issues
1. **Forms not loading configuration**
   - Check JWT token validity
   - Verify X-Healthcare-Entity-ID header
   - Check browser network tab for API errors

2. **Configuration changes not saving**
   - Verify admin role in JWT token
   - Check backend logs for validation errors
   - Verify database migrations applied

3. **Core fields being disabled**
   - Check field definitions in database
   - Verify is_core flag is set correctly
   - Check backend validation logic

### Debug Commands
```bash
# Check form configuration data
docker exec healthcare-user-db psql -U postgres -d user_service_db -c "
SELECT ft.name, fd.display_name, fd.is_core, efc.is_enabled, efc.is_required 
FROM form_types ft
JOIN field_definitions fd ON ft.id = fd.form_type_id
LEFT JOIN entity_field_configurations efc ON fd.id = efc.field_definition_id
WHERE efc.healthcare_entity_id = 1 OR efc.healthcare_entity_id IS NULL
ORDER BY ft.name, fd.sort_order;
"

# Check user roles
docker exec healthcare-user-db psql -U postgres -d user_service_db -c "
SELECT email, role, healthcare_entity_id FROM users WHERE role = 'admin';
"
```

## Final Validation

Before considering the feature complete, ensure:
- [ ] All backend tests pass
- [ ] All frontend functionality works
- [ ] Security measures are in place
- [ ] Multi-tenancy is properly implemented
- [ ] Performance is acceptable
- [ ] Documentation is complete
- [ ] Error handling is comprehensive

The form configuration system is ready for production use when all test criteria are met.