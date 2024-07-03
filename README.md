# TODO API Project

## Overview

This project involves developing a TODO API using Golang and ScyllaDB. The API supports basic CRUD operations, pagination, and filtering based on the TODO item status. Additionally, it includes sorting options for the list endpoint based on the creation date.

## Features

- Create, read, update, and delete TODO items.
- Pagination support for listing TODO items.
- Filtering by status (pending, completed).
- Sorting options for the list based on the creation date.

## Requirements

- Golang
- ScyllaDB

## Project Structure

```bash
.
├── README.md
├── main.go
├── handlers.go
├── models.go
├── go.mod
├── go.sum
```

## Installation
```bash
git clone https://github.com/Swarnim01/SamespaceAssignment.git
cd todo-api
```

```bash
go mod download
```

```bash
go run main.go
```

## Design Decisions:
- Golang and ScyllaDB: Golang was chosen for its performance and simplicity. ScyllaDB is used for its high throughput and low latency, making it suitable for handling a    large number of TODO items efficiently.
- Pagination and Filtering: Implementing pagination and filtering ensures efficient data retrieval and enhances the user experience.