package main

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"time"
	
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db                 *sql.DB
	formConfigService  *FormConfigService
	translationService *TranslationService
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(user *User) error {
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	
	now := time.Now()
	user.IsActive = true
	user.CreatedAt = now
	user.UpdatedAt = now

	err := s.db.QueryRow(
		query,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Role,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

// GetUserByEmail gets user by email
func (s *UserService) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, 
		       COALESCE(healthcare_entity_id, 0) as healthcare_entity_id,
		       COALESCE(license_number, '') as license_number, 
		       COALESCE(specialization, '') as specialization, 
		       COALESCE(preferred_locale, '') as preferred_locale, 
		       is_active, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_active = true
	`

	err := s.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.HealthcareEntityID,
		&user.LicenseNumber,
		&user.Specialization,
		&user.PreferredLocale,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// GetUserByID gets user by ID
func (s *UserService) GetUserByID(id int) (*User, error) {
	user := &User{}
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, 
		       COALESCE(healthcare_entity_id, 0) as healthcare_entity_id,
		       COALESCE(license_number, '') as license_number, 
		       COALESCE(specialization, '') as specialization, 
		       COALESCE(preferred_locale, '') as preferred_locale, 
		       is_active, created_at, updated_at
		FROM users
		WHERE id = $1 AND is_active = true
	`

	err := s.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.HealthcareEntityID,
		&user.LicenseNumber,
		&user.Specialization,
		&user.PreferredLocale,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(user *User) error {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, role = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND is_active = true
		RETURNING updated_at
	`

	err := s.db.QueryRow(
		query,
		user.FirstName,
		user.LastName,
		user.Role,
		user.ID,
	).Scan(&user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	return nil
}

// GetUsers gets all users (admin only)
func (s *UserService) GetUsers(limit, offset int) ([]User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at
		FROM users
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.FirstName,
			&user.LastName,
			&user.Role,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetDoctorsByEntity gets all doctors for a specific healthcare entity
func (s *UserService) GetDoctorsByEntity(healthcareEntityID int) ([]User, error) {
	// Query doctors filtered by healthcare entity for proper multi-tenant isolation
	query := `
		SELECT id, first_name, last_name, email, role, 
			   COALESCE(specialization, '') as specialization
		FROM users
		WHERE is_active = true AND role = 'doctor' AND healthcare_entity_id = $1
		ORDER BY first_name, last_name
	`
	
	rows, err := s.db.Query(query, healthcareEntityID)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()
	
	var doctors []User
	for rows.Next() {
		var doctor User
		
		err := rows.Scan(
			&doctor.ID,
			&doctor.FirstName,
			&doctor.LastName,
			&doctor.Email,
			&doctor.Role,
			&doctor.Specialization,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		
		doctors = append(doctors, doctor)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	
	return doctors, nil
}

// EmailExists checks if email already exists
func (s *UserService) EmailExists(email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = $1`
	
	err := s.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// DeactivateUser soft deletes a user
func (s *UserService) DeactivateUser(id int) error {
	query := `
		UPDATE users
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// GetHealthcareEntityByID gets healthcare entity by ID
func (s *UserService) GetHealthcareEntityByID(id int) (*HealthcareEntity, error) {
	entity := &HealthcareEntity{}
	query := `
		SELECT id, name, type, country_id, address, state_id, city_id, country, state, city,
		       postal_code, phone, email, website, license, tax_id, timezone, locale, currency, 
		       require_room_assignment, is_active, created_at, updated_at
		FROM healthcare_entities
		WHERE id = $1
	`

	// Use temporary variables for nullable fields
	var countryID, stateID, cityID sql.NullInt32
	
	err := s.db.QueryRow(query, id).Scan(
		&entity.ID,
		&entity.Name,
		&entity.Type,
		&countryID,
		&entity.Address,
		&stateID,
		&cityID,
		&entity.Country,
		&entity.State,
		&entity.City,
		&entity.PostalCode,
		&entity.Phone,
		&entity.Email,
		&entity.Website,
		&entity.License,
		&entity.TaxID,
		&entity.Timezone,
		&entity.Language,  // Maps to locale in DB
		&entity.Currency,
		&entity.RequireRoomAssignment,
		&entity.IsActive,
		&entity.CreatedAt,
		&entity.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("healthcare entity not found")
		}
		return nil, err
	}

	// Convert sql.NullInt32 to *int
	fmt.Printf("DEBUG: countryID.Valid=%t, countryID.Int32=%d\n", countryID.Valid, countryID.Int32)
	fmt.Printf("DEBUG: stateID.Valid=%t, stateID.Int32=%d\n", stateID.Valid, stateID.Int32)
	fmt.Printf("DEBUG: cityID.Valid=%t, cityID.Int32=%d\n", cityID.Valid, cityID.Int32)
	
	if countryID.Valid {
		countryVal := int(countryID.Int32)
		entity.CountryID = &countryVal
		fmt.Printf("DEBUG: Set entity.CountryID to %d\n", countryVal)
	} else {
		fmt.Printf("DEBUG: countryID is not valid, leaving nil\n")
	}
	if stateID.Valid {
		stateVal := int(stateID.Int32)
		entity.StateID = &stateVal
	}
	if cityID.Valid {
		cityVal := int(cityID.Int32)
		entity.CityID = &cityVal
	}

	return entity, nil
}

// CreateDoctorWithTempPassword creates a doctor with a temporary password
func (s *UserService) CreateDoctorWithTempPassword(req AdminCreateDoctorRequest, healthcareEntityID int) (*AdminCreateDoctorResponse, error) {
	// Generate temporary password (8 characters with letters and numbers)
	tempPassword := generateTempPassword()
	
	// Hash the temporary password
	hashedPassword, err := hashPassword(tempPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}
	
	// Set temp password expiry (7 days from now)
	expiryTime := time.Now().Add(7 * 24 * time.Hour)
	
	query := `
		INSERT INTO users (
			email, password_hash, first_name, last_name, role, 
			healthcare_entity_id, license_number, specialization, 
			preferred_locale, is_active, is_temp_password, temp_password_expires,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		) RETURNING id, created_at, updated_at
	`
	
	now := time.Now()
	var user User
	user.Email = req.Email
	user.Password = hashedPassword
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Role = "doctor"
	user.HealthcareEntityID = healthcareEntityID
	user.LicenseNumber = req.LicenseNumber
	user.Specialization = req.Specialization
	user.PreferredLocale = req.PreferredLocale
	user.IsActive = true
	user.IsTempPassword = true
	user.TempPasswordExpires = &expiryTime

	err = s.db.QueryRow(
		query,
		user.Email, hashedPassword, user.FirstName, user.LastName, user.Role,
		user.HealthcareEntityID, user.LicenseNumber, user.Specialization,
		user.PreferredLocale, user.IsActive, user.IsTempPassword, user.TempPasswordExpires,
		now, now,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create doctor: %v", err)
	}

	return &AdminCreateDoctorResponse{
		Doctor:          user.ToUserResponse(),
		TempPassword:    tempPassword,
		PasswordExpires: expiryTime,
		Message:         "Doctor created successfully. Please provide the temporary password to the doctor.",
	}, nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(userID int, currentPassword, newPassword string) error {
	// Get user's current password hash
	var currentHash string
	var isTempPassword bool
	
	query := "SELECT password_hash, is_temp_password FROM users WHERE id = $1"
	err := s.db.QueryRow(query, userID).Scan(&currentHash, &isTempPassword)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}
	
	// Verify current password
	if !checkPasswordHash(currentPassword, currentHash) {
		return errors.New("invalid current password")
	}
	
	// Hash new password
	newHash, err := hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %v", err)
	}
	
	// Update password and clear temp password flags
	updateQuery := `
		UPDATE users 
		SET password_hash = $1, is_temp_password = false, temp_password_expires = NULL, updated_at = $2
		WHERE id = $3
	`
	
	_, err = s.db.Exec(updateQuery, newHash, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}
	
	return nil
}

// generateTempPassword generates a secure 8-character temporary password
func generateTempPassword() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8
	
	password := make([]byte, length)
	for i := range password {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		password[i] = charset[num.Int64()]
	}
	
	return string(password)
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// checkPasswordHash checks if a password matches its hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}