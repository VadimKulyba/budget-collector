# Budget Collector

A Go-based tool for processing and categorizing bank transaction reports from CSV files. This tool helps analyze monthly spending patterns by automatically categorizing transactions and generating organized output reports.

## Features

- **CSV Processing**: Reads bank transaction reports in CSV format with Windows-1251 encoding
- **Automatic Categorization**: Maps bank transaction categories to user-defined spending categories
- **Multi-Payment Method Support**: Handles transactions from different payment cards/accounts
- **Period-Based Reporting**: Processes reports for specific monthly periods (MM.YYYY format)
- **Smart Filtering**: Excludes internal bank transfers and focuses on actual spending
- **Structured Output**: Generates clean, categorized CSV reports for budget analysis

## Requirements

- Go 1.24.5 or higher
- CSV files with Windows-1251 encoding
- Bank reports in the specified format

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd budget-collector
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o budget-collector cmd/main.go
```

## Usage

### Basic Usage

1. Place your bank report CSV files in a `reports/` directory
2. Run the application:
```bash
./budget-collector
```

3. Enter the report period in MM.YYYY format (e.g., "12.2024")
4. The tool will:
   - Find the matching report file
   - Process all transactions
   - Generate an `output.csv` file with categorized spending

### Input File Structure

Your CSV files should contain the following columns:
- **Операция** (Operation): Transaction description
- **Сумма** (Amount): Transaction amount
- **Дата операции по счету** (Transaction Date): Date of the transaction
- **Категория операции** (Operation Category): Bank's category classification

### Output Format

The generated `output.csv` contains:
- Transaction date
- Payment type (card/cash)
- Spending category
- Subcategory (currently empty)
- Amount
- Empty column (for additional data)
- Transaction description

### File Paths

- **Input Reports**: `reports/*.csv`
- **Output File**: `output.csv` (generated in the current directory)

## Dependencies

- `golang.org/x/text`: For handling Windows-1251 character encoding
- Standard Go libraries for CSV processing and file operations

## Error Handling

The application includes error handling for:
- Invalid period formats
- Missing report files
- CSV parsing errors
- File I/O operations

## TODO

- Currently supports only BYN currency
- Subcategory field is not populated
- Support only one input file
