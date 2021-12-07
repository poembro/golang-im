package model

import (
	"encoding/json"
	"errors"
)

type Int64 int64

type User struct {
	UserId Int64 `json:"user_id"`
	RoomId   string `json:"room_id"`
	DeviceId string `json:"device_id"`
	Nickname string `json:"nickname"`
	Face     string `json:"face"`
	ShopName string `json:"shop_name"`
	ShopId   string `json:"shop_id"`
	ShopFace string `json:"shop_face"`
	Platform string `json:"platform"`
	Suburl   string `json:"suburl"`
	Pushurl  string `json:"pushurl"`
	Unread      Int64    `json:"unread"`       // 未读
	LastMessage []string `json:"last_message"` //最后一条消息
}

func (u *Int64) UnmarshalJSON(bs []byte) error {
	var i int64
	if err := json.Unmarshal(bs, &i); err == nil {
		*u = Int64(i)
		return nil
	}
	var s string
	if err := json.Unmarshal(bs, &s); err != nil {
		return errors.New("expected a string or an integer")
	}
	if err := json.Unmarshal([]byte(s), &i); err != nil {
		return err
	}
	*u = Int64(i)
	return nil
}
