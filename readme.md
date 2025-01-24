# Data Availability Sampling Simulator

A Monte Carlo simulation tool for analyzing Data Availability Sampling (DAS) reconstruction probabilities in blockchain systems. This implementation was used to derive the empirical results discussed in the research on light node requirements for Celestia's DAS mechanism.

## Purpose

This simulator helps determine the minimum number of light nodes needed for successful data reconstruction with a specified probability (default 99%). It validates that fewer samples than theoretically predicted (k(3k-2)) can achieve reliable reconstruction.

## Key Features

- Simulates data reconstruction using random sampling
- Supports variable matrix sizes (16x16 to 256x256)
- Configurable parameters for samples per node and iteration count
- Calculates success probability across different configurations

## Usage

```go
config := NewDefaultConfig()
RunSimulation(config)
```

### Configuration Parameters

- `SamplesPerIteration`: Number of samples per light node (default: 16)
- `Iterations`: Number of Monte Carlo iterations (default: 1000)
- `InitialSize`: Starting matrix size k (default: 16)
- `MaxSize`: Maximum matrix size k (default: 256)
- `TargetProbability`: Required success rate (default: 0.99)

## Key Findings

The simulation helped establish that:
- Reconstruction is possible with ~2.5-2.8x fewer light nodes than originally estimated
- The relationship between parameters follows: (n*s)/(2k)² ≈ 0.6
    - n: number of light nodes
    - s: samples per node
    - k: matrix size

## References

- Original Research Paper: [Data Availability Proofs](http://arxiv.org/abs/1809.09044)
- Full Analysis: [Celestia DAS: Reconstruction Simulations and Light Node Requirements](https://forum.celestia.org/t/celestia-das-reconstruction-simulations-and-light-node-requirements/1891)
