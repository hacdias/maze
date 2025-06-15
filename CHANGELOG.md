# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

### Added

### Changed

### Deprecated

### Removed

### Fixed

### Security

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
