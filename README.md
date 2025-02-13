A terminal-based user interface (TUI) application for managing and tracking subscriptions. Built with Go and the `tview` library. The main objective was to create an application where data remains local to help keep track of recurring payments, their costs, and when they're due.

## 
<img width="249" alt="image" src="https://github.com/user-attachments/assets/2a6e646f-76d6-409d-af89-c0c43638ec79" />
<img width="605" alt="image" src="https://github.com/user-attachments/assets/0b8b185f-ae3e-4c1f-a21e-a7f87056a48b" />
<img width="809" alt="image" src="https://github.com/user-attachments/assets/f1c19972-7df2-41ad-a580-4716d62d9a0f" />


## Usage

Run the application:
```bash
./subscription-tracker
```
or
```bash
go run main.go
```

### Navigation

- Use arrow keys to navigate through menus
- Press the highlighted character (e.g., 'a' for Add) for quick actions
- Use Tab/Shift+Tab to move between form fields
- Press Enter to select/confirm
- Use ESC to go back/cancel in most contexts

### Available Actions

- **Add Subscription (a)**: Create a new subscription entry
- **List Subscriptions (l)**: View and manage existing subscriptions
- **Quit (q)**: Exit the application

## Dependencies

- [github.com/rivo/tview](https://github.com/rivo/tview): Terminal UI library

## Coming up
~~There's no persistance at the moment and data gets wiped on application close, which is fine for a preview but defeats the practical purpose. I haven't decided yet whether to stick with a simple solution like storing data in a JSON file or use something more industry-standard like sqlite. Both have their advantages, so TBD.~~ Opted for JSON-based storage since the data model is relatively simple and there would be no need for additional dependencies compared to what SQLite would require (e.g. database drivers). There's also negligible performance difference between both for such a simple and small data model (most people would not even reach double-digit subscriptions).

Next up would be to continue extensively testing for logical/validation errors and possibly re-write the UI to use [github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) which was the original plan but opted for `tview` for the preview due to familiarity
