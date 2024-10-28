package internal

import (
	"log"
	"math"
	"os"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	types2 "github.com/go-echarts/go-echarts/v2/types"
)

const (
	c1 = "#149b8e"
	c2 = "#2e86e3"
)

func RenderFiles(code *CompileResult) {
	var nodes []opts.GraphNode
	var links []opts.GraphLink

	var d2Lines []string

	for _, i := range code.interfaces {
		nodes = append(nodes, opts.GraphNode{
			Name: i.name,
			ItemStyle: &opts.ItemStyle{
				Color: c1,
			},
			Value:      float32(i.NumMethods()),
			SymbolSize: math.Log2(float64(i.NumMethods()+1)) * 4,
			Tooltip: &opts.Tooltip{
				Trigger:   "item",
				TriggerOn: "mousemove|click",
				Formatter: nodeTooltip(i.name, i.Methods()),
				Enterable: opts.Bool(true),
			},
		})

		d2Lines = append(d2Lines, dNode(i.name, i.Methods()))
	}

	for _, s := range code.structs {
		nodes = append(nodes, opts.GraphNode{
			Name: s.name,
			ItemStyle: &opts.ItemStyle{
				Color: c2,
			},
			Value:      float32(s.NumMethods()),
			SymbolSize: math.Log2(float64(s.NumMethods()+1)) * 4,
			Symbol:     "roundRect",
			Tooltip: &opts.Tooltip{
				Formatter: nodeTooltip(s.name, s.Methods()),
				Enterable: opts.Bool(true),
			},
		})

		d2Lines = append(d2Lines, dNode(s.name, s.Methods()))

		for _, i := range s.implements {
			links = append(links, opts.GraphLink{
				Source: i.name,
				Target: s.name,
			})

			d2Lines = append(d2Lines, dLink(i.name, s.name))
		}

		for _, i := range s.includes {
			links = append(links, opts.GraphLink{
				Source: s.name,
				Target: i,
			})

			d2Lines = append(d2Lines, dLink(s.name, i))
		}
	}

	sk := charts.NewGraph()
	sk.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: "1000px",
		}),
		charts.WithTitleOpts(opts.Title{
			Title: "Dependencies",
		}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show:   opts.Bool(true),
			Orient: "horizontal",
			Left:   "right",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show: opts.Bool(true), Title: "Save as image"},
				Restore: &opts.ToolBoxFeatureRestore{Show: opts.Bool(true), Title: "Reset"},
			},
		}),
	)

	sk.AddSeries("graph", nodes, links,
		charts.WithGraphChartOpts(opts.GraphChart{
			Layout: "force",
			Force: &opts.GraphForce{
				InitLayout: "circular",
				Repulsion:  10,
				Gravity:    0.01,
				EdgeLength: 0,
			},
			Roam:               opts.Bool(true),
			FocusNodeAdjacency: opts.Bool(true),
			Draggable:          opts.Bool(true),
		}),
		charts.WithItemStyleOpts(opts.ItemStyle{
			GapWidth: 100,
		}),
	)

	page := components.NewPage()

	page.AddCharts(sk)
	page.SetLayout(components.PageNoneLayout)

	err := os.WriteFile("result.html", page.RenderContent(), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	mermaidCode := "classDiagram\n\n" + strings.Join(d2Lines, "\n")
	err = os.WriteFile("result_mermaid.txt", []byte(mermaidCode), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func dNode(name string, methods []string) string {
	for i, m := range methods {
		methods[i] = "	" + m + "()"
	}

	methodsCode := strings.Join(methods, "\n")
	if methodsCode == "" {
		methodsCode = "empty"
	}

	return "class " + dClean(name) + "{\n" +
		methodsCode +
		"\n}"
}

func nodeTooltip(name string, methods []string) types2.FuncStr {
	for i, m := range methods {
		methods[i] = m + "()"
	}

	return types2.FuncStr("<b>" + name + "</b></br></br>" + strings.Join(methods, "</br>"))
}

func dLink(source, target string) string {
	return dClean(source) + " --|> " + dClean(target)
}

func dClean(s string) string {
	return strings.ReplaceAll(s, "/", "_")
}
