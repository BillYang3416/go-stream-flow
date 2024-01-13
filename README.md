# GoStreamFlow

## Overview
GoStreamFlow is a streamlined system designed for efficient processing of files and emails using Go (Golang). This platform excels in managing file uploads and automating email distribution tasks, offering a comprehensive workflow from user interface to execution.

## Components:
1. **Frontend App - FileFlowUI**: 
   - A modern Angular application where users can upload files and enter email information.
   - Directly interacts with GoFlowGateway for initiating file and email processing.

2. **Golang Application - GoFlowGateway**: 
   - A versatile application that incorporates the functionalities of file processing and email distribution.
   - Handles files and email data from FileFlowUI, manages records, interfaces with RabbitMQ for messaging, and oversees the email distribution process.

## Workflow:
- Users upload a file and input email details on FileFlowUI.
- FileFlowUI sends this data to GoFlowGateway.
- GoFlowGateway processes the data, storing records and forwarding the information to RabbitMQ.
- It then takes charge of message consumption from RabbitMQ, ensuring the file reaches its intended recipient.
- Following this, GoFlowGateway updates the job status on a RabbitMQ queue.
- Users can track the status of their file uploads and email dispatches on FileFlowUI.
