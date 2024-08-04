package main

import (
	"context"
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/kapparina/warframe-tui/cmd/tui"
	wapi "github.com/kapparina/warframe-tui/pkg/warframestatapi"
	"os"
	"time"
)

const (
	apiUrl string = "https://api.warframestat.us"
)

func main() {
	var apiClient wapi.ClientWithResponsesInterface
	var worldState *wapi.Ws

	apiClient, err := wapi.NewClientWithResponses(apiUrl)
	worldState, err = getData(apiClient)
	if err != nil {
		log.Fatal("Initialisation error", "error", err)
	}
	tabNames, data := buildStrings(worldState)
	m := tui.ViewportModel{Tabs: tabNames, TabContent: data}
	if _, err = tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func getData(client wapi.ClientWithResponsesInterface) (*wapi.Ws, error) {
	resp, err := client.GetWorldstateByPlatformWithResponse(
		context.TODO(),
		wapi.Pc,
		&wapi.GetWorldstateByPlatformParams{
			Language: wapi.En,
		},
	)
	if resp.JSON200 == nil || err != nil {
		return nil, errors.Join(errors.New("potential nil Worldstate value"), err)
	} else {
		return resp.JSON200, nil
	}
}

func buildStrings(worldState *wapi.Ws) ([]string, []string) {
	var newsStatusString, eventsString, cambionCycleString string

	for i, newsEntry := range worldState.News {
		entryDate, parseErr := time.Parse(time.RFC3339, newsEntry.Date)
		if parseErr == nil {
			newsEntry.Date = entryDate.Format(time.RFC1123)
		}
		switch i {
		case 0:
			newsStatusString += fmt.Sprintf("Date: %s\tMessage: %s\n", newsEntry.Date, newsEntry.Message)
		default:
			newsStatusString += "\n" + fmt.Sprintf("Date: %s\tMessage: %s\n", newsEntry.Date, newsEntry.Message)
		}
	}
	for i, event := range worldState.Events {
		switch i {
		case 0:
			eventsString += fmt.Sprintf(
				"Expiry: %s\tEvent: %s\n",
				event.Expiry.Format(time.RFC1123),
				*event.Description,
			)
		default:
			eventsString += "\n" + fmt.Sprintf(
				"Expiry: %s\tEvent: %s\n",
				event.Expiry.Format(time.RFC1123),
				*event.Description,
			)
		}
	}
	cambionCycleString = fmt.Sprintf(
		"State: %s\tTime Left: %s\tExpiry: %s",
		worldState.CambionCycle.State,
		*worldState.CambionCycle.TimeLeft,
		worldState.CambionCycle.Expiry,
	)
	return []string{"News", "Events", "Cambion Cycle"}, []string{newsStatusString, eventsString, cambionCycleString}
}
