![image](https://github.com/user-attachments/assets/dd3c48dd-022c-49ba-9c2e-49422f039401)

# Tiny Ping 🟢

TinyPing is a lightweight status checking service that provides simple and efficient monitoring of various service endpoints.

## Features

- Service status monitoring
- Response time (latency) tracking
- Clean and intuitive dashboard UI

## Prerequisites
- Go 1.23 or higher
- AWS DynamoDB
  - partition key : service(s)
  - sort key : timestamp(s)

## Environment Variables
Required environment variables:

```env
AWS_REGION=ap-northeast-2
DYNAMODB_TABLE_NAME=tinyping-test
```

## Getting Started

1. Clone the repository
```bash
git clone https://github.com/yourusername/tinyping.git
cd tinyping
```

2. Set environment variables
```bash
export AWS_REGION=ap-northeast-2
export DYNAMODB_TABLE_NAME=tinyping-test
```


## AWS Setup

1. Install and configure AWS CLI
2. Create DynamoDB table
3. Set up appropriate IAM permissions

## Contributing

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Acknowledgments
- Thanks to all contributors
