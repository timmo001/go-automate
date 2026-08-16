import QtQuick

QtObject {
  readonly property var bar: ({
    label: "Home Assistant",
    icon: "󰟐",
    fontSize: 10,
    horizontalMargin: 6,
    rowSpacing: 10,
    colors: {
      quiet: "",
      active: "#2bb3b1",
      warning: "#e7ad63",
      critical: "#e06c75",
      unavailable: "#a55555"
    }
  })

  readonly property var panel: ({
    title: "Home Assistant",
    width: 430,
    maxHeight: 670,
    rowPadding: 12,
    rowSpacing: 10
  })
}
