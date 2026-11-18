package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Primoz200/ps_naloga5/vreme"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "ps_naloga5",
		Usage: "izpis meritev vremenskih postaj",
		Flags: []cli.Flag{
			&cli.Float64Flag{Name: "minTemp", Usage: "najnižja dovoljena temperatura meritve", Value: -50},
			&cli.Float64Flag{Name: "maxTemp", Usage: "najvišja dovoljena temperatura meritve", Value: 50},
			&cli.IntFlag{Name: "stMeritev", Usage: "najmanjše število meritev za statistiko", Value: 3},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			minTemp, maxTemp, stMeritev := cmd.Float64("minTemp"), cmd.Float64("maxTemp"), cmd.Int("stMeritev")
			if minTemp > maxTemp {
				return fmt.Errorf("minTemp ne sme biti večji od maxTemp")
			}
			if stMeritev < 1 {
				return fmt.Errorf("stMeritev mora biti vsaj 1")
			}

			postaje := make(map[string]vreme.Postaja)
			vreme.DodajPostajo(postaje, "LJ01", "sentvid", "ljubljana")
			vreme.DodajPostajo(postaje, "LJ02", "fri", "ljubljana")
			vreme.DodajPostajo(postaje, "KR01", "center", "kranj")
			vreme.DodajPostajo(postaje, "MB01", "center", "maribor")

			začetne := map[string][]float64{
				"LJ01": {10, 15, 12, 19}, "LJ02": {20, 22, 14, 14},
				"KR01": {18, 10, -2, 5}, "MB01": {18, 10},
			}
			for oznaka, meritve := range začetne {
				for _, temperatura := range meritve {
					if err := vreme.DodajMeritev(postaje, oznaka, temperatura, minTemp, maxTemp); err != nil {
						return err
					}
				}
			}
			if err := vreme.DodajMeritev(postaje, "KR01", 25, minTemp, maxTemp); err != nil {
				fmt.Fprintln(cmd.ErrWriter, err)
			}

			vreme.IzpisMeritev(cmd.Writer, postaje)
			fmt.Fprintln(cmd.Writer)
			vreme.IzpisPorocila(cmd.Writer, postaje, stMeritev)
			return nil
		},
	}
	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "Napaka:", err)
		os.Exit(1)
	}
}
