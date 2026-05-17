// --- Menu types ---

/** Identifies a top-level TUI view for navigation */
export type ViewId = "main" | "submenu";

// --- Menu action types ---

/** Action that suspends the TUI and runs a command with inherited stdio */
export interface CommandAction {
  readonly type: "command";
  readonly cmd: string;
  /** When true, show "Press any key to continue" before resuming the TUI */
  readonly wait: boolean;
}

/** Action that runs a command in the background without suspending */
export interface SilentAction {
  readonly type: "silent";
  readonly cmd: string;
}

/** Toast notification config for notify actions */
export interface NotifyConfig {
  /** Stable ID for grouping — a new notification with the same ID replaces the previous one */
  readonly id: string;
  /** Message shown while the command is running */
  readonly progress: string;
  /** Message shown on success */
  readonly success: string;
}

/** Action that runs a command silently and shows toast notifications for progress/result */
export interface NotifyAction {
  readonly type: "notify";
  readonly cmd: string;
  readonly notify: NotifyConfig;
}

/** Toast variant controlling display colour */
export type ToastVariant = "info" | "success" | "error";

/** Action that navigates to a sub-view within the TUI */
export interface ViewAction {
  readonly type: "view";
  readonly viewId: ViewId;
}

/** Action that opens a nested submenu */
export interface SubmenuAction {
  readonly type: "submenu";
  readonly menuId: string;
}

/** Action that replaces the TUI process with a command (no return) */
export interface ReplaceAction {
  readonly type: "replace";
  readonly cmd: string;
}

/** Action that opens a flag popup form to collect values before running */
export interface FlagPopupAction {
  readonly type: "flagPopup";
  /** Base command without flags (e.g. "go-automate ha light toggle") */
  readonly baseCmd: string;
  /** Display title for the popup */
  readonly title: string;
  /** Flag field definitions */
  readonly fields: readonly FlagField[];
  /** Fields hidden behind "Advanced options" toggle (by index into fields) */
  readonly advancedFieldIndices?: readonly number[];
}

/** Action that exits the TUI */
export interface QuitAction {
  readonly type: "quit";
}

/** Discriminated union of all possible menu item actions */
export type MenuAction =
  | CommandAction
  | SilentAction
  | NotifyAction
  | ViewAction
  | SubmenuAction
  | ReplaceAction
  | FlagPopupAction
  | QuitAction;

/** A selectable variant for a menu item offering an alternative action */
export interface MenuVariant {
  /** Short display label (e.g. "Quick", "Full") */
  readonly label: string;
  /** Optional description shown below the label in the popup */
  readonly description?: string;
  /** The action to execute when this variant is selected */
  readonly action: MenuAction;
}

/** A single entry in the TUI menu system */
export interface MenuItem {
  /** Stable dot-separated identifier (e.g. "ha.light.toggle") */
  readonly id: string;
  /** Primary display text */
  readonly title: string;
  /** Secondary text shown below the title */
  readonly description: string;
  /** Nerd Font icon character */
  readonly icon: string;
  /** What happens when this item is selected */
  readonly action: MenuAction;
  /** Optional alternative actions shown in a popup selector when present */
  readonly variants?: readonly MenuVariant[];
  /** Optional search aliases for fuzzy filter matching */
  readonly keywords?: readonly string[];
}

// --- Flag popup types ---

/** Definition of a single flag field in the flag popup form */
export interface FlagField {
  /** Flag name as passed on the CLI (e.g. "entity_id", "message") */
  readonly name: string;
  /** Display label in the form (e.g. "Entity ID", "Message") */
  readonly label: string;
  /** Field type: string input, boolean toggle, or select dropdown */
  readonly type: "string" | "bool" | "select";
  /** Whether the field must have a value before running */
  readonly required: boolean;
  /** When true, value is appended as a positional arg instead of --name 'value' */
  readonly positional?: boolean;
  /** Pre-filled default value (string for string/select, "true"/"false" for bool) */
  readonly defaultValue?: string;
  /** Available options for "select" type fields */
  readonly options?: readonly string[];
  /** Placeholder text shown in empty string inputs */
  readonly placeholder?: string;
}
