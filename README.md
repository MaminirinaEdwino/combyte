# Combyte 
Combyte is a cli tools that allow you to compress txt file encoded in UTF-8 like .txt, .sql or something else

## How to install : 
### Dependecies : 
#### Golang : 
You will to install golang programming for building the cli 

### Installation : 
1. Method 1: 
- clone the repo
- build it manually 
- and add its path into your environnement variable(it depend on your OS)

2. Method 2: 
use the following command : 
```bash
go install github.com/MaminirinaEdwino/combyte
```

### Usage : 
1. Compress a file : 
Using the option `--compress` or `-c`

```bash
combyte --compress --filename="file.txt"
```

2. Extract compressed a file : 
Using the option `--extract` or `-e`
```bash
combyte --extract --filename="file.combyte"
```




