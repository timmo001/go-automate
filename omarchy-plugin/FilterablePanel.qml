import QtQuick
import qs.Commons

Item {
  id: root

  property var model: []
  property string filterText: ""
  property int cursorIndex: 0

  readonly property var filteredModel: filterModel(model, filterText)
  readonly property int count: filteredModel.length

  signal activateRequested(var entry)
  signal closeRequested()
  signal refreshRequested()
  signal tabRequested(int direction)
  signal revealRequested()

  focus: true
  Keys.priority: Keys.BeforeItem

  onFilteredModelChanged: {
    cursorIndex = Math.max(0, Math.min(cursorIndex, Math.max(0, count - 1)))
    revealRequested()
  }

  function filterModel(entries, query) {
    var term = String(query || "").trim().toLowerCase()
    if (!term) return entries || []
    return (entries || []).filter(function(entry) {
      return [entry.label, entry.group, entry.value].join(" ").toLowerCase().indexOf(term) >= 0
    })
  }

  function reset() {
    filterText = ""
    cursorIndex = 0
  }

  function moveCursor(delta) {
    if (count <= 0) return
    cursorIndex = Math.max(0, Math.min(cursorIndex + delta, count - 1))
    revealRequested()
  }

  Keys.onPressed: function(event) {
    if (event.key === Qt.Key_Escape) {
      if (root.filterText) root.reset()
      else root.closeRequested()
      event.accepted = true
    } else if (event.key === Qt.Key_Tab || event.key === Qt.Key_Backtab) {
      root.tabRequested((event.modifiers & Qt.ShiftModifier)
        || event.key === Qt.Key_Backtab ? -1 : 1)
      event.accepted = true
    } else if (Util.editsFilter(event, root.filterText)) {
      root.filterText = Util.editedFilter(event, root.filterText)
      root.cursorIndex = 0
      event.accepted = true
    } else if (event.key === Qt.Key_Up) {
      root.moveCursor(-1)
      event.accepted = true
    } else if (event.key === Qt.Key_Down) {
      root.moveCursor(1)
      event.accepted = true
    } else if (event.key === Qt.Key_Return || event.key === Qt.Key_Enter) {
      if (root.cursorIndex < root.count)
        root.activateRequested(root.filteredModel[root.cursorIndex])
      event.accepted = true
    } else if (event.key === Qt.Key_R
        && event.modifiers === Qt.ControlModifier) {
      root.refreshRequested()
      event.accepted = true
    } else if (event.text && event.text.length === 1
        && event.text.charCodeAt(0) >= 32
        && event.text.charCodeAt(0) !== 127
        && (event.modifiers === Qt.NoModifier
          || event.modifiers === Qt.ShiftModifier)) {
      root.filterText += event.text
      root.cursorIndex = 0
      event.accepted = true
    }
  }
}
