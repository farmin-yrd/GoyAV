# GoyAV
<div align="center">
  <img src="https://i.postimg.cc/VLwhjXm9/logo.webp" height="300">
</div>

GoyAV is a versatile REST API-based service for uploading documents and conducting antivirus scanning to ensure the security and integrity of files. It's designed to be adaptable with various antivirus analyzers, databases, and binary data repositories.

## Key Features

- **File Upload and Management**: Securely upload and manage documents.
- **Antivirus Analysis**: Perform antivirus scanning on uploaded documents.
- **Flexible Adapters**: Compatible with different antivirus analyzers, databases, and binary data repositories.

## API Specification

For detailed information about the API endpoints, parameters, and responses, refer to the API specification:

- [GoyAV API Specification](./resources/api/swagger.yml)

This link will lead to the Swagger YAML file containing the full API documentation. It's useful for understanding the available endpoints, their functionalities, and expected payloads.

## Architecture

GoyAV uses an hexagonal architecture pattern, consisting of:

- [**core**](/src/internal/core): Contains domain logic and port interfaces.
- [**adapters**](/src/internal/adapter): Facilitates interaction with various external components like antivirus services (e.g., ClamAV or others), databases (e.g., PostgreSQL, MySQL, Redis), and binary data repositories (e.g., file storage, object storage buckets like Minio).
- [**services**](/src/inernal/service): Implements business logic: this includes managing file uploads, initiating antivirus scans, and securely storing the analysis results.

## Implementation Details

The current implementation includes adapters for:
- **Antivirus Analyzer**: ClamAV (configurable for other analyzers).
- **Database**: PostgreSQL (adaptable for other databases like MySQL, Redis).
- **Binary Data Repository**: Minio (or other file/object storage solutions).


## Environment Setup

Configure the GoyAV application using the following environment variables:

### General Configuration
- `GOYAV_DEBUG_MODE`: Debug mode (true/false).
- `GOYAV_HOST`: Host for the API server.
- `GOYAV_PORT`: Port for the API server.

### File Upload Configuration
- `GOYAV_MAX_UPLOAD_SIZE`: Maximum upload size in bytes.
- `GOYAV_UPLOAD_TIMEOUT`: Upload timeout in seconds.

### Document Management
- `GOYAV_DOCUMENT_TTL`: Time-To-Live for documents in the system (format: s for seconds, m for minutes, h for hours, e.g., `2h50m10s`, `24h`, `30m`).

### Service Information
- `GOYAV_VERSION`: Version of the GoyAV service.
- `GOYAV_INFORMATION`: Additional information about the GoyAV service.

### Minio Configuration
- `GOYAV_MINIO_HOST`: Host for the Minio server.
- `GOYAV_MINIO_PORT`: Port for the Minio server.
- `GOYAV_MINIO_ACCESS_KEY`: Access key for Minio.
- `GOYAV_MINIO_SECRET_KEY`: Secret key for Minio.
- `GOYAV_MINIO_BUCKET_NAME`: Bucket name in Minio.
- `GOYAV_MINIO_USE_SSL`: Enable SSL for Minio connection (true/false).

### PostgreSQL Configuration
- `GOYAV_POSTGRES_HOST`: Host for the PostgreSQL database.
- `GOYAV_POSTGRES_PORT`: Port for the PostgreSQL database.
- `GOYAV_POSTGRES_USER`: User for the PostgreSQL database.
- `GOYAV_POSTGRES_USER_PASSWORD`: Password for the PostgreSQL user.
- `GOYAV_POSTGRES_DB`: Database name in PostgreSQL.
- `GOYAV_POSTGRES_SCHEMA`: Schema name in PostgreSQL.

### ClamAV Configuration
- `GOYAV_CLAMAV_HOST`: Host for the ClamAV service.
- `GOYAV_CLAMAV_PORT`: Port for the ClamAV service.
- `GOYAV_CLAMAV_TIMEOUT`: Timeout for ClamAV operations.


## Best Practices in GoyAV Implementation

GoyAV implements several best practices in its antivirus analysis to ensure efficiency and robustness:

1. **Asynchronous Analysis**: GoyAV performs the antivirus analysis asynchronously, allowing the upload process to complete without waiting for the analysis to finish.

2. **File Size Limit**: A predefined limit on the file size for uploads helps in managing resource usage effectively.

3. **Timeout for Antivirus Service**: A specific timeout is set for the antivirus service to ensure system responsiveness.

4. **Hash Preservation**: To optimize performance, GoyAV checks for existing hashes of files. If a file with the same hash is found, and the analysis is already complete, it reuses the existing result.

5. **Parallel Process Limitation (Semaphore)**: A semaphore mechanism limits the number of parallel analysis processes, preventing overloading and ensuring efficient resource allocation.

## Contributing: Extending GoyAV with New Adapters

The community is welcomed and encouraged  to contribute to GoyAV by developing new adapters. The hexagonal architecture of GoyAV is designed to be flexible and extendable, allowing for easy integration of various external systems and services.

### Areas for Extension

You can contribute by creating new adapters that implement the following interfaces:

- [AntivirusAnalyser](/src/internal/core/port/antivirus_analyser.go): Develop an adapter to integrate different antivirus scanning services.
- [ByteRepository](/src/internal/core/port/byte_repository.go): Create an adapter for alternative binary data storage solutions.
- [DocumentRepository](/src/internal/core/port/document_repository.go): Implement an adapter for various database systems to manage document metadata.
- [DocumentService](/src/internal/core/port/document_service.go): Enhance the application by developing additional document processing services.

### How to Contribute

1. **Fork the Repository**: Start by forking the GoyAV repository.
2. **Develop Your Adapter**: Follow the existing coding standards and project structure.
3. **Test Your Code**: Ensure your adapter works as expected and passes all tests.
4. **Submit a Pull Request**: Once you're satisfied with your adapter, submit a pull request for review.

Your contributions will help make GoyAV more versatile and robust, catering to a wider range of use cases and scenarios.
