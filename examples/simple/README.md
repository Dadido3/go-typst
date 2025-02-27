# Simple example

This example shows how to render Typst documents directly from strings in Go.

## The pros and cons

The main advantage of this method is that it's really easy to set up.
In the most simple case you build your Typst markup by concatenating strings, or by using `fmt.Sprintf`.

The downside is that the final Typst markup is only generated on demand.
This means that you can't easily use the existing Typst tooling to write, update or debug your Typst markup.
Especially as your your documents get more complex, you should switch to other methods.
