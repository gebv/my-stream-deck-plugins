PLUGIN_UUID ?= com.github.gebv.my-stream-deck-plugins

logs:
	tail -f ~/Library/Logs/StreamDeck/*

build:
	mkdir -p dist
	go build -o ./${PLUGIN_UUID}.sdPlugin/my-sd-plugins-backend.darwin-arm64.bin ./
	rm -rf dist/*
	./DistributionTool -b -i ${PLUGIN_UUID}.sdPlugin -o dist
	open ./dist/com.github.gebv.my-stream-deck-plugins.streamDeckPlugin

# release:
# 	rm -rf "$(HOME)/Library/Application Support/com.elgato.StreamDeck/Plugins/${PLUGIN_UUID}.sdPlugin"

# 	mv dist "$(HOME)/Library/Application Support/com.elgato.StreamDeck/Plugins/${PLUGIN_UUID}.sdPlugin"

all: build;

# open "/Users/gebv/Library/Application Support/com.elgato.StreamDeck/Plugins"
