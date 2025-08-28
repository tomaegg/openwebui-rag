#!/bin/sh

set -exu

ollama serve &

sleep 10

ollama pull qwen3:4b
ollama pull bge-m3:567m

wait
