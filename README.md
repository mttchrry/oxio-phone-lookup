# Phone-Number formatter

## Prerequisites
Need go installed, and get dependencies with `go get ./...` from the root dir

## To Build
build with `go build -o ./server ./cmd/main.go` from the root dir

## To Run
run `./server` from the root after building

## To Test

### Manual
Listing on `http://localhost:8081/v1/phone-numbers`, can send get commands with JSON body of the format
```
{
    "phoneNumber": "971 419 508 9397",
    "countryCode": "AE"
}
```

for example, with a header `"Content-Type": "application/json"`

### Unit tests
from the root, simply run `go test ./...`

## Explanations

### Tech Stack
I am most familiar with Golang and is seems suited to the task.  I found a library for country codes to be able to translate back and forth without maintaining my own DB or data structure so that seemed like the best approach to no re-invent the wheel.  

For checking that the number was validly formated, I used a regex as it vastly simplifies a somewhat complicated string parsing.  

### Deploying
We would set up a kubernetes and docker framework and deploy to a cloud platform like GCP or AWS.  Would need more infra to keep it from being DDOSed though, such as rate limiting on IP or having a sign-up/login workflow

### Assumptions and improvements
I am assuming that for the country number, if there isn't an associated country code passed in, returning any two letter country code that matches up (i.e. returning CA or US are both valid for '+1') is adaquate.  But for productization we'd likely want to find a better library or way to look this up. 

I also took the liberty of allowing parenthesis or dashes along with spaces where they are typically expected.  I figured this would be the next feature improvement and through them in as options in the Regex. 

This solution is wanting for a better way to log errors like Loggly, and standard different kinds of errors and returning 400s/500s and other standardized error codes I didn't get to in 2 hours.

The testing is also not super thorough, especially the logic around matching country codes to their numeric equivalent and vice versa. 