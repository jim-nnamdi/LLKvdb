# **LLKVdb** 🌟  
LLKVdb is an acronym for low latency key value database. Filesystems are major considerations during system designs of any kind. Proper handling of files and data either through Relational database management systems, NoSQL Databases or even Key-Value stores cannot be overemphasized. LLKvdb is a network available key value store database that is built based on the design principles of _Log Structured Merge Trees_ in this project my approach conveys a wholistic and Lean structure. LLKvdb is designed to solve some of the issues traditional operating systems face in File systems modules which are high latency and high memory usages and sometimes Runtime complexities.

## **Table of Contents**  
1. [Description](#description)  
2. [Features](#features)  
3. [Prerequisites](#Prerequisites)
3. [Installation](#installation)  
4. [Usage](#usage) 
5. [Notes](#notes)

## **Description**  
**LLKvdb** is a network available key value store database that is built based on the design principles of _Log Structured Merge Trees_ which can be incorporated in any system to serve as a datastore. LLKvdb provides low latency during writes and reads operations on any system it is incorporated on. it also provides a fault tolerant system which helps to preserve data incase of unforseen circumstances or a system crash. LLKvdb pays attention to the ACID properties of database engineering to enhance data safety and durability, It also uses algorithms like Quicksort (efficient for sorting large datasets), Binary search and Log Merge trees for internal file operations.

### **Key Highlights:**  
- **Problem it solves**: High Latency and High Memory utilisation issues in mordern Database systems  
- **Where can it be used**: Embedded systems Datastores, Web systems that require any form of Datastore etc.   


## **Features**  
| Feature |  Summary | Status     |  
|-------------|-------------|------------|  
| 🌟 **Write Ahead Log**   |  Takes in data before sent to memtables      |  ✅ |  
| 📝 **In-Mem Disk Memtable**   | Holds writes data request from server      |  ✅ |  
| ⚙️ **Sorted Tables**   | Memtables Flush on threshold limit to disk     |  ✅ |  
| 🚀 **Data Recovery**   | Replay mechanism done on WAL     |  ✅ |  
| 📮 **Compaction for SSTables**   | Merge multiple sstables for Memory optimization     |  🔺 |  
| 📊 **Data Replication**   | Data is split inside sstables and also in WAL     |  ✅ |  
| 🌐 **Network Availability**   | Public endpoints exposed for end users     |  ✅ |  

## **Prerequisites**

To run this project, ensure you have the following installed and configured on your system:

### **General Requirements**
- [Go (Golang) 1.21](https://golang.org/dl/)
- Git for version control ([Download Git](https://git-scm.com/))
- Internet connection to fetch dependencies via `go mod tidy`

## **Installation**
#### **Windows**
1. **Install Go**  
   - Download the Windows installer from [Go Downloads](https://golang.org/dl/).  
   - Run the installer and follow the instructions.  
   - Ensure the Go binary (`go.exe`) is added to your `PATH`.

2. **Install Git**  
   - Download Git from [Git for Windows](https://git-scm.com/).  
   - Follow the installation wizard to complete the setup.

3. **Environment Variables**  
   - Ensure your `GOPATH` and `GOROOT` environment variables are set correctly:
     - `GOPATH` (default): `C:\Users\<YourUsername>\go`
     - `GOROOT` (default): `C:\Go`

#### **Linux**
1. **Install Go**  
   - Use your system package manager or download the tarball from [Go Downloads](https://golang.org/dl/).  
   - Example for Ubuntu/Debian:
     ```bash
     sudo apt update
     sudo apt install golang
     ```
   - Example for Fedora/RHEL:
     ```bash
     sudo dnf install golang
     ```
#### **MacOS**
```bash 
brew install go
go --version
```

## **Usage**  
```bash
# Clone this repository
git clone https://github.com/jim-nnamdi/LLKvdb.git

# Navigate to the project directory
cd LLKvdb

# Install dependencies
go mod tidy

# Build the binary
go build -o llkvdb ./cmd/kvs

# Run the binary
./llkvdb start & 

# See running Process
ps aux | grep llkvdb

# To run unit tests
go test ./... -v
```

## **Make Usage**  
```bash
# Ensure you're in the project directory already
# compiles the project  
make compile

# starts the binary as a process
make run

# runs the unit tests
make test

# removes the binary and disk files
make clean
```

## **Endpoints**
```shell

# After running the binary, you can access the endpoints via postman or insomnia too
# remember how to change the port on the **Notes** section

curl --location 'http://localhost:7009/put' \
--form 'key="9"' \
--form 'value="Samuel"'

curl --location 'http://localhost:7009/read/1'
```

## **Notes**  
Inside ```pkg/command/start.go``` file you will see the port on line 22. Now since this would be a binary running on a local system. it's fair to assume that the provided port in the project might be used or is currently in used for another running process. this value can be changed to a free port and then after that rebuild the binary and run as normal.

```bash
go build -o llkvdb ./cmd/kvs

# Run the binary
./llkvdb start & 

# See running Process
ps aux | grep llkvdb
```