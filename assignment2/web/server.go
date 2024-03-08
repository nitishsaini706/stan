package web

import (
	"net/http"
    "github.com/gin-gonic/gin"
    "github.com/nitishsaini706/stan/assignment2/store"
	"github.com/nitishsaini706/stan/assignment2/models"
    "strconv"
)

func SetupRouter(s *store.Store) *gin.Engine {
    r := gin.Default()

    // Create user
    r.POST("/users", func(c *gin.Context) {
        var user models.User
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        err := s.CreateUser(user) // Modified to handle error from CreateUser
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, user)
    })

    // GET handler to read a user
    r.GET("/users/:id", func(c *gin.Context) {
        id, err := strconv.Atoi(c.Param("id"))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
            return
        }

        user, err := s.GetUser(id) // Modified to match the corrected function name
        if err == store.ErrNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        } else if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, user)
    })

    // PUT handler to update a user
    r.PUT("/users/:id", func(c *gin.Context) {
        id, err := strconv.Atoi(c.Param("id"))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
            return
        }

        var updateUser models.User
        if err := c.ShouldBindJSON(&updateUser); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        err = s.UpdateUser(id, updateUser) // Assuming UpdateUser returns an error to handle
        if err == store.ErrNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        } else if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, updateUser)
    })

    // DELETE handler to delete a user
    r.DELETE("/users/:id", func(c *gin.Context) {
        id, err := strconv.Atoi(c.Param("id"))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
            return
        }

        err = s.DeleteUser(id) // Assuming DeleteUser returns an error to handle
        if err == store.ErrNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        } else if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
    })

    return r
}
