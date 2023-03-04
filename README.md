# <img src="https://what3words.com/assets/images/w3w_square_red.png" width="64" height="64" alt="what3words">&nbsp;What3Words Go Wrapper
[![Go version](https://img.shields.io/github/go-mod/go-version/henrwal/w3w-go-wrapper/main?label=Go%20Version)](https://github.com/henrwal/w3w-go-wrapper)
[![GoReportCard Badge](https://goreportcard.com/badge/github.com/henrwal/w3w-go-wrapper)](https://goreportcard.com/report/github.com/henrwal/w3w-go-wrapper)

A Go module that wraps the [what3words REST API](https://docs.what3words.com/api/v3/).

# Overview

The what3words Go module gives you programmatic access to:
* Convert a 3 word address to coordinates.
* Convert coordinates to a 3 word address.
* Auto-suggest functionality which takes a slightly incorrect 3 word address, and suggests a list of valid 3 word addresses.
* Obtain a section of the 3m x 3m what3words grid for a bounding box.
* Determine the languages currently supported by what3words.

# Getting Started

## Authentication

To use this module you’ll need an API key, please visit [https://what3words.com/select-plan](https://what3words.com/select-plan) and sign up for an account.

# Usage

Quick example:
```go
import w3w "github.com/henrwal/w3w-go-wrapper"

func main() {
    key, ok := os.LookupEnv("W3W_API_KEY")
    if !ok {
		panic("env not set")
    }
	
	w := w3w.NewClient(key)
    coordinates, _ := w.ConvertToCoordinates("filled.count.soap")
	fmt.Println(coordinates)
}
```

## Convert To Coordinates

This function takes the words parameter as a string of 3 words `'table.book.chair'`

The returned payload from the `convert-to-coordinates` method is described in the [what3words REST API documentation](https://docs.what3words.com/api/v3/#convert-to-coordinates).

## Convert To 3 Word Address

This function takes the latitude and longitude:
- 2 parameters:  `lat=0.1234`, `lng=1.5678`

The returned payload from the `convert-to-3wa` method is described in the [what3words REST API documentation](https://docs.what3words.com/api/v3/#convert-to-3wa).


## AutoSuggest

Returns a list of 3 word addresses based on user input and other parameters.

This method provides corrections for the following types of input error:
* typing errors
* spelling errors
* misremembered words (e.g. singular vs. plural)
* words in the wrong order

The `autosuggest` method determines possible corrections to the supplied 3 word address string based on the probability of the input errors listed above and returns a ranked list of suggestions. This method can also take into consideration the geographic proximity of possible corrections to a given location to further improve the suggestions returned.

### Input 3 word address

You will only receive results back if the partial 3 word address string you submit contains the first two words and at least the first character of the third word; otherwise an error message will be returned.

### Clipping and Focus

We provide various `clip` policies to allow you to specify a geographic area that is used to exclude results that are not likely to be relevant to your users. We recommend that you use the `clip` parameter to give a more targeted, shorter set of results to your user. If you know your user’s current location, we also strongly recommend that you use the `focus` to return results which are likely to be more relevant.

In summary, the `clip` policy is used to optionally restrict the list of candidate AutoSuggest results, after which, if focus has been supplied, this will be used to rank the results in order of relevancy to the focus.

https://docs.what3words.com/api/v3/#autosuggest

The returned payload from the `autosuggest` method is described in the [what3words REST API documentation](https://docs.what3words.com/api/v3/#autosuggest).

## Grid Section

Returns a section of the 3m x 3m what3words grid for a bounding box.

## Available Languages

Retrieves a list of the currently loaded and available 3 word address languages.

The returned payload from the `available-languages` method is described in the [what3words REST API documentation](https://docs.what3words.com/api/v3/#available-languages).

## Code examples

### Convert to coordinates
```go
import w3w "github.com/henrwal/w3w-go-wrapper"

func main() {
	w := w3w.NewClient("SECRET_API_KEY")
	words := "filled.count.soap"
    coordinates, _ := w.ConvertToCoordinates(context.Background(), words)
	fmt.Println(coordinates)
}
```

### Convert to 3 word address
```go
import w3w "github.com/henrwal/w3w-go-wrapper"
    
func main() {
    coordinates := Coordinates{
		Lat: 51.520847,
		Lng: -0.195521,
	}
	locationResponse, _ := w.ConvertTo3wa(context.Background(), coordinates)
	fmt.Println(locationResponse)
}
```

## Issues

Find a bug or want to request a new feature? Please let me know by submitting an issue.

If any issues/bugs are found please raise an issue on GitHub here: [Issue Tracker](https://github.com/henrwal/w3w-go-wrapper/issues)