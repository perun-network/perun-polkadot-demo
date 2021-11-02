// Copyright 2021 - See NOTICE file for copyright holders.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package demo

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/pkg/errors"
)

type run struct {
	// data are the recorded times for sendUpdate in nanoseconds
	data []float64
	// start of the last run
	start time.Time
}

func (r *run) Start() {
	r.start = time.Now()
}

func (r *run) Stop() {
	r.data = append(r.data, float64(time.Since(r.start).Nanoseconds()))
}

func (r *run) String() string {
	functions := []func(stats.Float64Data) (float64, error){stats.Sum, stats.Min, stats.Max, stats.Median, stats.StdDevP}
	var str string
	values := make([]float64, len(functions))

	for i, f := range functions {
		values[i], _ = f(r.data)
		str += (time.Duration(values[i]) * time.Nanosecond).Round(time.Microsecond).String() + "\t"
	}

	freq := (float64(len(r.data)) / values[0]) * float64(time.Second.Nanoseconds())
	return fmt.Sprintf("N\ttx/s\tSum\tMin\tMax\tMedian\tStddev\t\n%d\t%.1f\t", len(r.data), freq) + str
}

// Benchmark updates the channel with a `peer` `n` times and measures the of every update.
// A statistic is then printed with run.String()
func (n *node) Benchmark(args []string) error {
	n.mtx.Lock()
	defer n.mtx.Unlock()
	peer := n.peers[args[0]]
	totalAmountDot, _ := strconv.Atoi(args[1])
	txCount, _ := strconv.Atoi(args[2])
	var r run

	if txCount < 1 {
		return errors.New("Number of runs cant be less than 1")
	} else if peer == nil {
		return errors.New("Peer not found")
	} else if peer.ch == nil {
		return errors.New("Open a state channel first")
	}

	totalAmountPlank := dotToPlank(big.NewFloat(float64(totalAmountDot)))[0]
	txAmount := new(big.Int).Div(totalAmountPlank, big.NewInt(int64(txCount)))
	for i := 0; i < txCount; i++ {
		r.Start()
		if err := peer.ch.sendMoney(txAmount); err != nil {
			return errors.WithMessage(err, "could not send update")
		}
		r.Stop()
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
	fmt.Fprintln(w, r.String())
	return w.Flush()
}
