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
	"os"
	"golang.org/x/net/proxy"
)

func SetProxy(url string) {
	os.Setenv("HQPROXY", url)
}

func HQVerify(number string) (*HQVerification, error) {
	verification := HQVerification{}
	verificationError := HQError{}

	body := `{"method":"sms","phone":"` + number + `"}`
	req, _ := http.NewRequest("POST", "https://api-quiz.hype.space/verifications", strings.NewReader(body))

	req.Header.Add("x-hq-client", "Android/1.6.2")
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("content-length", strconv.Itoa(len(body)))
	req.Header.Add("user-agent", "okhttp/3.8.0")

	t := &http.Transport{}
	if os.Getenv("HQPROXY") != "" {
		url, _ := url.Parse(os.Getenv("HQPROXY"))
		t.Proxy = http.ProxyURL(url)
	}
	res, err := http.Client{Transport: t}.Do(req)
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
func HQConfirm(verification *HQVerification, code string) (*HQAuth, error) {
	authInfo := HQAuth{}
	authErr := HQError{}

	body := `{"code":"` + code + `"}`
	req, _ := http.NewRequest("POST", "https://api-quiz.hype.space/verifications/"+verification.VerificationID, strings.NewReader(body))

	req.Header.Add("x-hq-client", "Android/1.6.2")
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("content-length", strconv.Itoa(len(body)))
	req.Header.Add("user-agent", "okhttp/3.8.0")

	t := &http.Transport{}
	if os.Getenv("HQPROXY") != "" {
		url, _ := url.Parse(os.Getenv("HQPROXY"))
		t.Proxy = http.ProxyURL(url)
	}
	res, err := http.Client{Transport: t}.Do(req)
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
func HQCreate(verification *HQVerification, username, referrer, region string) (*HQInfo, error) {
	info := HQInfo{}
	createError := HQError{}

	body := `{"country":"` + region + `","language":"en","referringUsername":"` + referrer + `","username":"` + username + `","verificationId":"` + verification.VerificationID + `"}`
	req, _ := http.NewRequest("POST", "https://api-quiz.hype.space/users", strings.NewReader(body))

	req.Header.Add("x-hq-client", "Android/1.6.2")
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("content-length", strconv.Itoa(len(body)))
	req.Header.Add("user-agent", "okhttp/3.8.0")

	t := &http.Transport{}
	if os.Getenv("HQPROXY") != "" {
		url, _ := url.Parse(os.Getenv("HQPROXY"))
		t.Proxy = http.ProxyURL(url)
	}
	res, err := http.Client{Transport: t}.Do(req)
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
func HQWeekly(info *HQInfo) (error) {
	authErr := HQError{}

	body := `{}`
	req, _ := http.NewRequest("POST", "https://api-quiz.hype.space/easter-eggs/makeItRain", strings.NewReader(body))

	req.Header.Add("x-hq-client", "Android/1.6.2")
	req.Header.Add("authorization", "Bearer "+info.AccessToken)
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("content-length", strconv.Itoa(len(body)))
	req.Header.Add("user-agent", "okhttp/3.8.0")

	t := &http.Transport{}
	if os.Getenv("HQPROXY") != "" {
		url, _ := url.Parse(os.Getenv("HQPROXY"))
		t.Proxy = http.ProxyURL(url)
	}
	res, err := http.Client{Transport: t}.Do(req)
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

func Schedule(bearer string) (HQSchedule) {
	req, _ := http.NewRequest("GET", "https://api-quiz.hype.space/shows/now?type=hq", nil)
	req.Header.Set("authorization", bearer)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return HQSchedule{}
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	schedule := HQSchedule{}
	json.Unmarshal(bytes, &schedule)

	return schedule
}
func HQConnect(id int, bearer string) (*HQSocket, error) {
	var u = url.URL{Scheme: "wss", Host: "ws-quiz.hype.space", Path: "/ws/" + strconv.Itoa(id)}

	request := http.Header{}
	request.Add("Authorization", bearer)
	var dialer *websocket.Dialer

	if os.Getenv("HQPROXY") != "" {
		netDialer, _ := proxy.SOCKS5("tcp", os.Getenv("HQPROXY"), nil, proxy.Direct)
		dialer = &websocket.Dialer{NetDial: netDialer.Dial}
	} else {
		dialer = websocket.DefaultDialer
	}

	c, _, err := dialer.Dial(u.String(), request)
	return &HQSocket{c}, err
}
func HQDebug() (*HQSocket, error) {
	var u = url.URL{Scheme: "wss", Host: "hqecho.herokuapp.com"}
	var dialer *websocket.Dialer

	if os.Getenv("HQPROXY") != "" {
		netDialer, _ := proxy.SOCKS5("tcp", os.Getenv("HQPROXY"), nil, proxy.Direct)
		dialer = &websocket.Dialer{NetDial: netDialer.Dial}
	} else {
		dialer = websocket.DefaultDialer
	}

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
