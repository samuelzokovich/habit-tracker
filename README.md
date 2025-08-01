# habit-tracker
![image](https://github.com/user-attachments/assets/9aa3a0c1-4154-45b6-8552-63b822155c33)

# 📈 Habit Tracker

**Habit Tracker** is a lightweight, Go-powered web application that helps users build and maintain healthy habits by tracking their daily progress and sending timely notifications and monthly reports via email.

---

## 🌟 Features

- 📬 **Email Notifications** — Remind users to log their habits each day.
- 🧾 **Web UI for Input** — A clean and simple web interface for submitting habits.
- 📊 **Monthly Reports** — Summary of habit progress sent via email.
- 🔐 **Secure and Lightweight** — Built in Go for high performance and minimal dependencies.

---

## 🚀 How It Works

1. **Daily Notification**  
   The app sends a reminder email to users prompting them to log their habits.

2. **Habit Input**  
   Users visit the web page (linked in the email) to input:
   - What habit they followed
   - Any optional notes or progress details

3. **Data Storage**  
   Habit entries are securely stored in a database with timestamps.

4. **Monthly Report**  
   At the end of each month, a summarized report of the user's habit activity is emailed automatically.

---

## 🛠️ Getting Started

### ⚙️ Prerequisites

- Go 1.20+
- PostgreSQL (or your preferred database)
- SMTP email credentials

### 📦 Installation

```bash
git clone https://github.com/your-org/habit-tracker.git
cd habit-tracker
go build -o habit-tracker
