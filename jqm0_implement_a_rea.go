package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// AutomationScript represents a single automation script
type AutomationScript struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Script    string `json:"script"`
	Interval  int    `json:"interval"` // in seconds
	LastRun   int64  `json:"last_run"` // unix timestamp
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"` // unix timestamp
	UpdatedAt int64  `json:"updated_at"` // unix timestamp
}

// AutomationScriptDashboard represents the dashboard data
type AutomationScriptDashboard struct {
	Scripts   []AutomationScript `json:"scripts"`
	Error    string             `json:"error"`
	Message  string             `json:"message"`
	UpdatedAt int64            `json:"updated_at"` // unix timestamp
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	scripts = []AutomationScript{
		{ID: 1, Name: "Script 1", Script: "echo 'Hello World!' > output.txt", Interval: 60},
		{ID: 2, Name: "Script 2", Script: "curl -s https://example.com", Interval: 300},
	}
)

func main() {
	r := gin.Default()

	r.GET("/dashboard", dashboardHandler)
	r.GET("/dashboard/ws", dashboardWSHandler)

	log.Fatal(r.Run(":8080"))
}

func dashboardHandler(c *gin.Context) {
	dashboard := AutomationScriptDashboard{Scripts: scripts}
	c.JSON(http.StatusOK, dashboard)
}

func dashboardWSHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		switch msgType {
		case websocket.TextMessage:
			var message struct {
				Action string `json:"action"`
			}
			err := json.Unmarshal(msg, &message)
			if err != nil {
				log.Println(err)
				continue
			}

			switch message.Action {
			case "get_scripts":
				ws.WriteJSON(scripts)
			case "run_script":
				// implement script running logic here
				log.Println("Script running logic not implemented")
			default:
				log.Println("Unknown action:", message.Action)
			}
		default:
			log.Println("Unexpected message type:", msgType)
		}
	}
}