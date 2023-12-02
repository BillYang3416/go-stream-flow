# go-stream-flow

## Components:
1. **Frontend App - FileFlowUI**: 
   - A sleek Angular application where users upload files and enter email information.
   - It interacts with the backend to initiate the file and email processing workflow.

2. **First Golang Application - GoFileGate**: 
   - This app acts as the primary gateway for processing user uploads. 
   - It receives files and email data from FileFlowUI, stores the records, and then sends this information to a RabbitMQ broker.

3. **Second Golang Application - GoMailDistributor**: 
   - This application is responsible for consuming messages from the RabbitMQ broker.
   - It sends the file to the specified email address and posts success or failure messages to another queue on the broker.

## Workflow:
- The user uploads a file and inputs email details on FileFlowUI.
- FileFlowUI sends this data to GoFileGate, which records it and forwards it to RabbitMQ.
- GoMailDistributor picks up the message, processes it (i.e., sends the file to the intended recipient), and updates the job status on another RabbitMQ queue.
- GoFileGate receives the status update and modifies the record accordingly.
- The user is then able to see the status of their file upload and email dispatch on FileFlowUI.

## Why These Names?
- **GoStreamFlow** encapsulates the core functionality of the project — streaming and processing data (files and emails) using Golang.
- **FileFlowUI** reflects the user-facing aspect of the project, emphasizing file uploads and user interactions.
- **GoFileGate** suggests a gateway role, handling initial data intake and queuing.
- **GoMailDistributor** accurately describes the app's role in distributing emails and handling message queue responses.

