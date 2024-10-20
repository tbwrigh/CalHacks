# AutoLock

AutoLock is a tool designed to automatically secure your code by setting up CodeQL and providing automatic fixes for vulnerabilities.

*Please note, because this project invovles automatically making PRs for insecure code, there may be activiity after the submission period.*

## Features

- **Automated CodeQL Setup**: Simplifies the process of integrating CodeQL into your project.
- **Automatic Vulnerability Fixes**: Identifies and fixes vulnerabilities in your codebase.

## Installation

To install AutoLock, follow these steps:

1. Clone the repository:
    ```sh
    git clone https://github.com/tbwrigh/CalHacks.git
    ```
2. Navigate to the project directory:
    ```sh
    cd CalHacks
    ```
3. Cd to the Frontend 
    ```sh
    cd frontend
    ```
4. Install the required dependencies:
    ```sh
    npm install
    ```
5. Run the frontend
    ```sh
    npm run dev
    ```
6. Change to backend
    ```sh
    cd ../backend
    ```
7. Run Go
    ```sh
    go run main.go
    ```
8. Go to Root Dir
    ```sh
    cd ..
    ```
9. Run docker compose
    ```
    docker compose up
    ```

## Usage

To use AutoLock, run the following command:
```sh
npm start
```

## License

This project is licensed under the MIT License. 

## Contact

For any questions or feedback, please open an issue on GitHub.
