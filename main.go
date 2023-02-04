package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/valyala/fastjson"
	"meow.tf/streamdeck/sdk"
)

func main() {
	f, errf := os.OpenFile("./plugin-backend.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errf != nil {
		log.Println(errf)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(time.Now().Format(time.RFC3339Nano))

	// Initialize handlers for events
	sdk.RegisterAction("com.github.gebv.my-stream-deck-plugins.dosomething1", doSomethingHandler)
	sdk.RegisterAction("com.github.gebv.my-stream-deck-plugins.toggle-on-off", doSomethingHandler)

	sdk.AddHandler(func(e *sdk.SendToPluginEvent) {
		log.Println("PI send to Plugin Event", e.Action, e.Payload.String())

		time.Sleep(time.Second)
		sdk.SendToPropertyInspector(e.Context, map[string]string{"result": "123"})
	})

	log.Println("input args", os.Args)

	// Open and connect the SDK
	err := sdk.Open()

	log.Println("plugin UUID", sdk.PluginUUID)

	if err != nil {
		log.Println("Open failed:", err)
		log.Fatalln(err)
	}
	log.Println("Open OK")

	go pool()
	go func() {
		log.Println("send log message")
		sdk.Log("Polling mic state..." + sdk.PluginUUID)
		log.Println("success sent log message")
	}()

	log.Println("Wait")
	// Wait until the socket is closed, or SIGTERM/SIGINT is received
	sdk.Wait()
	log.Println("Bye Bye")
}

var i = 0

func doSomethingHandler(action, context string, payload *fastjson.Value, deviceId string) {
	// Do something as a result of an action (keyDown)
	log.Println(">>hande:", map[string]interface{}{
		"action":   action,
		"deviceId": deviceId,
		"context":  context,
		"payload":  payload.String(),
	})

	i++
	sdk.SetTitle(context, fmt.Sprintf("title %d", i), 0)
	// sdk.SetSettings(context, map[string]any{
	// 	"context_settings": context,
	// 	"abc":              fmt.Sprintf("title %d", i),
	// 	action:             fmt.Sprintf("title %d", i),
	// })
	// sdk.SetGlobalSettings(map[string]any{
	// 	"global": true,
	// 	"abc":    fmt.Sprintf("title %d", i),
	// 	action:   fmt.Sprintf("title %d", i),
	// })
	// sdk.Log("test message" + fmt.Sprintf("title %d", i))
	// sdk.SetState()

}

func drawPng(context string, img image.Image) {
	var buff bytes.Buffer
	_ = png.Encode(&buff, img)

	str := base64.StdEncoding.EncodeToString(buff.Bytes())
	sdk.SetImage(
		context,
		"data:image/png;base64,"+str,
		0,
	)
}

func pool() {
	for {
		// sdk.Log("Polling mic state..." + sdk.PluginUUID)
		log.Println("Pooling")

		time.Sleep(time.Second * 60)
	}
}
