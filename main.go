package main

import (
	"fmt"
	"os"
	"time"

	"github.com/tsuzu/wautorle/runner"
	"github.com/tsuzu/wautorle/wordle"
)

func main() {
	a, err := wordle.NewAutomator()

	if err != nil {
		panic(err)
	}

	defer a.Close()

	var statsFile string

	if len(os.Args) > 1 {
		statsFile = os.Args[1]

		stats, err := os.ReadFile(statsFile)

		if err != nil {
			panic(err)
		}

		if err := a.SetStateStats(string(stats)); err != nil {
			panic(err)
		}
	}

	r, err := runner.New()

	if err != nil {
		panic(err)
	}

	defer r.Close()

	idx := 0
	for {
		func() {
			word, err := r.NextWord()

			if err != nil {
				panic(err)
			}
			fmt.Println("entering", word)

			if err := a.Enter(word); err != nil {
				panic(err)
			}
			time.Sleep(3 * time.Second)

			line, err := a.Line(idx)

			if err != nil {
				panic(err)
			}
			fmt.Println(line)

			result := runner.ParseResult(line)

			if result.Finished() {
				for {
					finished, err := a.Finished()

					if err != nil {
						panic(err)
					}

					if finished {
						break
					}
					time.Sleep(1 * time.Second)
				}
				fmt.Println("finished")

				time.Sleep(1 * time.Second)
				result, err := a.CopyResult()

				if err != nil {
					panic(err)
				}
				fmt.Println(result)

				ss, err := a.Screenshot()

				if err != nil {
					panic(err)
				}

				if err := os.WriteFile(os.Args[2], []byte(result), 0774); err != nil {
					panic(err)
				}

				if err := os.WriteFile(os.Args[3], ss, 0774); err != nil {
					panic(err)
				}

				if statsFile != "" {
					stats, err := a.GetStateStats()

					if err != nil {
						panic(err)
					}

					if err := os.WriteFile(statsFile, []byte(stats), 0644); err != nil {
						panic(err)
					}
				}

				os.Exit(0)
			}

			if err := r.WriteResult(result); err != nil {
				panic(err)
			}

			idx++
		}()
	}

}
