# **Triple-S (Simple Storage Service)**

## **Overview**
Triple-S is a simplified S3-like object storage service that allows users to create and manage buckets, upload, retrieve, and delete objects via a REST API. This project demonstrates the principles of RESTful API design, basic networking, and data management.

---

## **Features**
### **1. Bucket Operations**
- Create a bucket
- List all buckets
- Delete a bucket

### **2. Object Operations**
- Upload an object to a bucket
- Retrieve an object from a bucket
- Delete an object from a bucket

### **3. Metadata Management**
- Store bucket metadata in `buckets.csv`
- Store object metadata in `objects.csv` for each bucket

---

## **Project Structure**
```
triple-s
├── cmd
│   └── main.go
├── data
│   └── buckets.csv
├── flags
│   └── flags.go
├── go.mod
├── handlers
│   ├── bucket.go
│   └── object.go
├── metadata
│   └── metadata.go
├── README.md
└── utils
    └── utils.go
```

---

## **API Endpoints**

### **Bucket Operations**
| HTTP Method | Endpoint                 | Description              |
|------------|-------------------------|-------------------------|
| `PUT`       | `/{BucketName}`         | Create a new bucket     |
| `GET`       | `/`                     | List all buckets        |
| `DELETE`    | `/{BucketName}`         | Delete a bucket         |

### **Object Operations**
| HTTP Method | Endpoint                     | Description             |
|------------|------------------------------|-------------------------|
| `PUT`       | `/{BucketName}/{ObjectKey}`  | Upload an object       |
| `GET`       | `/{BucketName}/{ObjectKey}`  | Retrieve an object     |
| `DELETE`    | `/{BucketName}/{ObjectKey}`  | Delete an object       |

---

## **Installation and Setup**

### **1. Clone the Repository**
```bash
git clone git@github.com:SULTE4/triple-s.git
cd triple-s
```

### **2. Build the Project**
```bash
go build -o triple-s .
```

### **3. Run the Server**
```bash
./triple-s --port <PORT> --dir <BaseDir>
```
**Example:**
```bash
./triple-s --port 8080 --dir ./data
```

---

## **Error Handling**

| Error Scenario                 | HTTP Status Code | Response Message               |
|--------------------------------|-----------------|--------------------------------|
| Invalid bucket name            | `400`           | `Bucket name is required`      |
| Bucket already exists          | `409`           | `Bucket already exists`        |
| Bucket not found               | `404`           | `Bucket not found`             |
| Object not found               | `404`           | `Object not found`             |
| Bucket is not empty            | `409`           | `Bucket is not empty`          |
| Invalid method for endpoint    | `405`           | `Method not allowed`           |
| Internal server error          | `500`           | `Failed to perform operation`  |

---
