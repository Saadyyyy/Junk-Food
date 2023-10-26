package controller

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
)

type MakanRespon struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

type MakanUsecase interface {
	RecommendMakan(userInput, openAIKey string) (string, error)
}

type makan struct{}

func NewMakanUsecase() MakanUsecase {
	return &makan{}
}

func (uc *makan) RecommendMakan(userInput, openAIKey string) (string, error) {
	ctx := context.Background()
	client := openai.NewClient(openAIKey)
	model := openai.GPT3Dot5Turbo
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "Hai, Saya akan membantu anda untuk rekomendasi makanan ",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: userInput,
		},
	}

	resp, err := uc.getCompletionFromMessages(ctx, client, messages, model)
	if err != nil {
		return "", err
	}
	answer := resp.Choices[0].Message.Content
	return answer, nil
}

func (uc *makan) getCompletionFromMessages(
	ctx context.Context,
	client *openai.Client,
	messages []openai.ChatCompletionMessage,
	model string,
) (openai.ChatCompletionResponse, error) {
	if model == "" {
		model = openai.GPT3Dot5Turbo
	}

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    model,
			Messages: messages,
		},
	)
	return resp, err
}

func RecommendMakan(c echo.Context, makan MakanUsecase) error {
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Authorization token is missing"})
	}

	// Memeriksa apakah header Authorization mengandung token Bearer
	if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Invalid token format. Use 'Bearer [token]'"})
	}

	// Ekstrak token dari header
	tokenString = tokenString[7:]

	var requestData map[string]interface{}
	err := c.Bind(&requestData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": "Invalid JSON format"})
	}

	userInput, ok := requestData["message"].(string)
	if !ok || userInput == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": "Invalid or missing 'message' in the request"})
	}

	userInput = fmt.Sprintf("Rekomendasi Makanan: %s", userInput)

	answer, err := makan.RecommendMakan(userInput, os.Getenv("OPENAI_KEY"))
	if err != nil {
		errorMessage := "Failed to generate hotel recommendations"
		if strings.Contains(err.Error(), "rate limits exceeded") {
			errorMessage = "Rate limits exceeded. Please try again later."
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": errorMessage})
	}

	responseData := MakanRespon{
		Status: "success",
		Data:   answer,
	}

	return c.JSON(http.StatusOK, responseData)
}
