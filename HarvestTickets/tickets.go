package main

import (
	"os"
	"fmt"
	"math"
	"time"

	"github.com/adlio/harvest"
)

func main() {

	client := harvest.NewTokenAPI(os.Getenv("HarvestAccountID"), os.Getenv("HarvestAccessToken"))

	clients, err := client.GetClients(harvest.Defaults())
	if err != nil {
		println(err)
		return
	}

	projects, err := client.GetProjects(harvest.Defaults())
	if err != nil {
		println(err)
		return
	}

	for _, cl := range clients {

		fmt.Println("____________")
		fmt.Printf("[%s]\n", cl.Name)
		fmt.Println("____________")

		for _, project := range projects {

			if project.Client.Name == cl.Name && project.IsActive {
				fmt.Printf("[%d] %s\n", project.ID, project.Name)

				timeEntries, _ := client.GetTimeEntriesForProjectBetween(project.ID,
					time.Date(project.HintEarliestRecordAt.Year(),
						project.HintEarliestRecordAt.Month(),
						project.HintEarliestRecordAt.Day(), 0, 0, 0, 0, time.Local),
					time.Now(), harvest.Defaults())

				total := time.Duration(0)
				var startDate harvest.Date

				lastEntryDuration := time.Duration(0)

				for index, entry := range timeEntries {
					start, _ := time.Parse("15:04", entry.StartedTime)
					end, _ := time.Parse("15:04", entry.EndedTime)
					duration := end.Sub(start)

					if index == len(timeEntries)-1 {
						startDate = entry.SpentDate
					}

					fmt.Printf("    %s, %.0fh:%02.0fm (%02.2f ore), %s\n",
						entry.SpentDate.Format("02-01-2006"),
						math.Floor(duration.Hours()),
						math.Mod(duration.Minutes(), 60),
						duration.Minutes()/60,
						entry.Notes)

					if index == 0 {
						lastEntryDuration = duration
					}

					total = total + duration
				}

				budget := *project.CostBudget
				fmt.Printf("\nUltimo intervento: %.0fh:%02.0fm (%02.2f ore)\n",
					math.Floor(lastEntryDuration.Hours()),
					math.Mod(lastEntryDuration.Minutes(), 60),
					lastEntryDuration.Minutes()/60)

				fmt.Println("=========================")

				fmt.Printf("Piano di supporto: %2.0f ore\n", budget/110)
				elapsedMonths := time.Now().Sub(startDate.Time).Hours() / 744
				elapsedDays := int(math.Mod(time.Now().Sub(startDate.Time).Hours(), 31))

				if elapsedMonths > 0 {
					fmt.Printf("Data inizio: %s (circa %2.0f mesi fa)\n", startDate.Format("02-01-2006"), elapsedMonths)
				} else {
					fmt.Printf("Data inizio: %s (circa %d giorni fa)\n", startDate.Format("02-01-2006"), elapsedDays)
				}
				fmt.Printf("Totale ore utilizate: %2.0fh:%2.0fm\n\n", math.Floor(total.Hours()), math.Mod(total.Minutes(), 60))

				spent := total.Hours() * 110.0
				leftPercent := (1 - spent/budget)
				budgetHours := budget / 110.0
				budgetLeft := budgetHours * leftPercent * 60
				hoursLeft := int(budgetLeft) / 60
				minutesLeft := math.Mod(budgetLeft, 60)
				if budgetLeft > 0 {
					fmt.Printf("Rimangono: %d ore e %2.0f minuti\n", hoursLeft, minutesLeft)
				} else if budgetLeft < 0 {
					fmt.Printf("Piano di supporto esaurito!\n%d ore e %2.0f minuti in eccesso!!\n", -hoursLeft, -minutesLeft)
				} else {
					fmt.Printf("Piano di supporto esaurito!")
				}

				fmt.Println("=========================")
			}
		}
	}
}
