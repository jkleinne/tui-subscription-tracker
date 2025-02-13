A terminal-based user interface (TUI) application for managing and tracking subscriptions. Built with Go and the `tview` library. The main objective was to create an application where data remains local to help keep track of recurring payments, their costs, and when they're due.

## 
<img width="249" alt="image" src="https://github.com/user-attachments/assets/2a6e646f-76d6-409d-af89-c0c43638ec79" />
<img width="486" alt="image" src="https://github.com/user-attachments/assets/c2eaa3e5-54f0-4330-9f18-382cd74e4a7c" />
<img width="559" alt="image" src="https://github.com/user-attachments/assets/0b87917e-acee-4b34-a9e1-b3e95732a7bf" />

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
There's no persistance at the moment and data gets wiped on application close, which is fine for a preview but defeats the practical purpose. I haven't decided yet whether to stick with a simple solution like storing data in a JSON file or use something more industry-standard like sqlite. Both have their advantages, so TBD.
