# telemoney

Bot that provides tg interface for personal accounting and use GoogleSheets as DB.

## Questions
Sheets:
* What date format is used by Sheets?
* Does Shee

Functions:
* What categories should be defined? How?
* How to group tags?


`(?P<amount>\d+[\.,]?\d*) (?P<category>\w*) (?:\((?P<tags>[\w, ]*)\))?(?P<comment>.*$)?`
to parse "9,5 lunch (grenka, dumplings) I need foood!"