package token

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	pb "auth_serice/genproto"
	"github.com/form3tech-oss/jwt-go"
)

var AccessTokenKey = "salom"
var RefreshTokenKey = "salom"

func GenerateJwtToken(request *pb.Claims) (*pb.LoginResponse, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	accessClaims["user_name"] = request.UserName
	accessClaims["email"] = request.Email
	accessClaims["password"] = request.PasswordHash
	accessClaims["user_type"] = request.UserType
	accessClaims["iat"] = time.Now().Unix()
	accessClaims["exp"] = time.Now().Add(1 * time.Hour).Unix()

	access, err := accessToken.SignedString([]byte(AccessTokenKey))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["user_name"] = request.UserName
	refreshClaims["email"] = request.Email
	refreshClaims["password"] = request.PasswordHash
	refreshClaims["user_type"] = request.UserType
	refreshClaims["iat"] = time.Now().Unix()
	refreshClaims["exp"] = time.Now().Add(24 * time.Hour).Unix()

	refresh, err := refreshToken.SignedString([]byte(RefreshTokenKey))
	if err != nil {
		return nil, err
	}
	refreshExp := refreshClaims["exp"].(float64)
	refreshIat := refreshClaims["iat"].(float64)
	expired := int(refreshExp - refreshIat)
	return &pb.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    strconv.Itoa(expired),
	}, nil
}

func UpdateToken(request *pb.UpdateTokenRequest) (*pb.LoginResponse, error) {
	claims, err := ExtractClaim(request.RefreshToken)
	if err != nil {
		return nil, err
	}

	userName, ok := claims["user_name"].(string)
	if !ok {
		return nil, errors.New("user_name not found in token claims")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("email not found in token claims")
	}

	passwordHash, ok := claims["password"].(string)
	if !ok {
		return nil, errors.New("password not found in token claims")
	}

	userType, ok := claims["user_type"].(string)
	if !ok {
		return nil, errors.New("user_type not found in token claims")
	}
	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	accessClaims["user_name"] = userName
	accessClaims["email"] = email
	accessClaims["password"] = passwordHash
	accessClaims["user_type"] = userType
	accessClaims["iat"] = time.Now().Unix()
	accessClaims["exp"] = time.Now().Add(1 * time.Hour).Unix()

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["user_name"] = userName
	refreshClaims["email"] = email
	refreshClaims["password"] = passwordHash
	refreshClaims["user_type"] = userType
	refreshClaims["iat"] = time.Now().Unix()
	refreshClaims["exp"] = time.Now().Add(24 * time.Hour).Unix()

	refreshExp := claims["exp"].(float64)
	refreshIat := claims["iat"].(float64)
	expired := int(refreshExp - refreshIat)

	return &pb.LoginResponse{
		AccessToken:  AccessTokenKey,
		RefreshToken: RefreshTokenKey,
		ExpiresIn:    strconv.Itoa(expired),
	}, nil
}

func ExtractClaim(tokenStr string) (jwt.MapClaims, error) {
	var (
		token *jwt.Token
		err   error
	)

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(AccessTokenKey), nil
	}

	token, err = jwt.Parse(tokenStr, keyFunc)
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
