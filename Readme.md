# allcaps
AllCaps is a project that allows users to book classes for various activities. It provides endpoints for booking classes, canceling bookings, adding users. This project is designed to manage class bookings efficiently.

## Features

- **Class Booking**: Users can book classes for different activities.
- **Cancellation**: Users can cancel their bookings if needed.
- **User Management**: Users can be added to the system.
- **Waitlisting**: If a class is full, users are added to a waitlist.

## Setup

To run the AllCaps project locally,make sure  go1.22.4 is installed and then follow these steps:

1. **Clone the repository**:
   ```bash
   git clone https://github.com/dhiranj/allcaps.git

2. **change directory**:
   ```bash
   cd allcaps

3. **build the binary**:
   ```bash
   go build allcaps

4. **run the project**:
   ```bash
   ./allcaps

5. **To test open a new terminal**:
   ```bash
   cd tests
   go test
