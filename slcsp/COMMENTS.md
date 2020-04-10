# Assumptions

I found the requirements description of the model a bit confusing and tried to get clarity from Ad Hoc but was not able to get a response to my question which was "I'm working on SLCSP and the requirements mention counties. But you can only look up plans by state and rate area so from my understanding the counties don't come into the solution at all. Is that correct?"
If I had a product person around I would whiteboard this whole model till I accurately understood Plans and their Areas.
Meanwhile I'm assuming from the data that given a zip code I can look up the state and rate area tuple and with that look up all relevant plans. And that if that zip matches multiple rate areas the result is ambiguous.

# Notes for compiling and running SLCSP

See requirements in data/README.md

## Intall Go 1.14 (may work with earlier versions)
https://golang.org/doc/install

## Run
```
go run slcsp.go
```
