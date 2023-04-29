package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"io"
	"log"
	"os"
)

var openaiClient *openai.Client

// write init function to initialize openaiClient
func init() {

	openaiKey := os.Getenv("OPENAI_APIKEY")

	if openaiKey == "" {
		log.Fatal("Error: OPENAI_APIKEY environment variable is not set.")
	}

	openaiClient = openai.NewClient(openaiKey)
}

func askGpt4(
	command string,
	data string,
) string {
	// This function is assumed to be implemented already
	// see https://github.com/sashabaranov/go-openai
	resp, err := openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Temperature: 0.1,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a CLI tool. Please process the following data + instruction. Your output should not include any explanations or markdown or such. Only raw cli output, so that it could be used for example inside a cli script or unix pipe (|) chain of operations.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: data,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: command,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		os.Exit(1)
	}
	return resp.Choices[0].Message.Content
}

func main() {

	// Define permission flags
	var fs string
	var network bool
	flag.StringVar(&fs, "fs", "", "File system permissions: ro (read-only) or rw (read-write)")
	flag.BoolVar(&network, "netw", false, "Enable networking permissions")
	flag.Parse()

	// Validate permission flags
	if fs != "" && fs != "ro" && fs != "rw" {
		_, _ = fmt.Fprintln(os.Stderr, "Invalid --fs flag value, must be 'ro' or 'rw'")
		os.Exit(1)
	}

	// Read command line arguments after flag parsing
	args := flag.Args()
	if len(args) == 0 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [--fs=ro|--fs=rw] [--netw] <command>\n", os.Args[0])
		os.Exit(1)
	}

	command := args[0]

	// Read input from stdin
	inputData := ""
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			_, _ = fmt.Fprintln(os.Stderr, "Error reading from stdin:", err)
			os.Exit(2)
		}
		inputData += line
		if err == io.EOF {
			break
		}
	}

	// Process input with GPT-4
	processedData := askGpt4(inputData, command)

	// Write output to stdout
	fmt.Print(processedData)
}
