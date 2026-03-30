package middleware

// TODO: Role-based access control middleware.
// Runs after auth middleware — requires authenticated user.
// Checks user record fields (e.g., role, isAdmin) to authorize access.
// Returns 403 if the user lacks the required role.
