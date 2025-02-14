## Installation Project IWO KWIO


```bash
go run resource/install/main.go 
```

```bash
go mod tidy
```

### auto generate resource

please makse folder and json file ex: folder warga

warga/warga.json
```bash
{
  "name": "warga",
  "model": "Models/Unit.go"
}
```

```bash
go run resource/main.go warga/warga.json
```

## Installation Air

### Via `go install` (Recommended)

With go 1.23 or higher:

```bash
go install github.com/air-verse/air@latest
```


## Installation Air

### Via `go install` (Recommended)

With go 1.23 or higher:

```bash
go install github.com/air-verse/air@latest
```

configuration file to the current directory with the default settings running the following command.

```shell
air init
```
