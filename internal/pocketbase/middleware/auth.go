package middleware

// TODO: JWT validation middleware for custom routes.
// Checks Authorization header or pb_auth cookie.
// Can wrap apis.RequireAuth() or implement custom token validation logic.
// Used by routes package via Bind() on protected endpoints.
