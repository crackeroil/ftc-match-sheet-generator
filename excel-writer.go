package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/xuri/excelize/v2"
)

func GenerateMatchsheets(season string, eventCode string, apiKey string) {
	client := &http.Client{}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	teamsResp, err := GetTeams(client, season, eventCode, 1, apiKey)
	if err != nil {
		panic(err)
	}

	teams, err := ParseTeamsResponse(teamsResp)
	if err != nil {
		panic(err)
	}

	for teams.PageCurrent != teams.PageTotal {
		teamsResp, err = GetTeams(client, season, eventCode, teams.PageCurrent+1, apiKey)
		if err != nil {
			panic(err)
		}

		teams2, err := ParseTeamsResponse(teamsResp)
		if err != nil {
			panic(err)
		}

		teams.Teams = append(teams.Teams, teams2.Teams...)
		teams.PageCurrent = teams2.PageCurrent
	}

	for _, team := range teams.Teams {
		// Set up sheet for every team
		teamName := team.NameShort
		if len(teamName) > 31 {
			teamName = teamName[:31]
		}
		if strings.ContainsAny(teamName, ":\\/?*[]") {
			teamName = strings.ReplaceAll(teamName, ":", "")
			teamName = strings.ReplaceAll(teamName, "\\", "")
			teamName = strings.ReplaceAll(teamName, "/", "")
			teamName = strings.ReplaceAll(teamName, "?", "")
			teamName = strings.ReplaceAll(teamName, "*", "")
			teamName = strings.ReplaceAll(teamName, "[", "")
			teamName = strings.ReplaceAll(teamName, "]", "")
		}

		_, err := f.NewSheet(teamName)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Set up title
		f.MergeCell(teamName, "B2", "E2")
		f.SetCellValue(teamName, "B2", fmt.Sprintf("Team #%d | %s", team.TeamNumber, team.NameShort))

		style, err := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Size: 24,
			},
			Border: []excelize.Border{
				{Type: "left", Style: 2, Color: "000000"},
				{Type: "top", Style: 2, Color: "000000"},
				{Type: "bottom", Style: 2, Color: "000000"},
				{Type: "right", Style: 2, Color: "000000"},
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellStyle(teamName, "B2", "E2", style)
		f.SetRowHeight(teamName, 2, 36.75)
		f.SetColWidth(teamName, "B", "E", 30)

		// Set up column headers
		f.SetCellValue(teamName, "B3", "Match")
		f.SetCellValue(teamName, "C3", "Partner")
		f.SetCellValue(teamName, "D3", "Opponent 1")
		f.SetCellValue(teamName, "E3", "Opponent 2")

		style, err = f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Size: 14,
			},
			Border: []excelize.Border{
				{Type: "left", Style: 1, Color: "000000"},
				{Type: "top", Style: 1, Color: "000000"},
				{Type: "bottom", Style: 1, Color: "000000"},
				{Type: "right", Style: 1, Color: "000000"},
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetCellStyle(teamName, "B3", "E3", style)

		// Get schedule
		scheduleResp, err := GetTeamSchedule(client, season, eventCode, fmt.Sprint(team.TeamNumber), apiKey)
		if err != nil {
			fmt.Println(err)
			return
		}

		schedule, err := ParseScheduleResponse(scheduleResp)
		if err != nil {
			fmt.Println(err)
			return
		}

		for i, match := range schedule.Schedule {
			f.SetCellValue(teamName, fmt.Sprintf("B%d", 2*i+4), fmt.Sprintf("Qualification %d", match.MatchNumber))

			// Get partner and opponents numbers
			var friendlyStation string
			var partnerNum int
			var opponent1Num int
			var opponent2Num int

			for _, matchTeam := range match.Teams {
				if matchTeam.TeamNumber == team.TeamNumber {
					friendlyStation = matchTeam.Station[:len(matchTeam.Station)-1]
				}
			}

			f.SetCellValue(teamName, fmt.Sprintf("B%d", 2*i+5), friendlyStation)

			for _, matchTeam := range match.Teams {
				if matchTeam.TeamNumber == team.TeamNumber {
					continue
				}

				if matchTeam.Station[:len(matchTeam.Station)-1] == friendlyStation {
					partnerNum = matchTeam.TeamNumber
				} else {
					if opponent1Num == 0 {
						opponent1Num = matchTeam.TeamNumber
					} else {
						opponent2Num = matchTeam.TeamNumber
					}
				}
			}

			// Find partner and opponents names based on numbers
			var partnerName string
			var opponent1Name string
			var opponent2Name string

			for _, team := range teams.Teams {
				if partnerNum == team.TeamNumber {
					partnerName = team.NameShort
				} else if opponent1Num == team.TeamNumber {
					opponent1Name = team.NameShort
				} else if opponent2Num == team.TeamNumber {
					opponent2Name = team.NameShort
				}
			}

			f.SetCellValue(teamName, fmt.Sprintf("C%d", 2*i+4), partnerNum)
			f.SetCellValue(teamName, fmt.Sprintf("C%d", 2*i+5), partnerName)
			f.SetCellValue(teamName, fmt.Sprintf("D%d", 2*i+4), opponent1Num)
			f.SetCellValue(teamName, fmt.Sprintf("D%d", 2*i+5), opponent1Name)
			f.SetCellValue(teamName, fmt.Sprintf("E%d", 2*i+4), opponent2Num)
			f.SetCellValue(teamName, fmt.Sprintf("E%d", 2*i+5), opponent2Name)

			redStyleTop, err := f.NewStyle(&excelize.Style{
				Font: &excelize.Font{
					Size:  14,
					Color: "FF0000",
				},
				Border: []excelize.Border{
					{Type: "left", Style: 1, Color: "000000"},
					{Type: "top", Style: 1, Color: "000000"},
					{Type: "bottom", Style: 1, Color: "FFFFFF"},
					{Type: "right", Style: 1, Color: "000000"},
				},
				Alignment: &excelize.Alignment{
					Horizontal: "center",
					Vertical:   "center",
				},
			})
			if err != nil {
				fmt.Println(err)
				return
			}

			redStyleBot, err := f.NewStyle(&excelize.Style{
				Font: &excelize.Font{
					Size:  14,
					Color: "FF0000",
				},
				Border: []excelize.Border{
					{Type: "left", Style: 1, Color: "000000"},
					{Type: "top", Style: 1, Color: "FFFFFF"},
					{Type: "bottom", Style: 1, Color: "000000"},
					{Type: "right", Style: 1, Color: "000000"},
				},
				Alignment: &excelize.Alignment{
					Horizontal: "center",
					Vertical:   "center",
				},
			})
			if err != nil {
				fmt.Println(err)
				return
			}

			blueStyleTop, err := f.NewStyle(&excelize.Style{
				Font: &excelize.Font{
					Size:  14,
					Color: "0000FF",
				},
				Border: []excelize.Border{
					{Type: "left", Style: 1, Color: "000000"},
					{Type: "top", Style: 1, Color: "000000"},
					{Type: "bottom", Style: 1, Color: "FFFFFF"},
					{Type: "right", Style: 1, Color: "000000"},
				},
				Alignment: &excelize.Alignment{
					Horizontal: "center",
					Vertical:   "center",
				},
			})
			if err != nil {
				fmt.Println(err)
				return
			}

			blueStyleBot, err := f.NewStyle(&excelize.Style{
				Font: &excelize.Font{
					Size:  14,
					Color: "0000FF",
				},
				Border: []excelize.Border{
					{Type: "left", Style: 1, Color: "000000"},
					{Type: "top", Style: 1, Color: "FFFFFF"},
					{Type: "bottom", Style: 1, Color: "000000"},
					{Type: "right", Style: 1, Color: "000000"},
				},
				Alignment: &excelize.Alignment{
					Horizontal: "center",
					Vertical:   "center",
				},
			})
			if err != nil {
				fmt.Println(err)
				return
			}

			if friendlyStation == "Red" {
				f.SetCellStyle(teamName, fmt.Sprintf("B%d", 2*i+4), fmt.Sprintf("C%d", 2*i+4), redStyleTop)
				f.SetCellStyle(teamName, fmt.Sprintf("B%d", 2*i+5), fmt.Sprintf("C%d", 2*i+5), redStyleBot)
				f.SetCellStyle(teamName, fmt.Sprintf("D%d", 2*i+4), fmt.Sprintf("E%d", 2*i+4), blueStyleTop)
				f.SetCellStyle(teamName, fmt.Sprintf("D%d", 2*i+5), fmt.Sprintf("E%d", 2*i+5), blueStyleBot)
			} else {
				f.SetCellStyle(teamName, fmt.Sprintf("B%d", 2*i+4), fmt.Sprintf("C%d", 2*i+4), blueStyleTop)
				f.SetCellStyle(teamName, fmt.Sprintf("B%d", 2*i+5), fmt.Sprintf("C%d", 2*i+5), blueStyleBot)
				f.SetCellStyle(teamName, fmt.Sprintf("D%d", 2*i+4), fmt.Sprintf("E%d", 2*i+4), redStyleTop)
				f.SetCellStyle(teamName, fmt.Sprintf("D%d", 2*i+5), fmt.Sprintf("E%d", 2*i+5), redStyleBot)
			}
		}
	}

	if err := f.DeleteSheet("Sheet1"); err != nil {
		fmt.Println(err)
	}

	if err := f.SaveAs("FTCCMP1FRAN.xlsx"); err != nil {
		fmt.Println(err)
	}
}
