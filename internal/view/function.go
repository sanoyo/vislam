package view

import (
	"fmt"
	"log/slog"

	lambdatype "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/rivo/tview"
)

type functionView struct {
	view
	functions []lambdatype.FunctionConfiguration
}

func newFunctionView(functions []lambdatype.FunctionConfiguration, app *App) *functionView {
	keys := append(basicKeyInputs, []keyDescriptionPair{
		hotKeyMap["n"],
	}...)
	return &functionView{
		view: *newView(app, keys, secondaryPageKeyMap{
			DescriptionKind: describePageKeys,
		}),
		functions: functions,
	}
}

func (app *App) showFunctionsPage(reload bool) error {
	app.kind = ClusterKind
	if switched := app.switchPage(reload); switched {
		return nil
	}

	funcitions, err := app.Store.ListFunctions()
	if err != nil {
		slog.Error("failed to load funcitions", "region", app.Region, "error", err.Error())
		return err
	}

	if len(funcitions) == 0 {
		m := fmt.Sprintf("there is no valid funcitions in %s region", app.Region)
		slog.Warn("failed start", "reason", m)
		return fmt.Errorf(m)
	}

	view := newFunctionView(funcitions, app)
	page := buildAppPage(view)
	app.addAppPage(page)
	view.table.Select(app.rowIndex, 0)
	return nil
}

// Build table for function page
func (f *functionView) bodyBuilder() *tview.Pages {
	title, headers, dataBuilder := f.tableParam()
	f.buildTable(title, headers, dataBuilder)
	f.tableHandler()
	return f.bodyPages
}
