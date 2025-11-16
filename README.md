# Simrail Side (Go Rewrite)

This folder contains a Go rewrite of the Simrail Side project. All logic is ported from the original PHP files, but the originals remain unchanged.

## Structure
- `main.go`: Entry point, starts HTTP server
- `api.go`: API endpoints for layout and train data
- `update.go`: Update logic for fetching and processing data
- `convert.go`: Data conversion logic (to be expanded)
- `config.go`: Configuration constants
- `utils.go`: Utility functions

## Usage
1. Install Go (https://golang.org/dl/)
2. Run `go run main.go` in this folder
3. Access endpoints:
   - `/api?type=layout&layout=0`
   - `/api?type=trains&server=en1&layout=0`
   - `/update?task=1`

## Notes
- Data files are read from `../files/` and `../layouts/`.
- Extend `convert.go` for custom train/station/timetable logic.
