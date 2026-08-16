import QtQuick
import QtQuick.Controls
import qs.Commons
import qs.Ui

Panel {
  id: root
  moduleName: "timmo.home-assistant"

  Config { id: config }

  property var anchorItem: null
  property var hostWidget: null
  property var service: null
  readonly property var barIdentity: hostWidget || root
  readonly property var panelConfig: config.panel
  readonly property color foreground: bar ? bar.foreground : Color.foreground
  readonly property string fontFamily: bar ? bar.fontFamily : Style.font.family
  readonly property var panelRows: buildRows()

  function rowValue(row) {
    if (!row.available) return "Unavailable"
    var value = String(row.text || "").trim()
    if (row.icon && value.indexOf(row.icon) === 0)
      value = value.slice(String(row.icon).length).trim()
    if (value !== "") return value
    return row.severity === "active" ? row.activeText : row.inactiveText
  }

  function buildRows() {
    var entries = []
    var rows = service ? service.rows : []
    for (var i = 0; i < rows.length; i++) {
      entries.push({
        id: rows[i].id,
        group: rows[i].group,
        label: rows[i].label,
        icon: rows[i].icon,
        action: rows[i].action,
        available: rows[i].available,
        severity: rows[i].severity,
        color: rows[i].color,
        value: rowValue(rows[i])
      })
    }
    return entries
  }

  function open() {
    filterController.reset()
    if (service) service.refresh()
    controller.show()
    Qt.callLater(function() {
      list.contentY = 0
      filterController.forceActiveFocus()
    })
  }

  function close() { controller.hide() }
  function toggle() { if (opened) close(); else open() }

  function activate(entry) {
    if (!service || !entry || !entry.action) return
    service.activate(entry.id)
  }

  KeyboardPanel {
    id: panel
    anchorItem: root.anchorItem
    owner: root.barIdentity
    bar: root.bar
    open: root.opened
    focusTarget: filterController
    contentWidth: panel.fittedContentWidth(Style.space(root.panelConfig.width))
    contentHeight: panel.fittedContentHeight(content.implicitHeight,
      Style.space(root.panelConfig.maxHeight))

    FilterablePanel {
      id: filterController
      anchors.fill: parent
      model: root.panelRows
      onActivateRequested: function(entry) { root.activate(entry) }
      onCloseRequested: root.close()
      onTabRequested: function(direction) {
        if (root.bar && typeof root.bar.switchPanelFrom === "function")
          root.bar.switchPanelFrom(root.barIdentity, direction)
      }
      onRefreshRequested: if (root.service) root.service.refresh()

      Flickable {
        id: list
        anchors.fill: parent
        contentWidth: width
        contentHeight: content.implicitHeight
        clip: true
        boundsBehavior: Flickable.StopAtBounds
        ScrollBar.vertical: ScrollBar { policy: ScrollBar.AsNeeded }

        Column {
          id: content
          width: list.width

          PanelHero {
            width: parent.width
            title: filterController.filterText || root.panelConfig.title
            foreground: root.foreground
            fontFamily: root.fontFamily
            iconComponent: Component {
              Text {
                text: config.bar.icon
                color: config.bar.colors.active
                font.family: root.fontFamily
                font.pixelSize: Style.font.display
              }
            }
          }

          Repeater {
            model: filterController.filteredModel

            CursorSurface {
              required property int index
              required property var modelData
              width: content.width
              implicitHeight: rowContent.implicitHeight
                + Style.space(root.panelConfig.rowPadding)
              hasCursor: index === filterController.cursorIndex
              foreground: root.foreground
              accent: modelData.color || root.foreground

              Row {
                id: rowContent
                anchors.left: parent.left
                anchors.right: parent.right
                anchors.verticalCenter: parent.verticalCenter
                anchors.margins: Style.space(root.panelConfig.rowPadding)
                spacing: Style.space(root.panelConfig.rowSpacing)

                Text {
                  width: Style.space(22)
                  text: modelData.icon
                  color: modelData.color || root.foreground
                  font.family: root.fontFamily
                  font.pixelSize: Style.font.icon
                  horizontalAlignment: Text.AlignHCenter
                }

                Column {
                  width: Math.max(0, parent.width - Style.space(32))
                  Text {
                    width: parent.width
                    text: modelData.label
                    color: root.foreground
                    font.family: root.fontFamily
                    font.pixelSize: Style.font.body
                    elide: Text.ElideRight
                  }
                  Text {
                    width: parent.width
                    text: modelData.value
                    color: Qt.darker(root.foreground, 1.35)
                    font.family: root.fontFamily
                    font.pixelSize: Style.font.caption
                    elide: Text.ElideRight
                  }
                }
              }

              MouseArea {
                anchors.fill: parent
                hoverEnabled: true
                cursorShape: modelData.action ? Qt.PointingHandCursor : Qt.ArrowCursor
                onEntered: filterController.cursorIndex = index
                onClicked: root.activate(modelData)
              }
            }
          }

          Text {
            visible: filterController.filterText && filterController.count === 0
            width: parent.width
            text: "No matches for " + filterController.filterText
            color: Qt.darker(root.foreground, 1.35)
            font.family: root.fontFamily
            font.pixelSize: Style.font.body
            horizontalAlignment: Text.AlignHCenter
          }
        }
      }
    }
  }
}
