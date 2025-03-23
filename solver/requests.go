package solver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

var geetestResponseRegex = regexp.MustCompile(`^geetest_\d+\((.*)\)$`)

const geetestCaptchaImageUrl = "https://static.geetest.com/"

type V4PuzzleCaptchaInfo struct {
	Status string `json:"status"`
	Data   struct {
		LotNumber   string `json:"lot_number"`
		CaptchaType string `json:"captcha_type"`
		Slice       string `json:"slice"`
		Bg          string `json:"bg"`
		Ypos        int    `json:"ypos"`
		Arrow       string `json:"arrow"`
		Js          string `json:"js"`
		Css         string `json:"css"`
		StaticPath  string `json:"static_path"`
		GctPath     string `json:"gct_path"`
		ShowVoice   bool   `json:"show_voice"`
		Feedback    string `json:"feedback"`
		Logo        bool   `json:"logo"`
		Pt          string `json:"pt"`
		CaptchaMode string `json:"captcha_mode"`
		Guard       bool   `json:"guard"`
		CheckDevice bool   `json:"check_device"`
		Language    string `json:"language"`
		LangReverse bool   `json:"lang_reverse"`
		CustomTheme struct {
			Style      string `json:"_style"`
			Color      string `json:"_color"`
			Gradient   string `json:"_gradient"`
			Hover      string `json:"_hover"`
			Brightness string `json:"_brightness"`
			Radius     string `json:"_radius"`
		} `json:"custom_theme"`
		PowDetail struct {
			Version  string `json:"version"`
			Bits     int    `json:"bits"`
			Datetime string `json:"datetime"`
			Hashfunc string `json:"hashfunc"`
		} `json:"pow_detail"`
		Payload         string `json:"payload"`
		ProcessToken    string `json:"process_token"`
		PayloadProtocol int    `json:"payload_protocol"`
	} `json:"data"`
}

type V4PuzzleCaptchaResponse struct {
	Status string `json:"status"`
	Data   struct {
		LotNumber string `json:"lot_number"`
		Result    string `json:"result"`
		FailCount int    `json:"fail_count"`
		Seccode   struct {
			CaptchaId     string `json:"captcha_id"`
			LotNumber     string `json:"lot_number"`
			PassToken     string `json:"pass_token"`
			GenTime       string `json:"gen_time"`
			CaptchaOutput string `json:"captcha_output"`
		} `json:"seccode"`
		Score           string `json:"score"`
		Payload         string `json:"payload"`
		ProcessToken    string `json:"process_token"`
		PayloadProtocol int    `json:"payload_protocol"`
	} `json:"data"`
}

func (c *GeetestSolverConfig) getV4PuzzleCaptchaInfo() (*V4PuzzleCaptchaInfo, error) {
	currentTimestamp := time.Now().UnixMilli()

	reqUrl := fmt.Sprintf("https://gcaptcha4.geetest.com/load?callback=geetest_%d&captcha_id=%s&client_type=web&pt=1&lang=eng",
		currentTimestamp,
		c.CaptchaId,
	)

	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)

	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"pragma":                   {"no-cache"},
		"cache-control":            {"no-cache"},
		"sec-ch-ua-platform":       {"\"Windows\""},
		"user-agent":               {c.UserAgent},
		"sec-ch-ua":                {"\"Chromium\";v=\"134\", \"Not:A-Brand\";v=\"24\", \"Google Chrome\";v=\"134\""},
		"sec-ch-ua-mobile":         {"?0"},
		"accept":                   {"*/*"},
		"sec-fetch-site":           {"cross-site"},
		"sec-fetch-mode":           {"no-cors"},
		"sec-fetch-dest":           {"script"},
		"sec-fetch-storage-access": {"active"},
		"referer":                  {c.WebsiteUrl},
		"accept-encoding":          {"gzip, deflate, br, zstd"},
		"accept-language":          {"en-US,en;q=0.9"},
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	parsedResponse, err := parseGeetestResponse(body)

	if err != nil {
		return nil, err
	}

	var captchaInfo V4PuzzleCaptchaInfo

	if err := json.Unmarshal(parsedResponse, &captchaInfo); err != nil {
		return nil, err
	}

	return &captchaInfo, nil
}

func (c *GeetestSolverConfig) verifyV4PuzzleCaptchaInfo(captchaInfo *V4PuzzleCaptchaInfo) (*V4PuzzleCaptchaResponse, error) {
	currentTimestamp := time.Now().UnixMilli()

	payload, err := c.generateV4PuzzlePayload(captchaInfo)

	if err != nil {
		return nil, err
	}

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}

	encryptedPayload, err := encryptPayload(jsonPayload, encryptionPublicKey)

	if err != nil {
		return nil, err
	}

	reqUrl := fmt.Sprintf("https://gcaptcha4.geetest.com/verify?callback=geetest_%d&captcha_id=%s&client_type=web&lot_number=%s&payload=%s&process_token=%s&payload_protocol=1&pt=1&w=%s",
		currentTimestamp,
		c.CaptchaId,
		captchaInfo.Data.LotNumber,
		captchaInfo.Data.Payload,
		captchaInfo.Data.ProcessToken,
		encryptedPayload,
	)

	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)

	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"pragma":                   {"no-cache"},
		"cache-control":            {"no-cache"},
		"sec-ch-ua-platform":       {"\"Windows\""},
		"user-agent":               {c.UserAgent},
		"sec-ch-ua":                {"\"Chromium\";v=\"134\", \"Not:A-Brand\";v=\"24\", \"Google Chrome\";v=\"134\""},
		"sec-ch-ua-mobile":         {"?0"},
		"accept":                   {"*/*"},
		"sec-fetch-site":           {"cross-site"},
		"sec-fetch-mode":           {"no-cors"},
		"sec-fetch-dest":           {"script"},
		"sec-fetch-storage-access": {"active"},
		"referer":                  {c.WebsiteUrl},
		"accept-encoding":          {"gzip, deflate, br, zstd"},
		"accept-language":          {"en-US,en;q=0.9"},
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	parsedResponse, err := parseGeetestResponse(body)

	if err != nil {
		return nil, err
	}

	var captchaResponse V4PuzzleCaptchaResponse

	if err := json.Unmarshal(parsedResponse, &captchaResponse); err != nil {
		return nil, err
	}

	return &captchaResponse, nil
}

func (c *GeetestSolverConfig) FetchImage(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"pragma":             {"no-cache"},
		"cache-control":      {"no-cache"},
		"sec-ch-ua-platform": {"\"Windows\""},
		"user-agent":         {c.UserAgent},
		"sec-ch-ua":          {"\"Chromium\";v=\"134\", \"Not:A-Brand\";v=\"24\", \"Google Chrome\";v=\"134\""},
		"sec-ch-ua-mobile":   {"?0"},
		"accept":             {"detection/avif,detection/webp,detection/apng,detection/svg+xml,detection/*,*/*;q=0.8"},
		"sec-fetch-site":     {"same-site"},
		"sec-fetch-mode":     {"no-cors"},
		"sec-fetch-dest":     {"detection"},
		"referer":            {c.WebsiteUrl},
		"accept-encoding":    {"gzip, deflate, br, zstd"},
		"accept-language":    {"en-US,en;q=0.9"},
		"priority":           {"i"},
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	imageData, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return imageData, nil
}

func parseGeetestResponse(body []byte) ([]byte, error) {
	matches := geetestResponseRegex.FindSubmatch(body)

	if len(matches) < 2 {
		return nil, fmt.Errorf("could not parse geetest response")
	}

	return matches[1], nil
}
