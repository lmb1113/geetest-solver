package solver

import (
	"fmt"
	"net/http"
	"net/url"
)

const captchaSolvedSuccessFlag = "success"

type GeetestSolver interface {
	SolveV4Puzzle() (*V4PuzzleSolution, error)
}

type GeetestSolverConfig struct {
	WebsiteUrl string
	CaptchaId  string
	UserAgent  string
	Proxy      string
	httpClient *http.Client
}

type V4PuzzleSolution struct {
	Solution *V4PuzzleCaptchaResponse
}

func NewGeetestSolver(websiteUrl string, captchaId string, userAgent string, proxy string) (GeetestSolver, error) {
	httpClient := &http.Client{}

	if proxy != "" {
		parsedProxy, err := url.Parse(proxy)

		if err != nil {
			return nil, err
		}

		httpClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(parsedProxy),
		}
	}

	return &GeetestSolverConfig{
		WebsiteUrl: websiteUrl,
		CaptchaId:  captchaId,
		UserAgent:  userAgent,
		Proxy:      proxy,
		httpClient: httpClient,
	}, nil
}

func (c *GeetestSolverConfig) SolveV4Puzzle() (*V4PuzzleSolution, error) {
	captchaInfo, err := c.getV4PuzzleCaptchaInfo()

	if err != nil {
		return nil, err
	}

	captchaResponse, err := c.verifyV4PuzzleCaptchaInfo(captchaInfo)

	if err != nil {
		return nil, err
	}

	if captchaResponse.Status != captchaSolvedSuccessFlag {
		return nil, fmt.Errorf("failed solving captcha")
	}

	return &V4PuzzleSolution{
		Solution: captchaResponse,
	}, nil
}
