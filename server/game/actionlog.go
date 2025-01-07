package game

import "encoding/json"

type ActionLog struct {
	PlayerID string `json:"playerId"`
	Message  string `json:"message"`
}

type ActionLogResponse struct {
	Event string    `json:"event"`
	Data  ActionLog `json:"data"`
}

func NewActionLog(playerID string, message string) ActionLog {
	return ActionLog{
		PlayerID: playerID,
		Message:  message,
	}
}

func (log *ActionLog) getEncodedActionLogResponse() ([]byte, error) {
	response := ActionLogResponse{
		Event: "actionlog",
		Data:  *log,
	}
	encodedResp, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return encodedResp, nil
}
