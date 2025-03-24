package geetest_solver

import (
	"errors"
	"fmt"
	"github.com/brianxor/geetest-solver/solver"
	"time"
)

type Options struct {
	WebsiteUrl string `json:"websiteUrl"`
	CaptchaId  string `json:"captchaId"`
	UserAgent  string `json:"userAgent"`
	Proxy      string `json:"proxy"`
	UserInfo   string `json:"user_info"`
}
type Resp struct {
	Solution  any    `json:"solution"`
	SolveTime string `json:"solveTime"`
	Success   bool   `json:"success"`
}

func V4PuzzleSolveHandler(o *Options) (*Resp, error) {
	if o.WebsiteUrl == "" || o.CaptchaId == "" || o.UserAgent == "" {
		return nil, errors.New("websiteUrl or userAgent is empty")
	}

	captchaSolver, err := solver.NewGeetestSolver(o.WebsiteUrl, o.CaptchaId, o.UserAgent, o.Proxy, o.UserInfo)

	if err != nil {
		return nil, errors.New("failed creating solver")
	}
	start := time.Now()
	solution, err := captchaSolver.SolveV4Puzzle()
	if err != nil {
		return nil, errors.New("failed solving captcha")
	}

	solveTime := time.Since(start)

	return &Resp{
		Success:   true,
		SolveTime: fmt.Sprintf("%s", solveTime),
		Solution:  solution.Solution,
	}, nil
}
