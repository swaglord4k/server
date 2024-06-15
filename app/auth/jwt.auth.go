package auth

import (
	"fmt"

	"de.amplifonx/app/model"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret-key")

func CreateToken(role model.Role) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"role": role,
		})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func GetRoleFromToken(tokenString string) (*model.Role, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claimsString, err := token.Claims.GetAudience()
	if err != nil {
		return nil, fmt.Errorf("claims not present")
	}
	var bytes []byte
	claimsString.UnmarshalJSON(bytes)
	role := model.Role(string(bytes))
	return &role, nil
}
