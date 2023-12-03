# clean-links

_A small experiment to test out [Golang](https://go.dev/) myself._

An application that:

1. Walks through all HTML files
1. Looks for `<a>` elements
1. Add `rel="noreferrer"` attribute
   - [`rel="noopener"` is already implied](https://github.com/jsx-eslint/eslint-plugin-react/issues/2022) for most modern browsers

Future improvements:

- [ ] Customize `rel` values
- [ ] Handle other elements
  - See <https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Referrer-Policy#integration_with_html>
- [ ] Customize the `referrerpolicy` values
- [ ] Exclude elements based on condition
