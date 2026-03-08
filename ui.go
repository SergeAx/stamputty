package main

import (
	"fmt"

	"github.com/tailscale/walk"
	. "github.com/tailscale/walk/declarative"
	"github.com/tailscale/win"
)

type AppUI struct {
	mainWindow      *walk.MainWindow
	sessionsList    *walk.ListBox
	settingsTable   *walk.TableView
	showAllCheckBox *walk.CheckBox
	selectAllButton *walk.PushButton
	clearButton     *walk.PushButton
	applyButton     *walk.PushButton
	statusLabel     *walk.Label
	splitter        *walk.Splitter
	sessionsModel   *SessionsModel
	settingsModel   *SettingsModel
	selectedSession *Session
}

type SessionsModel struct {
	walk.ListModelBase
	sessions []Session
}

func (m *SessionsModel) ItemCount() int {
	return len(m.sessions)
}

func (m *SessionsModel) Value(index int) interface{} {
	return m.sessions[index].DisplayName
}

func (m *SessionsModel) SetSessions(sessions []Session) {
	m.sessions = sessions
	m.PublishItemsReset()
}

type SettingsModel struct {
	walk.TableModelBase
	settings         []Setting
	onCheckedChanged func()
}

func (m *SettingsModel) RowCount() int {
	return len(m.settings)
}

func (m *SettingsModel) Value(row, col int) interface{} {
	setting := m.settings[row]

	switch col {
	case 0:
		return setting.Name
	case 1:
		return setting.GetDefaultValueString()
	case 2:
		return setting.GetCurrentValueString()
	}
	return nil
}

func (m *SettingsModel) Checked(row int) bool {
	return m.settings[row].IsChecked
}

func (m *SettingsModel) SetChecked(row int, checked bool) error {
	if checked && !m.settings[row].IsDifferent {
		return nil
	}

	m.settings[row].IsChecked = checked
	if m.onCheckedChanged != nil {
		m.onCheckedChanged()
	}
	return nil
}

func (m *SettingsModel) SetSettings(settings []Setting) {
	oldLen := len(m.settings)
	m.settings = settings
	for i := range m.settings {
		if !m.settings[i].IsDifferent {
			m.settings[i].IsChecked = false
		}
	}
	newLen := len(m.settings)

	if oldLen > newLen {
		m.PublishRowsRemoved(newLen, oldLen-1)
	}
	if newLen > oldLen {
		m.PublishRowsInserted(oldLen, newLen-1)
	}
	if newLen > 0 {
		m.PublishRowsChanged(0, newLen-1)
	}

	if m.onCheckedChanged != nil {
		m.onCheckedChanged()
	}
}

func (m *SettingsModel) GetCheckedSettings() []Setting {
	var checked []Setting
	for _, setting := range m.settings {
		if setting.IsChecked {
			checked = append(checked, setting)
		}
	}
	return checked
}

func (m *SettingsModel) CheckedSettingNames() map[string]bool {
	checkedNames := make(map[string]bool)
	for _, setting := range m.settings {
		if setting.IsChecked {
			checkedNames[setting.Name] = true
		}
	}

	return checkedNames
}

func (m *SettingsModel) RestoreCheckedSettingNames(checkedNames map[string]bool) {
	if len(checkedNames) == 0 {
		return
	}

	for i := range m.settings {
		if !m.settings[i].IsDifferent {
			continue
		}

		if checkedNames[m.settings[i].Name] {
			m.settings[i].IsChecked = true
		}
	}

	if len(m.settings) > 0 {
		m.PublishRowsChanged(0, len(m.settings)-1)
	}

	if m.onCheckedChanged != nil {
		m.onCheckedChanged()
	}
}

func (m *SettingsModel) UncheckAll() {
	for i := range m.settings {
		m.settings[i].IsChecked = false
	}
	if len(m.settings) > 0 {
		m.PublishRowsChanged(0, len(m.settings)-1)
	}
	if m.onCheckedChanged != nil {
		m.onCheckedChanged()
	}
}

func (m *SettingsModel) CheckAll() {
	for i := range m.settings {
		m.settings[i].IsChecked = m.settings[i].IsDifferent
	}
	if len(m.settings) > 0 {
		m.PublishRowsChanged(0, len(m.settings)-1)
	}
	if m.onCheckedChanged != nil {
		m.onCheckedChanged()
	}
}

func createUI() (*AppUI, error) {
	ui := &AppUI{
		sessionsModel: &SessionsModel{},
		settingsModel: &SettingsModel{},
	}
	ui.settingsModel.onCheckedChanged = ui.updateApplyButtonState

	sessionListPane := Composite{
		StretchFactor: 1,
		Layout:        VBox{MarginsZero: true},
		MinSize:       Size{Width: 200},
		Children: []Widget{
			Label{
				Text: "Sessions:",
			},
			ListBox{
				AssignTo:              &ui.sessionsList,
				Model:                 ui.sessionsModel,
				OnCurrentIndexChanged: ui.onSessionSelected,
			},
		},
	}

	buttonBar := Composite{
		Layout: HBox{MarginsZero: true},
		Children: []Widget{
			CheckBox{
				AssignTo: &ui.showAllCheckBox,
				Text:     "Show unchanged settings",
				OnCheckedChanged: func() {
					ui.onShowAllToggled()
				},
			},
			HSpacer{},
			PushButton{
				AssignTo: &ui.selectAllButton,
				Text:     "Select All",
				OnClicked: func() {
					ui.settingsModel.CheckAll()
				},
			},
			PushButton{
				AssignTo: &ui.clearButton,
				Text:     "Clear",
				OnClicked: func() {
					ui.onClear()
				},
			},
			PushButton{
				AssignTo: &ui.applyButton,
				Text:     "Apply",
				Enabled:  false,
				OnClicked: func() {
					ui.onApply()
				},
			},
		},
	}

	settingsPane := Composite{
		StretchFactor: 3,
		Layout:        VBox{MarginsZero: true},
		MinSize:       Size{Width: 650},
		Children: []Widget{
			Label{
				Text: "Difference with Default Settings:",
			},
			TableView{
				AssignTo:            &ui.settingsTable,
				CheckBoxes:          true,
				ColumnsOrderable:    false,
				AlternatingRowBG:    true,
				StyleCell:           ui.styleSettingsCell,
				LastColumnStretched: true,
				Columns: []TableViewColumn{
					{Title: "Setting Name", Width: 200},
					{Title: "Default Value", Width: 200},
					{Title: "Current Value", Width: 200},
				},
				Model: ui.settingsModel,
			},
			Label{
				AssignTo: &ui.statusLabel,
				Text:     "",
			},
			buttonBar,
		},
	}

	err := MainWindow{
		AssignTo: &ui.mainWindow,
		Title:    "StamPuTTY - PuTTY Session Sync",
		MinSize:  Size{Width: 800, Height: 600},
		Layout:   VBox{SpacingZero: true},
		Children: []Widget{
			HSplitter{
				AssignTo: &ui.splitter,
				Children: []Widget{
					sessionListPane,
					settingsPane,
				},
			},
		},
	}.Create()

	if err != nil {
		return nil, err
	}

	return ui, nil
}

func (ui *AppUI) onSessionSelected() {
	index := ui.sessionsList.CurrentIndex()
	if index < 0 || index >= len(ui.sessionsModel.sessions) {
		return
	}

	ui.selectedSession = &ui.sessionsModel.sessions[index]
	ui.refreshSettings(true)
}

func (ui *AppUI) onShowAllToggled() {
	if ui.selectedSession == nil {
		return
	}

	ui.refreshSettings(false)
}

func (ui *AppUI) refreshSettings(showNoDifferencePopup bool) {
	showAll := ui.showAllCheckBox != nil && ui.showAllCheckBox.Checked()
	checkedSettingNames := ui.settingsModel.CheckedSettingNames()

	var (
		settings []Setting
		err      error
	)

	if showAll {
		settings, err = computeAllSettings(ui.selectedSession.EncodedName)
	} else {
		settings, err = computeDiff(ui.selectedSession.EncodedName)
	}

	if err != nil {
		showTaskDialog(ui.mainWindow, "Error", fmt.Sprintf("Failed to compute differences: %v", err))
		return
	}

	if len(settings) == 0 {
		ui.settingsModel.SetSettings(nil)
		if ui.statusLabel != nil {
			_ = ui.statusLabel.SetText("No difference detected")
		}
		if showNoDifferencePopup && !showAll {
			showTaskDialog(ui.mainWindow, "Info", "No difference detected")
		}
		return
	}

	ui.settingsModel.SetSettings(settings)
	ui.settingsModel.RestoreCheckedSettingNames(checkedSettingNames)
	if ui.statusLabel != nil {
		_ = ui.statusLabel.SetText("")
	}
}

func (ui *AppUI) onApply() {
	if ui.selectedSession == nil {
		showTaskDialog(ui.mainWindow, "Error", "No session selected")
		return
	}

	checkedSettings := ui.settingsModel.GetCheckedSettings()
	if len(checkedSettings) == 0 {
		showTaskDialog(ui.mainWindow, "Info", "No settings selected to apply")
		return
	}

	for _, setting := range checkedSettings {
		err := writeSettingToSession(ui.selectedSession.EncodedName, setting.Name, setting.DefaultValue, setting.Type)
		if err != nil {
			showTaskDialog(ui.mainWindow, "Error", fmt.Sprintf("Failed to write setting %s: %v", setting.Name, err))
			return
		}
	}

	showTaskDialog(ui.mainWindow, "Success", fmt.Sprintf("Successfully updated %d setting(s)", len(checkedSettings)))
	ui.refreshSettings(false)
}

func (ui *AppUI) onClear() {
	ui.settingsModel.UncheckAll()
}

func (ui *AppUI) updateApplyButtonState() {
	if ui.applyButton == nil {
		return
	}

	ui.applyButton.SetEnabled(len(ui.settingsModel.GetCheckedSettings()) > 0)
}

func (ui *AppUI) styleSettingsCell(style *walk.CellStyle) {
	row := style.Row()
	if row < 0 || row >= len(ui.settingsModel.settings) {
		return
	}

	if !ui.settingsModel.settings[row].IsDifferent {
		style.TextColor = walk.Color(win.GetSysColor(int(walk.SysColorGrayText)))
		style.BackgroundColor = walk.Color(win.GetSysColor(int(walk.SysColorInactiveBorder)))
	}
}

func (ui *AppUI) loadSessions() error {
	sessions, err := getSessions()
	if err != nil {
		return err
	}

	if len(sessions) == 0 {
		return fmt.Errorf("no sessions found, check if PuTTY is installed")
	}

	ui.sessionsModel.SetSessions(sessions)
	return nil
}

func (ui *AppUI) Run() {
	ui.mainWindow.Show()
	walk.App().Run()
}

func showTaskDialog(owner walk.Form, title, content string) {
	taskDialog := walk.NewTaskDialog()
	_, _ = taskDialog.Show(walk.TaskDialogOpts{
		Owner:         owner,
		Title:         title,
		Content:       content,
		CommonButtons: win.TDCBF_OK_BUTTON,
	})
}
