package posthandler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func (h *PostsHandler) authMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Получаем токен из заголовка
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// 2. Отправляем запрос к users_service
			userInfo, err := fetchUserInfo(r.Context(), h.usersUrl, token, h.usersClient)
			if err != nil {
				handleAuthError(w, err)
				return
			}

			// 3. Добавляем userInfo в контекст
			newCtx := context.WithValue(r.Context(), "userInfo", userInfo)
			log.Println("add: ", userInfo)

			// 4. Продолжаем цепочку middleware
			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}

// fetchUserInfo делает запрос к users_service
func fetchUserInfo(ctx context.Context, baseURL, token string, client *http.Client) (*userInfo, error) {
	// 1. Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/whoami", nil)
	if err != nil {
		return nil, err
	}

	// 2. Устанавливаем заголовки
	req.Header.Set("Authorization", token)
	req.Header.Set("Accept", "application/json")

	// 3. Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 4. Обрабатываем ответ
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("user service returned non-200 status")
	}

	// 5. Парсим ответ
	var userInfo userInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// handleAuthError обрабатывает ошибки аутентификации
func handleAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		http.Error(w, "User service timeout", http.StatusGatewayTimeout)
	case err.Error() == "user service returned non-200 status":
		http.Error(w, "Invalid token", http.StatusUnauthorized)
	default:
		http.Error(w, "Authentication failed", http.StatusInternalServerError)
	}
}
