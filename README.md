# Modbus Scanner Go
This program is a Go application that scans a Modbus device and reads various data types.

## Features:

Reads Coils, Discrete Inputs, Holding Registers, and Input Registers.
Scans a range of addresses (0 to 65534).
Handles potential "illegal" address errors.
Outputs successful reads and errors to separate JSON files.

## How to use:

Run the program: 
Follow the prompts to enter the serial port, baud rate, and slave ID of your Modbus device.
The program will scan the device and write the results to output.json and errors.json.

## Output:

 - `output.json`: Contains a JSON object with successful reads for each function code (Coils, Discrete Inputs, Holding Registers, Input Registers). The keys are function code names, and the values are objects mapping addresses to their corresponding read values.
 - `errors.json`: Contains a JSON object listing any errors encountered during the scan. The keys are formatted as "function_code:address" and the values are the error messages.


## Notes:

This program is for educational purposes and may require adjustments for specific Modbus devices.
Ensure proper communication settings (serial port, baud rate, slave ID) match your Modbus device.
