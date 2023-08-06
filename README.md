# himawari-dl

[![Star](https://img.shields.io/github/stars/ej-agas/himawari-dl.svg?style=flat-square)](https://github.com/ej-agas/himawari-dl/stargazers) [![License](https://img.shields.io/github/license/ej-agas/himawari-dl.svg?style=flat-square)](https://github.com/ej-agas/himawari-dl/blob/main/LICENSE) [![Release](https://img.shields.io/github/release/ej-agas/himawari-dl.svg?style=flat-square)](https://github.com/ej-agas/himawari-dl/releases)

CLI application that downloads Himawari-8/Himawari-9 satellite images from [Meteorological Satellite Center of JMA](https://www.data.jma.go.jp/mscweb/en/index.html).

## 📖 Usage

### Help
```shell
himawari-dl help
```

### Downloading Images
```shell
himawari-dl download --dir img --parallel 5
```
### Flags
**--dir** path where to save downloaded images (default: img) \
**--parallel** number of consecutive downloads (default: 5) \