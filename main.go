package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"strconv"
)

var openaiClient *openai.Client
var gptVersionEnum = openai.GPT3Dot5Turbo

// write init function to initialize openaiClient
func init() {

	openaiKey := os.Getenv("OPENAI_APIKEY")
	if openaiKey == "" {
		log.Fatal("Error: OPENAI_APIKEY environment variable is not set.")
	}
	openaiClient = openai.NewClient(openaiKey)

	// init default gpt version from env var GPT_VERSION, either 3 or 4, convert from string to int
	gptVersionStr := os.Getenv("GPT_VERSION")
	if gptVersionStr != "" {
		// convert string gptVersion to int using atoi
		num, err := strconv.Atoi(gptVersionStr)
		if err != nil {
			log.Fatal("Error: GPT_VERSION environment variable is not set to a valid integer.")
		}
		// if num is 3 or 4, set gptVersionEnum to num, use a switch statement
		switch num {
		case 3:
			gptVersionEnum = openai.GPT3Dot5Turbo
		case 4:
			gptVersionEnum = openai.GPT4
		default:
			log.Fatal("Error: GPT_VERSION environment variable is not set to a valid value (must be 3 or 4)")
		}
	}
}

func askGpt(
	command string,
	data string,
	gptVersionOvrd int,
) string {

	if gptVersionOvrd == 4 {
		gptVersionEnum = openai.GPT4
	}

	resp, err := openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       gptVersionEnum,
			Temperature: 0.1,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: "You are a CLI tool. Please process the following data + instruction.\n" +
						"Your output should not include any explanations or markdown or such.\n" +
						"Only raw cli output, so that it could be used for example inside a cli script or unix pipe (|) chain of operations.\n",
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
		log.Fatalf("ChatCompletion error: %v\n", err)
	}
	return resp.Choices[0].Message.Content
}

func readAllStdIn() string {
	// Read all stdin
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(os.Stdin)
	if err != nil {
		log.Fatal("Could not read std in")
	}
	return buf.String()
}

func main() {

	var fs string
	var network bool
	var gptVersion int
	flag.StringVar(&fs, "fs", "", "File system permissions: ro (read-only) or rw (read-write)")
	flag.BoolVar(&network, "netw", false, "Enable networking permissions")
	flag.IntVar(&gptVersion, "gpt", 3, "GPT version, 3 or 4")
	flag.Parse()

	// Validate flags
	if gptVersion != 3 && gptVersion != 4 {
		log.Fatal("Invalid --gpt flag value, must be 3 or 4")
	}
	if fs != "" && fs != "ro" && fs != "rw" {
		log.Fatal("Invalid --fs flag value, must be 'ro' or 'rw'")
	}

	// Read command line arguments after flag parsing
	args := flag.Args()
	if len(args) == 0 {
		log.Fatalf("Usage: %s [--fs=ro|--fs=rw] [--gpt=4|--gpt=3] [--netw] <command>\n", os.Args[0])
	}

	command := args[0]

	// Read input from stdin if not connected to a terminal
	inputData := ""
	if !terminal.IsTerminal(int(os.Stdin.Fd())) {
		inputData = readAllStdIn()
	}

	result := askGpt(inputData, command, gptVersion)
	fmt.Print(result)
}
