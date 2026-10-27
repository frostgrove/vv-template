package product

import "github.com/frostgrove/vv/auth"

// The permissions a caller needs to reach the product routes. An API key (or a
// token) carries them; the routes declare and check them (product.http-handler.go).
const (
	PermRead   auth.Permission = "product:read"
	PermWrite  auth.Permission = "product:write"
	PermDelete auth.Permission = "product:delete"
)
