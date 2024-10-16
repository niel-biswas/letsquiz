
```
 /$$                   /$$                     /$$$$$$            /$$           /$$
| $$                  | $$                    /$$__  $$          |__/          | $$
| $$        /$$$$$$  /$$$$$$   /$$$$$$$      | $$  \ $$ /$$   /$$ /$$ /$$$$$$$$| $$
| $$       /$$__  $$|_  $$_/  /$$_____/      | $$  | $$| $$  | $$| $$|____ /$$/| $$
| $$      | $$$$$$$$  | $$   |  $$$$$$       | $$  | $$| $$  | $$| $$   /$$$$/ |__/
| $$      | $$_____/  | $$ /$$\____  $$      | $$/$$ $$| $$  | $$| $$  /$$__/      
| $$$$$$$$|  $$$$$$$  |  $$$$//$$$$$$$/      |  $$$$$$/|  $$$$$$/| $$ /$$$$$$$$ /$$
|________/ \_______/   \___/ |_______/        \____ $$$ \______/ |__/|________/|__/
                                                   \__/
```

![Build Status](https://img.shields.io/github/actions/workflow/status/niel-biswas/letsquiz/main.yml?branch=main)
![License](https://img.shields.io/badge/license-GNU%20GPLv3-blue)

## Description

`LetsQuiz` is a quiz platform built as a Text User Interface (TUI) application in Go, utilizing the Bubble Tea framework, which is inspired by the Elm Architecture. It combines Go's efficiency and scalability to deliver a robust and responsive user experience.

## Table of Contents

- [Getting Started](#getting-started)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)
- [Authors and Acknowledgment](#authors-and-acknowledgment)
- [Project Status](#project-status)

## Getting Started

To make it easy for you to get started with Gihub, here's a list of recommended next steps.

## Add Your Files

To add your files, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/niel-biswas/letsquiz.git
   ```

2. Change into the repository directory:

   ```bash
   cd letsquiz
   ```

3. Add your files and commit the changes:

   ```bash
   git add .
   git commit -m "Added my quiz files"
   ```

4. Push your changes to GitHub:

   ```bash
   git push origin main
   ```

## Installation

To install the project, follow these steps:
```sh
git clone https://github.com/niel-biswas/letsquiz.git
cd letsquiz
go mod tidy
go build
```

## Database ERD

![Database ERD](./images/quiz-ERD.png)


## Application Flow

![App Flow Wire Diagram](./images/AppFlow-WireDiagram.png)


## Usage

To run the `LetsQuiz` application, make sure you start the backend server before launching the application.

### Step 1: Start the Backend Server

The backend server handles all API requests, database interactions, and user management. To start the backend:

1. Navigate to the `server` directory:

    ```bash
    cd server
    ```

2. Build and run the backend server:

    ```bash
    go run main.go
    ```

   This will start the backend server, which listens for incoming requests from the quiz application.

### Step 2: Start the Application

Once the backend server is up and running, return to the project root directory and start the main quiz application:

1. Navigate back to the root directory:

    ```bash
    cd ..
    ```

2. Run the application:

    ```bash
    go run main.go
    ```

Now your quiz application will communicate with the backend server to handle user sessions, quizzes, and leaderboards.

## Contributing

Please refer to the contribution guidelines for more details.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](./LICENSE) file for details.

## Authors and Acknowledgment

- Author: niel-biswas

## Project Status

This project is currently in active development. Contributions are welcome!
