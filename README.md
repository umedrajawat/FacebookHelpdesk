## FacebookHelpDesk
FacebookHelDesk is a web application designed for efficiently managing Facebook pages' help desks. It utilizes React.js in the frontend and Golang with the Gin framework in the backend.

## Purpose
The primary purpose of this project is to streamline the management of help desk operations for Facebook pages. When a user sends a message to the page, it should seamlessly reflect in the help desk interface, enabling prompt responses and efficient user interaction. The application ensures that messages sent by clients are visible, and any replies from the help desk are promptly reflected in Facebook Messenger and vice versa.

## Setup
To set up Facebook integration, follow these steps:

- Create a Developers Account: Begin by creating a developers account on Facebook and create a new project.
- Configure Login URL and Webhook: In your Facebook project settings, configure the login URL and webhook.
- Note: Facebook mandates HTTPS URLs for login and webhook configurations. To work locally, set up Ngrok to expose localhost ports as HTTPS URLs.
- Flows
The application follows these essential flows:

## User Management:
- New users can be created at the /createUser page.
Existing users can log in at the /login page and subsequently connect to Facebook.
- Message Handling:
Client messages sent on Facebook are received via webhooks in the backend.
- These messages are seamlessly transferred to the UI in real-time using websockets.
- Upon opening a chat in the UI, a websocket connection is established to facilitate message reception.
- Data Storage
MongoDB serves as the primary database for storing relevant application data.

## API endpoints

-- /auth/login

-- /webhook/

-- /create_user/

-- /messages/getAllMessages/

-- /auth/get-user/

-- /messages/sendMessage

This README provides a concise overview of FacebookHelDesk, elucidating its purpose, setup instructions, fundamental flows, and data storage mechanisms. For more comprehensive documentation and support, kindly refer to the project's detailed documentation or reach out to the project maintainers.
