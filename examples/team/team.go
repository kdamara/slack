package main

import (
	"fmt"

	"github.com/kdamara/slack"
)

func main() {
	api := slack.New("xoxp-2165961063-2187024495-91143250439-f7260fba7101be25673e3d402da10c25")
	billingActive, err := api.GetBillableInfo("U025H0QEK")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Println(billingActive)

  billingActiveTeam, err := api.GetBillableInfoForTeam()
  fmt.Println(billingActiveTeam)
}
