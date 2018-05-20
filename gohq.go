package gohq

import (
	"net/http"
	"strings"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"errors"
	"net/url"
	"github.com/gorilla/websocket"
	"fmt"
	"time"
)

func HQVerify(number string, transport *http.Transport) (*HQVerification, error) {
	if transport == nil {
		transport = &http.Transport{}
	}

	verification := HQVerification{}
	verificationError := HQError{}

	body := `{"method":"sms","phone":"` + number + `"}`
	req, _ := http.NewRequest("POST", "https://api-quiz.hype.space/verifications", strings.NewReader(body))

	req.Header.Add("x-hq-client", "Android/1.6.2")
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("content-length", strconv.Itoa(len(body)))
	req.Header.Add("user-agent", "okhttp/3.8.0")

	client := http.Client{Transport:transport, Timeout:time.Second * 10}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bytes, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	json.Unmarshal(bytes, &verification)
	if verification.VerificationID == "" {
		json.Unmarshal(bytes, &verificationError)
		if verificationError.Error != "" {
			return nil, errors.New(verificationError.Error)
		} else {
			return nil, errors.New("unknown error")
		}
	}

	return &verification, nil
}
func HQConfirm(verification *HQVerification, code string, transport *http.Transport) (*HQAuth, error) {
	if transport == nil {
		transport = &http.Transport{}
	}

	authInfo := HQAuth{}
	authErr := HQError{}

	body := `{"code":"` + code + `"}`
	req, _ := http.NewRequest("POST", "https://api-quiz.hype.space/verifications/"+verification.VerificationID, strings.NewReader(body))

	req.Header.Add("x-hq-client", "Android/1.6.2")
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("content-length", strconv.Itoa(len(body)))
	req.Header.Add("user-agent", "okhttp/3.8.0")

	client := http.Client{Transport:transport, Timeout:time.Second * 10}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bytes, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	json.Unmarshal(bytes, &authInfo)
	if authInfo.Auth.AccessToken == "" {
		json.Unmarshal(bytes, &authErr)
		if authErr.Error != "" {
			return nil, errors.New(authErr.Error)
		} else {
			if authInfo.Auth == (HQInfo{}) {
				return nil, nil
			}
			return nil, errors.New("unknown error")
		}
	}

	return &authInfo, nil
}
func HQCreate(verification *HQVerification, username, referrer, region string, transport *http.Transport) (*HQInfo, error) {
	if transport == nil {
		transport = &http.Transport{}
	}

	info := HQInfo{}
	createError := HQError{}

	body := `{"country":"` + region + `","language":"en","referringUsername":"` + referrer + `","username":"` + username + `","verificationId":"` + verification.VerificationID + `"}`
	req, _ := http.NewRequest("POST", "https://api-quiz.hype.space/users", strings.NewReader(body))

	req.Header.Add("x-hq-client", "Android/1.6.2")
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("content-length", strconv.Itoa(len(body)))
	req.Header.Add("user-agent", "okhttp/3.8.0")

	client := http.Client{Transport:transport, Timeout:time.Second * 10}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bytes, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	json.Unmarshal(bytes, &info)
	if info.AccessToken == "" {
		json.Unmarshal(bytes, &createError)
		if createError.Error != "" {
			return nil, errors.New(createError.Error)
		} else {
			return nil, errors.New("unknown error")
		}
	}

	return &info, nil
}
func HQWeekly(info *HQInfo, transport *http.Transport) (error) {
	if transport == nil {
		transport = &http.Transport{}
	}

	authErr := HQError{}

	body := `{}`
	req, _ := http.NewRequest("POST", "https://api-quiz.hype.space/easter-eggs/makeItRain", strings.NewReader(body))

	req.Header.Add("x-hq-client", "Android/1.6.2")
	req.Header.Add("authorization", "Bearer "+info.AccessToken)
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("content-length", strconv.Itoa(len(body)))
	req.Header.Add("user-agent", "okhttp/3.8.0")

	client := http.Client{Transport:transport, Timeout:time.Second * 10}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	bytes, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	json.Unmarshal(bytes, &authErr)
	if authErr.Error != "" {
		return errors.New(authErr.Error)
	} else {
		return nil
	}
}

func Schedule(bearer string, transport *http.Transport) (HQSchedule) {
	if transport == nil {
		transport = &http.Transport{}
	}

	req, _ := http.NewRequest("GET", "https://api-quiz.hype.space/shows/now?type=hq", nil)
	req.Header.Set("authorization", bearer)

	client := http.Client{Transport:transport, Timeout:time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return HQSchedule{}
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	schedule := HQSchedule{}
	json.Unmarshal(bytes, &schedule)

	return schedule
}
func HQConnect(id int, bearer string, dialer *websocket.Dialer) (*HQSocket, error) {
	if dialer == nil {
		dialer = websocket.DefaultDialer
	}
	var u = url.URL{Scheme: "wss", Host: "ws-quiz.hype.space", Path: "/ws/" + strconv.Itoa(id)}

	request := http.Header{}
	request.Add("Authorization", bearer)

	c, _, err := dialer.Dial(u.String(), request)
	return &HQSocket{c}, err
}

func HQDebug(dialer *websocket.Dialer) (*HQSocket, error) {
	if dialer == nil {
		dialer = websocket.DefaultDialer
	}

	var u = url.URL{Scheme: "wss", Host: "hqecho.herokuapp.com"}

	c, _, err := dialer.Dial(u.String(), nil)
	return &HQSocket{c}, err
}

func (ws *HQSocket) SendSocketSubscribe(broadcastID int) error {
	return ws.WriteMessage(websocket.TextMessage, []byte(`{"type":"subscribe","broadcastId":`+strconv.Itoa(broadcastID)+`}`))
}
func (ws *HQSocket) SendSocketAnswer(broadcastID, questionID, answerID int) error {
	return ws.WriteMessage(websocket.TextMessage, []byte(`{"type":"answer","broadcastId":`+strconv.Itoa(broadcastID)+`,"questionId":`+strconv.Itoa(questionID)+`,"answerId":`+strconv.Itoa(answerID)+`}`))
}
func (ws *HQSocket) SendSocketExtraLife(broadcastID, questionID int) error {
	return ws.WriteMessage(websocket.TextMessage, []byte(`{"type":"useExtraLife","broadcastId":`+strconv.Itoa(broadcastID)+`,"questionId":`+strconv.Itoa(questionID)+`}`))
}

func (ws *HQSocket) Read() (message []byte, err error) {
	_, message, err = ws.ReadMessage()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from a crash, fix this lazy boi")
			message = nil
			err = errors.New("panic")
		}
	}()

	if err != nil {
		return nil, err
	} else {
		return message, nil
	}
}
func (ws *HQSocket) ParseQuestion(message []byte) (*HQQuestion) {
	var question *HQQuestion
	json.Unmarshal(message, &question)
	if question != nil && len(question.Answers) != 0 {
		return question
	}
	return nil
}
func (ws *HQSocket) ParseStats(message []byte) (*HQStats) {
	var stats *HQStats
	json.Unmarshal(message, &stats)
	if stats != nil && stats.ViewerCounts.Playing != 0 {
		return stats
	}

	return nil
}
func (ws *HQSocket) ParseQuestionSummary(message []byte) (*HQQuestionSummary) {
	var summary *HQQuestionSummary
	json.Unmarshal(message, &summary)
	if summary != nil && summary.QuestionID != 0 {
		return summary
	}

	return nil
}