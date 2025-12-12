package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/smtp"
)

type EmailService interface {
	SendVerificationCode(ctx context.Context, email, code, language string) error
}

type emailService struct {
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	fromEmail    string
	fromName     string
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

type VerificationEmailData struct {
	Code     string
	Language string
	FromName string
}

func NewEmailService(config EmailConfig) EmailService {
	return &emailService{
		smtpHost:     config.SMTPHost,
		smtpPort:     config.SMTPPort,
		smtpUsername: config.SMTPUsername,
		smtpPassword: config.SMTPPassword,
		fromEmail:    config.FromEmail,
		fromName:     config.FromName,
	}
}

func (e *emailService) SendVerificationCode(ctx context.Context, email, code, language string) error {
	// Generate HTML content
	htmlContent, err := e.generateVerificationHTML(code, language)
	if err != nil {
		return fmt.Errorf("failed to generate email HTML: %w", err)
	}

	// Generate plain text content
	textContent := e.generateVerificationText(code, language)

	// Create email message
	subject := e.getSubject(language)
	message := e.createEmailMessage(email, subject, htmlContent, textContent)

	// Send email
	auth := smtp.PlainAuth("", e.smtpUsername, e.smtpPassword, e.smtpHost)
	addr := fmt.Sprintf("%s:%d", e.smtpHost, e.smtpPort)

	err = smtp.SendMail(addr, auth, e.fromEmail, []string{email}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (e *emailService) generateVerificationHTML(code, language string) (string, error) {
	data := VerificationEmailData{
		Code:     code,
		Language: language,
		FromName: e.fromName,
	}

	var htmlTemplate string
	if language == "tr" {
		htmlTemplate = turkishHTMLTemplate
	} else {
		htmlTemplate = englishHTMLTemplate
	}

	tmpl, err := template.New("verification").Parse(htmlTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (e *emailService) generateVerificationText(code, language string) string {
	if language == "tr" {
		return fmt.Sprintf(`
Doğrulama Kodu: %s

Merhaba,

E-posta adresinizi doğrulamak için aşağıdaki 6 haneli kodu kullanın:

%s

Bu kod 10 dakika içinde geçerliliğini yitirecektir.

Eğer bu işlemi siz yapmadıysanız, bu e-postayı görmezden gelebilirsiniz.

Saygılarımızla,
%s Ekibi
`, code, code, e.fromName)
	}

	return fmt.Sprintf(`
Verification Code: %s

Hello,

Please use the following 6-digit code to verify your email address:

%s

This code will expire in 10 minutes.

If you didn't request this verification, you can safely ignore this email.

Best regards,
%s Team
`, code, code, e.fromName)
}

func (e *emailService) getSubject(language string) string {
	if language == "tr" {
		return "E-posta Doğrulama Kodu"
	}
	return "Email Verification Code"
}

func (e *emailService) createEmailMessage(to, subject, htmlContent, textContent string) string {
	boundary := "boundary123456789"

	message := fmt.Sprintf(`From: %s <%s>
To: %s
Subject: %s
MIME-Version: 1.0
Content-Type: multipart/alternative; boundary="%s"

--%s
Content-Type: text/plain; charset=UTF-8
Content-Transfer-Encoding: 7bit

%s

--%s
Content-Type: text/html; charset=UTF-8
Content-Transfer-Encoding: 7bit

%s

--%s--
`, e.fromName, e.fromEmail, to, subject, boundary, boundary, textContent, boundary, htmlContent, boundary)

	return message
}

// Turkish HTML Template
const turkishHTMLTemplate = `
<!DOCTYPE html>
<html lang="tr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>E-posta Doğrulama</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f4f4f4;
        }
        .container {
            background-color: #ffffff;
            padding: 40px;
            border-radius: 10px;
            box-shadow: 0 0 20px rgba(0,0,0,0.1);
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        .logo {
            font-size: 28px;
            font-weight: bold;
            color: #2c3e50;
            margin-bottom: 10px;
        }
        .title {
            color: #34495e;
            font-size: 24px;
            margin-bottom: 20px;
        }
        .code-container {
            background-color: #ecf0f1;
            padding: 20px;
            border-radius: 8px;
            text-align: center;
            margin: 30px 0;
            border-left: 4px solid #3498db;
        }
        .verification-code {
            font-size: 36px;
            font-weight: bold;
            color: #2c3e50;
            letter-spacing: 8px;
            margin: 10px 0;
        }
        .message {
            font-size: 16px;
            line-height: 1.8;
            color: #555;
            margin-bottom: 20px;
        }
        .warning {
            background-color: #fff3cd;
            border: 1px solid #ffeaa7;
            color: #856404;
            padding: 15px;
            border-radius: 5px;
            margin: 20px 0;
        }
        .footer {
            text-align: center;
            margin-top: 40px;
            padding-top: 20px;
            border-top: 1px solid #eee;
            color: #777;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">{{.FromName}}</div>
            <h1 class="title">E-posta Doğrulama</h1>
        </div>
        
        <div class="message">
            <p>Merhaba,</p>
            <p>E-posta adresinizi doğrulamak için aşağıdaki 6 haneli kodu kullanın:</p>
        </div>
        
        <div class="code-container">
            <div class="verification-code">{{.Code}}</div>
        </div>
        
        <div class="warning">
            <strong>Önemli:</strong> Bu kod 10 dakika içinde geçerliliğini yitirecektir.
        </div>
        
        <div class="message">
            <p>Eğer bu işlemi siz yapmadıysanız, bu e-postayı görmezden gelebilirsiniz.</p>
        </div>
        
        <div class="footer">
            <p>Saygılarımızla,<br>{{.FromName}} Ekibi</p>
        </div>
    </div>
</body>
</html>
`

// English HTML Template
const englishHTMLTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Verification</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f4f4f4;
        }
        .container {
            background-color: #ffffff;
            padding: 40px;
            border-radius: 10px;
            box-shadow: 0 0 20px rgba(0,0,0,0.1);
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        .logo {
            font-size: 28px;
            font-weight: bold;
            color: #2c3e50;
            margin-bottom: 10px;
        }
        .title {
            color: #34495e;
            font-size: 24px;
            margin-bottom: 20px;
        }
        .code-container {
            background-color: #ecf0f1;
            padding: 20px;
            border-radius: 8px;
            text-align: center;
            margin: 30px 0;
            border-left: 4px solid #3498db;
        }
        .verification-code {
            font-size: 36px;
            font-weight: bold;
            color: #2c3e50;
            letter-spacing: 8px;
            margin: 10px 0;
        }
        .message {
            font-size: 16px;
            line-height: 1.8;
            color: #555;
            margin-bottom: 20px;
        }
        .warning {
            background-color: #fff3cd;
            border: 1px solid #ffeaa7;
            color: #856404;
            padding: 15px;
            border-radius: 5px;
            margin: 20px 0;
        }
        .footer {
            text-align: center;
            margin-top: 40px;
            padding-top: 20px;
            border-top: 1px solid #eee;
            color: #777;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">{{.FromName}}</div>
            <h1 class="title">Email Verification</h1>
        </div>
        
        <div class="message">
            <p>Hello,</p>
            <p>Please use the following 6-digit code to verify your email address:</p>
        </div>
        
        <div class="code-container">
            <div class="verification-code">{{.Code}}</div>
        </div>
        
        <div class="warning">
            <strong>Important:</strong> This code will expire in 10 minutes.
        </div>
        
        <div class="message">
            <p>If you didn't request this verification, you can safely ignore this email.</p>
        </div>
        
        <div class="footer">
            <p>Best regards,<br>{{.FromName}} Team</p>
        </div>
    </div>
</body>
</html>
`
