# shrt

An alternative command line tool to shorten URLs using the Shlink URL shortening service.

`shrt` allows the configuration of a self-hosted instance for you to use.

## Installation

Download the binary for your platform from [here](https://ithub.com/arenzana/shrt/releases).

## Usage

`shrt` allows listing all shortened URLs and shortening one or more URLs.

```bash
An alternative Shlink service application to interact with a Shlink instance without having to install the official client and server.

Usage:
  shrt [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  list        Lists all the URLs shortened in the configured Shlink service
  short       Shortens a long URL using the Shlink API

Flags:
      --config string   config file (default is $HOME/.shrt.yaml)
  -h, --help            help for shrt
  -q, --quiet           config file (default is $HOME/.shrt.yaml) (default true)
  -t, --toggle          Help message for toggle
  -v, --version         version for shrt

Use "shrt [command] --help" for more information about a command.
```

## Examples

List all shortened URLs on an instance:

```bash
./shrt ls                                                    (base)
CREATED         SHORT URL                       LONG URL                                                          VISITS
2023-11-03      https://l.isma.rip/lolz         https://pkg.go.dev/golang.design/x/clipboard                      1
2023-11-03      https://l.isma.rip/p3aMM        https://www.themoviedb.org/tv/62317                               2
2023-11-03      https://l.isma.rip/VKeeU        https://appiculous.com                                            0
2023-11-02      https://l.isma.rip/rsCWY        https://www.macrumors.com/2023/11/02/ios-17-1-1-is-coming/        0
2023-11-02      https://l.isma.rip/YBmzA        https://x.com/ryanair/status/1720046153444610100                  2
```

or shorten a URL (copies shortened URL to clipboard automatically):

```bash
./shrt s https://tailscale.com/changelog/ -g tail           (base)
https://l.isma.rip/tail
```

The `-g` option allows the user to set a custom slug. Otherwise Shlink will generate one automatically.

The `-q` option will toggle quiet mode to know what `shrt` is doing.
