package authApi

import (
	"github.com/dtekltd/common/jwt"
	"github.com/gofiber/fiber/v2"
)

func RegisterHandlers(router fiber.Router, keyManager *jwt.KeyManager) {
	router.Post("/login", jwt.WithKeyManager(keyManager, login))
	router.Post("/register", jwt.WithKeyManager(keyManager, register))
	router.Post("/reset-password", resetPassword)
	router.Post("/reset-password-request", resetPasswordRequest)

	// auth
	auth := router.Group("/auth")
	auth.Get("/account", account)

	// auth / enforcer
	handler := NewEnforcerHandler()
	auth.Post("/reload", handler.reloadPolicy)
	auth.Post("/enforce", handler.enforce)

	// auth / policy
	// router = router.Group("/policy", auth.HasRoleMiddleware("manager", "admin"))
	// should control by policy instead of roles
	policy := auth.Group("/policy")
	policy.Get("", handler.getPolicy)
	policy.Post("", handler.addPolicy)
	policy.Delete("", handler.deletePolicy)
	policy.Get("/user", handler.getUserPolicies)
	policy.Get("/user/ext", handler.getUserPolicyEx)
	policy.Delete("/user", handler.deleteUserPolicy)
	policy.Get("/all", handler.getAllPolicy)
	// policy.Get("/ext", handler.getPolicyEx)

	// auth / google
	google := router.Group("/google")
	google.Post("", jwt.WithKeyManager(keyManager, googleLogin))
	google.Get("/callback", jwt.WithKeyManager(keyManager, googleCallback))
}
