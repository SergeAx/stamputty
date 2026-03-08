# PuTTY Session Sync

A lightweight Windows utility built in **Go** to help you stop repeating yourself. If you’ve ever changed your "Default Settings" in PuTTY (like switching to a high-contrast theme or a larger font) and realized none of your saved sessions inherited those changes, this tool is for you.

## ⚠️ The "Don't Break My SSH" Guarantee

Directly editing the Windows Registry can be spicy. To keep your workflow safe, this app **automatically excludes** sensitive session-specific keys from being overwritten:

* **Hostnames & IP Addresses**
* **Usernames**
* **SSH Private Key paths**
* **Proxy Credentials**

---

## 🚀 Features

* **Session Discovery:** Automatically finds all saved PuTTY sessions in your Registry.
* **The "Diff" Engine:** Select a session on the left to see exactly how it differs from your **Default Settings**.
* **Granular Control:** Use checkboxes to pick exactly which settings (colors, fonts, scrolling behavior) you want to sync.
* **Batch Update:** Apply changes to the Registry with a single click.
* **Safety First:** A "Cancel" button to reset your pending changes before you commit.

## 🛠 Prerequisites

* **Windows OS** (Obviously, since we're poking the Registry).
* **PuTTY** installed and configured with at least one saved session.
* **Go 1.18+** (If building from source).
* **rsrc:** Required to embed the manifest file for the `lxn/walk` UI.

---

## 🏗 Installation & Building

Since `lxn/walk` requires a Windows Manifest to use modern common controls, you’ll need to generate a `.syso` file before building.

1. **Install the manifest tool:**
```bash
go install github.com/akavel/rsrc@latest

```


2. **Generate the resource file:**
```bash
rsrc -manifest test.manifest -o rsrc.syso

```


3. **Build the app:**
```bash
go build -ldflags="-H windowsgui"

```



---

## 📖 How to Use

1. **Launch the App:** You'll see your sessions listed in the left sidebar.
2. **Select a Session:** Click "Production-DB" or whatever you've named your connection.
3. **Compare:** The right pane will populate with a list of settings where your session differs from the "Default Settings."
4. **Mark for Change:** Check the boxes for the settings you want to update (e.g., `FontHeight` or `Colour2`).
5. **Save:** Hit **Save** to write the Default values into that specific session's Registry entry.

---

## 📂 Technical Details

* **Registry Root:** `HKEY_CURRENT_USER\Software\SimonTatham\PuTTY\Sessions`
* **UI Framework:** [lxn/walk](https://github.com/lxn/walk)
* **Encoding:** Handles PuTTY’s URL-encoded session names (e.g., converting `%20` back to spaces for the UI).

> **Note:** This tool modifies `HKEY_CURRENT_USER`. No Administrator privileges are required, but it's always a good idea to export a backup of your PuTTY registry key before performing bulk updates.
