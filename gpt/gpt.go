package gpt

import (
	"bytes"
	"deouy/wechatbot/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const BASEURL = "https://api.openai.com/v1/"

type OpenAIRequest struct {
	Model            string    `json:"model"`
	Messages         []Message `json:"messages"`
	Temperature      float64   `json:"temperature,omitempty"`
	TopP             float64   `json:"top_p,omitempty"`
	FrequencyPenalty float64   `json:"frequency_penalty,omitempty"`
	PresencePenalty  float64   `json:"presence_penalty,omitempty"`
	MaxTokens        int       `json:"max_tokens,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Created int64            `json:"created"`
	Choices []ResponseChoice `json:"choices"`
	Usage   UsageInfo        `json:"usage"`
}

type ResponseChoice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type UsageInfo struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func Completions(messages []Message) (string, error) {
	req := OpenAIRequest{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
		//Temperature: 0.5,
		//MaxTokens:   4096,
	}

	client := &http.Client{}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	log.Printf("request gtp json string : %v", messages)

	request, err := http.NewRequest("POST", BASEURL+"chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.LoadConfig().ApiKey))

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var resp OpenAIResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", err
	}
	var reply string
	if len(resp.Choices) > 0 {
		reply = resp.Choices[0].Message.Content
	}

	log.Printf("gpt response text: %s \n", reply)
	return reply, nil
}

func CreateMessage(content string) Message {
	return Message{
		Role:    "user",
		Content: content,
	}
}
