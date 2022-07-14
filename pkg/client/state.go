package client

type State string

const(
	OFFLINE State = "Offline"
	MENU State = "In Menu"
	PREGAME State = "In Pregame"
	INGAME State = "In Game"
)