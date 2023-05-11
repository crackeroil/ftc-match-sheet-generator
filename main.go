package main

func main() {
	goodFlags, seasonFlag, eventCodeFlag, apiKeyFlag := CheckFlags()

	if goodFlags {
		GenerateMatchsheets(seasonFlag, eventCodeFlag, apiKeyFlag)
	}
}
