package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mohammadrezajavid-lab/goauth/goauthapp/service/goauth"
	"github.com/mohammadrezajavid-lab/goauth/pkg/logger"
	"log/slog"
	"strings"
)

// UserRepository struct handles database operations for users.
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new user repository.
func NewUserRepository(db *pgxpool.Pool) goauth.UserRepository {
	return &UserRepository{
		db: db,
	}
}

// FindByPhoneNumber finds a user by their phone number.
// It returns goauth.ErrUserNotFound if the user does not exist.
func (r *UserRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*goauth.User, error) {
	log := logger.L().With(slog.String("phone_number", phoneNumber))
	log.Info("Finding user by phone number in database")

	query := `
		SELECT id, phone_number, created_at, updated_at 
		FROM users 
		WHERE phone_number = $1
	`
	row := r.db.QueryRow(ctx, query, phoneNumber)

	var user goauth.User
	err := row.Scan(&user.ID, &user.PhoneNumber, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("User not found in database")
			// Return a specific, known error for the service layer to handle.
			return nil, goauth.ErrUserNotFound
		}
		// For any other database error, log it and return a generic error.
		log.Error("Failed to find user by phone number", slog.String("error", err.Error()))
		return nil, err
	}

	log.Debug("User found successfully")
	return &user, nil
}

// Create inserts a new user into the database.
// It updates the user struct passed by reference with the database-generated values (ID, CreatedAt, UpdatedAt).
func (r *UserRepository) Create(ctx context.Context, user *goauth.User) error {
	log := logger.L().With(slog.String("phone_number", user.PhoneNumber))
	log.Info("Creating a new user in the database")

	query := `
		INSERT INTO users (phone_number) 
		VALUES ($1) 
		RETURNING id, created_at, updated_at
	`
	// Using QueryRow because RETURNING clause returns a single row with the new data.
	err := r.db.QueryRow(ctx, query, user.PhoneNumber).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		// Check for unique constraint violation.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 is the code for unique_violation
			log.Warn("Attempted to create a user with a duplicate phone number")
			return goauth.ErrUserAlreadyExists
		}

		log.Error("Failed to create user", slog.String("error", err.Error()))
		return err
	}

	log.Debug("User created successfully", slog.Int64("user_id", user.ID))
	return nil
}

// FindByID finds a single user by their unique ID.
func (r *UserRepository) FindByID(ctx context.Context, id int64) (*goauth.User, error) {
	log := logger.L().With(slog.Int64("user_id", id))
	log.Info("Finding user by ID in database")

	query := `
		SELECT id, phone_number, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)

	var user goauth.User
	err := row.Scan(&user.ID, &user.PhoneNumber, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("User with given ID not found")
			return nil, goauth.ErrUserNotFound
		}
		log.Error("Failed to find user by ID", slog.String("error", err.Error()))
		return nil, err
	}

	log.Debug("User found successfully by ID")
	return &user, nil
}

// List returns a paginated and searchable list of users.
func (r *UserRepository) List(ctx context.Context, params goauth.ListUsersParams) ([]goauth.User, int, error) {
	log := logger.L().With(slog.Any("params", params))
	log.Info("Listing users from database")

	// --- 1. Build dynamic query for search ---
	var args []interface{}
	var conditions []string
	argID := 1

	if params.Search != "" {
		conditions = append(conditions, fmt.Sprintf("phone_number LIKE $%d", argID))
		args = append(args, "%"+params.Search+"%")
		argID++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// --- 2. Execute query to get the total count ---
	countQuery := "SELECT COUNT(*) FROM users " + whereClause
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		log.Error("Failed to count users", slog.String("error", err.Error()))
		return nil, 0, err
	}

	if total == 0 {
		return []goauth.User{}, 0, nil
	}

	// --- 3. Build the main query for fetching the user list with pagination ---
	query := fmt.Sprintf(`
		SELECT id, phone_number, created_at, updated_at
		FROM users
		%s
		ORDER BY created_at DESC
		LIMIT $%d
		OFFSET $%d
	`, whereClause, argID, argID+1)

	offset := (params.Page - 1) * params.PageSize
	args = append(args, params.PageSize, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		log.Error("Failed to list users", slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer rows.Close()

	// --- 4. Scan the results ---
	var users []goauth.User
	for rows.Next() {
		var user goauth.User
		if err := rows.Scan(&user.ID, &user.PhoneNumber, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Error("Failed to scan user row", slog.String("error", err.Error()))
			return nil, 0, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Error("Error during user list iteration", slog.String("error", err.Error()))
		return nil, 0, err
	}

	log.Debug("Users listed successfully", slog.Int("count", len(users)), slog.Int("total", total))
	return users, total, nil
}
