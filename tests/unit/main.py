# main.py

import os
import logging
import yaml
from typing import Dict

# Set up logging
logging.basicConfig(level=logging.INFO)

# Load configuration from YAML file
with open("config.yaml", "r") as file:
    config: Dict[str, str] = yaml.safe_load(file)

# Define a function to perform the main operation
def main() -> None:
    # Check if the necessary environment variables are set
    if not os.environ.get("API_KEY") or not os.environ.get("SECRET_KEY"):
        logging.error("Environment variables 'API_KEY' and 'SECRET_KEY' are required.")
        exit(1)

    # Perform the main operation using the loaded configuration and environment variables
    # Replace this with the actual operation
    logging.info("Performing main operation...")
    logging.info(f"Using API key: {os.environ.get('API_KEY')}")
    logging.info(f"Using secret key: {os.environ.get('SECRET_KEY')}")

if __name__ == "__main__":
    main()