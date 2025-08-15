# ngomap
Ngomap is a simple homemade tool for network scanning.

### ⚙️ Features:
- Scan a single host or a whole network
- Rich and user-friendly CLI
 
### Installation:
1) Clone the repository:
- ```git clone https://github.com/BetterLuckNextLife/ngomap```
2) Build the project:
- ```cd ngomap```
- ```go build```
3) Run the binary:
- ```./ngomap```

### 📁 Structure:
```
ngomap
├── cmd              # CLI commands
│   ├── network.go
│   ├── root.go
│   └── single.go
├── go.mod           # Go modules 
├── go.sum           # Go modules
├── main.go          # Calls the CLI manager
└── scanners
    ├── scanner.go   # Main scanning logic
    ├── synSender.go # Raw SYN packet sending logic
    └── utils.go     # Utilities
```

### 🤝 Contribution:
How can you contribute?

1) Fork the repository.
2) Create your feature branch (git checkout -b feature/YourFeature).
3) Commit your changes (git commit -m 'Added some features').
4) Push to the branch (git push origin feature/YourFeature).
5) Open a pull request.

If you find any bugs, please report them.

### License:
This project is licensed under the [MIT](./LICENSE) License.
Author: BetterLuckNextLife (2025).
You can use, modify, and distribute the project, as long as you mention the author.

Anyways, I had fun developing this little tool. This project really boosted my understanding of Go.
*Not a single line of code was written by AI*
