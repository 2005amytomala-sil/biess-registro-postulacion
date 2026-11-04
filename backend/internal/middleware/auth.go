package middleware

import (
	"strings"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

const (
	postulanteIDKey = "postulante_id"
	cedulaKey       = "cedula"
)

// TokenValidator interface para validar tokens
type TokenValidator interface {
	ValidarToken(token string) (int, string, error)
}

// AuthMiddleware protege rutas con validación JWT
type AuthMiddleware struct {
	tokenValidator TokenValidator
}

func NewAuthMiddleware(tokenValidator TokenValidator) *AuthMiddleware {
	return &AuthMiddleware{
		tokenValidator: tokenValidator,
	}
}

// RequireAuth verifica el token JWT y establece el usuario en el contexto
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extraer token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(domain.ErrNoAutorizado)
			c.Abort()
			return
		}

		// Formato: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(domain.ErrNoAutorizado)
			c.Abort()
			return
		}

		token := parts[1]

		// Validar token
		postulanteID, cedula, err := m.tokenValidator.ValidarToken(token)
		if err != nil {
			c.Error(domain.ErrNoAutorizado)
			c.Abort()
			return
		}

		// Guardar en el contexto
		c.Set(postulanteIDKey, postulanteID)
		c.Set(cedulaKey, cedula)

		c.Next()
	}
}

// PostulanteIDFrom extrae el postulante_id del contexto
func PostulanteIDFrom(c *gin.Context) (int, bool) {
	val, exists := c.Get(postulanteIDKey)
	if !exists {
		return 0, false
	}
	postulanteID, ok := val.(int)
	return postulanteID, ok
}

// CedulaFrom extrae la cédula del contexto
func CedulaFrom(c *gin.Context) (string, bool) {
	val, exists := c.Get(cedulaKey)
	if !exists {
		return "", false
	}
	cedula, ok := val.(string)
	return cedula, ok
}
