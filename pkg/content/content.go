package content

import (
	"github.com/muesli/termenv"
	"github.com/sovietscout/valorank/pkg/models"
)

const (
	NAME    = "Valorank"
	VERSION = "v1.0" // β
	AUTHOR  = "SovietScout"
)

var (
	// Using structs to avoid looping
	preAscendantSeasonID = map[string]struct{}{
		"0df5adb9-4dcb-6899-1306-3e9860661dd3": {},
		"3f61c772-4560-cd3f-5d3f-a7ab5abda6b3": {},
		"0530b9c4-4980-f2ee-df5d-09864cd00542": {},
		"46ea6166-4573-1128-9cea-60a15640059b": {},
		"fcf2c8f4-4324-e50b-2e23-718e4a3ab046": {},
		"97b6e739-44cc-ffa7-49ad-398ba502ceb0": {},
		"ab57ef51-4e59-da91-cc8d-51a5a2b9b8ff": {},
		"52e9749a-429b-7060-99fe-4595426a0cf7": {},
		"71c81c67-4fae-ceb1-844c-aab2bb8710fa": {},
		"2a27e5d2-4d30-c9e2-b15a-93b8909a442c": {},
		"4cb622e1-4244-6da3-7276-8daaf1c01be2": {},
		"a16955a5-4ad0-f761-5e9e-389df1c892fb": {},
		"97b39124-46ce-8b55-8fd1-7cbf7ffe173f": {},
		"573f53ac-41a5-3a7d-d9ce-d6a6298e5704": {},
		"d929bc38-4ab6-7da4-94f0-ee84f8ac141e": {},
		"3e47230a-463c-a301-eb7d-67bb60357d4f": {},
		"808202d6-4f2b-a8ff-1feb-b3a0590ad79f": {},
	}

	gamePodStrings = map[string]string{
		"aresqa.aws-rclusterprod-use1-1.dev1-gp-ashburn-1":                 "Ashburn",
		"aresriot.aws-mes1-prod.eu-gp-bahrain-1":                           "Bahrain",
		"aresriot.aws-mes1-prod.ext1-gp-bahrain-1":                         "Bahrain",
		"aresriot.aws-rclusterprod-mes1-1.eu-gp-bahrain-awsedge-1":         "Bahrain",
		"aresriot.aws-rclusterprod-mes1-1.ext1-gp-bahrain-awsedge-1":       "Bahrain",
		"aresriot.aws-rclusterprod-mes1-1.tournament-gp-bahrain-awsedge-1": "Bahrain",
		"aresriot.aws-rclusterprod-bog1-1.latam-gp-bogota-1":               "Bogotá",
		"aresriot.aws-rclusterprod-bog1-1.tournament-gp-bogota-1":          "Bogotá",
		"aresriot.aws-rclusterprod-usw2-1.tournament-gp-cmob-1":            "CMOB 1",
		"aresriot.aws-rclusterprod-usw2-1.tournament-gp-cmob-2":            "CMOB 2",
		"aresriot.aws-rclusterprod-usw2-1.tournament-gp-cmob-3":            "CMOB 3",
		"aresriot.aws-rclusterprod-usw2-1.tournament-gp-cmob-4":            "CMOB 4",
		"aresriot.mtl-riot-ord2-3.ext1-gp-chicago-1":                       "Chicago",
		"aresriot.mtl-riot-ord2-3.latam-gp-chicago-1":                      "Chicago",
		"aresqa.aws-rclusterprod-dfw1-1.dev1-gp-dallas-1":                  "Dallas",
		"aresqa.aws-rclusterprod-euc1-1.dev1-gp-frankfurt-1":               "Frankfurt",
		"aresqa.aws-rclusterprod-euc1-1.stage1-gp-frankfurt-1":             "Frankfurt",
		"aresriot.aws-euc1-prod.eu-gp-frankfurt-1":                         "Frankfurt",
		"aresriot.aws-euc1-prod.ext1-gp-eu1":                               "Frankfurt",
		"aresriot.aws-euc1-prod.ext1-gp-frankfurt-1":                       "Frankfurt",
		"aresriot.aws-rclusterprod-euc1-1.ext1-gp-eu1":                     "Frankfurt",
		"aresriot.aws-rclusterprod-euc1-1.tournament-gp-frankfurt-1":       "Frankfurt",
		"aresriot.aws-rclusterprod-euc1-1.eu-gp-frankfurt-1":               "Frankfurt 1",
		"aresriot.aws-rclusterprod-euc1-1.eu-gp-frankfurt-awsedge-1":       "Frankfurt 2",
		"aresriot.aws-ape1-prod.ap-gp-hongkong-1":                          "Hong Kong",
		"aresriot.aws-ape1-prod.ext1-gp-hongkong-1":                        "Hong Kong",
		"aresriot.aws-rclusterprod-ape1-1.ext1-gp-hongkong-1":              "Hong Kong",
		"aresriot.aws-rclusterprod-ape1-1.tournament-gp-hongkong-1":        "Hong Kong",
		"aresriot.aws-rclusterprod-ape1-1.ap-gp-hongkong-1":                "Hong Kong 1",
		"aresriot.aws-rclusterprod-ape1-1.ap-gp-hongkong-awsedge-1":        "Hong Kong 2",
		"aresriot.mtl-riot-ist1-2.eu-gp-istanbul-1":                        "Istanbul",
		"aresriot.mtl-riot-ist1-2.tournament-gp-istanbul-1":                "Istanbul",
		"aresriot.aws-euw2-prod.eu-gp-london-1":                            "London",
		"aresriot.aws-rclusterprod-euw2-1.eu-gp-london-awsedge-1":          "London",
		"aresriot.aws-rclusterprod-euw2-1.tournament-gp-london-awsedge-1":  "London",
		"aresriot.aws-rclusterprod-mad1-1.eu-gp-madrid-1":                  "Madrid",
		"aresriot.aws-rclusterprod-mad1-1.tournament-gp-madrid-1":          "Madrid",
		"aresriot.mtl-tmx-mex1-1.ext1-gp-mexicocity-1":                     "Mexico City",
		"aresriot.mtl-tmx-mex1-1.latam-gp-mexicocity-1":                    "Mexico City",
		"aresriot.mtl-tmx-mex1-1.tournament-gp-mexicocity-1":               "Mexico City",
		"aresriot.mia1.latam-gp-miami-1":                                   "Miami",
		"aresriot.mia1.tournament-gp-miami-1":                              "Miami",
		"aresriot.aws-aps1-prod.ap-gp-mumbai-1":                            "Mumbai",
		"aresriot.aws-rclusterprod-aps1-1.ap-gp-mumbai-awsedge-1":          "Mumbai",
		"aresriot.aws-rclusterprod-aps1-1.tournament-gp-mumbai-awsedge-1":  "Mumbai",
		"aresriot.aws-rclusterprod-usw2-1.tournament-gp-offline-1":         "Offline 1",
		"aresriot.aws-rclusterprod-usw2-1.tournament-gp-offline-2":         "Offline 2",
		"aresriot.aws-rclusterprod-usw2-1.tournament-gp-offline-3":         "Offline 3",
		"aresriot.aws-rclusterprod-usw2-1.tournament-gp-offline-4":         "Offline 4",
		"aresriot.aws-rclusterprod-usw2-1.tournament-gp-offline-5":         "Offline 5",
		"aresriot.aws-rclusterprod-usw2-1.tournament-gp-offline-6":         "Offline 6",
		"aresriot.aws-rclusterprod-usw2-1.tournament-gp-offline-7":         "Offline 7",
		"aresriot.aws-rclusterprod-usw2-1.tournament-gp-offline-8":         "Offline 8",
		"aresriot.aws-euw3-prod.eu-gp-paris-1":                             "Paris",
		"aresriot.aws-rclusterprod-euw3-1.tournament-gp-paris-1":           "Paris",
		"aresriot.aws-rclusterprod-euw3-1.eu-gp-paris-1":                   "Paris 1",
		"aresriot.aws-rclusterprod-euw3-1.eu-gp-paris-awsedge-1":           "Paris 2",
		"aresriot.mtl-ctl-scl2-2.ext1-gp-santiago-1":                       "Santiago",
		"aresriot.mtl-ctl-scl2-2.latam-gp-santiago-1":                      "Santiago",
		"aresriot.mtl-ctl-scl2-2.tournament-gp-santiago-1":                 "Santiago",
		"aresriot.aws-rclusterprod-sae1-1.ext1-gp-saopaulo-1":              "Sao Paulo",
		"aresriot.aws-rclusterprod-sae1-1.tournament-gp-saopaulo-1":        "Sao Paulo",
		"aresriot.aws-sae1-prod.br-gp-saopaulo-1":                          "Sao Paulo",
		"aresriot.aws-sae1-prod.ext1-gp-saopaulo-1":                        "Sao Paulo",
		"aresriot.aws-rclusterprod-sae1-1.br-gp-saopaulo-1":                "Sao Paulo 1",
		"aresriot.aws-rclusterprod-sae1-1.br-gp-saopaulo-awsedge-1":        "Sao Paulo 2",
		"aresriot.aws-apne2-prod.ext1-gp-seoul-1":                          "Seoul",
		"aresriot.aws-apne2-prod.kr-gp-seoul-1":                            "Seoul",
		"aresriot.aws-rclusterprod-apne2-1.ext1-gp-seoul-1":                "Seoul",
		"aresriot.aws-rclusterprod-apne2-1.tournament-gp-seoul-1":          "Seoul",
		"aresriot.aws-rclusterprod-apne2-1.kr-gp-seoul-1":                  "Seoul 1",
		"aresriot.aws-apse1-prod.ap-gp-singapore-1":                        "Singapore",
		"aresriot.aws-apse1-prod.ext1-gp-singapore-1":                      "Singapore",
		"aresriot.aws-rclusterprod-apse1-1.ext1-gp-singapore-1":            "Singapore",
		"aresriot.aws-rclusterprod-apse1-1.tournament-gp-singapore-1":      "Singapore",
		"aresriot.aws-rclusterprod-apse1-1.ap-gp-singapore-1":              "Singapore 1",
		"aresriot.aws-rclusterprod-apse1-1.ap-gp-singapore-awsedge-1":      "Singapore 2",
		"aresriot.aws-eun1-prod.eu-gp-stockholm-1":                         "Stockholm",
		"aresriot.aws-rclusterprod-eun1-1.tournament-gp-stockholm-1":       "Stockholm",
		"aresriot.aws-rclusterprod-eun1-1.eu-gp-stockholm-1":               "Stockholm 1",
		"aresriot.aws-rclusterprod-eun1-1.eu-gp-stockholm-awsedge-1":       "Stockholm 2",
		"aresriot.aws-apse2-prod.ap-gp-sydney-1":                           "Sydney",
		"aresriot.aws-apse2-prod.ext1-gp-sydney-1":                         "Sydney",
		"aresriot.aws-rclusterprod-apse2-1.ext1-gp-sydney-1":               "Sydney",
		"aresriot.aws-rclusterprod-apse2-1.tournament-gp-sydney-1":         "Sydney",
		"aresriot.aws-rclusterprod-apse2-1.ap-gp-sydney-1":                 "Sydney 1",
		"aresriot.aws-rclusterprod-apse2-1.ap-gp-sydney-awsedge-1":         "Sydney 2",
		"aresriot.aws-apne1-prod.ap-gp-tokyo-1":                            "Tokyo",
		"aresriot.aws-apne1-prod.eu-gp-tokyo-1":                            "Tokyo",
		"aresriot.aws-apne1-prod.ext1-gp-kr1":                              "Tokyo",
		"aresriot.aws-apne1-prod.ext1-gp-tokyo-1":                          "Tokyo",
		"aresriot.aws-rclusterprod-apne1-1.eu-gp-tokyo-1":                  "Tokyo",
		"aresriot.aws-rclusterprod-apne1-1.ext1-gp-kr1":                    "Tokyo",
		"aresriot.aws-rclusterprod-apne1-1.tournament-gp-tokyo-1":          "Tokyo",
		"aresriot.aws-rclusterprod-apne1-1.ap-gp-tokyo-1":                  "Tokyo 1",
		"aresriot.aws-rclusterprod-apne1-1.ap-gp-tokyo-awsedge-1":          "Tokyo 2",
		"aresqa.aws-usw2-dev.main1-gp-tournament-2":                        "Tournament",
		"aresriot.aws-rclusterprod-atl1-1.na-gp-atlanta-1":                 "US Central (Georgia)",
		"aresriot.aws-rclusterprod-atl1-1.tournament-gp-atlanta-1":         "US Central (Georgia)",
		"aresriot.mtl-riot-ord2-3.na-gp-chicago-1":                         "US Central (Illinois)",
		"aresriot.mtl-riot-ord2-3.tournament-gp-chicago-1":                 "US Central (Illinois)",
		"aresriot.aws-rclusterprod-dfw1-1.na-gp-dallas-1":                  "US Central (Texas)",
		"aresriot.aws-rclusterprod-dfw1-1.tournament-gp-dallas-1":          "US Central (Texas)",
		"aresriot.aws-rclusterprod-use1-1.na-gp-ashburn-1":                 "US East (N. Virginia 1)",
		"aresriot.aws-rclusterprod-use1-1.na-gp-ashburn-awsedge-1":         "US East (N. Virginia 2)",
		"aresriot.aws-rclusterprod-use1-1.ext1-gp-ashburn-1":               "US East (N. Virginia)",
		"aresriot.aws-rclusterprod-use1-1.pbe-gp-ashburn-1":                "US East (N. Virginia)",
		"aresriot.aws-rclusterprod-use1-1.tournament-gp-ashburn-1":         "US East (N. Virginia)",
		"aresriot.aws-use1-prod.ext1-gp-ashburn-1":                         "US East (N. Virginia)",
		"aresriot.aws-use1-prod.na-gp-ashburn-1":                           "US East (N. Virginia)",
		"aresriot.aws-use1-prod.pbe-gp-ashburn-1":                          "US East (N. Virginia)",
		"aresriot.aws-rclusterprod-usw1-1.na-gp-norcal-1":                  "US West (N. California 1)",
		"aresriot.aws-rclusterprod-usw1-1.na-gp-norcal-awsedge-1":          "US West (N. California 2)",
		"aresriot.aws-rclusterprod-usw1-1.ext1-gp-na2":                     "US West (N. California)",
		"aresriot.aws-rclusterprod-usw1-1.pbe-gp-norcal-1":                 "US West (N. California)",
		"aresriot.aws-rclusterprod-usw1-1.tournament-gp-norcal-1":          "US West (N. California)",
		"aresriot.aws-usw1-prod.ext1-gp-na2":                               "US West (N. California)",
		"aresriot.aws-usw1-prod.ext1-gp-norcal-1":                          "US West (N. California)",
		"aresriot.aws-usw1-prod.na-gp-norcal-1":                            "US West (N. California)",
		"aresriot.aws-rclusterprod-usw2-1.na-gp-oregon-1":                  "US West (Oregon 1)",
		"aresriot.aws-rclusterprod-usw2-1.na-gp-oregon-awsedge-1":          "US West (Oregon 2)",
		"aresriot.aws-rclusterprod-usw2-1.pbe-gp-oregon-1":                 "US West (Oregon)",
		"aresriot.aws-rclusterprod-usw2-1.tournament-gp-oregon-1":          "US West (Oregon)",
		"aresriot.aws-usw2-prod.na-gp-oregon-1":                            "US West (Oregon)",
		"aresriot.aws-usw2-prod.pbe-gp-oregon-1":                           "US West (Oregon)",
		"aresqa.aws-usw2-dev.main1-gp-1":                                   "US West 1",
		"aresqa.aws-usw2-dev.stage1-gp-1":                                  "US West 1",
		"aresqa.aws-usw2-dev.main1-gp-4":                                   "US West 2",
		"aresriot.aws-rclusterprod-waw1-1.eu-gp-warsaw-1":                  "Warsaw",
		"aresriot.aws-rclusterprod-waw1-1.tournament-gp-warsaw-1":          "Warsaw",
	}

	agentIDs = map[string]*Data{
		"dade69b4-4f5a-8528-247b-219e5a1facd6": {Name: "Fade", Colour: "#5c5c5e"},
		"5f8d3a7f-467b-97f3-062c-13acf203c006": {Name: "Breach", Colour: "#d97a2e"},
		"f94c3b30-42be-e959-889c-5aa313dba261": {Name: "Raze", Colour: "#d97a2e"},
		"22697a3d-45bf-8dd7-4fec-84a9e28c69d7": {Name: "Chamber", Colour: "#ffce6f"},
		"601dbbe7-43ce-be57-2a40-4abd24953621": {Name: "KAY/O", Colour: "#85929c"},
		"6f2a04ca-43e0-be17-7f36-b3908627744d": {Name: "Skye", Colour: "#c0e69e"},
		"117ed9e3-49f3-6512-3ccf-0cada7e3823b": {Name: "Cypher", Colour: "#d18a5b"},
		"320b2a48-4d9b-a075-30f1-1f93a9b638fa": {Name: "Sova", Colour: "#258fcc"},
		"1e58de9c-4950-5125-93e9-a0aee9f98746": {Name: "Killjoy", Colour: "#ffd91f"},
		"707eab51-4836-f488-046a-cda6bf494859": {Name: "Viper", Colour: "#30ba87"},
		"eb93336a-449b-9c1b-0a54-a891f7921d69": {Name: "Phoenix", Colour: "#fe8266"},
		"41fb69c1-4189-7b37-f117-bcaf1e96f1bf": {Name: "Astra", Colour: "#712ae8"},
		"9f0d8ba9-4140-b941-57d3-a7ad57c6b417": {Name: "Brimstone", Colour: "#d97a2e"},
		"bb2a4828-46eb-8cd1-e765-15848195d751": {Name: "Neon", Colour: "#1c45a1"},
		"7f94d92c-4234-0a36-9646-3a87eb8b5c89": {Name: "Yoru", Colour: "#344ccf"},
		"569fdd95-4d10-43ab-ca70-79becc718b46": {Name: "Sage", Colour: "#5ae6d5"},
		"a3bfb853-43b2-7238-a4f1-ad90e9e46bcc": {Name: "Reyna", Colour: "#b565b5"},
		"8e253930-4c05-31dd-1b6c-968525494517": {Name: "Omen", Colour: "#47508f"},
		"add6443a-41bd-e414-f6ad-e58d267f4e95": {Name: "Jett", Colour: "#9adeff"},
		// "": 									{Name: "Unknown", Colour: "#646464"},
	}

	ranks = []*Data{
		{Name: "Unranked", Colour: "#2e2e2e"},
		{Name: "Unranked", Colour: "#2e2e2e"},
		{Name: "Unranked", Colour: "#2e2e2e"},
		{Name: "Iron 1", Colour: "#48453e"},
		{Name: "Iron 2", Colour: "#48453e"},
		{Name: "Iron 3", Colour: "#48453e"},
		{Name: "Bronze 1", Colour: "#bb8f5a"},
		{Name: "Bronze 2", Colour: "#bb8f5a"},
		{Name: "Bronze 3", Colour: "#bb8f5a"},
		{Name: "Silver 1", Colour: "#aeb2b2"},
		{Name: "Silver 2", Colour: "#aeb2b2"},
		{Name: "Silver 3", Colour: "#aeb2b2"},
		{Name: "Gold 1", Colour: "#c5ba3f"},
		{Name: "Gold 2", Colour: "#c5ba3f"},
		{Name: "Gold 3", Colour: "#c5ba3f"},
		{Name: "Platinum 1", Colour: "#18a7b9"},
		{Name: "Platinum 2", Colour: "#18a7b9"},
		{Name: "Platinum 3", Colour: "#18a7b9"},
		{Name: "Diamond 1", Colour: "#d864c7"},
		{Name: "Diamond 2", Colour: "#d864c7"},
		{Name: "Diamond 3", Colour: "#d864c7"},
		{Name: "Ascendant 1", Colour: "#189452"},
		{Name: "Ascendant 2", Colour: "#189452"},
		{Name: "Ascendant 3", Colour: "#189452"},
		{Name: "Immortal 1", Colour: "#dd4444"},
		{Name: "Immortal 2", Colour: "#dd4444"},
		{Name: "Immortal 3", Colour: "#dd4444"},
		{Name: "Radiant", Colour: "#fffdcd"},
	}

	p = termenv.ColorProfile()

	// Name Colours
	teamColour = p.Color("#4c97ed")
	oppColour  = p.Color("#ee4d4d")

	// State Colours
	offlineColour = p.Color("#ffffff")
	menuColour    = p.Color("#67ed4c")
	pregameColour = p.Color("#eef136")
	ingameColour  = p.Color("#67ed4c")

	ClientPlatform = "ew0KCSJwbGF0Zm9ybVR5cGUiOiAiUEMiLA0KCSJwbGF0Zm9ybU9TIjog" +
		"IldpbmRvd3MiLA0KCSJwbGF0Zm9ybU9TVmVyc2lvbiI6ICIxMC4wLjE5" +
		"MDQyLjEuMjU2LjY0Yml0IiwNCgkicGxhdGZvcm1DaGlwc2V0IjogIlVua25vd24iDQp9"

	ClientVersion   string
	CurrentSeasonID string
)

func PreAscendantSeasonID(seasonID string) bool {
	_, found := preAscendantSeasonID[seasonID]
	return found
}

func AgentFromID(AgentID string) *Data {
	if agent, found := agentIDs[AgentID]; found {
		return agent
	}

	return nil
}

func StyleAgent(agent *Data) string {
	if agent != nil {
		return termenv.String(agent.Name).Foreground(p.Color(agent.Colour)).String()
	} else {
		return ""
	}
}

func RankFromID(RankID int) string {
	rank := ranks[RankID]
	return termenv.String(rank.Name).Foreground(p.Color(rank.Colour)).String()
}

func ServerFromGamePod(gamePod string) string {
	return gamePodStrings[gamePod]
}

func ColourFromPlayer(player *models.Player) string {
	name := termenv.String(player.Name)

	if player.Ally {
		return name.Foreground(teamColour).String()
	}

	return name.Foreground(oppColour).String()
}

func ColourFromState(state models.State) string {
	s := termenv.String(string(state))
	var stateCol termenv.Color

	switch state {
	case models.OFFLINE:
		stateCol = offlineColour
	case models.MENU:
		stateCol = menuColour
	case models.PREGAME:
		stateCol = pregameColour
	case models.INGAME:
		stateCol = ingameColour
	}

	return s.Foreground(stateCol).String()
}
