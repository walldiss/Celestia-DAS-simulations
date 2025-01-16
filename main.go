package main

import (
	"log"
	"math/rand"
)

// Sample represents a single point in the data square
type Sample struct {
	Row, Col int
}

// SampleSet maintains a collection of unique samples
type SampleSet struct {
	samples map[Sample]bool
}

// NewSampleSet creates a new initialized SampleSet
func NewSampleSet(capacity int) *SampleSet {
	return &SampleSet{
		samples: make(map[Sample]bool, capacity),
	}
}

// Clear removes all samples from the set
func (s *SampleSet) Clear() {
	clear(s.samples)
}

// FillUnique adds n unique random samples within the given size bounds
func (s *SampleSet) FillUnique(n, size int) {
	for n > 0 {
		row := rand.Intn(size * 2)
		col := rand.Intn(size * 2)
		sample := Sample{Row: row, Col: col}

		if !s.samples[sample] {
			s.samples[sample] = true
			n--
		}
	}
}

// DataSquare represents the main data structure for the recovery simulation
type DataSquare struct {
	Size          int
	Matrix        [][]int
	RowCounts     []int
	ColCounts     []int
	RecoveredRows map[int]bool
	RecoveredCols map[int]bool
	TotalCount    int
}

// NewDataSquare creates a new initialized DataSquare
func NewDataSquare(size int) *DataSquare {
	matrix := make([][]int, 2*size)
	for i := range matrix {
		matrix[i] = make([]int, 2*size)
	}

	return &DataSquare{
		Size:          size,
		Matrix:        matrix,
		RecoveredRows: make(map[int]bool),
		RecoveredCols: make(map[int]bool),
	}
}

// Reset clears all data in the DataSquare
func (ds *DataSquare) Reset() {
	ds.RowCounts = make([]int, ds.Size*2)
	ds.ColCounts = make([]int, ds.Size*2)
	clear(ds.RecoveredRows)
	clear(ds.RecoveredCols)
	ds.TotalCount = 0

	for i := range ds.Matrix {
		for j := range ds.Matrix[i] {
			ds.Matrix[i][j] = 0
		}
	}
}

// AddSamples adds all samples from the given set to the DataSquare
func (ds *DataSquare) AddSamples(samples *SampleSet) {
	for s := range samples.samples {
		if ds.Matrix[s.Row][s.Col] == 0 {
			ds.AddSample(s.Row, s.Col)
		}
	}
}

// AddSample adds a single sample to the DataSquare
func (ds *DataSquare) AddSample(row, col int) bool {
	if ds.Matrix[row][col] > 0 {
		return false
	}

	ds.Matrix[row][col] = 1
	ds.RowCounts[row]++
	ds.ColCounts[col]++
	ds.TotalCount++
	return true
}

// TryRecoverRow attempts to recover a row if it meets the criteria
func (ds *DataSquare) TryRecoverRow(row int) bool {
	if ds.RecoveredRows[row] {
		return false
	}

	if ds.RowCounts[row] >= ds.Size {
		ds.RecoveredRows[row] = true
		for col := range ds.Matrix[row] {
			if ds.AddSample(row, col) {
				ds.TryRecoverCol(col)
			}
		}
		return true
	}
	return false
}

// TryRecoverCol attempts to recover a column if it meets the criteria
func (ds *DataSquare) TryRecoverCol(col int) bool {
	if ds.RecoveredCols[col] {
		return false
	}

	if ds.ColCounts[col] >= ds.Size {
		ds.RecoveredCols[col] = true
		for row := range ds.Matrix {
			if ds.AddSample(row, col) {
				ds.TryRecoverRow(row)
			}
		}
		return true
	}
	return false
}

// IsRecovered checks if the DataSquare is fully recovered
func (ds *DataSquare) IsRecovered() bool {
	return len(ds.RecoveredRows) >= ds.Size || len(ds.RecoveredCols) >= ds.Size
}

// Recover attempts to recover the entire DataSquare
func (ds *DataSquare) Recover() bool {
	if ds.TotalCount < ds.Size*ds.Size {
		return false
	}

	for {
		var rowRecovered, colRecovered bool
		for i := 0; i < ds.Size*2; i++ {
			rowRecovered = ds.TryRecoverRow(i) || rowRecovered
			colRecovered = ds.TryRecoverCol(i) || colRecovered
		}

		if ds.IsRecovered() {
			return true
		}
		if !rowRecovered && !colRecovered {
			return false
		}
	}
}

// SimulationConfig holds the configuration for running simulations
type SimulationConfig struct {
	// SamplesPerIteration is the number of unique samples to generate in each iteration
	// This represents how many points we try to recover in each step
	SamplesPerIteration int

	// Iterations is the number of times to run each simulation scenario
	// Higher values provide more accurate probability estimates but take longer to run
	Iterations int

	// InitialLights is the starting number of light sources for the simulation
	// This value may be overridden by LightsAt16 calculation
	InitialLights int

	// LightsAt16 is used to calculate InitialLights for different grid sizes
	// If non-zero, InitialLights is scaled proportionally to the grid size
	// Formula: InitialLights = LightsAt16 * (currentSize^2) / (16^2)
	LightsAt16 int

	// SizeIterFactor determines how much to increment the number of lights
	// in each iteration. The increment is calculated as: size / SizeIterFactor
	SizeIterFactor int

	// InitialSize is the starting size for the data square
	// The actual grid will be 2x this size in both dimensions
	InitialSize int

	// MaxSize is the largest size to test
	// The simulation will double the size until reaching this value
	MaxSize int

	// TargetProbability is the success rate we want to achieve
	// Once this probability is reached, we move to the next size
	// Value should be between 0 and 1 (e.g., 0.99 for 99%)
	TargetProbability float64
}

// NewDefaultConfig creates a SimulationConfig with default values
func NewDefaultConfig() *SimulationConfig {
	return &SimulationConfig{
		SamplesPerIteration: 16,
		Iterations:          1000,
		LightsAt16:          10,
		InitialLights:       7500,
		SizeIterFactor:      16,
		InitialSize:         16,
		MaxSize:             256,
		TargetProbability:   0.99,
	}
}

// RunSimulation executes the main simulation with the given configuration
func RunSimulation(config *SimulationConfig) {
	log.Printf("Starting simulation with target probability: %.2f%%\n", config.TargetProbability*100)

	for size := config.InitialSize; size <= config.MaxSize; size *= 2 {
		log.Printf("\nProcessing size: %d x %d\n", size*2, size*2)

		ds := NewDataSquare(size)
		samples := NewSampleSet(config.SamplesPerIteration)

		initialLights := config.InitialLights
		if config.LightsAt16 != 0 {
			initialLights = config.LightsAt16 * (size * size) / (16 * 16)
		}

		log.Printf("Initial lights: %d\n", initialLights)

		for lights := initialLights; ; lights += size / config.SizeIterFactor {
			successCount := 0

			for i := 0; i < config.Iterations; i++ {
				ds.Reset()

				for n := 0; n < lights; n++ {
					samples.FillUnique(config.SamplesPerIteration, size)
					ds.AddSamples(samples)
					samples.Clear()
				}

				if ds.Recover() {
					successCount++
				}
			}

			probability := float64(successCount) / float64(config.Iterations)
			log.Printf("Lights: %d, Success Rate: %.2f%% (%d/%d)\n",
				lights,
				probability*100,
				successCount,
				config.Iterations)

			if probability >= config.TargetProbability {
				log.Printf("Target probability reached for size %d with %d lights\n", size, lights)
				break
			}
		}
	}
}

func main() {
	rand.Seed(1)
	config := NewDefaultConfig()
	RunSimulation(config)
}
