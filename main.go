package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/valyala/fastjson"
	"meow.tf/streamdeck/sdk"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
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

	defer func() {
		if rec := recover(); rec != nil {
			log.Println("Recover panic:\n", rec)
			return
		}
	}()

	// Initialize handlers for events
	sdk.RegisterAction("com.github.gebv.my-stream-deck-plugins.dosomething1", doSomethingHandler)
	sdk.RegisterAction("com.github.gebv.my-stream-deck-plugins.toggle-on-off", doSomethingHandler)
	// sdk.RegisterAction("com.github.gebv.my-stream-deck-plugins.mem-info", memInfoHandler)

	sdk.AddHandler(func(e *sdk.WillAppearEvent) {
		log.Println("Active element:", e.Action, e.Context, e.Payload)
		registredActionsMux.Lock()
		registredActions[e.Context] = action{
			context:      e.Context,
			action:       e.Action,
			selectedSkin: string(e.Payload.Get("settings").GetStringBytes("selectedSkin")),
		}
		registredActionsMux.Unlock()
	})
	sdk.AddHandler(func(e *sdk.ReceiveSettingsEvent) {
		log.Println("Got Settings:", e.Action, e.Context, e.Settings)
		registredActionsMux.Lock()
		if action, exists := registredActions[e.Context]; exists {
			action.selectedSkin = string(e.Settings.GetStringBytes("selectedSkin"))
			registredActions[e.Context] = action
		}
		registredActionsMux.Unlock()

	})
	sdk.AddHandler(func(e *sdk.GlobalSettingsEvent) {
		log.Println("Got Global Settings:", e.Settings)
	})
	sdk.AddHandler(func(e *sdk.WillDisappearEvent) {
		log.Println("Hide element:", e.Action, e.Context, e.Payload)

		registredActionsMux.Lock()
		delete(registredActions, e.Context)
		registredActionsMux.Unlock()
	})
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
		log.Println()

		time.Sleep(time.Second * 1)

		registredActionsMux.RLock()
		for context := range registredActions {
			item := registredActions[context]

			log.Printf("Registred Action %q, context %q\n", item.action, item.context)
			switch item.action {
			case "com.github.gebv.my-stream-deck-plugins.mem-info":
				osCPU, _ := cpu.Percent(0, false)
				memInfo, _ := mem.VirtualMemory()
				skin := item.selectedSkin

				log.Println("  Skin:", skin)

				switch skin {
				case "cpu_usage_percent":
					if len(osCPU) > 0 {
						sdk.SetTitle(item.context, fmt.Sprintf("CPU\n%d%%", int(osCPU[0])), 0)
					}
				case "mem_total":
					sdk.SetTitle(item.context, fmt.Sprintf("MEM\nTotal\n%s", humanize.Bytes(memInfo.Total)), 0)
				case "mem_free":
					sdk.SetTitle(item.context, fmt.Sprintf("MEM\nFree\n%s", humanize.Bytes(memInfo.Available)), 0)
				case "mem_usage_percent":
					sdk.SetTitle(item.context, fmt.Sprintf("MEM\nUsage\n%d%%", int(memInfo.UsedPercent)), 0)
				default:
					sdk.SetTitle(item.context, fmt.Sprintf("MEM\nUsage\n%d%%", int(memInfo.UsedPercent)), 0)
				}
			}
		}
		registredActionsMux.RUnlock()
	}
}

type action struct {
	context      string
	action       string
	selectedSkin string
}

var registredActions = map[string]action{}
var registredActionsMux sync.RWMutex

func info() {
	cpuInfo, _ := cpu.Info()
	v, _ := mem.VirtualMemory()

	// almost every return value is a struct
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
	fmt.Printf("CPU: %v\n", cpuInfo)
}
