package go_hqtrivia

import (
	"time"
	"github.com/gorilla/websocket"
)

type UserInfo struct {
	UserID      int    `json:"userId"`
	Username    string `json:"username"`
	Admin       bool   `json:"admin"`
	Tester      bool   `json:"tester"`
	Guest       bool   `json:"guest"`
	AvatarURL   string `json:"avatarUrl"`
	LoginToken  string `json:"loginToken"`
	AccessToken string `json:"accessToken"`
	AuthToken   string `json:"authToken"`
}
type AuthInfo struct {
	Auth UserInfo `json:"auth"`
}

type Verification struct {
	CallsEnabled   bool      `json:"callsEnabled"`
	Expires        time.Time `json:"expires"`
	Phone          string    `json:"phone"`
	RetrySeconds   int       `json:"retrySeconds"`
	VerificationID string    `json:"verificationId"`
}

type HQError struct {
	Error     string `json:"error"`
	ErrorCode int    `json:"errorCode"`
}

type Schedule struct {
	Active        bool        `json:"active"`
	AtCapacity    bool        `json:"atCapacity"`
	ShowID        int         `json:"showId"`
	ShowType      string      `json:"showType"`
	StartTime     time.Time   `json:"startTime"`
	NextShowTime  interface{} `json:"nextShowTime"`
	NextShowPrize interface{} `json:"nextShowPrize"`
	Upcoming []struct {
		Time  time.Time `json:"time"`
		Prize string    `json:"prize"`
	} `json:"upcoming"`
	Prize int `json:"prize"`
	Broadcast struct {
		BroadcastID   int           `json:"broadcastId"`
		UserID        int           `json:"userId"`
		Title         string        `json:"title"`
		Status        int           `json:"status"`
		State         string        `json:"state"`
		ChannelID     int           `json:"channelId"`
		Created       time.Time     `json:"created"`
		Started       time.Time     `json:"started"`
		Ended         interface{}   `json:"ended"`
		Permalink     string        `json:"permalink"`
		ThumbnailData interface{}   `json:"thumbnailData"`
		Tags          []interface{} `json:"tags"`
		SocketURL     string        `json:"socketUrl"`
		Streams struct {
			Source      string `json:"source"`
			Passthrough string `json:"passthrough"`
			High        string `json:"high"`
			Medium      string `json:"medium"`
			Low         string `json:"low"`
		} `json:"streams"`
		StreamURL         string `json:"streamUrl"`
		StreamKey         string `json:"streamKey"`
		RelativeTimestamp int    `json:"relativeTimestamp"`
		Links struct {
			Self       string `json:"self"`
			Transcript string `json:"transcript"`
			Viewers    string `json:"viewers"`
		} `json:"links"`
	} `json:"broadcast"`
	GameKey       string `json:"gameKey"`
	BroadcastFull bool   `json:"broadcastFull"`
}

type HQSocket struct {
	*websocket.Conn
}

type HQQuestion struct {
	Type        string `json:"type"`
	TotalTimeMs int    `json:"totalTimeMs"`
	TimeLeftMs  int    `json:"timeLeftMs"`
	QuestionID  int    `json:"questionId"`
	Question    string `json:"question"`
	Category    string `json:"category"`
	Answers []struct {
		AnswerID int    `json:"answerId"`
		Text     string `json:"text"`
	} `json:"answers"`
	QuestionNumber int       `json:"questionNumber"`
	QuestionCount  int       `json:"questionCount"`
	Ts             time.Time `json:"ts"`
	Sent           time.Time `json:"sent"`
}