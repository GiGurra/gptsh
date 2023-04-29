# gptsh
Some silly saturday morning gpt shell stuff. Written by GPT4 + github copilot. Uses GPT3.5-turbo by default.
The results are pretty bad with gpt3.x, so I definitely recommend using gpt4. 

Unfortunately I don't have a gpt4 API key (the waitlist seems long), but I've verified using chatgpt with GPT that gpt4 produces much better results.

## Usage
```ll . | gptsh "reverse line order" | gptsh "output as prettified json with some useful field names"```

## Installation
```go install github.com/gigurra/gptsh```

## Environment variables
* OPENAI_APIKEY -  Required. Your OpenAI API key.
* GPT_VERSION - Optional. The GPT model to use. Defaults to "3". Optionally you can set it to "4".