package main

import (
  "encoding/csv"
  "fmt"
  "io"
  "strconv"
)

// statsFunc defines a generic statistical function
type statsFunc func(data []float64) float64

func sum(data []float64) float64 {
  sum := 0.0

  for _, v := range data {
    sum += v
  }

  return sum
}

func avg(data []float64) float64 {
  return sum(data) / float64(len(data))
}

func csv2float(r io.Reader, column int) ([]float64, error) {
  // Create the CSV Reader used to read in data from CSV files
  cr := csv.NewReader(r)
  cr.ReuseRecord = true

  // Adjusting for 0 based index
  column--

  var data []float64

  // Looping through all records
  for i := 0; ; i++ {
    row, err := cr.Read()
    if err == io.EOF {
      break
    }

    if err != nil {
      return nil, fmt.Errorf("Cannot read data from file: %s", err)
    }

    if i == 0 {
      continue
    }

    // Checking number of columns in CSV file
    if len(row) <= column {
      // File does not have that many columns
      return nil,
        fmt.Errorf("Invalid column #. File has only %d columns", len(row))
    }

    // Try to convert data read into a float number
    v, err := strconv.ParseFloat(row[column], 64)
    if err != nil {
      return nil, fmt.Errorf("Data is not numeric: %s", err)
    }

    data = append(data, v)
  }

  // Return the slice of float64 and nil error
  return data, nil
}
