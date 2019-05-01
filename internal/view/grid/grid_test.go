package grid_test

import (
	"testing"

	"github.com/slok/grafterm/internal/model"
	"github.com/slok/grafterm/internal/view/grid"
	"github.com/stretchr/testify/assert"
)

func TestGrid(t *testing.T) {
	tests := []struct {
		name   string
		grid   func() (*grid.Grid, error)
		exp    *grid.Grid
		expErr bool
	}{
		{
			name: "Empty Adaptive grid.",
			grid: func() (*grid.Grid, error) {
				widgets := []model.Widget{}
				maxWidth := 100

				return grid.NewAdaptiveGrid(maxWidth, widgets)
			},
			exp: &grid.Grid{
				MaxWidth: 100,
				Rows:     []*grid.Row{},
			},
			expErr: false,
		},
		{
			name: "Empty Fixed grid.",
			grid: func() (*grid.Grid, error) {
				widgets := []model.Widget{}
				maxWidth := 100

				return grid.NewFixedGrid(maxWidth, widgets)
			},
			exp: &grid.Grid{
				MaxWidth: 100,
				Rows:     []*grid.Row{},
			},
			expErr: false,
		},
		{
			name: "On adaptive grid widgets that exceed the total row size should pass to the next row.",
			grid: func() (*grid.Grid, error) {
				maxWidth := 100
				widgets := []model.Widget{
					model.Widget{GridPos: model.GridPos{W: 50}},
					model.Widget{GridPos: model.GridPos{W: 60}},
					model.Widget{GridPos: model.GridPos{W: 25}},
					model.Widget{GridPos: model.GridPos{W: 10}},
					model.Widget{GridPos: model.GridPos{W: 75}},
					model.Widget{GridPos: model.GridPos{W: 100}},
				}

				return grid.NewAdaptiveGrid(maxWidth, widgets)
			},
			exp: &grid.Grid{
				MaxWidth: 100,
				Rows: []*grid.Row{
					&grid.Row{
						PercentSize: 25,
						Elements: []*grid.Element{
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{W: 50}},
								PercentSize: 50,
							},
							&grid.Element{
								Empty:       true,
								PercentSize: 50,
							},
						},
					},
					&grid.Row{
						PercentSize: 25,
						Elements: []*grid.Element{
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{W: 60}},
								PercentSize: 60,
							},
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{W: 25}},
								PercentSize: 25,
							},
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{W: 10}},
								PercentSize: 10,
							},
							&grid.Element{
								Empty:       true,
								PercentSize: 5,
							},
						},
					},
					&grid.Row{
						PercentSize: 25,
						Elements: []*grid.Element{
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{W: 75}},
								PercentSize: 75,
							},
							&grid.Element{
								Empty:       true,
								PercentSize: 25,
							},
						},
					},
					&grid.Row{
						PercentSize: 25,
						Elements: []*grid.Element{
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{W: 100}},
								PercentSize: 100,
							},
						},
					},
				},
			},
			expErr: false,
		},
		{
			name: "Using adaptive grids of different maxWidth.",
			grid: func() (*grid.Grid, error) {
				maxWidth := 1000
				widgets := []model.Widget{
					model.Widget{GridPos: model.GridPos{W: 500}},
					model.Widget{GridPos: model.GridPos{W: 600}},
					model.Widget{GridPos: model.GridPos{W: 250}},
					model.Widget{GridPos: model.GridPos{W: 100}},
					model.Widget{GridPos: model.GridPos{W: 750}},
					model.Widget{GridPos: model.GridPos{W: 1000}},
				}

				return grid.NewAdaptiveGrid(maxWidth, widgets)
			},
			exp: &grid.Grid{
				MaxWidth: 1000,
				Rows: []*grid.Row{
					&grid.Row{
						PercentSize: 25,
						Elements: []*grid.Element{
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{W: 500}},
								PercentSize: 50,
							},
							&grid.Element{
								Empty:       true,
								PercentSize: 50,
							},
						},
					},
					&grid.Row{
						PercentSize: 25,
						Elements: []*grid.Element{
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{W: 600}},
								PercentSize: 60,
							},
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{W: 250}},
								PercentSize: 25,
							},
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{W: 100}},
								PercentSize: 10,
							},
							&grid.Element{
								Empty:       true,
								PercentSize: 5,
							},
						},
					},
					&grid.Row{
						PercentSize: 25,
						Elements: []*grid.Element{
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{W: 750}},
								PercentSize: 75,
							},
							&grid.Element{
								Empty:       true,
								PercentSize: 25,
							},
						},
					},
					&grid.Row{
						PercentSize: 25,
						Elements: []*grid.Element{
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{W: 1000}},
								PercentSize: 100,
							},
						},
					},
				},
			},
			expErr: false,
		},
		{
			name: "On fixed grids (and unsorted) the widgets that are not one after the other should have an empty element in between.",
			grid: func() (*grid.Grid, error) {
				maxWidth := 100
				widgets := []model.Widget{
					model.Widget{GridPos: model.GridPos{Y: 0, X: 0, W: 50}},
					model.Widget{GridPos: model.GridPos{Y: 2, X: 0, W: 10}},
					model.Widget{GridPos: model.GridPos{Y: 2, X: 10, W: 20}},
					model.Widget{GridPos: model.GridPos{Y: 0, X: 60, W: 25}},
					model.Widget{GridPos: model.GridPos{Y: 2, X: 35, W: 5}},
					model.Widget{GridPos: model.GridPos{Y: 1, X: 10, W: 90}},
					model.Widget{GridPos: model.GridPos{Y: 2, X: 45, W: 10}},
					model.Widget{GridPos: model.GridPos{Y: 2, X: 30, W: 5}},
				}

				return grid.NewFixedGrid(maxWidth, widgets)
			},
			exp: &grid.Grid{
				MaxWidth:  100,
				MaxHeight: 3,
				Rows: []*grid.Row{
					&grid.Row{
						PercentSize: 33,
						Elements: []*grid.Element{
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{Y: 0, X: 0, W: 50}},
								PercentSize: 50,
							},
							&grid.Element{
								Empty:       true,
								PercentSize: 10,
							},
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{Y: 0, X: 60, W: 25}},
								PercentSize: 25,
							},
							&grid.Element{
								Empty:       true,
								PercentSize: 15,
							},
						},
					},
					&grid.Row{
						PercentSize: 33,
						Elements: []*grid.Element{
							&grid.Element{
								Empty:       true,
								PercentSize: 10,
							},
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{Y: 1, X: 10, W: 90}},
								PercentSize: 90,
							},
						},
					},
					&grid.Row{
						PercentSize: 33,
						Elements: []*grid.
							Element{
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{Y: 2, X: 0, W: 10}},
								PercentSize: 10,
							},
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{Y: 2, X: 10, W: 20}},
								PercentSize: 20,
							},
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{Y: 2, X: 30, W: 5}},
								PercentSize: 5,
							},
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{Y: 2, X: 35, W: 5}},
								PercentSize: 5,
							},
							&grid.Element{
								Empty:       true,
								PercentSize: 5,
							},
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{Y: 2, X: 45, W: 10}},
								PercentSize: 10,
							},
							&grid.Element{
								Empty:       true,
								PercentSize: 45,
							},
						},
					},
				},
			},
			expErr: false,
		},
		{
			name: "On fixed grids with different maxwidth.",
			grid: func() (*grid.Grid, error) {
				maxWidth := 1000
				widgets := []model.Widget{
					model.Widget{GridPos: model.GridPos{Y: 0, X: 0, W: 500}},
					model.Widget{GridPos: model.GridPos{Y: 2, X: 0, W: 100}},
					model.Widget{GridPos: model.GridPos{Y: 2, X: 100, W: 200}},
					model.Widget{GridPos: model.GridPos{Y: 0, X: 600, W: 250}},
					model.Widget{GridPos: model.GridPos{Y: 2, X: 350, W: 50}},
					model.Widget{GridPos: model.GridPos{Y: 1, X: 100, W: 900}},
					model.Widget{GridPos: model.GridPos{Y: 2, X: 450, W: 100}},
					model.Widget{GridPos: model.GridPos{Y: 2, X: 300, W: 50}},
				}

				return grid.NewFixedGrid(maxWidth, widgets)
			},
			exp: &grid.Grid{
				MaxWidth:  1000,
				MaxHeight: 3,
				Rows: []*grid.Row{
					&grid.Row{
						PercentSize: 33,
						Elements: []*grid.Element{
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{Y: 0, X: 0, W: 500}},
								PercentSize: 50,
							},
							&grid.Element{
								Empty:       true,
								PercentSize: 10,
							},
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{Y: 0, X: 600, W: 250}},
								PercentSize: 25,
							},
							&grid.Element{
								Empty:       true,
								PercentSize: 15,
							},
						},
					},
					&grid.Row{
						PercentSize: 33,
						Elements: []*grid.Element{
							&grid.Element{
								Empty:       true,
								PercentSize: 10,
							},
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{Y: 1, X: 100, W: 900}},
								PercentSize: 90,
							},
						},
					},
					&grid.Row{
						PercentSize: 33,
						Elements: []*grid.
							Element{
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{Y: 2, X: 0, W: 100}},
								PercentSize: 10,
							},
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{Y: 2, X: 100, W: 200}},
								PercentSize: 20,
							},
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{Y: 2, X: 300, W: 50}},
								PercentSize: 5,
							},
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{Y: 2, X: 350, W: 50}},
								PercentSize: 5,
							},
							&grid.Element{
								Empty:       true,
								PercentSize: 5,
							},
							&grid.Element{
								Widget:      model.Widget{GridPos: model.GridPos{Y: 2, X: 450, W: 100}},
								PercentSize: 10,
							},
							&grid.Element{
								Empty:       true,
								PercentSize: 45,
							},
						},
					},
				},
			},
			expErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			got, err := test.grid()
			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.exp, got)
			}
		})
	}

}
