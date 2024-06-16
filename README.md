# IP Information Emailer

## Overview

This Go application fetches IP address information using the IP-API and sends it via email formatted as HTML. It was created out of necessity for a cost-effective solution to access a dynamic IP address remotely without the need for services like No-IP or other paid dynamic DNS services. The app allows the user to receive their IP address by email multiple times a day, enabling remote access without a fixed IP or additional services.

## Features

- **Data Retrieval:** Obtains comprehensive details about an IP address including location, ISP, and network information.
- **HTML Formatting:** Uses Handlebars templates to generate visually appealing email content.
- **Email Notifications:** Sends the formatted data as an HTML email to a predefined list menu of recipients.

## Prerequisites

To run this application, you will need:

- Go 1.15 or newer installed on your machine.
- SMTP server credentials for sending emails (Gmail or any other SMTP server).
- Access to modify and save YAML configuration files.

## Configuration

Before running the application, you need to set up a few configurations:

1. **SMTP Settings:** Configure your SMTP server details in the `config.yaml` file. Here's an example setup for Gmail:

   ```yaml
   email:
     from: "your-email@gmail.com"
     password: "your-app-password"

   smtp:
     host: "smtp.gmail.com"
     port: "587"

   cron:
     schedule: "0 6,12,18,0 * * *" # Executes at 6 AM, 12 PM, 6 PM, and Midnight
   emails: "email1@example.com,email2@example.com"
   ```

   Note: If you are using Gmail, you might need to generate an app-specific password.

2. API Configuration: No API key is needed as the application uses the free version of IP-API. However, ensure you adhere to the IP-API usage terms.

## Usage

Run the application using the following command in your terminal:

```bash
go run . -config path/to/your/config.yaml -cron "0 6,12,18,0 * * *" -emails "email1@example.com,email2@example.com"
```

### Command-line Arguments

- `-config`: Path to the configuration file.
- `-cron`: Cron schedule specification to override the default or configured value.
- `-emails`: Comma-separated list of email addresses to override the default or configured value.

## Downloading the Executable

If you prefer not to compile the code yourself, pre-compiled executables are available for download through GitHub Releases. To download:

1. Go to the Releases section of the repository.
2. Choose the release that suits your operating system and architecture.
3. Download the executable file from the assets.

## Building the Application

To compile the application into an executable, run the following command in the root directory of the project:

```bash
go build -o ipemailer
```

You can then execute the compiled binary directly:

```bash
./ipemailer -config path/to/your/config.yaml -emails "email1@example.com,email2@example.com"
```

## Troubleshooting

_Emails not sending?_ Check your SMTP settings and ensure that the email and password are correctly entered.
_HTML not rendering correctly?_ Ensure that your template in the _`renderTemplate`_ function is correctly formatted as HTML.

### Contributing

Feel free to fork the repository and submit pull requests. You can also open issues if you find bugs or have feature requests.

### License

This project is licensed under the MIT License. See the LICENSE file for more details.
