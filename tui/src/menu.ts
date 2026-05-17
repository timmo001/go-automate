import type { MenuItem, MenuVariant, NotifyConfig, FlagField, FlagPopupAction, ViewId } from "./types.js";

// --- Helpers ---

function item(
  id: string,
  icon: string,
  title: string,
  description: string,
  action: MenuItem["action"],
  variants?: readonly MenuVariant[],
  keywords?: readonly string[],
): MenuItem {
  return {
    id,
    icon,
    title,
    description,
    action,
    ...(variants && { variants }),
    ...(keywords && { keywords }),
  };
}

function cmd(command: string, wait = true): MenuItem["action"] {
  return { type: "command", cmd: command, wait };
}

function silent(command: string): MenuItem["action"] {
  return { type: "silent", cmd: command };
}

function notify(command: string, config: NotifyConfig): MenuItem["action"] {
  return { type: "notify", cmd: command, notify: config };
}

function view(viewId: ViewId): MenuItem["action"] {
  return { type: "view", viewId };
}

function submenu(menuId: string): MenuItem["action"] {
  return { type: "submenu", menuId };
}

function replace(command: string): MenuItem["action"] {
  return { type: "replace", cmd: command };
}

function flagPopup(
  baseCmd: string,
  title: string,
  fields: readonly FlagField[],
  advancedFieldIndices?: readonly number[],
): FlagPopupAction {
  return { type: "flagPopup", baseCmd, title, fields, advancedFieldIndices };
}

// --- Shared field definitions ---

const entityNameField: FlagField = {
  name: "entity_name",
  label: "Entity Name",
  type: "string",
  required: true,
  positional: true,
  placeholder: "e.g. bedroom_lamp",
};

const areaIdField: FlagField = {
  name: "area_id",
  label: "Area ID",
  type: "string",
  required: true,
  positional: true,
  placeholder: "e.g. living_room",
};

const messageField: FlagField = {
  name: "message",
  label: "Message",
  type: "string",
  required: true,
  positional: true,
  placeholder: "e.g. Dinner is ready",
};

const entityIdField: FlagField = {
  name: "entity_id",
  label: "Entity ID",
  type: "string",
  required: true,
  positional: true,
  placeholder: "e.g. sensor.temperature",
};

const notifySummaryField: FlagField = {
  name: "summary",
  label: "Summary",
  type: "string",
  required: true,
  positional: true,
  placeholder: "e.g. Build complete",
};

const notifyBodyField: FlagField = {
  name: "body",
  label: "Body",
  type: "string",
  required: false,
  positional: true,
  placeholder: "(optional)",
};

// --- Helpers for control submenus (light, switch, input_boolean) ---

function createControlSubmenuItems(
  domain: string,
  labelPrefix: string,
): readonly MenuItem[] {
  return [
    item(
      `ha.${domain}.turn-on`,
      "󰔌",
      "Turn On",
      `Turn on a ${labelPrefix}`,
      flagPopup(
        `go-automate ha ${domain} turn-on`,
        `${labelPrefix} › Turn On`,
        [entityNameField],
      ),
      undefined,
      ["on", "enable", "activate", "start"],
    ),
    item(
      `ha.${domain}.turn-off`,
      "󰔍",
      "Turn Off",
      `Turn off a ${labelPrefix}`,
      flagPopup(
        `go-automate ha ${domain} turn-off`,
        `${labelPrefix} › Turn Off`,
        [entityNameField],
      ),
      undefined,
      ["off", "disable", "deactivate", "stop"],
    ),
    item(
      `ha.${domain}.toggle`,
      "󰑐",
      "Toggle",
      `Toggle a ${labelPrefix}`,
      flagPopup(
        `go-automate ha ${domain} toggle`,
        `${labelPrefix} › Toggle`,
        [entityNameField],
      ),
      undefined,
      ["switch", "flip", "swap"],
    ),
  ];
}

// --- Main menu ---

const mainItems: readonly MenuItem[] = [
  item(
    "bridge-serve",
    "󰐊",
    "Bridge Serve",
    "Start the HA websocket bridge",
    replace("go-automate ha bridge serve"),
    undefined,
    ["start", "daemon", "run", "server", "launch", "bridge", ":run", ":serve", "go"],
  ),

  item(
    "ha",
    "󰋜",
    "Home Assistant",
    "Control HA entities",
    submenu("ha"),
    undefined,
    ["home", "assistant", "hass", "light", "switch", "entity", ":ha"],
  ),

  item(
    "watch",
    "󰈈",
    "Watch Entity",
    "Watch entity state changes via bridge",
    flagPopup(
      "go-automate ha bridge watch entity",
      "Watch Entity",
      [entityIdField],
    ),
    undefined,
    ["observe", "monitor", "state", "sensor", "bridge", ":w"],
  ),

  item(
    "notify",
    "󰍡",
    "Notify",
    "Send a desktop notification",
    flagPopup(
      "go-automate notify",
      "Send Notification",
      [notifySummaryField, notifyBodyField],
    ),
    undefined,
    ["notification", "alert", "message", "send", ":n"],
  ),

  item("quit", "󰩈", "Quit", "Exit the menu", { type: "quit" }, undefined, [
    ":q",
    ":wq",
    ":qa",
    "exit",
    "quit",
    "close",
    "bye",
  ]),
];

// --- Home Assistant submenu ---

const haItems: readonly MenuItem[] = [
  item(
    "ha.light",
    "󰌵",
    "Light",
    "Control lights",
    submenu("ha.light"),
    undefined,
    ["lamp", "bulb", "brightness", "illuminate", ":l"],
  ),
  item(
    "ha.switch",
    "󰔡",
    "Switch",
    "Control switches",
    submenu("ha.switch"),
    undefined,
    ["outlet", "plug", "relay", ":s"],
  ),
  item(
    "ha.input_boolean",
    "󰨚",
    "Input Boolean",
    "Toggle input booleans",
    submenu("ha.input_boolean"),
    undefined,
    ["boolean", "flag", "helper", "toggle", ":ib"],
  ),
  item(
    "ha.assist_satellite",
    "󱄡",
    "Assist Satellite",
    "Announce to areas",
    submenu("ha.assist_satellite"),
    undefined,
    ["announce", "satellite", "speaker", "voice", "tts", ":as"],
  ),
];

// --- Control submenus ---

const lightItems = createControlSubmenuItems("light", "Light");
const switchItems = createControlSubmenuItems("switch", "Switch");
const inputBooleanItems = createControlSubmenuItems("input_boolean", "Input Boolean");

// --- Assist Satellite submenu ---

const assistSatelliteItems: readonly MenuItem[] = [
  item(
    "ha.assist_satellite.announce",
    "󱄡",
    "Announce",
    "Announce a message to an area",
    flagPopup(
      "go-automate ha assist_satellite announce",
      "Assist Satellite › Announce",
      [areaIdField, messageField],
    ),
    undefined,
    ["message", "tts", "broadcast", "speak", "say"],
  ),
];

// --- Registries ---

/** Top-level main menu items */
export const mainMenuItems: readonly MenuItem[] = mainItems;

/** Map of submenu ID → items */
export const submenus: Map<string, readonly MenuItem[]> = new Map([
  ["ha", haItems],
  ["ha.light", lightItems],
  ["ha.switch", switchItems],
  ["ha.input_boolean", inputBooleanItems],
  ["ha.assist_satellite", assistSatelliteItems],
]);

/** Display titles for submenu breadcrumbs */
export const submenuTitles: Map<string, string> = new Map([
  ["ha", "Home Assistant"],
  ["ha.light", "Light"],
  ["ha.switch", "Switch"],
  ["ha.input_boolean", "Input Boolean"],
  ["ha.assist_satellite", "Assist Satellite"],
]);

/** Flat map of every menu item by its ID (main items + all submenu items) */
export const menuItemsById: Map<string, MenuItem> = new Map();

function registerItems(items: readonly MenuItem[]): void {
  for (const m of items) {
    menuItemsById.set(m.id, m);
  }
}

registerItems(mainItems);
registerItems(haItems);
registerItems(lightItems);
registerItems(switchItems);
registerItems(inputBooleanItems);
registerItems(assistSatelliteItems);
