# gptsh
Some silly saturday morning gpt shell stuff. Written by GPT4 + github copilot. Uses GPT3.5-turbo by default.

## Usage
```ll . | gptsh "reverse line order" | gptsh "output as prettified json with some useful field names"```

## Installation
```go install github.com/gigurra/gptsh```

## Environment variables
* OPENAI_APIKEY -  Required. Your OpenAI API key.
* GPT_VERSION - Optional. The GPT model to use. Defaults to "3". Optionally you can set it to "4".