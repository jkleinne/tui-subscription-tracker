# Plan to Improve the Subscription Tracker Application

## Overview

After reviewing the application's code, the following logical errors and bugs have been identified:

1. **Persistent Storage Not Utilized in UI**

   - While `JSONStorage` is implemented in `storage/json_storage.go`, the application currently uses in-memory storage (`storage.NewMemoryStorage()` in `ui/ui.go`), causing all subscription data to be lost when the application closes.

2. **Input Parsing Without Error Handling**

   - In `ui/saveSubscription()` and `ui/showEditForm()`, input parsing for `cost` and `totalPayments` lacks error handling. Invalid inputs may default to zero without notifying the user.

3. **Validation of Next Payment Date**

   - The application prevents users from entering a `NextPaymentDate` in the past, hindering the addition of existing subscriptions with past start dates.

4. **Subscription Status Updates**

   - There is no mechanism to automatically update subscriptions based on the current date. Subscriptions do not progress unless `ProcessPayment()` is invoked manually.

5. **Payment Frequency Input Validation**

   - The `PaymentFrequency` input field allows any string, leading to potential invalid inputs.

6. **User Feedback on Errors**

   - Error messages may lack clarity and sufficient context for the user.

## Proposed Solutions

1. **Integrate Persistent Storage into the UI**

   - Replace `storage.NewMemoryStorage()` with `storage.NewJSONStorage("data/subscriptions.json")` in `ui/ui.go`.

   - Ensure the application correctly loads existing data at startup and saves data on changes.

2. **Add Error Handling for Input Parsing**

   - After parsing `cost` and `totalPayments`, check for parsing errors and validate that the values are greater than zero.

   - Provide user feedback if inputs are invalid.

3. **Allow Past Dates for Next Payment Date**

   - Modify validation in `models.NewSubscription()` to allow `NextPaymentDate` in the past.

   - Alternatively, introduce a separate `StartDate` field for when the subscription originally began.

4. **Automate Subscription Status Updates**

   - Implement a function that updates subscriptions during listing, checking if payments are due based on the current date.

   - Automatically call `ProcessPayment()` as needed when the subscription is accessed.

5. **Improve Payment Frequency Input**

   - Change the `PaymentFrequency` input field to a drop-down selection with predefined valid options (`daily`, `weekly`, `monthly`, `yearly`).

6. **Enhance Error Messages**

   - Review and revise error messages to be more descriptive and provide actionable guidance.

## Implementation Steps

1. **Integrate Persistent Storage**

   - **Update Initialization in `ui/ui.go`:**

     - Replace the initialization of `storage` in `NewUI` from:

       ```go
       storage: storage.NewMemoryStorage(),
       ```

       to:

       ```go
       storage: initializeStorage(),
       ```

     - Add the `initializeStorage()` function in `ui/ui.go`:

       ```go
       func initializeStorage() storage.Storage {
           dataFilePath := filepath.Join("data", "subscriptions.json")
           jsonStorage, err := storage.NewJSONStorage(dataFilePath)
           if err != nil {
               // Handle the error appropriately
               panic(fmt.Sprintf("Failed to initialize storage: %v", err))
           }
           return jsonStorage
       }
       ```

     - Ensure you import `"path/filepath"` at the top of `ui/ui.go`.

   - **Ensure Data Directory Exists:**

     - The `NewJSONStorage` function already creates necessary directories; verify this behavior during testing.

   - **Test Data Persistence:**

     - Run the application and confirm that subscriptions are saved to and loaded from `data/subscriptions.json`.

2. **Improve Input Parsing and Validation**

   - **Cost Parsing:**

     - After parsing `costStr`, check for errors:

       ```go
       _, err := fmt.Sscanf(costStr, "%f", &cost)
       if err != nil || cost <= 0 {
           ui.showError("Invalid cost. Please enter a positive number.")
           return
       }
       ```

   - **Total Payments Parsing:**

     - After parsing `totalPaymentsStr`, check for errors:

       ```go
       _, err := fmt.Sscanf(totalPaymentsStr, "%d", &totalPayments)
       if err != nil || totalPayments <= 0 {
           ui.showError("Invalid total payments. Please enter a positive integer.")
           return
       }
       ```

3. **Adjust Date Validation**

   - **Allow Past Dates:**

     - In `models/subscription.go`, modify `NewSubscription` to allow past dates by removing or commenting out the check:

       ```go
       // if nextPayment.Before(time.Now()) {
       //     return nil, fmt.Errorf("Next payment date cannot be in the past")
       // }
       ```

   - **Optional Start Date:**

     - Consider adding a `StartDate` field to `Subscription` if tracking the original start date is necessary for future features.

4. **Automate Subscription Updates**

   - **Update Mechanism:**

     - In `ui/showSubscriptions()`, before displaying subscriptions, update each subscription's status:

       ```go
       for _, sub := range subs {
           for sub.NextPaymentDate.Before(time.Now()) {
               sub.ProcessPayment()
               ui.storage.UpdateSubscription(sub.Name, sub)
           }
       }
       ```

     - Ensure that `UpdateSubscription` saves changes to the storage.

5. **Update UI for Payment Frequency**

   - **Change to Drop-Down:**

     - Replace the input field for `Payment Frequency` with a drop-down:

       ```go
       .AddDropDown("Payment Frequency", []string{"daily", "weekly", "monthly", "yearly"}, 0, nil)
       ```

     - Adjust form submission handlers to retrieve the selected value:

       ```go
       _, frequency := ui.form.GetFormItem(2).(*tview.DropDown).GetCurrentOption()
       ```

6. **Enhance Error Messages**

   - **Consistent Error Handling:**

     - Ensure all user input errors are caught and displayed using `ui.showError()`.

   - **Descriptive Messages:**

     - Revise error messages to be clear and helpful, indicating what went wrong and how to fix it.

## Additional Recommendations

- **Form Field Validation:**

  - Implement real-time validation for form fields if possible.

- **Code Refactoring:**

  - Review the codebase for repeated patterns and consider refactoring to reduce duplication.

- **User Experience Enhancements:**

  - Add confirmations for actions like saving, updating, or deleting subscriptions.

  - Enhance navigation options to improve usability.

## Conclusion

Implementing these improvements will enhance the application's functionality, user experience, and reliability. Users will be able to retain their data between sessions, have a smoother interaction with the application, and encounter fewer errors.