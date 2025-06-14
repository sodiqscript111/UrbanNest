UrbanNest - Modern Accommodation Booking Platform
🏡 Project Overview
UrbanNest is a robust and scalable full-stack accommodation booking platform designed to connect property owners (Listers) with customers seeking short-term and long-term stays. This project showcases a comprehensive backend API built with Go, offering secure authentication, listing management, booking functionalities, and seamless integration with external services for enhanced user experience. The frontend is developed using React, providing a dynamic and intuitive interface for users.

This project demonstrates a strong understanding of modern web development principles, including API design, secure authentication, database management, and cloud-based file storage.

✨ Features
For Customers:

Browse Listings: Easily search and filter accommodation listings by location and availability dates.
Detailed Listing View: Access comprehensive information about each property, including descriptions, amenities, and media.
Secure Authentication: Register and log in securely to manage bookings.
Email Verification: Real-time email validation during signup using ZeroBounce API to ensure data quality.
Booking Management: Create and view personal booking history.
Secure Payments: Integrated with Paystack API for reliable and verified payment processing.
For Property Listers:

Lister Registration & Login: Secure accounts for managing properties.
Listing Creation: Add new accommodation listings with detailed information and media.
Listing Management: View, edit, and delete existing listings.
Media Upload: Securely upload listing images and videos directly to AWS S3 using presigned URLs.
🛠️ Technologies & Tools
This project leverages a modern and powerful tech stack to deliver a high-performance and maintainable application.

Backend (Go):

Go: The primary language for building a fast, efficient, and concurrent API.
Gin Web Framework: A high-performance HTTP web framework for Go, used for routing and middleware.
GORM: An excellent ORM library for Go, simplifying database interactions.
PostgreSQL: A powerful, open-source relational database used for storing all application data.
github.com/joho/godotenv: For loading environment variables from .env files.
github.com/go-resty/resty/v2: An easy-to-use HTTP client for making external API requests (e.g., to ZeroBounce and Paystack).
AWS SDK for Go v2 (github.com/aws/aws-sdk-go-v2/service/s3): For interacting with AWS S3 for media storage.
JSON Web Tokens (JWT): Used for secure, stateless authentication.
Bcrypt: For secure password hashing.
ZeroBounce API: Integrated for real-time email verification during user registration.
Paystack API: Integrated for secure and verified payment processing for bookings.
Frontend (React):

React: A declarative, component-based JavaScript library for building user interfaces.
(Likely) React Router: For client-side routing.
(Likely) State Management Library: (e.g., Context API, Redux, Zustand) for managing application state.
HTML5 & CSS3: For structuring and styling the web application.
JavaScript (ES6+): For interactive client-side logic.
Deployment & Infrastructure:

AWS S3: Cloud storage for listing media (images, videos).
Vercel: (Mentioned in CORS configuration) A platform for deploying frontend applications.
Render/Heroku/DigitalOcean (or similar, depending on where the Go backend is deployed): For hosting the Go API. (You can specify your actual deployment platform here if you know it).
Development Tools:

Git & GitHub: For version control and collaborative development.
Postman/Insomnia: For API testing and development.
VS Code: Preferred IDE for development.
🚀 Getting Started
Follow these instructions to get a copy of the project up and running on your local machine for development and testing purposes.

Prerequisites
Go (version 1.18 or higher recommended)
Node.js & npm/yarn (for the React frontend)
PostgreSQL
An AWS account (for S3 storage)
ZeroBounce API Key
Paystack Secret Key
Installation (Backend)
Clone the repository:

Bash

git clone https://github.com/your-username/UrbanNest.git
cd UrbanNest
Set up environment variables:
Create a .env file in the root of the backend directory and add the following:

Code snippet

DB_HOST=localhost
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=urbannest_db
DB_PORT=5432
JWT_SECRET=your_jwt_secret_key
ZERBOUNCE_API_KEY=your_zerobounce_api_key
PAYSTACK_SECRET_KEY=your_paystack_secret_key
AWS_REGION=your_aws_region # e.g., us-east-1
AWS_ACCESS_KEY_ID=your_aws_access_key_id
AWS_SECRET_ACCESS_KEY=your_aws_secret_access_key
PORT=8080 # Or any desired port
Replace placeholders with your actual credentials.

Initialize the database:
Ensure your PostgreSQL server is running. The db.InitDB() function in main.go will handle database migration based on your models definitions.

Run the application:

Bash

go mod tidy
go run main.go
The API server will start on the port specified in your .env file (default: 8080).

Installation (Frontend)
Navigate to your React frontend project directory.

Follow the instructions in the frontend's README.md (if separate) to install dependencies and run the development server. Typically:

Bash

cd ../UrbanNest-Frontend # Assuming your frontend is in a separate directory
npm install # or yarn install
npm start # or yarn start
📂 Project Structure (Backend)
The Go backend follows a modular structure for better organization and maintainability:

UrbanNest/
├── main.go               # Entry point of the application, sets up routes and starts the server.
├── db/                   # Database connection and migration logic.
│   └── db.go
├── middleware/           # Custom middleware (e.g., authentication).
│   └── auth.go
├── models/               # Defines database models and their associated CRUD operations.
│   ├── booking.go
│   ├── customer.go
│   ├── lister.go
│   ├── listing.go
│   └── media.go
└── utils/                # Utility functions (e.g., JWT token generation, password hashing).
    ├── auth.go
    └── helpers.go
💡 API Endpoints
Here's a summary of the key API endpoints:

Public Endpoints:

GET /listings: Retrieve all listings, with optional filters for location, check-in, and check-out dates.
GET /listings/:id: Retrieve a specific listing by ID.
POST /signup: Register a new customer.
POST /signuplister: Register a new property lister.
POST /login: Authenticate a customer.
POST /loginlister: Authenticate a property lister.
POST /booking: Create a new booking (includes Paystack payment verification).
GET /getmybookings/:id: Retrieve bookings for a specific customer.
GET /api/upload-url: Generate a presigned URL for direct S3 uploads.
POST /verify-email: Verify an email address using ZeroBounce.
Authenticated Endpoints (require JWT in Authorization header):

POST /listings: Create a new listing.
GET /mylisting: Retrieve listings associated with the authenticated lister.
PUT /listings/:id: Update an existing listing by ID.
DELETE /listings/:id: Delete a listing by ID.
🤝 Contributing
Contributions are welcome! If you have suggestions for improvements, new features, or bug fixes, please open an issue or submit a pull request.
