package packetType

const (
	ESTABLISH = "ESTABLISH"

	//To Create a game room
	CREATE_ROOM = "CREATE_ROOM"

	//To join a game room
	JOIN_ROOM = "JOIN_ROOM"

	//To leave a game room
	EXIT_ROOM = "EXIT_ROOM"

	//To Get room list
	ROOM_LIST_INFO = "ROOM_LIST_INFO"

	//To share game information
	GET_SESSION_INFO = "GET_SESSION_INFO"

	//Start the game
	START_GAME = "START_GAME"

	//Update game status
	UPDATE_GAME_STATUS = "UPDATE_GAME_STATUS"

	//End the game
	END_GAME = "END_GAME"
)
