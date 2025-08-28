#!/bin/sh

set -exu

ollama serve &

sleep 10

ollama pull qwen3:8b
ollama pull bge-m3:567m

wait
