# clean-links

_A small experiment to test out [Golang](https://go.dev/) myself._

A CLI application that:

1. Accepts one or more paths
1. Walks through all HTML files
1. Looks for `<a>` elements
1. Add `referrerpolicy="noreferrer"` attribute
   - `referrerpolicy="noreferrer"` will imply `referrerpolicy="noopener"` for most modern browsers ([ref](https://github.com/jsx-eslint/eslint-plugin-react/issues/2022))

Future improvements:

- [x] ~~Customize `rel` values~~
  - Moved to `referrerpolicy`
- [x] Handle other elements
  - See <https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Referrer-Policy#integration_with_html>
- [x] Customize the `referrerpolicy` values
- [ ] Exclude elements based on condition
  - [x] Exclude based on `class`

## Usage

1. Download the TGZ package from [releases](https://github.com/altbdoor/clean-links/releases)
1. Unpack the TGZ file
   - `tar xzf clean-links-*.tgz`
1. For Linux based distribution, just call `./clean-links`
1. For Windows based distribution, just call `./clean-links.exe`
1. For MacOS based distribution:
   1. Get the binary out of quarantine with `xattr -d com.apple.quarantine clean-links`
   1. Just call `./clean-links`
