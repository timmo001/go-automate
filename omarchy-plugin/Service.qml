import QtQuick
import Quickshell
import Quickshell.Io

Item {
  id: root

  Config { id: config }

  property string modulesFile: ""
  property var modules: []
  property var moduleStates: []
  property int revision: 0

  readonly property var rows: buildRows(revision)
  readonly property var barConfig: config.bar

  function configure(path) {
    var nextPath = String(path || "")
    if (nextPath === modulesFile) return
    modulesFile = nextPath
    if (!nextPath) clearModules()
  }

  function clearModules() {
    modules = []
    moduleStates = []
    revision++
  }

  function loadModules(content) {
    try {
      var parsed = JSON.parse(String(content || ""))
      var nextModules = Array.isArray(parsed) ? parsed : parsed.modules
      if (!Array.isArray(nextModules)) throw new Error("modules must be an array")
      modules = nextModules.filter(function(module) {
        return module && typeof module.id === "string"
          && typeof module.command === "string"
      })
      var initial = []
      for (var i = 0; i < modules.length; i++)
        initial.push({ text: "", tooltip: "", className: "", available: false })
      moduleStates = initial
      revision++
    } catch (error) {
      console.warn("timmo.home-assistant: invalid modules file", modulesFile, error)
      clearModules()
    }
  }

  function setState(index, state) {
    if (index < 0 || index >= moduleStates.length) return
    var next = moduleStates.slice()
    next[index] = state
    moduleStates = next
    revision++
  }

  function applyOutput(index, raw) {
    try {
      var payload = JSON.parse(String(raw || "").trim())
      var text = payload.text === undefined || payload.text === null
        ? "" : String(payload.text)
      var tooltip = payload.tooltip === undefined || payload.tooltip === null
        ? "" : String(payload.tooltip)
      var className = payload["class"] === undefined || payload["class"] === null
        ? String(payload.alt || "") : String(payload["class"])
      setState(index, {
        text: text,
        tooltip: tooltip,
        className: className,
        available: text !== "" || tooltip !== "" || className !== "hidden"
      })
    } catch (error) {
      setState(index, { text: "", tooltip: "", className: "", available: false })
    }
  }

  function hasClass(className, values) {
    var classes = String(className || "").split(/\s+/)
    for (var i = 0; i < (values || []).length; i++)
      if (classes.indexOf(values[i]) !== -1) return true
    return false
  }

  function severity(module, state) {
    if (!state.available) return "unavailable"
    var classes = module.severityClasses || {}
    if (hasClass(state.className, classes.critical || ["critical"])) return "critical"
    if (hasClass(state.className, classes.warning || ["warning"])) return "warning"
    if (hasClass(state.className, classes.active || ["active"])) return "active"
    return "quiet"
  }

  function buildRows(dependency) {
    var result = []
    for (var i = 0; i < modules.length; i++) {
      var module = modules[i]
      if (module.background === true) continue
      var state = moduleStates[i]
        || { text: "", tooltip: "", className: "", available: false }
      if (module.hideUnavailable === true && !state.available) continue
      if (hasClass(state.className, module.hideClasses || [])) continue
      var rowSeverity = severity(module, state)
      result.push({
        id: module.id,
        group: module.group || "Status",
        label: module.label || module.id,
        icon: module.icon || "",
        action: module.action || "",
        text: state.text,
        tooltip: state.tooltip,
        available: state.available,
        inactiveText: module.inactiveText || "Quiet",
        activeText: module.activeText || "Active",
        barIconOnly: module.barIconOnly === true,
        panelOnly: module.panelOnly === true,
        severity: rowSeverity,
        barActive: module.panelOnly !== true && state.available
          && !hasClass(state.className, module.barHideClasses || ["hidden", "inactive"]),
        color: (module.colors || {})[rowSeverity] || ""
      })
    }
    return result
  }

  function refresh() {
    for (var i = 0; i < processInstantiator.count; i++) {
      var runner = processInstantiator.objectAt(i)
      if (runner && typeof runner.poll === "function") runner.poll()
    }
  }

  function activate(rowId) {
    for (var i = 0; i < modules.length; i++) {
      if (modules[i].id === rowId && modules[i].action) {
        Quickshell.execDetached(["bash", "-lc", modules[i].action])
        return
      }
    }
  }

  FileView {
    path: root.modulesFile
    watchChanges: true
    printErrors: false
    onLoaded: root.loadModules(text())
    onLoadFailed: root.clearModules()
    onFileChanged: reload()
  }

  Instantiator {
    id: processInstantiator
    model: root.modules

    delegate: QtObject {
      id: runner
      required property int index
      required property var modelData

      function poll() {
        if (modelData.stream === true || process.running) return
        process.running = true
      }

      property Process process: Process {
        running: runner.modelData.stream === true
        command: ["bash", "-lc", runner.modelData.command]
        stdout: runner.modelData.stream === true ? streamOutput : pollOutput
        onExited: {
          if (runner.modelData.stream === true) {
            running = false
            root.setState(runner.index,
              { text: "", tooltip: "", className: "", available: false })
            restartTimer.restart()
          } else {
            root.applyOutput(runner.index, pollOutput.text)
          }
        }
      }

      property StdioCollector pollOutput: StdioCollector { waitForEnd: true }
      property SplitParser streamOutput: SplitParser {
        onRead: function(line) { root.applyOutput(runner.index, line) }
      }
      property Timer pollTimer: Timer {
        interval: Math.max(1000, Number(runner.modelData.interval || 60000))
        running: runner.modelData.stream !== true
        repeat: true
        triggeredOnStart: true
        onTriggered: runner.poll()
      }
      property Timer restartTimer: Timer {
        interval: 5000
        onTriggered: if (runner.modelData.stream === true
            && !runner.process.running) runner.process.running = true
      }
    }
  }
}
