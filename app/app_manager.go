package app

import "errors"

// WtfAppManager handles the instances of WtfApp, ensuring that they're displayed as requested
type WtfAppManager struct {
	WtfApps []*WtfApp

	selected int
}

// Current returns the currently-displaying instance of WtfApp
func (appMan *WtfAppManager) Current() (*WtfApp, error) {
	if appMan.selected < 0 || appMan.selected >= len(appMan.WtfApps) {
		return nil, errors.New("invalid app index selected")
	}

	return appMan.WtfApps[appMan.selected], nil
}
