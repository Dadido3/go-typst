#let template(data, customText) = {
  set page(paper: "a4")

  [= Example]

  customText

  [== List of items]

  show table.cell.where(y: 0): strong
  set table(
    stroke: (x, y) => if y == 0 {
      (bottom: 0.7pt + black)
    },
  )

  table(
    columns: 5,
    table.header(
      [Name],
      [Size],
      [Example box],
      [Created],
      [Numbers],
    ),
    ..for value in data {
      (
        [#value.Name],
        [#value.Size],
        box(fill: black, width: 0.1mm * value.Size.X, height: 0.1mm * value.Size.Y),
        value.Created.display(),
        [#list(..for num in value.Numbers {([#num],)})],
      )
    }
  )
}
