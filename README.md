# StamPuTTY: PuTTY "Default Session" settings sync

A lightweight Windows utility built in **Go** to help you stop repeating
yourself. If you’ve ever changed your "Default Settings" in [PuTTY][putty]
(like switching to a high-contrast theme or a larger font) and realized none
of your saved sessions inherited those changes, this tool is for you.

[putty]: https://www.chiark.greenend.org.uk/~sgtatham/putty/

> **Note:** This tool modifies `HKEY_CURRENT_USER`. No Administrator privileges
are required, but it's always a good idea to export a backup of your PuTTY
registry branch before performing updates:

```cmd
reg export "HKEY_CURRENT_USER\Software\SimonTatham\PuTTY\Sessions" "PuTTYSessionsBackup.reg"
```

## ️The "Don't Break My SSH" Guarantee

Directly editing the Windows Registry can be spicy. To keep your workflow safe,
this app **automatically excludes** sensitive session-specific keys from being
overwritten:

* **Hostnames & IP Addresses**
* **Usernames**
* **SSH Private Key paths**
* **Proxy Credentials**

---

## Features

* **Session Discovery:** Automatically finds all saved PuTTY sessions in your
  Registry.
* **The "Diff" Engine:** Select a session on the left to see exactly how it
  differs from your **Default Settings**.
* **Granular Control:** Use checkboxes to pick exactly which settings (colors,
  fonts, scrolling behavior) you want to sync.
* **Batch Update:** Apply changes to the Registry with a single click.
* **Safety First:** A "Cancel" button to reset your pending changes before you
  commit.

---

## How to Use

1. **Launch the App:** You'll see your sessions listed in the left sidebar.
2. **Select a Session:** Click "Production-DB" or whatever you've named your
   connection.
3. **Compare:** The right pane will populate with a list of settings where your
   session differs from the "Default Settings."
4. To see all the settings instead of only the changed ones, just in case, check
   the "Show unchanged settings" checkbox.
5. **Mark for Change:** Check the boxes for the settings you want to update
   (e.g., `FontHeight` or `Colour2`), or hit **Select All** to select all
   diferent settings (not recommended).
6. **Apply changes:** Hit **Apply** to write the selected Default Settings
   values into that specific session's Registry entry.

---

## Building yourself

You need to have a Go 1.25+ environment set up.

Since `tailscale/walk` requires a Windows Manifest to use modern common controls,
you’ll need to generate a `.syso` file before building:

```cmd
go tool rsrc -manifest stamputty.manifest -o rsrc.syso
```

Then build the app:

```cmd
go build -ldflags="-H windowsgui" -o StamPuTTY.exe
```

**Cross-platform build:** if you know what cross-platform building in Go is,
you know how to do it. You will need to generate platform-specific `.syso`
files for each target platform, ex.:

```cmd
go tool rsrc -manifest stamputty.manifest -arch arm64 -o rsrc_windows_arm64.syso
```

(and, of course, remove the catch-all `rsrc.syso` file)

---

## Acknowledgments

* **Simon Tatham:** For creating and maintaining **PuTTY**, the indispensable
  tool this utility was built to support.
* **The PuTTY Team:** For their decades of work on the suite of tools we all
  rely on.
* This project is not affiliated with, or endorsed by, the official PuTTY
  project.


## License

[MIT](LICENSE)
