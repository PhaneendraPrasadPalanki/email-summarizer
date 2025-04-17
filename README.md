# 📬 Smart Email Summarizer (GoLang)

This project connects to your Gmail, fetches recent **Primary tab** emails, summarizes them, and:
- Displays summaries on a clean dashboard
- (Optional) Emails you a daily PDF report

## 🛠 Tech Stack
- GoLang
- Gmail API (OAuth2)
- HTML Dashboard
- PDF Generator (`gofpdf`)
- Email Sender (`jordan-wright/email`)

## 🚀 How to Run

See [How to Run](How_to_run.txt) for step-by-step setup instructions.

## 🔐 Secrets
- `.token.json` and `credentials.json` are ignored for safety.
- Use environment variables to store keys (e.g., `OPENAI_API_KEY`).

---

## 📧 Coming Soon
- Automatic daily summary emails
- Real summarization via OpenAI
- Better UI filters & styling
