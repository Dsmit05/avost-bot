package nclient

import (
	"encoding/json"
	"io"
	"net/http"
	neturl "net/url"
	"time"

	"github.com/Dsmit05/avost-bot/pkg/animevost/models"
	"github.com/mailru/easyjson"
	"github.com/pkg/errors"

	"strconv"
)

type Client struct {
	cli     http.Client
	baseUrl string
}

func NewClient(baseUrl string) (*Client, error) {
	client := http.Client{
		Timeout: time.Second * 3,
	}

	return &Client{baseUrl: baseUrl, cli: client}, nil
}

// GetPage получает последние страницы с количеством аниме.
func (c *Client) GetPage(page, quantity int) (models.AnimeSpecs, error) {
	url := c.baseUrl + "/last?page=" + strconv.Itoa(page) + "&quantity=" + strconv.Itoa(quantity)

	resp, err := c.cli.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "c.cli.Get(%s) error", url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("error: status code not ok, code: %d, url: %s", resp.StatusCode, url)
	}

	var m models.Resp
	if err = easyjson.UnmarshalFromReader(resp.Body, &m); err != nil {
		return nil, errors.Wrap(err, "easyjson.UnmarshalFromReader() error")
	}

	return m.Items, nil
}

// SearchForName по имени ищем название.
func (c *Client) SearchForName(name string) (models.AnimeSpecs, error) {
	url := c.baseUrl + "/search"

	args := neturl.Values{}
	args.Add("name", name)

	resp, err := c.cli.PostForm(url, args)
	if err != nil {
		return nil, errors.Wrapf(err, "c.cli.PostForm(%s) with params: %s error", url, name)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("error: status code not ok, code: %d, url: %s", resp.StatusCode, url)
	}

	var m models.Resp
	if err = easyjson.UnmarshalFromReader(resp.Body, &m); err != nil {
		return nil, errors.Wrap(err, "easyjson.UnmarshalFromReader() error")
	}

	return m.Items, nil
}

// GetPlayList получаем плейлист по id.
func (c *Client) GetPlayList(id int) (models.Playlists, error) {
	url := c.baseUrl + "/playlist"

	args := neturl.Values{}
	args.Add("id", strconv.Itoa(id))

	resp, err := c.cli.PostForm(url, args)
	if err != nil {
		return nil, errors.Wrapf(err, "c.cli.PostForm(%s) with params: %d error", url, id)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("error: status code not ok, code: %d, url: %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "io.ReadAll() with params: %d error", id)
	}

	var m models.Playlists
	if err = json.Unmarshal(body, &m); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal() error")
	}

	return m, nil
}

// GetInfo получаем models.AnimeSpec по id
func (c *Client) GetInfo(id int) (models.AnimeSpec, error) {
	url := c.baseUrl + "/info"

	args := neturl.Values{}
	args.Add("id", strconv.Itoa(id))

	resp, err := c.cli.PostForm(url, args)
	if err != nil {
		return models.AnimeSpec{}, errors.Wrapf(err, "c.cli.PostForm(%s) with params: %d error", url, id)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.AnimeSpec{}, errors.Errorf("error: status code not ok, code: %d, url: %s", resp.StatusCode, url)
	}

	var m models.Resp
	if err = easyjson.UnmarshalFromReader(resp.Body, &m); err != nil {
		return models.AnimeSpec{}, errors.Wrap(err, "easyjson.UnmarshalFromReader() error")
	}

	if len(m.Items) == 0 {
		return models.AnimeSpec{}, errors.Errorf("error: Not found, empty items with params: %d", id)
	}

	return m.Items[0], nil
}
