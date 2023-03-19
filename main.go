package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	gpt3 "github.com/PullRequestInc/go-gpt3"
)

func main() {
	apiKey := getApiKey()
	data, err := readfile("./data.csv")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("APIkey is :" + apiKey)
	response, err := callChatCompletionEndpoint(apiKey, data)
	if err != nil {
		log.Fatal(err)
	}
	// Combine the responses into a single string
	fmt.Println(response.Choices[0].Message.Content)

}

func readfile(filepath string) (string, error) {
	// Open the file at the given file path
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the contents of the file into a byte slice
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	// Convert the byte slice to a string and return it
	return string(bytes), nil
}

func getApiKey() string {
	return os.Getenv("OPENAI_API_KEY")
}

func callChatCompletionEndpoint(apiKey string, data string) (*gpt3.ChatCompletionResponse, error) {
	systemText := "You are taking rows in a spreadsheet NHL sports data document. Generate prompt/completion lines with questions based on the data, that can be used to train a data model. Each line should be a json dictionary. Here is an example with unrelated data {\"prompt\":\"Item is a handbag. Colour is army green. Price is midrange. Size is small.->\", \"completion\":\" This stylish small green handbag will add a unique touch to your look, without costing you a fortune.\"}"

	ctx := context.Background()

	client := gpt3.NewClient(apiKey)
	lines := strings.Split(data, "\n")

	first20 := strings.Join(lines[:20], "\n")
	data = first20

	chatResp, err := client.ChatCompletion(ctx, gpt3.ChatCompletionRequest{
		Model: gpt3.GPT3Dot5Turbo0301,
		Messages: []gpt3.ChatCompletionRequestMessage{
			{
				Role:    "system",
				Content: systemText,
			},
			{
				Role:    "user",
				Content: data,
			},
		},
	})
	return chatResp, err

}
