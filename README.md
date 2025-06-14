# DropLink 🔐

A modern, minimalistic and secure file sharing app with end-to-end encryption. Built for users who want privacy without complexity.

## Features ✨

- 🔐 **Zero-Knowledge** – Client-side file encryption, ensuring only the recipient can decrypt them.
- 🔑 **Password-based** – Encrypt and Decrypt files with custom password
- ⚡ **No Sign-Up Required** – Share files instantly without creating an account.
- 📁 **Cloud-agnostic (WIP)** – Self-host it or simply deploy it to a Cloud Provider.
- 🧩 **Minimalistic Design** – Clean, fast, and intuitive interface built with React and TypeScript.
- 🐳 **Containerized (WIP)** – Easily deployable backend and frontend with Docker.
- ⏳ **Expiration Support (WIP)** – Set time-based expiration for shared files.


## Components

### Backend (Golang) 🖥️
[Start Backend](./be/README.md)

### Frontend (Typescript/React/Vite) 🧑‍💻
[Start Frontend](./fe/README.md)

## File structure 🗂️
```
droplink/
├── be/                     // Backend
│   ├── Dockerfile
│   ├── handler
│   ├── main.go
│   ├── model
│   ├── storage
│   ├── uploads
│   └── utils
├── fe/                     // Frontend
│   ├── index.html
│   └── src/ 
├── infra/                  // Infra
│   ├── terraform/
│   └── k8s/   
├── LICENSE
└── README.md
```

## To-Do 🛠️
 - ☁️ **Add support for GCP and Azure**
 - ☸️ **Helmify as much as possible**

## Contribution 🤝

 1. Fork the repository
 2. Create a feature branch
 3. Commit & push to branch
 4. Open a Pull Request