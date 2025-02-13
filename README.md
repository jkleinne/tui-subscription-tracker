A terminal-based user interface (TUI) application for managing and tracking subscriptions. Built with Go and the `tview` library. The main motivation was to create a privacy-first application where data remains local to help keep track of recurring payments, their costs, and when they're due.

## 
https://github.com/user-attachments/assets/d56e9a84-fd6b-4d17-bd25-3e679e9163d8



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
~~There's no persistance at the moment and data gets wiped on application close, which is fine for a preview but defeats the practical purpose. I haven't decided yet whether to stick with a simple solution like storing data in a JSON file or use something more industry-standard like sqlite. Both have their advantages, so TBD.~~ 

Opted for JSON-based storage since the data model is relatively simple and there would be no need for additional dependencies compared to what SQLite would require (e.g. database drivers). There's also negligible performance difference between both for such a simple and small data model (most people would not even reach double-digit subscriptions).

Next up would be to continue fleshing out the current implementation with more features (including some logical checks that was ommitted for the sake of the preview). For example, could maybe even add an optional background task for notifications/reminders, validation is also minimal at the moment and the automatic payment processing (when the next payment date occurs) hasn't been implemented yet, etc.
