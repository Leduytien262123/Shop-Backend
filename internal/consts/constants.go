package consts

const (
	// JWT
	JWT_SECRET_KEY   = "your-secret-key-here"
	JWT_EXPIRE_HOURS = 24

	// User roles
	ROLE_ADMIN = "admin"
	ROLE_USER  = "user"

	// Response messages
	MSG_SUCCESS             = "Success"
	MSG_INVALID_CREDENTIALS = "Invalid credentials"
	MSG_USER_NOT_FOUND      = "User not found"
	MSG_EMAIL_EXISTS        = "Email already exists"
	MSG_USERNAME_EXISTS     = "Username already exists"
	MSG_UNAUTHORIZED        = "Unauthorized"
	MSG_FORBIDDEN           = "Forbidden"
	MSG_INTERNAL_ERROR      = "Internal server error"
	MSG_VALIDATION_ERROR    = "Validation error"
)
