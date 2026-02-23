package email

import "fmt"

// SendWelcomeEmail sends a welcome email with email verification link and API key to a new user.
func (s *Sender) SendWelcomeEmail(to, username, verificationURL string) error {
	subject := "Welcome to Safecast - Verify Your Email"

	displayName := username
	if displayName == "" {
		displayName = to
	}

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .button {
            display: inline-block;
            padding: 12px 24px;
            background-color: #4caf50;
            color: white !important;
            text-decoration: none;
            border-radius: 4px;
            margin: 16px 0;
        }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Welcome to Safecast!</h1>
        <p>Hi %s,</p>
        <p>Thank you for registering with Safecast. We're excited to have you join our community of radiation monitoring enthusiasts.</p>
        <p>Please verify your email address to complete your registration and start uploading radiation measurement data:</p>
        <p><a href="%s" class="button">Verify Email Address</a></p>
        <p>Or copy and paste this link into your browser:</p>
        <p style="word-break: break-all; color: #666;">%s</p>
        <p>This verification link will expire in 24 hours.</p>
        <p>If you didn't create this account, please ignore this email.</p>
        <div class="footer">
            <p>This is an automated message from Safecast. Please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
	`, displayName, verificationURL, verificationURL)

	textBody := fmt.Sprintf(`
Welcome to Safecast!

Hi %s,

Thank you for registering with Safecast. We're excited to have you join our community of radiation monitoring enthusiasts.

Please verify your email address to complete your registration and start uploading radiation measurement data:

%s

This verification link will expire in 24 hours.

If you didn't create this account, please ignore this email.

---
This is an automated message from Safecast. Please do not reply to this email.
	`, displayName, verificationURL)

	return s.Send(to, subject, htmlBody, textBody)
}

// SendWelcomeEmailWithAPIKey sends a welcome email with email verification link and API key to a new user.
func (s *Sender) SendWelcomeEmailWithAPIKey(to, username, verificationURL, apiKey string) error {
	subject := "Welcome to Safecast - Verify Your Email"

	displayName := username
	if displayName == "" {
		displayName = to
	}

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .button {
            display: inline-block;
            padding: 12px 24px;
            background-color: #4caf50;
            color: white !important;
            text-decoration: none;
            border-radius: 4px;
            margin: 16px 0;
        }
        .api-key-box {
            background-color: #f5f5f5;
            border-left: 4px solid #4caf50;
            padding: 16px;
            margin: 16px 0;
            font-family: 'Courier New', monospace;
            word-break: break-all;
        }
        .info-box {
            background-color: #e3f2fd;
            border-left: 4px solid #2196f3;
            padding: 12px;
            margin: 16px 0;
        }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Welcome to Safecast!</h1>
        <p>Hi %s,</p>
        <p>Thank you for registering with Safecast. We're excited to have you join our community of radiation monitoring enthusiasts.</p>

        <h2>Your API Key</h2>
        <p>Your unique API key for programmatic access is:</p>
        <div class="api-key-box">%s</div>

        <div class="info-box">
            <strong>Important:</strong> Store this API key securely. You can use it to:
            <ul>
                <li>Log in to your account (along with your email)</li>
                <li>Upload measurement data programmatically</li>
                <li>Access protected API endpoints</li>
            </ul>
            If you lose this key, an admin can regenerate it for you.
        </div>

        <h2>Verify Your Email</h2>
        <p>Please verify your email address to complete your registration and start uploading radiation measurement data:</p>
        <p><a href="%s" class="button">Verify Email Address</a></p>
        <p>Or copy and paste this link into your browser:</p>
        <p style="word-break: break-all; color: #666;">%s</p>
        <p>This verification link will expire in 24 hours.</p>
        <p>If you didn't create this account, please ignore this email.</p>
        <div class="footer">
            <p>This is an automated message from Safecast. Please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
	`, displayName, apiKey, verificationURL, verificationURL)

	textBody := fmt.Sprintf(`
Welcome to Safecast!

Hi %s,

Thank you for registering with Safecast. We're excited to have you join our community of radiation monitoring enthusiasts.

YOUR API KEY
============
%s

IMPORTANT: Store this API key securely. You can use it to:
- Log in to your account (along with your email)
- Upload measurement data programmatically
- Access protected API endpoints

If you lose this key, an admin can regenerate it for you.

VERIFY YOUR EMAIL
=================
Please verify your email address to complete your registration and start uploading radiation measurement data:

%s

This verification link will expire in 24 hours.

If you didn't create this account, please ignore this email.

---
This is an automated message from Safecast. Please do not reply to this email.
	`, displayName, apiKey, verificationURL)

	return s.Send(to, subject, htmlBody, textBody)
}

// SendPasswordSetupEmail sends a password setup email to existing users being migrated.
func (s *Sender) SendPasswordSetupEmail(to, username, setupURL string) error {
	subject := "Set Up Your Safecast Account Password"

	displayName := username
	if displayName == "" {
		displayName = to
	}

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .button {
            display: inline-block;
            padding: 12px 24px;
            background-color: #4caf50;
            color: white !important;
            text-decoration: none;
            border-radius: 4px;
            margin: 16px 0;
        }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Welcome to the New Safecast Map!</h1>
        <p>Hi %s,</p>
        <p>We've migrated your account to our new system. To continue uploading radiation measurements and accessing your account, please create a password:</p>
        <p><a href="%s" class="button">Create Password</a></p>
        <p>Or copy and paste this link into your browser:</p>
        <p style="word-break: break-all; color: #666;">%s</p>
        <p>This link will expire in 7 days.</p>
        <p>Your existing uploads and data are safe and will be linked to your account once you set up your password.</p>
        <div class="footer">
            <p>This is an automated message from Safecast. Please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
	`, displayName, setupURL, setupURL)

	textBody := fmt.Sprintf(`
Welcome to the New Safecast Map!

Hi %s,

We've migrated your account to our new system. To continue uploading radiation measurements and accessing your account, please create a password:

%s

This link will expire in 7 days.

Your existing uploads and data are safe and will be linked to your account once you set up your password.

---
This is an automated message from Safecast. Please do not reply to this email.
	`, displayName, setupURL)

	return s.Send(to, subject, htmlBody, textBody)
}

// SendPasswordResetEmail sends a password reset email.
func (s *Sender) SendPasswordResetEmail(to, resetURL string) error {
	subject := "Reset Your Safecast Password"

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .button {
            display: inline-block;
            padding: 12px 24px;
            background-color: #4caf50;
            color: white !important;
            text-decoration: none;
            border-radius: 4px;
            margin: 16px 0;
        }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; font-size: 12px; color: #666; }
        .warning { background-color: #fff3cd; border-left: 4px solid #ffc107; padding: 12px; margin: 16px 0; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Reset Your Password</h1>
        <p>We received a request to reset the password for your Safecast account.</p>
        <p><a href="%s" class="button">Reset Password</a></p>
        <p>Or copy and paste this link into your browser:</p>
        <p style="word-break: break-all; color: #666;">%s</p>
        <p>This link will expire in 1 hour.</p>
        <div class="warning">
            <strong>Security Note:</strong> If you didn't request this password reset, please ignore this email. Your password will not be changed.
        </div>
        <div class="footer">
            <p>This is an automated message from Safecast. Please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
	`, resetURL, resetURL)

	textBody := fmt.Sprintf(`
Reset Your Password

We received a request to reset the password for your Safecast account.

%s

This link will expire in 1 hour.

SECURITY NOTE: If you didn't request this password reset, please ignore this email. Your password will not be changed.

---
This is an automated message from Safecast. Please do not reply to this email.
	`, resetURL)

	return s.Send(to, subject, htmlBody, textBody)
}

// SendPasswordChangedEmail sends a notification email when password is changed.
func (s *Sender) SendPasswordChangedEmail(to string) error {
	subject := "Your Safecast Password Was Changed"

	htmlBody := `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; font-size: 12px; color: #666; }
        .warning { background-color: #f8d7da; border-left: 4px solid #dc3545; padding: 12px; margin: 16px 0; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Password Changed</h1>
        <p>This email confirms that your Safecast account password was successfully changed.</p>
        <p>If you made this change, you can safely ignore this email.</p>
        <div class="warning">
            <strong>Important:</strong> If you did not make this change, your account may have been compromised. Please contact Safecast support immediately.
        </div>
        <div class="footer">
            <p>This is an automated message from Safecast. Please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
	`

	textBody := `
Password Changed

This email confirms that your Safecast account password was successfully changed.

If you made this change, you can safely ignore this email.

IMPORTANT: If you did not make this change, your account may have been compromised. Please contact Safecast support immediately.

---
This is an automated message from Safecast. Please do not reply to this email.
	`

	return s.Send(to, subject, htmlBody, textBody)
}

// SendAPIKeyRegeneratedEmail sends a notification email when an API key is regenerated.
func (s *Sender) SendAPIKeyRegeneratedEmail(to, username, newAPIKey string) error {
	subject := "Your Safecast API Key Has Been Regenerated"

	displayName := username
	if displayName == "" {
		displayName = to
	}

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .api-key-box {
            background-color: #f5f5f5;
            border-left: 4px solid #4caf50;
            padding: 16px;
            margin: 16px 0;
            font-family: 'Courier New', monospace;
            word-break: break-all;
        }
        .warning { background-color: #fff3cd; border-left: 4px solid #ffc107; padding: 12px; margin: 16px 0; }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <h1>API Key Regenerated</h1>
        <p>Hi %s,</p>
        <p>Your Safecast API key has been regenerated. Your new API key is:</p>
        <div class="api-key-box">%s</div>
        <div class="warning">
            <strong>Important:</strong>
            <ul>
                <li>Your old API key is no longer valid</li>
                <li>Update any scripts or applications that use your API key</li>
                <li>Store this new key securely</li>
            </ul>
        </div>
        <p>You can use this API key to:</p>
        <ul>
            <li>Log in to your account (along with your email)</li>
            <li>Upload measurement data programmatically</li>
            <li>Access protected API endpoints</li>
        </ul>
        <div class="footer">
            <p>This is an automated message from Safecast. Please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
	`, displayName, newAPIKey)

	textBody := fmt.Sprintf(`
API Key Regenerated

Hi %s,

Your Safecast API key has been regenerated. Your new API key is:

%s

IMPORTANT:
- Your old API key is no longer valid
- Update any scripts or applications that use your API key
- Store this new key securely

You can use this API key to:
- Log in to your account (along with your email)
- Upload measurement data programmatically
- Access protected API endpoints

---
This is an automated message from Safecast. Please do not reply to this email.
	`, displayName, newAPIKey)

	return s.Send(to, subject, htmlBody, textBody)
}
