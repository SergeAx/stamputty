# Requirements Specification: PuTTY Session Manager (Go + walk)

## 1. Project Overview

The objective is to create a Windows desktop application in Go using the `lxn/walk` library. The app allows users to synchronize specific configuration settings from the PuTTY "Default Settings" to individual saved sessions by reading and writing directly to the Windows Registry.

## 2. Technical Stack

* **Language:** Go (latest stable).
* **UI Framework:** `github.com/lxn/walk` and `github.com/lxn/win`.
* **Target OS:** Windows 10/11.
* **Permissions:** The app must be able to read/write to `HKEY_CURRENT_USER`.

## 3. Data & Registry Logic

* **Registry Root Path:** `HKEY_CURRENT_USER\Software\SimonTatham\PuTTY\Sessions`
* **Default Settings Path:** `HKEY_CURRENT_USER\Software\SimonTatham\PuTTY\Sessions\Default%20Settings`
* **Session Naming:** Note that PuTTY URL-encodes session names in the registry (e.g., a space is `%20`). The app should decode these for display and encode them when accessing the registry.
* **Sensitive Keys (EXCLUSION LIST):** Do not display or allow copying of the following keys:
* `UserName`
* `PublicKeyFile`
* `ProxyUsername`
* `ProxyPassword`
* `LocalProxyCommand`
* `Hostname`

## 4. UI Layout & Component Requirements

### A. Left Pane (Session List)

* **Component:** `ListBox` or `TableView`.
* **Content:** List all sub-keys under the `Sessions` path.
* **Filter:** Exclude the entry named `Default%20Settings`.
* **Behavior:** Selecting a session triggers the "Diff Logic" to populate the Right Pane.

### B. Right Pane (Comparison View)

* **Component:** `TableView` with Checkboxes.
* **Columns:** 1.  **Setting Name** (e.g., `FontHeight`).
2.  **Default Value** (Value found in `Default%20Settings`).
3.  **Current Value** (Value found in the selected session).
* **Logic:** Only show rows where the `Default Value` differs from the `Current Value`.
* **Visual State:** When a checkbox is checked, the row should be visually marked (e.g., bold text or a status column saying "To be changed").

### C. Action Buttons (Footer)

* **Save Button:** * Iterate through all checked items in the Right Pane.
* Write the values from `Default Settings` into the selected session's registry key.
* Refresh the view upon completion.


* **Cancel Button:** * Uncheck all selections in the Right Pane.
* Clear the "To be changed" status.



## 5. Functional Workflow

1. **Initialization:** On startup, the app scans the Registry. If no PuTTY sessions are found, display a message box and exit.
2. **Selection:** User clicks "Work Server" in the left pane.
3. **Comparison:** * App reads all values from `Default%20Settings`.
* App reads all values from `Work%20Server`.
* App filters out "Sensitive Keys".
* App identifies keys where `Value_Default != Value_Session`.


4. **Modification:** User checks the box for `Colour0` (Background color). The "Save" button becomes enabled.
5. **Commit:** Upon "Save," the app performs a `Registry.SetStringValue` or `SetDWordValue` for the chosen keys.

## 6. Non-Functional Requirements

* **Error Handling:** Provide `walk.MsgBox` alerts if the Registry is locked or permissions are denied.
* **Type Safety:** Properly handle Registry types (`REG_SZ` vs `REG_DWORD`). Most PuTTY settings are strings, but some are integers.
* **Clean Exit:** Ensure registry handles are closed after each read/write operation.
