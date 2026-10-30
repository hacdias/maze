# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

### Added

### Changed

- Upgraded to Go 1.27.
- `NewMaze(nil)` now uses an HTTP client with a 10 second timeout, instead of none.

### Deprecated

### Removed

- Dependency on `github.com/paulmach/go.geojson`.

### Fixed

- `Search` now fills `CountryCode` and `PostalCode`, like `Reverse`.
- HTTP errors from Photon and aviowiki are now reported as such, instead of as missing results.
- `Reverse` and `Search` no longer fail when `lang` is empty.
- `Location.String` and `ParseLocation` now keep `CountryCode` and `PostalCode`.
- `ReverseGeoURI` now keeps a `?name=` query parameter, and `ParseLocation` now reads a `;name=` geo URI parameter.

### Security

## [0.6.1]

### Changed

- Preserve geo URI `name` when calling `ReverseGeoURI`.

## [0.6.0]

### Added

- Added `CountryCode` and `PostalCode` to the `Location` struct.

## [0.5.0]

### Added

- Adding a geo URI representation for the `Location` struct, which can be used to store it as a string.

## [0.4.0]

### Added

- When searching for an airport, the fields `IATA` and `ICAO` will be filled.

### Changed

- When searching for an airport, the name of the airport will be used in the returned `Location`, instead of the search parameter.

## [0.3.0]

### Changed

- Upgraded to Go 1.23.
