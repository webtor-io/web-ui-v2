# web-ui-v2

The future version of [webtor.io](https://webtor.io)

Some features to mention:

1. Lightweight - less JavaScript code, no frontend frameworks, fewer bytes sent to the client.
2. Based on [webtor REST-API](https://github.com/webtor-io/rest-api).

## Roadmap

- [x] Torrent/magnet upload
- [x] Torrent listing
- [x] Direct file download
- [x] Direct folder download as ZIP-archive
- [x] Picture preview
- [x] Audio streaming
- [x] Video streaming
  - [x] Base player
  - [x] Subtitles support
  - [x] OpenSubtitles support
  - [ ] Subtitle uploading support
  - [ ] Chromecast support
  - [ ] Subtitle size control
  - [ ] Embed control
- [x] Authentication
  - [x] Passwordless authentication
  - [x] Patreon account linking
- [x] Ads and statistic integration support
- [ ] Tools
  - [ ] Torrent => DDL
  - [ ] Magnet => DDL
  - [ ] Magnet => Torrent
- [ ] Misc
  - [ ] Feedback form
  - [ ] Allow magnet-url as query string
- [x] Chrome extension integration
- [x] Embed support
  - [x] Base version
  - [x] Extended version
- [ ] 🚀Switch webtor.io to web-ui-v2


## Setting up connection to Webtor RestAPI

You have to set up connection to [Webtor RestAPI](https://github.com/webtor-io/rest-api) before using WebUI.

If you have already installed [backend part](https://github.com/webtor-io/helm-charts),
then you have to proxy rest-api from your k8s instance to your local machine with [kubefwd](https://github.com/txn2/kubefwd):
```
sudo kubefwd svc -f metadata.name=rest-api -n webtor
```
or with [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl):
```shell
kubectl port-forward svc/rest-api 9090:80 -n webtor

# you have to setup additional environment variables before starting application
export REST_API_SERVICE_PORT=9090
export REST_API_SERVICE_HOST=localhost
```

If you have [RapidAPI subscription](https://rapidapi.com/paveltatarsky-Dx4aX7s_XBt/api/webtor/)
you can just do the following:

```shell
export RAPIDAPI_KEY={YOUR_RAPIDAPI_KEY}
export RAPIDAPI_HOST={YOUR_RAPIDAPI_HOST}
```

## Usage

```shell
./web-ui-v2 help
NAME:
   web-ui-v2 - runs webtor web ui v2

USAGE:
   web-ui-v2 [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   serve, s  Serves web server
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Development

```
npm install
npm start
```

## Building

```shell
make build
```