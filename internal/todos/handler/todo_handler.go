package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	custom "github.com/mesh-dell/todo-list-API/internal/errors"
	"github.com/mesh-dell/todo-list-API/internal/todos"
	"github.com/mesh-dell/todo-list-API/internal/todos/dtos"
	"github.com/mesh-dell/todo-list-API/internal/todos/service"
	"github.com/mesh-dell/todo-list-API/internal/utils"
)

type TodoHandler struct {
	svc service.TodoService
}

func (h *TodoHandler) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}
	id := userID.(uint)
	var req dtos.TodoItemRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	item, err := h.svc.Create(c.Request.Context(), id, req)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, dtos.TodoItemResponseDto{
		Id:          item.ID,
		Title:       item.Title,
		Description: item.Description,
	})
}

func (h *TodoHandler) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, gin.H{"error": "invalid id"})
		return
	}
	item, err := h.svc.FindByID(c.Request.Context(), uint(id))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "todo item not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, dtos.TodoItemResponseDto{
		Id:          item.ID,
		Title:       item.Title,
		Description: item.Description,
	})
}

func (h *TodoHandler) Delete(c *gin.Context) {
	// get todo id
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, gin.H{"error": "invalid id"})
		return
	}
	// get userId
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDStr.(uint)

	err = h.svc.Delete(c.Request.Context(), uint(id), userID)
	if err != nil {
		if errors.Is(err, custom.ErrItemNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if errors.Is(err, custom.ErrCannotDeleteItem) {
			c.IndentedJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusNoContent, gin.H{"error": "item deleted"})
}

func (h *TodoHandler) Update(c *gin.Context) {
	var req dtos.TodoItemRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDStr.(uint)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, gin.H{"error": "invalid id"})
		return
	}
	item, err := h.svc.Update(c.Request.Context(), uint(id), userID, req)
	if err != nil {
		if errors.Is(err, custom.ErrItemNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if errors.Is(err, custom.ErrCannotUpdateItem) {
			c.IndentedJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.IndentedJSON(http.StatusOK, dtos.TodoItemResponseDto{
		Id:          item.ID,
		Title:       item.Title,
		Description: item.Description,
	})

}

// TODO add pagination
func (h *TodoHandler) FindAllForUser(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDStr.(uint)

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	pagination := &utils.Pagination{
		Page:  page,
		Limit: limit,
	}

	todoItems, err := h.svc.FindAllForUser(c.Request.Context(), userID, pagination)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var todoItemsRes []dtos.TodoItemResponseDto
	d := todoItems.Data.(*[]todos.TodoItem)
	for _, item := range *d {
		todoItemsRes = append(todoItemsRes, dtos.TodoItemResponseDto{
			Id:          item.ID,
			Title:       item.Title,
			Description: item.Description,
		})
	}

	todoItems.Data = todoItemsRes
	c.IndentedJSON(http.StatusOK, todoItems)
}
func NewTodoHandler(s service.TodoService) *TodoHandler {
	return &TodoHandler{
		svc: s,
	}
}
