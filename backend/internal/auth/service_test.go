package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	t.Run("should retrieve a valid token when request is valid", func(t *testing.T) {
		tokenConfig := TokenConfig{
			SecretKey: "key",
			Issuer:    "test",
		}
		userID := uint(1001)
		service := NewService(tokenConfig)

		_, err := service.GenerateToken(userID)
		assert.NoError(t, err)
	})
}

func TestValidateToken(t *testing.T) {
	t.Run("should not return any error when token is valid", func(t *testing.T) {
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIxMjM0IiwiZXhwIjoxOTE5ODg5Mjk0LCJpYXQiOjE2NDc0MzAwOTQsImlzcyI6InRlc3QifQ.te-HZfMijZQ-z2fMQSd5CUjGW8Yv0iRwtkmOdJPQ_i4"
		tokenConfig := TokenConfig{
			SecretKey: "key",
			Issuer:    "test",
		}
		service := NewService(tokenConfig)

		_, err := service.ValidateToken(token)
		assert.NoError(t, err)
	})
	t.Run("should return an error when token is not valid(generated with different secret)", func(t *testing.T) {
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImpvaG4uZG9lQGV4YW1wbGUuY29tIiwiZXhwIjoxNjQ3Njg5ODM2LCJpYXQiOjE2NDc0MzA2MzYsImlzcyI6InRlc3QifQ.8fM04RvYTBMo-aEX4Ugvbvp5oxJGdhdlYzdauly3eGA"
		tokenConfig := TokenConfig{
			SecretKey: "key",
			Issuer:    "test",
		}
		service := NewService(tokenConfig)

		_, err := service.ValidateToken(token)
		assert.Error(t, err)
	})
	t.Run("should return an error when token is not valid(uses different signing method)", func(t *testing.T) {
		token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImpvaG4uZG9lQGV4YW1wbGUuY29tIiwiZXhwIjoxNjQ3Njg5OTEwLCJpYXQiOjE2NDc0MzA3MTAsImlzcyI6InRlc3QifQ.ScpcU6JToMMA1ZoiU8GezzdUA2rpDjvh-lEImRoXCrMt1tZ3hh7itszY8oKF6QO-yWVN3zCYyZY2tX0wL3ykqKV4QHH_ZKcdjbDD5bgSmYPbv06txX_Df655tRx0mdFaOFgpvYfC2a6zJvcpcKXGbKMgKmXbbANVTKYQQGkxPR6ITx2-Pyuu2LE3Mg6A-pjHcvLjK3rofoxwymlgoQ9EhDxs3sMVJl0RBoIPwsF1qjTvcnUDF_YBlCkgAYZBTARbRuYq6cMIYXfJvDuNYih2s12a_hI5Gmay4y_8TSgs91wo1GCe-yeebVxdR--Kql-SGZuoJWeZCPE1j3SVfYa8dg"
		tokenConfig := TokenConfig{
			SecretKey: "key",
			Issuer:    "test",
		}
		service := NewService(tokenConfig)

		_, err := service.ValidateToken(token)
		assert.Error(t, err)
	})
}
