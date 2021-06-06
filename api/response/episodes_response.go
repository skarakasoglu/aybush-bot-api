package response

import "time"

type Episode struct{
	Id int `json:"id"`
	Name string `json:"name"`
	StartTimestamp time.Time `json:"start_timestamp"`
	EndTimestamp time.Time `json:"end_timestamp"`
}

type EpisodesResponse struct {
	Response
	Episodes []Episode `json:"episodes"`
}
