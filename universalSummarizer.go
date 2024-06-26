package kagi

import (
	"fmt"
	"github.com/httpjamesm/kagigo/types"
)

type SummaryType string

const (
	SummaryTypeSummary   SummaryType = "summary"
	SummaryTypeTakeaways SummaryType = "takeaway"
)

type SummaryEngine string

const (
	SummaryEngineCecil  SummaryEngine = "cecil"
	SummaryEngineAgnes  SummaryEngine = "agnes"
	SummaryEngineDaphne SummaryEngine = "daphne"
	SummaryEngineMuriel SummaryEngine = "muriel"
)

type UniversalSummarizerParams struct {
	URL         string        `json:"url"`
	SummaryType SummaryType   `json:"summary_type"`
	Engine      SummaryEngine `json:"engine"`
}

type UniversalSummarizerResponse struct {
	Meta struct {
		ID   string `json:"id"`
		Node string `json:"node"`
		Ms   int    `json:"ms"`
	} `json:"meta"`
	Data struct {
		Output string `json:"output"`
		Tokens int    `json:"tokens"`
	} `json:"data"`
	Errors []types.Error `json:"error"`
}

func (c *Client) UniversalSummarizerCompletion(params UniversalSummarizerParams) (res UniversalSummarizerResponse, err error) {
	if params.URL == "" {
		err = fmt.Errorf("url is required")
		return
	}

	err = c.SendRequest("POST", "/summarize", params, &res)
	if err != nil {
		return
	}

	if len(res.Errors) != 0 {
		errObj := res.Errors[0]
		err = fmt.Errorf("api returned error: %v", fmt.Sprintf("[code: %d, msg: %s, ref: %v]", errObj.Code, errObj.Msg, errObj.Ref))
		return
	}

	return
}
