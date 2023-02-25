package services

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/urfave/cli"

	ra "github.com/webtor-io/rest-api/services"

	"github.com/dgrijalva/jwt-go"
)

const (
	apiKeyFlag       = "webtor-key"
	apiSecretFlag    = "webtor-secret"
	apiSecureFlag    = "webtor-rest-api-secure"
	apiHostFlag      = "webtor-rest-api-host"
	apiPortFlag      = "webtor-rest-api-port"
	rapidApiKeyFlag  = "rapidapi-key"
	rapidApiHostFlag = "rapidapi-host"
)

func RegisterApiFlags(f []cli.Flag) []cli.Flag {
	return append(f,
		cli.StringFlag{
			Name:   apiHostFlag,
			Usage:  "webtor rest-api host",
			EnvVar: "REST_API_SERVICE_HOST",
		},
		cli.IntFlag{
			Name:   apiPortFlag,
			Usage:  "webtor rest-api port",
			EnvVar: "REST_API_SERVICE_PORT",
			Value:  80,
		},
		cli.BoolFlag{
			Name:   apiSecureFlag,
			Usage:  "webtor rest-api secure (https)",
			EnvVar: "REST_API_SECURE",
		},
		cli.StringFlag{
			Name:   apiKeyFlag,
			Usage:  "webtor api key",
			Value:  "",
			EnvVar: "WEBTOR_API_KEY",
		},
		cli.StringFlag{
			Name:   apiSecretFlag,
			Usage:  "webtor api secret",
			Value:  "",
			EnvVar: "WEBTOR_API_SECRET",
		},
		cli.StringFlag{
			Name:   rapidApiHostFlag,
			Usage:  "RapidAPI host",
			Value:  "",
			EnvVar: "RAPIDAPI_HOST",
		},
		cli.StringFlag{
			Name:   rapidApiKeyFlag,
			Usage:  "RapidAPI key",
			Value:  "",
			EnvVar: "RAPIDAPI_KEY",
		},
	)
}

type EventData struct {
	Total     int64 `json:"total"`
	Completed int   `json:"completed"`
	Peers     int   `json:"peers"`
	Status    int   `json:"status"`
	Pieces    []struct {
		Position int  `json:"position"`
		Complete bool `json:"complete"`
		Priority int  `json:"priority"`
	} `json:"pieces"`
	Seeders  int `json:"seeders"`
	Leechers int `json:"leechers"`
}

type Claims struct {
	jwt.StandardClaims
	Rate          string `json:"rate"`
	Role          string `json:"role"`
	SessionID     string `json:"sessionID"`
	Domain        string `json:"domain"`
	Agent         string `json:"agent"`
	RemoteAddress string `json:"remoteAddress"`
}

type Api struct {
	url            string
	prepareRequest func(r *http.Request, c *Claims) (*http.Request, error)
	cl             *http.Client
}

type ListResourceContentOutputType string

const (
	OutputList ListResourceContentOutputType = "list"
	OutputTree ListResourceContentOutputType = "tree"
)

type ListResourceContentArgs struct {
	Limit  uint
	Offset uint
	Path   string
	Output ListResourceContentOutputType
}

func (s *ListResourceContentArgs) ToQuery() url.Values {
	q := url.Values{}
	limit := uint(10)
	offset := s.Offset
	path := "/"
	output := OutputList
	if s.Limit > 0 {
		limit = s.Limit
	}
	if s.Path != "" {
		path = s.Path
	}
	if s.Output != "" {
		output = s.Output
	}
	q.Set("limit", strconv.Itoa(int(limit)))
	q.Set("offset", strconv.Itoa(int(offset)))
	q.Set("path", path)
	q.Set("output", string(output))
	return q
}

func NewApi(c *cli.Context, cl *http.Client) *Api {
	host := c.String(apiHostFlag)
	port := c.Int(apiPortFlag)
	secure := c.Bool(apiSecureFlag)
	secret := c.String(apiSecretFlag)
	key := c.String(apiKeyFlag)
	rapidApiHost := c.String(rapidApiHostFlag)
	rapidApiKey := c.String(rapidApiKeyFlag)
	if rapidApiHost != "" {
		host = rapidApiHost
		port = 443
		secure = true
	}
	protocol := "http"
	if secure {
		protocol = "https"
	}
	u := fmt.Sprintf("%v://%v:%v", protocol, host, port)
	prepareRequest := func(r *http.Request, cl *Claims) (*http.Request, error) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, cl)
		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			return nil, err
		}
		r.Header.Set("X-Token", tokenString)
		r.Header.Set("X-Api-Key", key)
		return r, nil
	}
	if rapidApiHost != "" && rapidApiKey != "" {
		log.Info("using RapidAPI")
		prepareRequest = func(r *http.Request, cl *Claims) (*http.Request, error) {
			r.Header.Set("X-RapidAPI-Host", rapidApiHost)
			r.Header.Set("X-RapidAPI-Key", rapidApiKey)
			return r, nil
		}
	}
	log.Infof("api endpoint %v", u)
	return &Api{
		url:            u,
		cl:             cl,
		prepareRequest: prepareRequest,
	}
}

func (s *Api) GetResource(ctx context.Context, c *Claims, infohash string) (e *ra.ResourceResponse, err error) {
	u := s.url + "/resource/" + infohash
	e = &ra.ResourceResponse{}
	err = s.doRequest(ctx, c, u, "GET", nil, e)
	if e.ID == "" {
		e = nil
	}
	return
}

func (s *Api) StoreResource(ctx context.Context, c *Claims, resource []byte) (e *ra.ResourceResponse, err error) {
	u := s.url + "/resource"
	e = &ra.ResourceResponse{}
	err = s.doRequest(ctx, c, u, "POST", resource, e)
	if e.ID == "" {
		e = nil
	}
	return
}

func (s *Api) ListResourceContent(ctx context.Context, c *Claims, infohash string, args *ListResourceContentArgs) (e *ra.ListResponse, err error) {
	u := s.url + "/resource/" + infohash + "/list?" + args.ToQuery().Encode()
	e = &ra.ListResponse{}
	err = s.doRequest(ctx, c, u, "GET", nil, e)
	return
}

func (s *Api) doRequest(ctx context.Context, c *Claims, url string, method string, data []byte, v any) error {
	var payload io.Reader

	if data != nil {
		payload = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, payload)

	if err != nil {
		return err
	}

	req, err = s.prepareRequest(req, c)

	if err != nil {
		return err
	}

	res, err := s.cl.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, v)
		if err != nil {
			return err
		}
		return nil
	} else if res.StatusCode == http.StatusNotFound {
		return nil
	} else if res.StatusCode == http.StatusForbidden {
		return errors.Errorf("access is forbidden url=%v", url)
	} else {
		var e ra.ErrorResponse
		err = json.Unmarshal(body, &e)
		if err != nil {
			return errors.Wrapf(err, "failed to parse status=%v body=%v url=%v", res.StatusCode, body, req.URL)
		}
		return errors.New(e.Error)
	}
}

func (s *Api) ExportResourceContent(ctx context.Context, c *Claims, infohash string, itemID string) (e *ra.ExportResponse, err error) {
	u := s.url + "/resource/" + infohash + "/export/" + itemID
	e = &ra.ExportResponse{}
	err = s.doRequest(ctx, c, u, "GET", nil, e)
	// if e.Source.ID == nil
	// 	e = nil
	// }
	return
}

func (s *Api) Download(ctx context.Context, u string) (io.ReadCloser, error) {
	return s.DownloadWithRange(ctx, u, 0, -1)
}

func (s *Api) DownloadWithRange(ctx context.Context, u string, start int, end int) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if start != 0 || end != -1 {
		startStr := strconv.Itoa(start)
		endStr := ""
		if end != -1 {
			endStr = strconv.Itoa(end)
		}
		req.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", startStr, endStr))
	}
	if err != nil {
		log.WithError(err).Error("failed to make new request")
		return nil, err
	}
	res, err := s.cl.Do(req)
	if err != nil {
		log.WithError(err).Error("failed to do request")
		return nil, err
	}
	b := res.Body
	return b, nil
}

func (s *Api) Stats(ctx context.Context, u string) (chan EventData, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make new request")
	}
	res, err := s.cl.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do request")
	}
	ch := make(chan EventData)
	go func() {
		b := res.Body
		defer func() {
			close(ch)
			_ = b.Close()
		}()
		scanner := bufio.NewScanner(b)
		scanner.Split(bufio.ScanLines)

		t := ""
		for scanner.Scan() {
			if ctx.Err() != nil {
				log.WithError(ctx.Err()).Error("context error")
				break
			}
			if scanner.Err() != nil {
				log.WithError(scanner.Err()).Error("scanner error")
				break
			}
			line := scanner.Text()
			if strings.HasPrefix(line, "event: ") {
				t = strings.TrimSpace(strings.TrimPrefix(line, "event: "))
				continue
			}
			if t == "statupdate" && strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")
				var event EventData
				err := json.Unmarshal([]byte(data), &event)
				if err != nil {
					log.WithError(err).Errorf("failed to unmarshal data=%v line=%v", data, line)
					continue
				}
				ch <- event
				continue
			}
		}
	}()
	return ch, nil
}
