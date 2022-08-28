package rooms

import "errors"

type QuestionAnswer struct {
	Question string;
	Answers []Answer
}

type Answer struct {
	Data string;
	Correct bool;
}

type Room struct {
	Code string;
	QuestionAnswers []QuestionAnswer;
}

var Rooms []Room

func CreateRoom(room Room) (*Room, error) {
	if r := FindRoomByCode(room.Code); r != nil {
		return nil, errors.New("Room by this code already exists");
	}

	newRoom := Room{room.Code, room.QuestionAnswers}

	_ = append(Rooms, newRoom)

	return &newRoom, nil
}

func FindRoomByCode(code string) *Room {
	for _, r := range Rooms {
		if code == r.Code {
			return &r
		}
	}

	return nil
}

func JoinRoom(code string) bool {
	room := FindRoomByCode(code)

	return room != nil
}