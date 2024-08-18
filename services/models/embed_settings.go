package models

type EmbedSettings struct {
	StreamSettings
	Version    string `json:"version"`
	Magnet     string `json:"magnet"`
	TorrentURL string `json:"torrentUrl"`
	Referer    string `json:"referer"`
}
