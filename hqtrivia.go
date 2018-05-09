package go_hqtrivia

import (
	"net/http"
	"strings"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"errors"
)

func Verify(number string) (*Verification, error) {
	verification := Verification{}
	verificationError := HQError{}

	body := `{"method":"sms","phone":"` + number + `"}`
	req, _ := http.NewRequest("POST", "https://api-quiz.hype.space/verifications", strings.NewReader(body))

	req.Header.Add("x-hq-client", "Android/1.6.2")
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("content-length", strconv.Itoa(len(body)))
	req.Header.Add("user-agent", "okhttp/3.8.0")

	res, err := http.DefaultClient.Do(req)
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
func Confirm(verification *Verification, code string) (*AuthInfo, error) {
	authInfo := AuthInfo{}
	authErr := HQError{}

	body := `{"code":"` + code + `"}`
	req, _ := http.NewRequest("POST", "https://api-quiz.hype.space/verifications/"+verification.VerificationID, strings.NewReader(body))

	req.Header.Add("x-hq-client", "Android/1.6.2")
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("content-length", strconv.Itoa(len(body)))
	req.Header.Add("user-agent", "okhttp/3.8.0")

	res, err := http.DefaultClient.Do(req)
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
			if authInfo.Auth == (UserInfo{}) {
				return nil, nil
			}
			return nil, errors.New("unknown error")
		}
	}

	return &authInfo, nil
}
func Create(verification *Verification, username, referrer, region string) (*UserInfo, error) {
	info := UserInfo{}
	createError := HQError{}

	body := `{"country":"` + region + `","language":"en","referringUsername":"` + referrer + `","username":"` + username + `","verificationId":"` + verification.VerificationID + `"}`
	req, _ := http.NewRequest("POST", "https://api-quiz.hype.space/users", strings.NewReader(body))

	req.Header.Add("x-hq-client", "Android/1.6.2")
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("content-length", strconv.Itoa(len(body)))
	req.Header.Add("user-agent", "okhttp/3.8.0")

	res, err := http.DefaultClient.Do(req)
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
