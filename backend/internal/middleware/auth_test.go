package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// MockTokenValidator para testing
type MockTokenValidator struct {
	validToken      bool
	postulanteID    int
	cedula          string
	shouldReturnErr bool
}

func (m *MockTokenValidator) ValidarToken(token string) (int, string, error) {
	if m.shouldReturnErr {
		return 0, "", domain.ErrTokenInvalido
	}
	if !m.validToken {
		return 0, "", domain.ErrTokenInvalido
	}
	return m.postulanteID, m.cedula, nil
}

func TestRequireAuth_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ErrorHandler())

	mockValidator := &MockTokenValidator{validToken: false}
	authMiddleware := NewAuthMiddleware(mockValidator)

	router.GET("/protected", authMiddleware.RequireAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestRequireAuth_InvalidTokenFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ErrorHandler())

	mockValidator := &MockTokenValidator{validToken: false}
	authMiddleware := NewAuthMiddleware(mockValidator)

	router.GET("/protected", authMiddleware.RequireAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestRequireAuth_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ErrorHandler())

	mockValidator := &MockTokenValidator{validToken: false}
	authMiddleware := NewAuthMiddleware(mockValidator)

	router.GET("/protected", authMiddleware.RequireAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestRequireAuth_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockValidator := &MockTokenValidator{
		validToken:   true,
		postulanteID: 123,
		cedula:       "0101010101",
	}
	authMiddleware := NewAuthMiddleware(mockValidator)

	router.GET("/protected", authMiddleware.RequireAuth(), func(c *gin.Context) {
		postulanteID, ok := PostulanteIDFrom(c)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no postulante_id"})
			return
		}
		cedula, ok := CedulaFrom(c)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no cedula"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"postulante_id": postulanteID,
			"cedula":        cedula,
		})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if body != `{"cedula":"0101010101","postulante_id":123}` {
		t.Errorf("Unexpected body: %s", body)
	}
}

func TestPostulanteIDFrom_NotSet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	postulanteID, ok := PostulanteIDFrom(c)
	if ok {
		t.Error("Expected ok to be false")
	}
	if postulanteID != 0 {
		t.Errorf("Expected postulanteID to be 0, got %d", postulanteID)
	}
}

func TestCedulaFrom_NotSet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	cedula, ok := CedulaFrom(c)
	if ok {
		t.Error("Expected ok to be false")
	}
	if cedula != "" {
		t.Errorf("Expected cedula to be empty, got %s", cedula)
	}
}
