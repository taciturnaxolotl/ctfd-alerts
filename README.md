# ⛳ CTFd alerts

Sends alerts for any arbitrary [CTFd](https://ctfd.io/) instance via [ntfy](https://ntfy.sh/)

## Install

You can download a pre-built binary from the releases or you can use the following options

### Go

```bash
# Go
go install github.com/taciturnaxolotl/ctfd-alerts@latest
```

### Nix

```bash
# Direct installation with flakes enabled
nix profile install github:taciturnaxolotl/ctfd-alerts
```

For use in your own flake:

```nix
# In your flake.nix
{
  inputs.akami.url = "github:taciturnaxolotl/ctfd-alerts";

  outputs = { self, nixpkgs, akami, ... }: {
    # Access the package as:
    # ctfd-alerts.packages.${system}.default
  };
}
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
