# PubGolf.co

## Setup

Bootstrap the dev tools by [installing the latest version of Go](https://go.dev/dl/) and running the following:

```sh
go install ./tools/cmd/pubgolf-devctrl && pubgolf-devctrl update
```

After this the `pubgolf-devctrl` binary will be available in your $PATH to run development tasks. Run the following to install all dev dependencies:

```sh
pubgolf-devctrl install
```
