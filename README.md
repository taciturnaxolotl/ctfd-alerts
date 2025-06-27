# ⛳ CTFd alerts

Sends alerts for any arbitrary [CTFd](https://ctfd.io/) instance via [ntfy](https://ntfy.sh/)

![vhs gif of the command being run](https://github.com/taciturnaxolotl/ctfd-alerts/blob/main/.github/images/out.gif?raw=true)

## Install

You can download a pre-built binary from the releases or you can use the following options

### Go

```bash
# Go
go install github.com/taciturnaxolotl/ctfd-alerts@latest
```

If you need a systemd service file there is one in `ctfd-alerts.service`

### Nix

```bash
# Direct installation with flakes enabled
nix profile install github:taciturnaxolotl/ctfd-alerts
```

For use in your own flake:

```nix
# In your flake.nix
{
  inputs.ctfd-alerts.url = "github:taciturnaxolotl/ctfd-alerts";

  outputs = { self, nixpkgs, ctfd-alerts, ... }: {
    # Access the package as:
    # ctfd-alerts.packages.${system}.default
  };
}
```

## Config

The config for the bot is quite simple. Create a `config.toml` file in the same directory as the binary (or link to the config location with `-c ./path/to/config/config.toml`) with the following format:

```toml
debug = true
interval = 100 # defaults to 300 if unset
user = "echo_kieran"

[ctfd]
api_base = "http://163.11.237.79/api/v1"
api_key = "ctfd_10698fd44950bf7556bc3f5e1012832dae5bddcffb1fe82191d8dd3be3641393"

[ntfy]
api_base = "https://ntfy.sh/"
acess_token = ""
topic = "youralert"
```

Written in go. If you have any suggestions or issues feel free to open an issue on my [tangled](https://tangled.sh/@dunkirk.sh/ctfd-alerts) knot

<p align="center">
	<img src="https://raw.githubusercontent.com/taciturnaxolotl/carriage/master/.github/images/line-break.svg" />
</p>

<p align="center">
	<i><code>&copy 2025-present <a href="https://github.com/taciturnaxolotl">Kieran Klukas</a></code></i>
</p>

<p align="center">
	<a href="https://github.com/taciturnaxolotl/ctfd-alerts/blob/master/LICENSE.md"><img src="https://img.shields.io/static/v1.svg?style=for-the-badge&label=License&message=MIT&logoColor=d9e0ee&colorA=363a4f&colorB=b7bdf8"/></a>
</p>
