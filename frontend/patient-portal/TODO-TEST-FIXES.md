# Test Fixes Todo List

## Overview
This todo list tracks the remaining test failures that need to be fixed in the healthcare platform frontend. 

**Current Status**: 161/209 tests passing (~77% pass rate)  
**Target**: 100% pass rate

## Test Failures by Category

### 🔌 API Mocking Issues (Priority: High)

- [ ] **Task #21**: Fix API mocking in useFormConfig tests - formsApi export missing
  - **Issue**: `No "formsApi" export is defined on the "@/api/forms" mock`
  - **Files**: `src/composables/__tests__/useFormConfig.test.js`
  - **Solution**: Fix mock export structure in vi.mock() calls

- [ ] **Task #26**: Fix useLocationData API mocking and response handling
  - **Issue**: Mock functions not properly configured for locationsApi
  - **Files**: `src/composables/__tests__/useLocationData.test.js`
  - **Solution**: Update mock setup to handle async responses correctly

- [ ] **Task #34**: Add missing @pinia/testing dependency properly
  - **Issue**: Pinia testing utilities not properly configured
  - **Files**: `package.json`, test setup files
  - **Solution**: Configure Pinia testing with proper createTestingPinia setup

- [ ] **Task #35**: Fix unhandled promise rejections in test environment
  - **Issue**: 30 unhandled promise rejections causing test instability
  - **Files**: Various test files
  - **Solution**: Add proper async/await handling and error catching

- [ ] **Task #38**: Fix API client mocking for consistent behavior
  - **Issue**: Inconsistent API client mocking across test files
  - **Files**: `src/test/setup.js`, various test files
  - **Solution**: Create centralized API mock setup

### 📝 useFormState Issues (Priority: High)

- [ ] **Task #22**: Fix useFormState default value generation for different field types
  - **Issue**: `expected undefined to deeply equal ''`
  - **Files**: `src/composables/__tests__/useFormState.test.js`
  - **Solution**: Fix default value logic in useFormState.js

- [ ] **Task #23**: Fix useFormState form submission handling for missing record ID
  - **Issue**: `expected undefined to be null`
  - **Files**: `src/composables/__tests__/useFormState.test.js`
  - **Solution**: Add proper null handling for missing record IDs

- [ ] **Task #24**: Fix useFormState error management delegation methods
  - **Issue**: `expected 'undefined' to be 'function'`
  - **Files**: `src/composables/__tests__/useFormState.test.js`
  - **Solution**: Ensure error management methods are properly exposed

- [ ] **Task #25**: Fix useFormState null initial data handling
  - **Issue**: `Cannot read properties of null`
  - **Files**: `src/composables/useFormState.js`
  - **Solution**: Add null check before accessing initialData properties

### 🗺️ useLocationData Issues (Priority: Medium)

- [ ] **Task #27**: Fix state/city loading with proper mock responses
  - **Issue**: `expected [] to deeply equal [expected array]`
  - **Files**: `src/composables/__tests__/useLocationData.test.js`
  - **Solution**: Fix mock responses to return proper data arrays

- [ ] **Task #28**: Fix resolveCountryCode function logic
  - **Issue**: `expected '' to be 'CA'`
  - **Files**: `src/composables/__tests__/useLocationData.test.js`
  - **Solution**: Fix country code resolution logic in useLocationData.js

- [ ] **Task #29**: Fix form data clearing in country/state changes
  - **Issue**: Form fields not properly cleared when changing selections
  - **Files**: `src/composables/__tests__/useLocationData.test.js`
  - **Solution**: Fix form data mutation logic in change handlers

- [ ] **Task #30**: Fix concurrent loading and debouncing tests
  - **Issue**: `expected "spy" to be called 1 times, but got 0 times`
  - **Files**: `src/composables/__tests__/useLocationData.test.js`
  - **Solution**: Fix async timing issues in debounced function tests

### ✅ Component & Validation Issues (Priority: Medium)

- [ ] **Task #31**: Fix useValidation phone number validation regex
  - **Issue**: `'555.123.4567'` not recognized as valid phone number
  - **Files**: `src/composables/__tests__/useValidation.test.js`
  - **Solution**: Update phone validation regex to accept dot notation

- [ ] **Task #32**: Fix DynamicPatientForm country/city dependency test
  - **Issue**: City options not properly loaded after country selection
  - **Files**: `src/components/__tests__/DynamicPatientForm.spec.js`
  - **Solution**: Fix async timing and mock setup for location cascade

- [ ] **Task #33**: Fix virtualScrollMetrics scroll session tracking logic
  - **Issue**: `expected 1 to be 2` for scroll sessions
  - **Files**: `src/utils/__tests__/virtualScrollMetrics.test.js`
  - **Solution**: Implement session gap detection logic in trackScroll function

### 🧪 Test Infrastructure (Priority: Low)

- [ ] **Task #36**: Improve console mocking to prevent test interference
  - **Issue**: Console methods causing test output pollution
  - **Files**: `src/test/setup.js`
  - **Solution**: Add comprehensive console method mocking

- [ ] **Task #37**: Add proper test cleanup in afterEach hooks
  - **Issue**: Tests not properly cleaning up state between runs
  - **Files**: Various test files
  - **Solution**: Add consistent cleanup patterns

- [ ] **Task #39**: Add integration test for complete form submission workflow
  - **Issue**: Missing end-to-end form workflow testing
  - **Files**: New integration test file needed
  - **Solution**: Create comprehensive form workflow integration test

- [ ] **Task #40**: Add integration test for location cascade functionality
  - **Issue**: Missing integration test for country -> state -> city flow
  - **Files**: New integration test file needed
  - **Solution**: Create location cascade integration test

## Priority Order
1. **High Priority**: API Mocking & useFormState issues (Tasks #21-25, #26, #34-35, #38)
2. **Medium Priority**: useLocationData & Component issues (Tasks #27-33)
3. **Low Priority**: Test Infrastructure improvements (Tasks #36-37, #39-40)

## Testing Commands
```bash
# Run all tests
npm run test:run

# Run specific test file
npm run test:run -- src/composables/__tests__/useFormConfig.test.js

# Run tests with coverage
npm run test:coverage

# Run tests in watch mode (development)
npm run test
```

## Success Criteria
- [ ] All 209 tests passing (100% pass rate)
- [ ] No unhandled promise rejections
- [ ] Clean test output without console pollution
- [ ] All composables have comprehensive test coverage
- [ ] Integration tests verify complete workflows

---
*Last Updated: 2025-01-09*
*Current Pass Rate: 161/209 tests (~77%)*