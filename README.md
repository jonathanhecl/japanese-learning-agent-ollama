# 🏯 Japanese Learning Agent with Ollama
 
 Tiny CLI to learn Japanese using a local Ollama model (Mistral).

 ![Screenshot](screenshot.png)
 
 ## Requirements
 
 - Ollama installed and running: https://ollama.com/
 - Mistral model available in Ollama.
 - Windows, Linux, or macOS.
 
 ## Quick start

```sh
# Ensure the model is available
ollama pull mistral

# Run the app
japanese-learning-agent-ollama
```

### Options

- Environment variable `OLLAMA_MODEL` to select a different model at runtime (defaults to `mistral:latest`). For example:

```sh
export OLLAMA_MODEL=llama3.1
japanese-learning-agent-ollama
```

- Graceful exit: press `Ctrl+C` or type `exit`/`quit`.
 
 ## Download

Grab the latest binaries from [Releases](https://github.com/jonathanhecl/japanese-learning-agent-ollama/releases).
 
 ## 📝 License
  
 Apache License
