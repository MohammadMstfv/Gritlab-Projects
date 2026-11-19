# Tool Prompts

    This quest teaches the prompting side of it: getting clean JSON, getting the model to pick the right tool, and chaining a request into an answer. Everything runs in a browser by having the model produce and read JSON. No installation required.

## Why structured output matters for tool use:

 Structured JSON output from the model guarantees that tool names and arguments can be parsed deterministically by code, eliminating ambiguity and enabling safe, automated execution of external functions.

## What your code (not the model) is responsible for:

    Your orchestration code must validate the model's JSON, dispatch the correct tool, handle errors or malformed responses, and feed clean results back so the model can produce a final natural-language answer.