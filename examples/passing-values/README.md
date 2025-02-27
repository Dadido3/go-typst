# Passing values example

This example demonstrates how to pass values to Typst, which can be useful in rendering custom documents such as reports, invoices, and more.

## How it works

This example follows the [template pattern](https://typst.app/docs/tutorial/making-a-template/) described in the Typst documentation.
Here is a short overview of the files:

- [template.typ](template.typ) defines a Typst template function that constructs a document based on parameters.
- [main.go](main.go) shows how to convert/encode Go values into Typst markup, and how to call/render the template with these converted values.
- [template-preview.typ](template-preview.typ) also invokes the template while providing mock data.
    This is useful when you want to preview, update or debug the template.
