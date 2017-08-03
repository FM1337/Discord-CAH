When adding cards use the following syntax for the json:

Whitecards:

[{"CardText":"Text here"}, {"CardText":"Another whitecard"}]

Black Cards:

[{"CardText":"Text here _","Cards2Play":1}, {"CardText":"No underscores here.","Cards2Play":0}]

Note if there are no whitespaces just set Card2Play to 0 and the importer will handle it.

Make sure to put your black cards under /custom/BlackCards and white cards under /custom/WhiteCards and make sure the card's filename ends with .json

