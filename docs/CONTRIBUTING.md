# Contributing to Healthcare Platform

Thank you for your interest in contributing to the Healthcare Management Platform! This document provides guidelines and information for contributors.

## Table of Contents
- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing Requirements](#testing-requirements)
- [Pull Request Process](#pull-request-process)
- [Issue Guidelines](#issue-guidelines)

## Code of Conduct

### Our Pledge
We pledge to create a welcoming, inclusive environment for all contributors regardless of background, experience level, gender, gender identity, sexual orientation, disability, personal appearance, body size, race, ethnicity, age, religion, or nationality.

### Standards
- Use welcoming and inclusive language
- Respect differing viewpoints and experiences
- Accept constructive criticism gracefully
- Focus on what's best for the community
- Show empathy towards other community members

### Enforcement
Instances of abusive, harassing, or otherwise unacceptable behavior may be reported to the project maintainers.

## Getting Started

### Prerequisites
- Read the [Development Setup Guide](development/setup.md)
- Familiarize yourself with the [Architecture Documentation](architecture/README.md)
- Review the [API Documentation](api/README.md)

### First-Time Setup
1. Fork the repository
2. Clone your fork locally
3. Set up the development environment
4. Create a new branch for your changes
5. Make your changes and test thoroughly
6. Submit a pull request

## Development Workflow

### Branch Naming
Use descriptive branch names with prefixes:
- `feature/` - New features
- `bugfix/` - Bug fixes
- `hotfix/` - Critical production fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring
- `test/` - Test improvements

Examples:
- `feature/patient-search-filters`
- `bugfix/appointment-conflict-detection`
- `docs/api-authentication-guide`

### Commit Messages
Follow conventional commit format:
```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

Types:
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation changes
- `style` - Code formatting (no functional changes)
- `refactor` - Code refactoring
- `test` - Adding/updating tests
- `chore` - Maintenance tasks

Examples:
```
feat(patient): add advanced search filters

- Add filters for gender, insurance, and age range
- Implement server-side filtering logic
- Add frontend UI components for filters

Closes #123
```

```
fix(auth): resolve JWT token expiration issue

The token validation was not properly checking expiration time,
causing expired tokens to be accepted.

Fixes #456
```

## Coding Standards

### Go Backend Standards

#### Code Style
- Follow standard Go formatting (`gofmt`)
- Use `golangci-lint` for linting
- Follow Go naming conventions
- Keep functions focused and small
- Use meaningful variable names

#### Error Handling
```go
// Good - Proper error handling
if err != nil {
    log.Printf("Failed to create user: %v", err)
    return nil, fmt.Errorf("user creation failed: %w", err)
}

// Bad - Ignoring errors
user, _ := userService.CreateUser(userData)
```

#### Database Queries
```go
// Good - Parameterized queries
query := "SELECT * FROM users WHERE email = $1"
row := db.QueryRow(query, email)

// Bad - String concatenation (SQL injection risk)
query := "SELECT * FROM users WHERE email = '" + email + "'"
```

#### API Responses
```go
// Good - Consistent response format
c.JSON(http.StatusOK, gin.H{
    "data": user,
    "message": "User retrieved successfully",
})

// Error responses
c.JSON(http.StatusBadRequest, gin.H{
    "error": "Invalid input",
    "code": "VALIDATION_ERROR",
    "message": "Email format is invalid",
})
```

### Frontend Standards

#### Vue.js Components
```vue
<!-- Good - Composition API with proper structure -->
<template>
  <div class="patient-card">
    <h3>{{ patient.fullName }}</h3>
    <p>{{ patient.email }}</p>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  patient: {
    type: Object,
    required: true
  }
})

const fullName = computed(() => 
  `${props.patient.first_name} ${props.patient.last_name}`
)
</script>
```

#### CSS/Tailwind
```vue
<!-- Good - Semantic class combinations -->
<button class="btn btn-primary">
  Save Patient
</button>

<!-- Good - Responsive design -->
<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
  <!-- Content -->
</div>
```

### Security Standards

#### Input Validation
```go
// Always validate input
type UserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

if err := validator.Struct(userReq); err != nil {
    return c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}
```

#### Authentication
```go
// Always check authentication
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
            c.Abort()
            return
        }
        // Validate token...
    }
}
```

## Testing Requirements

### Backend Testing
All backend changes must include appropriate tests:

#### Unit Tests
```go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        user    *User
        wantErr bool
    }{
        {
            name: "valid user",
            user: &User{
                Email:     "test@example.com",
                FirstName: "John",
                LastName:  "Doe",
                Role:      "doctor",
            },
            wantErr: false,
        },
        {
            name: "invalid email",
            user: &User{
                Email: "invalid-email",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := userService.CreateUser(tt.user)
            if (err != nil) != tt.wantErr {
                t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

#### Integration Tests
Test API endpoints with real database:
```go
func TestUserHandler_Register(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer db.Close()
    
    // Create test request
    req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(payload))
    w := httptest.NewRecorder()
    
    // Execute request
    router.ServeHTTP(w, req)
    
    // Assert response
    assert.Equal(t, http.StatusCreated, w.Code)
}
```

### Frontend Testing
```javascript
// Component tests
import { mount } from '@vue/test-utils'
import PatientCard from '@/components/PatientCard.vue'

describe('PatientCard', () => {
  it('displays patient name correctly', () => {
    const patient = {
      first_name: 'John',
      last_name: 'Doe',
      email: 'john@example.com'
    }
    
    const wrapper = mount(PatientCard, {
      props: { patient }
    })
    
    expect(wrapper.text()).toContain('John Doe')
  })
})
```

### Test Coverage Requirements
- **Backend**: Minimum 80% coverage
- **Frontend**: Minimum 70% coverage
- **Critical paths**: 100% coverage (authentication, data validation)

## Pull Request Process

### Before Submitting
1. **Test thoroughly**
   - Run all tests: `go test ./...`
   - Test frontend: `npm test`
   - Manual testing of changed functionality

2. **Code quality checks**
   - Run linters: `golangci-lint run`
   - Format code: `gofmt -w .`
   - Check for vulnerabilities

3. **Documentation**
   - Update API documentation if needed
   - Add/update code comments
   - Update README if necessary

### PR Template
Use this template for pull requests:

```markdown
## Description
Brief description of changes and why they're needed.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed
- [ ] Test coverage maintained/improved

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Comments added to hard-to-understand areas
- [ ] Documentation updated
- [ ] No console.log or debug statements

## Related Issues
Closes #issue_number
```

### Review Process
1. **Automated checks** must pass (CI/CD pipeline)
2. **Peer review** by at least one maintainer
3. **Security review** for sensitive changes
4. **Documentation review** if docs are updated

### Merge Requirements
- All CI checks pass ✅
- At least one approval from maintainer ✅
- No requested changes ✅
- Branch is up to date with main ✅

## Issue Guidelines

### Bug Reports
Use the bug report template:
```markdown
**Describe the bug**
Clear description of the issue.

**To Reproduce**
Steps to reproduce:
1. Go to '...'
2. Click on '....'
3. See error

**Expected behavior**
What should happen.

**Screenshots**
If applicable, add screenshots.

**Environment:**
- OS: [e.g. macOS]
- Browser: [e.g. Chrome 91]
- Version: [e.g. 1.0.0]

**Additional context**
Any other context about the problem.
```

### Feature Requests
```markdown
**Feature Summary**
Brief description of the feature.

**Problem it Solves**
What problem does this solve?

**Proposed Solution**
Describe your proposed solution.

**Alternatives Considered**
Other solutions you considered.

**Additional Context**
Any other context or screenshots.
```

## Development Guidelines

### Performance Considerations
- Database queries should be optimized
- Use pagination for large datasets
- Implement proper caching strategies
- Monitor API response times

### Security Best Practices
- Never commit secrets or credentials
- Always validate and sanitize inputs
- Use parameterized database queries
- Implement proper authorization checks
- Follow OWASP security guidelines

### Accessibility
- Follow WCAG 2.1 guidelines
- Ensure keyboard navigation works
- Provide proper ARIA labels
- Test with screen readers
- Maintain good color contrast

## Getting Help

### Resources
- [Development Setup](development/setup.md)
- [Architecture Guide](architecture/README.md)
- [API Documentation](api/README.md)
- [Troubleshooting](troubleshooting.md)

### Communication
- Create issues for bugs and feature requests
- Join development discussions
- Ask questions in issue comments
- Reach out to maintainers for guidance

### Mentorship
New contributors are welcome! Maintainers are happy to:
- Help with first contributions
- Provide code review feedback
- Guide through the development process
- Answer questions about the codebase

Thank you for contributing to the Healthcare Platform! 🏥✨