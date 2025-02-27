// Prepare data that will be used as preview.
#let data = (
  (Name: "Bell", Size: (X: 80, Y: 40, Z: 40), Created: datetime(year: 2010, month: 12, day: 1, hour: 12, minute: 13, second: 14), Numbers: (1, 2, 3)),
  (Name: "Bell", Size: (X: 80, Y: 40, Z: 40), Created: datetime(year: 2010, month: 12, day: 1, hour: 12, minute: 13, second: 14), Numbers: (10, 12, 15)),
  (Name: "Bell", Size: (X: 80, Y: 40, Z: 40), Created: datetime(year: 2010, month: 12, day: 1, hour: 12, minute: 13, second: 14), Numbers: (100, 109, 199, 200)),
)

#let customText = "Hey, this is example data to test the template."

// Invoke the template with the preview data and replace this whole document with the result.
#import "template.typ": template
#show: doc => template(data, customText)
