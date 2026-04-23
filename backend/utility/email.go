package utility

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to string, otp string) error {
	from := os.Getenv("SMTP_EMAIL")        // sender "From" address (your Gmail)
	username := os.Getenv("SMTP_USERNAME") // Brevo SMTP login (a802b7001@smtp-brevo.com)
	password := os.Getenv("SMTP_PASSWORD") // Brevo SMTP key
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	fromName := os.Getenv("SMTP_FROM_NAME")

	if fromName == "" {
		fromName = "VitaTrack.AI"
	}
	// Fall back to from address if no dedicated username is set
	if username == "" {
		username = from
	}

	subject := "Subject: Verify Your Email - VitaTrack.AI\r\n"
	mime := "MIME-version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n"
	from_header := fmt.Sprintf("From: %s <%s>\r\n", fromName, from)
	to_header := fmt.Sprintf("To: %s\r\n", to)

	body := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"></head>
<body style="margin:0;padding:0;background-color:#f0f4f8;font-family:Arial,sans-serif;">
  <table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f0f4f8;padding:40px 20px;">
    <tr>
      <td align="center">
        <table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 4px 20px rgba(0,0,0,0.08);">

          <!-- Header -->
          <tr>
            <td style="background:linear-gradient(135deg,#2563eb,#1d4ed8);padding:36px 40px;text-align:center;">
              <div style="font-size:28px;font-weight:700;color:#ffffff;letter-spacing:1px;">
                + VitaTrack.AI
              </div>
              <div style="font-size:13px;color:rgba(255,255,255,0.75);margin-top:6px;">
                Your Personal Health Intelligence
              </div>
            </td>
          </tr>

          <!-- Body -->
          <tr>
            <td style="padding:40px;">
              <h2 style="margin:0 0 12px;font-size:22px;color:#1e293b;">Verify Your Email Address</h2>
              <p style="margin:0 0 24px;font-size:15px;color:#475569;line-height:1.6;">
                Thanks for signing up! Use the OTP below to verify your email and activate your VitaTrack.AI account.
              </p>

              <!-- OTP Box -->
              <div style="background:#f8fafc;border:2px dashed #2563eb;border-radius:10px;padding:28px;text-align:center;margin-bottom:28px;">
                <div style="font-size:11px;font-weight:600;color:#64748b;letter-spacing:2px;text-transform:uppercase;margin-bottom:12px;">
                  Your One-Time Password
                </div>
                <div style="font-size:42px;font-weight:700;letter-spacing:12px;color:#2563eb;font-family:monospace;">
                  %s
                </div>
                <div style="font-size:13px;color:#94a3b8;margin-top:12px;">
                  Expires in <strong>5 minutes</strong>
                </div>
              </div>

              <p style="margin:0;font-size:13px;color:#94a3b8;line-height:1.6;">
                If you didn't create an account with VitaTrack.AI, you can safely ignore this email.
              </p>
            </td>
          </tr>

          <!-- Footer -->
          <tr>
            <td style="background:#f8fafc;padding:24px 40px;border-top:1px solid #e2e8f0;text-align:center;">
              <p style="margin:0;font-size:12px;color:#94a3b8;">
                © 2025 VitaTrack.AI · Built with ❤️ by the ByIITians team
              </p>
            </td>
          </tr>

        </table>
      </td>
    </tr>
  </table>
</body>
</html>`, otp)

	message := []byte(from_header + to_header + subject + mime + "\r\n" + body)

	// Brevo authenticates with SMTP_USERNAME, but sends from SMTP_EMAIL
	auth := smtp.PlainAuth("", username, password, host)

	return smtp.SendMail(host+":"+port, auth, from, []string{to}, message)
}
