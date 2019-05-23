package room

import (
	"encoding/json"

	"github.com/Balhazraell/logger"
	"github.com/Balhazraell/tools"
)

// ------------------- Вспомогательные методы -------------------
func createMap() {
	var step = 100
	var y = 0
	var chunckIDCounter = 0

	for i := 0; i < 3; i++ {
		var x = 0
		for j := 0; j < 3; j++ {
			chunc := chunc{
				ID:    chunckIDCounter,
				State: ChuncStateEmpty,
			}
			chunc.Сoordinates = append(
				chunc.Сoordinates,
				[2]int{x, y},
				[2]int{x + step, y},
				[2]int{x + step, y + step},
				[2]int{x, y + step},
				// необходимо указывать 5 элементов, так как последняя точка замыкает фигуру,
				// это нужно для определения пересечения координат мышки и элементов.
				[2]int{x, y},
			)

			x += step

			Room.Map[chunckIDCounter] = &chunc
			chunckIDCounter++
		}
		y += step
	}
}

//---------------------------------------------------------------
func updateClientsMap(clientsIDs []int) {
	gameMap, err := json.Marshal(Room.Map)

	if err != nil {
		logger.WarningPrintf("При формировнии json при упаковке карты, для обновления у клиентов произошла ошибка %v.", err)
		return
	}

	updateMap := updateMapStruct{
		Map:        gameMap,
		ClientsIDs: clientsIDs,
	}

	data, err := json.Marshal(updateMap)

	if err != nil {
		logger.WarningPrintf("При формировнии json при создании сооббщения для обновления карт произошла ошибка %v.", err)
		return
	}

	newMessage := MessageRMQ{
		HandlerName: "UpdateClientsMap",
		Data:        string(data),
	}

	PublishMessage(newMessage)
}

func clientConnect(clientID int){
	logger.InfoPrintf("В комнату room_id:%v вошел новый пользователь user_id=%v.", Room.ID, clientID)
	updateClientsMap([]int{clientID})
}

func clientDisconnect(clientID int){
	logger.InfoPrintf("Комнату room_id:%v покидает пользователь user_id=%v.", Room.ID, clientID)
	Room.clients = tools.DeleElementFromArraByIndex(Room.clients, clientID)
}

func setChunckState(clientID int, chunkID int) {
	Room.Map[chunkID].State = Room.GameState

	if Room.GameState == ChuncStateCross {
		Room.GameState = ChuncStateZero
	} else {
		Room.GameState = ChuncStateCross
	}

	updateClientsMap(Room.clients)
}