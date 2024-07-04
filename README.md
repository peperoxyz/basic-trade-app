<p align="center"><img src="https://miro.medium.com/v2/resize:fit:600/1*i2skbfmDsHayHhqPfwt6pA.png" width="400" alt="BasicTrade Logo"></p>

# BasicTrade

BasicTrade is a backend application designed to manage products and their variants, with robust authorization, authentication, and validation mechanisms. It leverages Gin Gonic for routing, GORM for ORM, JWT for authentication, and Cloudinary for image uploads.

## Features

- **Admin Management**
  - Create, update, and delete products
  - Upload product photos to Cloudinary
  - Create, update, and delete product variants
- **Authorization & Authentication**
  - Only the user who created a product can modify it
  - JSON Web Token (JWT) based authentication
- **Validation**
  - Field validations for all tables
  - Image upload validations (format and size)
- **Search & Pagination**
  - Search products by name
  - Search variants by name
  - Pagination for products and variants

## Endpoints

### Admin

| Endpoint | Method | Description                              | Authorization |
| -------- | ------ | ---------------------------------------- | ------------- |
| /login   | POST   | Authenticate admin and receive JWT token | No            |

### Products

| Endpoint        | Method | Description                                                    | Authorization |
| --------------- | ------ | -------------------------------------------------------------- | ------------- |
| /products       | GET    | Retrieve all products with optional name search and pagination | No            |
| /products/:uuid | GET    | Retrieve product details by UUID                               | No            |
| /products       | POST   | Create a new product                                           | Yes           |
| /products/:uuid | PUT    | Update product details                                         | Yes           |
| /products/:uuid | DELETE | Delete a product                                               | Yes           |

### Variants

| Endpoint                 | Method | Description                                                    | Authorization |
| ------------------------ | ------ | -------------------------------------------------------------- | ------------- |
| /variants                | GET    | Retrieve all variants with optional name search and pagination | No            |
| /variants/:uuid          | GET    | Retrieve variant details by UUID                               | No            |
| /products/:uuid/variants | POST   | Create a new variant for a product                             | Yes           |
| /variants/:uuid          | PUT    | Update variant details                                         | Yes           |
| /variants/:uuid          | DELETE | Delete a variant                                               | Yes           |

## Tech Stack

- **Framework**: Gin Gonic
- **ORM**: GORM
- **Authentication**: JWT
- **Image Upload**: Cloudinary
- **Database**: PostgreSQL

## Required Libraries

- `github.com/golang-jwt/jwt/v5`
- `golang.org/x/crypto`
- `github.com/cloudinary/cloudinary-go/v2`

## Getting Started

### Prerequisites

- Go 1.16+
- PostgreSQL
- Cloudinary Account

## Usage

To use BasicTrade, follow these steps:

1. **Access the Application**

   - Access the application locally at [http://localhost:7070](http://localhost:7070).

2. **Interact with Endpoints**

   - Use tools like Postman or any REST client to interact with the endpoints.
   - Refer to the [Endpoints section](#endpoints) in the README for a list of available endpoints and their descriptions.

3. **Authentication**

   - Certain endpoints require authentication using JSON Web Tokens (JWT). Make sure to obtain a JWT token by logging in via the `/login` endpoint before accessing authenticated endpoints.

4. **Testing**

   - Test different scenarios by sending requests to various endpoints.
   - Ensure proper handling of validations, authorization, and error responses as per the API documentation.

5. **Deployment**

   - Deploy the application to your preferred environment for production use. Ensure all necessary configurations (database, Cloudinary credentials, etc.) are correctly set up.

6. **Feedback and Contributions**
   - Provide feedback or suggestions for improvement by raising an issue or submitting a pull request on the project's repository.
