package store

import (
    "testing"
    "../models"
    "github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
    store := New()
    user := models.User{ID: 1, Name: "John Doe", Email: "john@example.com"}
    store.CreateUser(user)
    assert.Equal(t, store.users[1], user)
}
