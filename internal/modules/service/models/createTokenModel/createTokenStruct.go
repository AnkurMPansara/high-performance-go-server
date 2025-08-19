package createTokenModel

import "time"

type ApiInput struct {
	UserId              int    `json:"user_id"`
	ExpirationTime      string `json:"expiration_time"`
	ValidationStartTime string `json:"validation_start_time"`
}

type ApiData struct {
	StartTime           time.Time
	UserId              int
	ExpirationTime      time.Time
	ValidationStartTime time.Time
	SessionId           string
	GeneratedToken      string
	Error               string
	Code                int
}

type ApiResponse struct {
	Code    int                `json:"code"`
	Status  string             `json:"status"`
	Respose CreateTokenRespose `json:"response"`
	Error   string             `json:"error"`
}

type CreateTokenRespose struct {
	Token  string `json:"token"`
	UserId int    `json:"user_id"`
}

type JwtTokenHeader struct {
	Algorithm string `json:"alg"`
	Type string `json:"typ"`
}

type JwtTokenPayLoad struct {
	Issuer         string `json:"iss"`
	Subject        string `json:"sub"`
	Audience       string `json:"aud"`
	ExpirationTime int    `json:"exp"`
	NotBeforeTime  int    `json:"nbf"`
	IssuedAt       int    `json:"iat"`
	JwtId          string `json:"jti"`
}
